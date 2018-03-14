// +build linux

package gpuinfo

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/rai-project/cudainfo"
	mconfig "github.com/rai-project/micro18-tools/pkg/config"
	nvml "github.com/rai-project/nvml-go"
	"github.com/spf13/cast"
	yaml "gopkg.in/yaml.v2"
)

type Device struct {
	mut     sync.Mutex
	index   uint
	handle  *cudainfo.NVMLDevice
	entries []cudainfo.DeviceStatus
}

type System struct {
	done       chan bool
	started    bool
	fmt        string
	fullOutput bool
	output     io.Writer
	devices    []*Device
}

func isVisible(ii uint) bool {
	visibleDevices := mconfig.Config.VisibleDevices
	for _, dev := range strings.Split(visibleDevices, ",") {
		idev := cast.ToUint(dev)
		if ii == idev {
			return true
		}
	}
	return false
}

func New() (*System, error) {
	devs, err := cudainfo.GetDeviceCount()
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get nvml device count")
	}
	devices := []*Device{}
	for ii := uint(0); ii < devs; ii++ {
		if !isVisible(ii) {
			continue
		}
		handle, err := cudainfo.NewNvmlDevice(uint(ii))
		if err != nil {
			return nil, errors.Wrapf(err, "cannot get device handle for %d", ii)
		}
		devices = append(devices, &Device{
			index:  ii,
			handle: handle,
		})
	}
	return &System{
		done:    make(chan bool),
		started: false,
		devices: devices,
		output:  os.Stdout,
	}, nil
}

func (m *System) Start(timestep time.Duration) error {
	if m.started == true {
		return errors.New("memory info already started")
	}
	ticker := time.NewTicker(timestep)
	go func() {
		for {
			select {
			case <-ticker.C:
				go func() {
					for _, dev := range m.devices {
						dev.recordInfo()
					}
				}()
			case <-m.done:
				ticker.Stop()
				break
			}
		}
	}()
	m.started = true
	return nil
}

func (m *System) Stop() {
	if m.started == false {
		return
	}
	close(m.done)
}

func (m *System) Print() {
	m.Write("table", os.Stdout)
}

func (m *System) SetOutputFormat(fmt string) {
	m.fmt = fmt
}

func (m *System) SetOutput(output io.Writer) {
	m.output = output
}

func (m *System) Write(fmt string, output io.Writer) error {
	fmt = strings.ToLower(fmt)

	m.SetOutputFormat(fmt)
	m.SetOutput(output)

	switch fmt {
	case "csv", "tsv":
		return m.writeCSV(fmt)
	case "table":
		return m.writeTable()
	case "json":
		bts, err := json.Marshal(m)
		if err != nil {
			return errors.Wrap(err, "failed to serialize gpu memory information to json")
		}
		_, err = output.Write(bts)
		return err
	case "yaml":
		bts, err := yaml.Marshal(m)
		if err != nil {
			return errors.Wrap(err, "failed to serialize gpu memory information to json")
		}
		_, err = output.Write(bts)
		return err
	}
	log.Errorf("the format %s is not a valid output format for gpu memory information", fmt)
	return errors.Errorf("the format %s is not a valid output format for gpu memory information", fmt)
}

func (m *System) dsvHeader() []string {
	return []string{
		"device_idx",
    "time_stamp",
    "power",
    "temperature",
    "utilization"
		"memory_used",
		"human_memory_used",
    "clock_core",
    "clock_memory",
    "pci_throughput_rx",
    "pci_throughput_tx",
    "num_processes",
	}
}

