// +build linux

package gpuinfo

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/rai-project/cudainfo"
	nvml "github.com/rai-project/nvml-go"
	"github.com/spf13/cast"
	yaml "gopkg.in/yaml.v2"
)

type Device struct {
	mut     sync.Mutex
	index   int
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

func New() (*System, error) {
	devs, err := cudainfo.GetDeviceCount()
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get nvml device count")
	}
	devices := make([]*Device, devs)
	for ii := range devices {
		handle, err := cudainfo.NewNvmlDevice(ii)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot get device handle for %d", ii)
		}
		devices[ii] = &Device{
			index:  ii,
			handle: handle,
		}
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
		"memory_used",
		"human_memory_used",
		"memory_free",
		"human_memory_free",
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
	}
	fullOutput := m.fullOutput
	rows := [][]string{}
	totalUsed := []uint64{}
	totalFree := []uint64{}
	totalEntries := []uint64{}
	for _, dev := range m.devices {
		currTotalUsed := uint64(0)
		currTotalFree := uint64(0)
		devIdx := strconv.Itoa(dev.index)
		for _, entry := range dev.entries {
			mem := entry.Memory
			memoryUsed := mem.Used
			memoryFree := mem.Free
			currTotalUsed += memoryUsed
			currTotalFree += memoryFree
			if fullOutput {
				rows = append(
					rows,
					[]string{
						devIdx,
						entry.timestamp.Format(time.RFC3339Nano),
						cast.ToString(memoryUsed),
						humanize.Bytes(memoryUsed),
						cast.ToString(memoryFree),
						humanize.Bytes(memoryFree),
					},
				)
			}
		}
		totalUsed = append(totalUsed, currTotalUsed)
		totalFree = append(totalFree, currTotalFree)
		totalEntries = append(totalEntries, uint64(len(dev.entries)))
	}
	if fullOutput {
		rows = append(
			rows,
			rowDivider,
		)
		for ii, dev := range m.devices {
			devIdx := strconv.Itoa(dev.index)
			rows = append(
				rows,
				[]string{
					devIdx,
					"total",
					cast.ToString(totalUsed[ii]),
					humanize.Bytes(totalUsed[ii]),
					cast.ToString(totalFree[ii]),
					humanize.Bytes(totalFree[ii]),
				},
			)
		}
		rows = append(
			rows,
			rowDivider,
		)
	}
	for ii, dev := range m.devices {
		devIdx := strconv.Itoa(dev.index)
		averageUsed := uint64(float64(totalUsed[ii]) / float64(totalEntries[ii]))
		averageFree := uint64(float64(totalFree[ii]) / float64(totalEntries[ii]))
		rows = append(
			rows,
			[]string{
				devIdx,
				"average",
				cast.ToString(averageUsed),
				humanize.Bytes(averageUsed),
				cast.ToString(averageFree),
				humanize.Bytes(averageFree),
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
	timestamp := time.Now()
	info, err := dev.handle.Status()
	if err != nil {
		log.WithError(err).Error("failed to get device memory information")
		return
	}

	meminfo, err := nvml.DeviceMemoryInformation(dev.handle)
	if err != nil {
		log.WithError(err).Error("failed to get device memory information")
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
