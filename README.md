# Micro 18 Tools

This repository includes a set of tools that are useful for performing experiments for the Micro18 papers.
The tools may be applicable for other types of projects which perform workload characterization and/or use the Chrome trace format.

## Installing

### Installing Go

The tool is developed using [golang](https://golang.org/) which needs to be installed for code to be compiled from source.
You can install Golang either through [Go Version Manager](https://github.com/moovweb/gvm)(recommended) or from the instructions on the [golang site](https://golang.org/). We recommend the Go Version Manager.

The following are instruction on how to install Go 1.8 through Go version manager.
Go version 1.8+ is required to compile RAI.

Download the [GVM](https://github.com/moovweb/gvm) using

```
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
```

Add the following line to your `.bashrc`(or `.zshrc` if using zsh) to set up the GVM environment.
This is sometimes done for you by default.

```
[[ -s "$HOME/.gvm/scripts/gvm" ]] && source "$HOME/.gvm/scripts/gvm"
```

Ran your updated `.bashrc` file

```
source .bashrc
```

You can then install the Go 1.8 binary and set it as the default

```
gvm install go1.10 -B
gvm use go1.10 --default
```

`gvm` will setup both your `$GOPATH` and `$GOROOT` and you can validate that the installation completed by invoking

```sh
$ go env
GOARCH="amd64"
GOBIN=""
GOEXE=""
GOHOSTARCH="amd64"
GOHOSTOS="linux"
GOOS="linux"
GOPATH="/home/abduld/.gvm/pkgsets/go1.8/global"
GORACE=""
GOROOT="/home/abduld/.gvm/gos/go1.8"
GOTOOLDIR="/home/abduld/.gvm/gos/go1.8/pkg/tool/linux_amd64"
GCCGO="gccgo"
CC="gcc"
GOGCCFLAGS="-fPIC -m64 -pthread -fmessage-length=0 -fdebug-prefix-map=/tmp/go-build917072201=/tmp/go-build -gno-record-gcc-switches"
CXX="g++"
CGO_ENABLED="1"
PKG_CONFIG="pkg-config"
CGO_CFLAGS="-g -O2"
CGO_CPPFLAGS=""
CGO_CXXFLAGS="-g -O2"
CGO_FFLAGS="-g -O2"
CGO_LDFLAGS="-g -O2"
```

### Install the package

Navigate to where Go will expect to find the source for this repo. Make the path if it does not exist.

    mkdir -p $GOPATH/src/github.com/rai-project
    cd $GOPATH/src/github.com/rai-project

Clone this repository there.

    git clone https://github.com/rai-project/micro18-tools.git
    cd micro18-tools

Install the dependencies

    go get -v ./...

You should now be able to build the micro18-tools

    go build main.go

## Configurations

The current client looks for the config in `~/.carml_config.yml`, but this can be overridden using the `--config=<<file_path>>` option.

### Defaults

```
micro18:
  build_timeout: 0
  polling_interval: 0
  src_path: $GOPATH/src/github.com/rai-project/mxnet-mirror
  base_path: $HOME/carml/data/mxnet
  server_relative_path: bin
  server_build_cmd: make
  server_run_cmd: bin/uprd
  client_relative_path: example/image-classification/predict-cpp
  client_build_cmd: make
  client_run_cmd: ./image-classification-predict
  base_bucket_url: http://s3.amazonaws.com/micro18profiles
  upload_bucket_name: traces
  profile_output_directory: $HOME//micro18_profiles
```

## Monitoring Memory Usage

... using the `monitor_memory` option.

This option is only supported on linux and uses the nvml library.

## Trace Tools

### Combining Traces

```
go run main.go trace combine /tmp/profiles/* -o /tmp/combined.json
```

There are options that can be passed in

```
      --adjust          Adjust the timeline to ignore categories, adjust event names, and zero out the trace (default true)
  -o, --output string   Ther output path to the combined trace (default "combined.json")
      --skip_first      Skip the first input argument
```

### Uploading Traces

...

## Client

### Examples

```
go run main.go client run --model_iterations=100 --iterations=100 --concurrent=15 --distribution=pareto
```

### Basic Run

```
micro18-tools client run opts...
```

### Comparing UPR vs Original

```
micro18-tools client compare opts...
```

### Client Options

#### Specifying Profile Output Directory

```
go run main.go client run --profile_output=<<SOMETHING>>
```

You can use the `--profile_output_overwrite=true` to delete and overwrite the existing profile output directory.

#### Specifying Experiment Description

```
go run main.go client run --experiment_description="description of your experiment"
```

#### Eager Initialization

Causes the CUDA initialization code to be run eagerly.
By default, MXNet runs the CUDA initialization code lazily.
The end-to-end time for inference is unchanged by this modification, and initialization will not be shown in the output profile.

This mode is enabled using the `--eager=true` option.

#### Eager Async Initialization

Causes the CUDA initialization code to be run eagerly in asynchronous mode within a background thread.
The CUDA initialization code will be run eagily in a background thread, and a thread wait is placed before the inference beginning.
If the user code (be it preprocessing or input reading) is long enough, then the initialization can be hidden by this technique.
The end-to-end time for inference will be by this modification, and initialization might be shown in the output profile.

This mode is enabled using the `--eager_async=true` option.

#### Specifying the Distribution

The model selection process can be specified by a distribution that simulates realistic workloads.
The following distributions are supported by the tool:

* Pareto: the `xm` and `alpha` prameters can be specified
* Uniform the `min` and `max` prameters can be specified
* Exponential: the `rate` parameter can be specified
* Weibull: the `k` and `lambda` prameters can be specified
* Poisson: the `lambda` prameter can be specified

Both the distribution as well as the distribution parameters can be specified using

```
micro18-tools client ... --distribution=<<dist_name>> --distribution_params=<<<param1,param2,...>>>
```

If neither the distribution is not specified then the client is run across all models in sequence multiple times (the number of times is specified by the `--iterations` option).

Models are selected from the distribution upto `--iteration` times.
If the `model_iterations` of iterations is set to a valid other than `-1` then each model is selected from the distribution upto `model_iterations` times (this means that the workload does not truly follow the distribution, but does guarantee that the all models get selected).

#### Specifying the Number of Concurrent Requests

The number of concurrent requests can be specified using the `--concurrent` option.
By default this is set to 1.

## Server

```
micro18-tools server run -d
```

...

## Other Tools

### Downloading Models

```
micro18-tools download assets
```

This will download it to the directory specified by `micro18.base_path` in the config file.
