# Micro 18 Tools

This repository includes a set of tools that are useful for performing experiments for the Micro18 papers.
The tools may be applicable for other types of projects which perform workload characterization and/or use the Chrome trace format.

## Config

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

...

### Uploading Traces

...

## Client

### Basic Run

```
micro18-tools client run opts...
```

### Comparing UPR vs Original

```
micro18-tools client compare opts...
```

### Client Options

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

...

## Other Tools

### Downloading Models

```
micro18-tools download assets
```

This will download it to the directory specified by `micro18.base_path` in the config file.
