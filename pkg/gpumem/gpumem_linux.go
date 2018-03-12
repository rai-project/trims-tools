// +build linux

package gpumem

import (
	"sync"
	"errors"
	"time"

	nvml "github.com/rai-project/nvml-go"
)

type Entry struct {
	timestamp time.Time
	MemoryUsed      int64
	MemoryFree      int64
}

type Device struct {
  mut sync.Mutex
  index int
  handle  nvml.DeviceHandle
	entries []*Entry
}

type Memory struct {
	done     chan bool
	started  bool
	devices  []*Device
}

func New() (*Memory, error) {
	devs, err := nvml.DeviceCount()
	if err != nil {
		return nil, err
	}
	devices := make([]Device, devs)
	for ii := range devices {
		handle, err := nvml.DeviceGetHandleByIndex(ii)
		if err != nil {
			return nil, err
		}
		devices[ii] = &Device{
      index: ii,
      handle: handle,
    }
	}
	return  &Memory{
		devices: devices,
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

func (*Memory) Stop() {
	if m.started == false {
		return
	}
	close(done)
}

func (*Memory) Print() {
  fmt.Println("todo.. implement the memory print function")
}

func (*Memory) Show(fmt string) {
  fmt.Println("todo.. implement the memory show function")
}

func (dev *DeviceMemory) recordInfo() {
  timestamp :=time.Now()
  info, err := nvml.DeviceMemoryInformation(dev.handle)
  if err != nil {
    return
  }
  dev.mut.Lock()
  defer dev.mut.Unlock()
  dev.entries = append(
    dev.entries,
    Entry{
      timestamp: timestamp,
      MemoryUsed: info.Used,
      MemoryFree: info.Free,
    }
  )
}

func init() {
  err := nvml.Init()
  if err != nil {
    panic("failed to initialize nvml " + err.Error())
  }
}