func (m *System) dsvRows() [][]string {
	rowDivider := []string{
		"---",
		"---",
		"---",
		"---",
		"---",
    "---",
    "---",
		"---",
		"---",
		"---",
		"---",
		"---",
	}
	fullOutput := m.fullOutput
  rows := [][]string{}
  totalPower :=  []uint64{}
  totalTemperature :=  []uint64{}
  totalUtilization := []uint64{}
  totalMemoryUsed := []uint64{}
  totalClockCore :=  []uint64{}
  totalClockMemory :=  []uint64{}
  totalPCIThroughputRX := []uint64{}
	totalPCIThroughputTX := []uint64{}
	totalNumProcesses := []uint64{}
	totalEntries := []uint64{}
	for _, dev := range m.devices {
    currTotalPower :=  uint64(0)
    currTotalTemperature :=  uint64(0)
    currTotalUtilization := uint64(0)
    currTotalMemoryUsed := uint64(0)
    currTotalClockCore := uint64(0)
    currTotalClockMemory := uint64(0)
    currTotalPCIThroughputRX := uint64(0)
    currTotalPCIThroughputTX := uint64(0)
    currTotalNumProcesses := uint64(0)

    devIdx := cast.ToString(dev.index)
		for _, entry := range dev.entries {
      power = entry.Power
      temperature = entry.Temperature
      utilization = entry.UtilizationInfo.GPU
      memoryUsed = entry.nvmlMemoryInfo.GlobalUsed
      clockCore = entry.ClockInfo.Core
      clockMemory = entry.ClockInfo.Memory
      PCIThroughputRX += entry.PCIStatusInfo.Throughput.RX
      PCIThroughputTX += entry.PCIStatusInfo.Throughput.TX
      numProcesses += len(entry.Processes)

      currTotalPower += power
      currTotalTemperature += temperature
      currTotalUtilization += utilization
      currTotalMemoryUsed += memoryUsed
      currTotalClockCore += clockCore
      currTotalClockMemory += clockMemory
      currTotalPCIThroughputRX += PCIThroughputRX
      currTotalPCIThroughputTX += PCIThroughputTX
      currTotalNumProcesses += numProcesses

			if fullOutput {
				rows = append(
					rows,
					[]string{
						devIdx,
            entry.TimeStamp.Format(time.RFC3339Nano),
            cast.ToString(power)
            cast.ToString(temperature)
            cast.ToString(utilization)
            cast.ToString(memoryUsed)
            humanize.Bytes(memoryUsed),
            cast.ToString(clockCore)
            cast.ToString(clockMemory)
            cast.ToString(PCIThroughputRX)
            cast.ToString(PCIThroughputTX)
            cast.ToString(numProcesses)
						cast.ToString(memoryUsed),
					},
				)
			}
		}
		totalPower = append(totalPower, currTotalPower)
		totalTemperature = append(totalTemperature, currTotalTemperature)
		totalUtilization = append(totalUtilization, currTotalUtilization)
		totalMemoryUsed = append(totalMemoryUsed, currTotalMemoryUsed)
		totalClockCore = append(totalClockCore, currTotalClockCore)
		totalClockMemory = append(totalClockMemory, currTotalClockMemory)
		totalPCIThroughputRX = append(totalPCIThroughputRX, currTotalPCIThroughputRX)
		totalPCIThroughputTX = append(totalPCIThroughputTX, currTotalPCIThroughputTX)
		totalNumProcesses = append(totalNumProcesses, currTotalNumProcesses)
    totalEntries = append(totalEntries, uint64(len(dev.entries)))

	}
	// if fullOutput {
	// 	rows = append(
	// 		rows,
	// 		rowDivider,
	// 	)
	// 	for ii, dev := range m.devices {
	// 		devIdx := cast.ToString(dev.index)
	// 		rows = append(
	// 			rows,
	// 			[]string{
	// 				devIdx,
	// 				"peak",
	// 			},
	// 		)
	// 	}
	// 	rows = append(
	// 		rows,
	// 		rowDivider,
	// 	)
	// }
	for ii, dev := range m.devices {
		devIdx := cast.ToString(dev.index)
		averagePower := uint64(float64(totalPower[ii]) / float64(totalEntries[ii]))
		averageTemperature := uint64(float64(totalTemperature[ii]) / float64(totalEntries[ii]))
		averageUtilization := uint64(float64(totalUtilization[ii]) / float64(totalEntries[ii]))
		averageMemoryUsed := uint64(float64(totalMemoryUsed[ii]) / float64(totalEntries[ii]))
		averageClockCore := uint64(float64(totalClockCore[ii]) / float64(totalEntries[ii]))
		averageClockMemory := uint64(float64(totalClockMemory[ii]) / float64(totalEntries[ii]))
		averagePCIThroughputRX := uint64(float64(totalPCIThroughputRX[ii]) / float64(totalEntries[ii]))
		averagePCIThroughputTX := uint64(float64(totalPCIThroughputTX[ii]) / float64(totalEntries[ii]))
		averageNumProcesses := uint64(float64(totalNumProcesses[ii]) / float64(totalEntries[ii]))
		rows = append(
			rows,
			[]string{
				devIdx,
				"average",
				cast.ToString(averagePower),
				cast.ToString(averageTemperature),
				cast.ToString(averageUtilization),
        cast.ToString(averageMemoryUsed),
        humanize.Bytes(averageMemoryUsed),
				cast.ToString(averageClockCore),
				cast.ToString(averageClockMemory),
				cast.ToString(averagePCIThroughputRX),
				cast.ToString(averagePCIThroughputTX),
				cast.ToString(averageNumProcesses),
			},
		)
	}
	return rows
}

func (m *System) writeCSV(fmt string) error {
	w := csv.NewWriter(m.output)
	if fmt == "tsv" {
		w.Comma = '\t'
	}
	w.Write(m.dsvHeader())
	w.WriteAll(m.dsvRows())
	w.Flush()
	return nil
}

func (m *System) writeTable() error {
	w := tablewriter.NewWriter(m.output)
	w.SetHeader(m.dsvHeader())
	w.AppendBulk(m.dsvRows())
	w.Render()
	return nil
}

func (dev *Device) recordInfo() {
	info, err := dev.handle.Status()
	if err != nil {
		log.WithError(err).Error("failed to get device status")
		return
	}
	dev.mut.Lock()
	defer dev.mut.Unlock()
	dev.entries = append(
		dev.entries,
		*info,
	)
}

func init() {
	err := nvml.Init()
	if err != nil {
		panic("failed to initialize nvml " + err.Error())
	}
}
