// +build linux

package gpumem

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	nvml "github.com/rai-project/nvml-go"
	"github.com/spf13/cast"
	yaml "gopkg.in/yaml.v2"
)

type Entry struct {
	timestamp  time.Time
	MemoryUsed uint64
	MemoryFree uint64
}

type Device struct {
	mut     sync.Mutex
	index   int
	handle  nvml.DeviceHandle
	entries []Entry
}

type Memory struct {
	done    chan bool
	started bool
	fmt     string
	output  io.Writer
	devices []*Device
}

func New() (*Memory, error) {
	devs, err := nvml.DeviceCount()
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get nvml device count")
	}
	devices := make([]*Device, devs)
	for ii := range devices {
		handle, err := nvml.DeviceGetHandleByIndex(ii)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot get device handle for %d", ii)
		}
		devices[ii] = &Device{
			index:  ii,
			handle: handle,
		}
	}
	return &Memory{
		devices: devices,
		output:  os.Stdout,
	}, nil
}

func (m *Memory) Start(timestep time.Duration) error {
	if m.started == true {
		return errors.New("memory info already started")
	}
	m.started = true
	ticker := time.NewTicker(timestep)
	for {
		select {
		case <-ticker.C:
			for _, dev := range m.devices {
				dev.recordInfo()
			}
		case <-m.done:
			ticker.Stop()
			break
		}
	}
	return nil
}

func (m *Memory) Stop() {
	if m.started == false {
		return
	}
	close(m.done)
}

func (m *Memory) Print() {
	m.Write("table", os.Stdout)
}

func (m *Memory) SetOutputFormat(fmt string) {
	m.fmt = fmt
}

func (m *Memory) SetOutput(output io.Writer) {
	m.output = output
}

func (m *Memory) Write(fmt string, output io.Writer) error {
	fmt = strings.ToLower(m.fmt)

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
	return errors.Errorf("the format %s is not a valid output format for gpu memory information", fmt)
}

func (m *Memory) dsvHeader() []string {
	firstDevice := m.devices[0]
	header := []string{"device_id"}
	for ii := range firstDevice.entries {
		header = append(
			header,
			fmt.Sprintf("timestamp_%d", ii),
			fmt.Sprintf("memory_used_%d", ii),
			fmt.Sprintf("memory_free_%d", ii),
		)
	}
	return header
}

func (m *Memory) dsvRows() [][]string {
	rows := [][]string{}
	for _, dev := range m.devices {
		row := []string{
			strconv.Itoa(dev.index),
		}
		for _, entry := range dev.entries {
			row = append(
				row,
				entry.timestamp.Format(time.RFC3339),
				cast.ToString(entry.MemoryUsed),
				cast.ToString(entry.MemoryFree),
			)
		}
		rows = append(rows, row)
	}
	return rows
}

func (m *Memory) writeCSV(fmt string) error {
	w := csv.NewWriter(m.output)
	if fmt == "tsv" {
		w.Comma = '\t'
	}
	w.Write(m.dsvHeader())
	w.WriteAll(m.dsvRows())
	w.Flush()
	return nil
}

func (m *Memory) writeTable() error {
	w := tablewriter.NewWriter(m.output)
	w.SetHeader(m.dsvHeader())
	w.AppendBulk(m.dsvRows())
	w.Render()
	return nil
}

func (dev *Device) recordInfo() {
	timestamp := time.Now()
	info, err := nvml.DeviceMemoryInformation(dev.handle)
	if err != nil {
		return
	}
	dev.mut.Lock()
	defer dev.mut.Unlock()
	dev.entries = append(
		dev.entries,
		Entry{
			timestamp:  timestamp,
			MemoryUsed: info.Used,
			MemoryFree: info.Free,
		},
	)
}

func init() {
	err := nvml.Init()
	if err != nil {
		panic("failed to initialize nvml " + err.Error())
	}
}
