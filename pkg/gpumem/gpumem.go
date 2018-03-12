// +build !linux

package gpumem

import (
	"errors"
	"time"
)

var (
	UnsupportedError = errors.New("unsupported device. nvml is only supported on linux hosts")
)

type Memory struct {
}

func New() (*Memory, error) {
	return nil, UnsupportedError
}

func (m *Memory) Start(timestep time.Duration) error {
	return UnsupportedError
}

func (*Memory) Stop() {
}

func (*Memory) Print() {
}

func (*Memory) Show(fmt string) {
}
