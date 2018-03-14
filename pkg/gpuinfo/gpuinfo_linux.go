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
	index   uint64
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

func isVisible(ii uint64) bool {
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
	for ii := uint64(0); ii < devs; ii++ {
		if !isVisible(ii) {
			continue
		}
		handle, err := cudainfo.NewNvmlDevice(uint64(ii))
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
		"gpu_utilization",
		"memory_utilization",
		"memory_used",
		"human_memory_used",
		"memory_free",
		"human_memory_free",
		"clock_core",
		"clock_memory",
		"pci_throughput_rx",
		"pci_throughput_tx",
		"num_processes",
	}
}

func max(a, b uint64) uint64 {
	if a <= b {
		return b
	}
	return a
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
		"---",
		"---",
		"---",
	}

	fullOutput := m.fullOutput
	rows := [][]string{}
	totalEntries := []uint64{}

	totalPower := []uint64{}
	totalTemperature := []uint64{}
	totalGPUUtilization := []uint64{}
	totalMemoryUtilization := []uint64{}
	totalMemoryUsed := []uint64{}
	totalMemoryFree := []uint64{}
	totalClockCore := []uint64{}
	totalClockMemory := []uint64{}
	totalPCIThroughputRX := []uint64{}
	totalPCIThroughputTX := []uint64{}
	totalNumProcesses := []uint64{}

	peakPower := []uint64{}
	peakTemperature := []uint64{}
	peakGPUUtilization := []uint64{}
	peakMemoryUtilization := []uint64{}
	peakMemoryUsed := []uint64{}
	peakMemoryFree := []uint64{}
	peakClockCore := []uint64{}
	peakClockMemory := []uint64{}
	peakPCIThroughputRX := []uint64{}
	peakPCIThroughputTX := []uint64{}
	peakNumProcesses := []uint64{}

	for _, dev := range m.devices {
		currTotalPower := uint64(0)
		currTotalTemperature := uint64(0)
		currTotalGPUUtilization := uint64(0)
		currTotalMemoryUtilization := uint64(0)
		currTotalMemoryUsed := uint64(0)
		currTotalMemoryFree := uint64(0)
		currTotalClockCore := uint64(0)
		currTotalClockMemory := uint64(0)
		currTotalPCIThroughputRX := uint64(0)
		currTotalPCIThroughputTX := uint64(0)
		currTotalNumProcesses := uint64(0)

		currPeakPower := uint64(0)
		currPeakTemperature := uint64(0)
		currPeakGPUUtilization := uint64(0)
		currPeaklMemoryUtilization := uint64(0)
		currPeakMemoryUsed := uint64(0)
		currPeakMemoryFree := uint64(0)
		currPeakClockCore := uint64(0)
		currPeakClockMemory := uint64(0)
		currPeakPCIThroughputRX := uint64(0)
		currPeakPCIThroughputTX := uint64(0)
		currPeakNumProcesses := uint64(0)

		devIdx := cast.ToString(dev.index)
		for _, entry := range dev.entries {
			power := uint64(entry.Power)
			temperature := uint64(entry.Temperature)
			gpuUtilization := uint64(entry.Utilization.GPU)
			memoryUtilization := uint64(entry.Utilization.Memory)
			memoryUsed := uint64(entry.Memory.Used)
			memoryFree := uint64(entry.Memory.Free)
			clockCore := uint64(entry.Clocks.Core)
			clockMemory := uint64(entry.Clocks.Memory)
			PCIThroughputRX := uint64(entry.PCI.Throughput.RX)
			PCIThroughputTX := uint64(entry.PCI.Throughput.TX)
			numProcesses := uint64(len(entry.Processes))

			currTotalPower += power
			currTotalTemperature += temperature
			currTotalGPUUtilization += gpuUtilization
			currTotalMemoryUtilization += memoryUtilization
			currTotalMemoryUsed += memoryUsed
			currTotalMemoryFree += memoryFree
			currTotalClockCore += clockCore
			currTotalClockMemory += clockMemory
			currTotalPCIThroughputRX += PCIThroughputRX
			currTotalPCIThroughputTX += PCIThroughputTX
			currTotalNumProcesses += numProcesses

			currPeakPower = max(currPeakPower, power)
			currPeakTemperature += max(currPeakTemperature, temperature)
			currPeakGPUUtilization += max(currPeakGPUUtilization, gpuUtilization)
			currPeakMemoryUtilization += max(currPeakMemoryUtilizationm, memoryUtilization)
			currPeakMemoryUsed += max(currPeakMemoryUsed, memoryUsed)
			currPeakMemoryFree += max(currPeakMemoryFree, memoryFree)
			currPeakClockCore += max(currPeakClockCore, clockCore)
			currPeakClockMemory += max(currPeakClockMemory, clockMemory)
			currPeakPCIThroughputRX += max(currPeakPCIThroughputRX, PCIThroughputRX)
			currPeakPCIThroughputTX += max(currPeakPCIThroughputTX, PCIThroughputTX)
			currPeakNumProcesses += max(currPeakNumProcesses, numProcesses)

			if fullOutput {
				rows = append(
					rows,
					[]string{
						devIdx,
						entry.TimeStamp.Format(time.RFC3339Nano),
						cast.ToString(power),
						cast.ToString(temperature),
						cast.ToString(gpuUtilization),
						cast.ToString(memoryUtilization),
						cast.ToString(memoryUsed),
						humanize.Bytes(memoryUsed),
						cast.ToString(memoryFree),
						humanize.Bytes(memoryFree),
						cast.ToString(clockCore),
						cast.ToString(clockMemory),
						cast.ToString(PCIThroughputRX),
						cast.ToString(PCIThroughputTX),
						cast.ToString(numProcesses),
					},
				)
			}
		}
		totalPower = append(totalPower, currTotalPower)
		totalTemperature = append(totalTemperature, currTotalTemperature)
		totalGPUUtilization = append(totalGPUUtilization, currTotalGPUUtilization)
		totalMemoryUtilization = append(totalMemoryUtilization, currTotalMemoryUtilization)
		totalMemoryUsed = append(totalMemoryUsed, currTotalMemoryUsed)
		totalMemoryFree = append(totalMemoryFree, currTotalMemoryFree)
		totalClockCore = append(totalClockCore, currTotalClockCore)
		totalClockMemory = append(totalClockMemory, currTotalClockMemory)
		totalPCIThroughputRX = append(totalPCIThroughputRX, currTotalPCIThroughputRX)
		totalPCIThroughputTX = append(totalPCIThroughputTX, currTotalPCIThroughputTX)
		totalNumProcesses = append(totalNumProcesses, currTotalNumProcesses)
		totalEntries = append(totalEntries, uint64(len(dev.entries)))

		peakPower = append(peakPower, currPeakPower)
		peakTemperature = append(peakTemperature, currPeakTemperature)
		peakGPUUtilization = append(peakGPUUtilization, currPeakGPUUtilization)
		peakMemoryUtilization = append(peakMemoryUtilization, currPeakMemoryUtilization)
		peakMemoryUsed = append(peakMemoryUsed, currPeakMemoryUsed)
		peakMemoryFree = append(peakMemoryFree, currPeakMemoryFree)
		peakClockCore = append(peakClockCore, currPeakClockCore)
		peakClockMemory = append(peakClockMemory, currPeakClockMemory)
		peakPCIThroughputRX = append(peakPCIThroughputRX, currPeakPCIThroughputRX)
		peakPCIThroughputTX = append(peakPCIThroughputTX, currPeakPCIThroughputTX)
		peakNumProcesses = append(peakNumProcesses, currPeakNumProcesses)
	}
	rows = append(
		rows,
		rowDivider,
	)

	for ii, dev := range m.devices {
		devIdx := cast.ToString(dev.index)
		rows = append(
			rows,
			[]string{
				devIdx,
				"peak",
				cast.ToString(peakPower[ii]),
				cast.ToString(averageTemperature[ii]),
				cast.ToString(averageGPUUtilization[ii]),
				cast.ToString(averageMemoryUtilization[ii]),
				cast.ToString(averageMemoryUsed[ii]),
				humanize.Bytes(averageMemoryUsed[ii]),
				cast.ToString(averageMemoryFree[ii]),
				humanize.Bytes(averageMemoryFree),
				cast.ToString(averageClockCore[ii]),
				cast.ToString(averageClockMemory[ii]),
				cast.ToString(averagePCIThroughputRX[ii]),
				cast.ToString(averagePCIThroughputTX[ii]),
				cast.ToString(averageNumProcesses[ii]),
			},
		)
	}

	rows = append(
		rows,
		rowDivider,
	)

	for ii, dev := range m.devices {
		devIdx := cast.ToString(dev.index)
		averagePower := uint64(float64(totalPower[ii]) / float64(totalEntries[ii]))
		averageTemperature := uint64(float64(totalTemperature[ii]) / float64(totalEntries[ii]))
		averageGPUUtilization := uint64(float64(totalGPUUtilization[ii]) / float64(totalEntries[ii]))
		averageMemoryUtilization := uint64(float64(totalMemoryUtilization[ii]) / float64(totalEntries[ii]))
		averageMemoryUsed := uint64(float64(totalMemoryUsed[ii]) / float64(totalEntries[ii]))
		averageMemoryFree := uint64(float64(totalMemoryFree[ii]) / float64(totalEntries[ii]))
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
				cast.ToString(averageGPUUtilization),
				cast.ToString(averageMemoryUtilization),
				cast.ToString(averageMemoryUsed),
				humanize.Bytes(averageMemoryUsed),
				cast.ToString(averageMemoryFree),
				humanize.Bytes(averageMemoryFree),
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
