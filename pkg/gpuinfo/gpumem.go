// +build !linux

package gpuinfo

import (
	"errors"
	"io"
	"time"
)

var (
	UnsupportedError = errors.New("unsupported device. nvml is only supported on linux hosts")
)

type System struct {
}

func New() (*System, error) {
	return nil, UnsupportedError
}

func (m *System) Start(timestep time.Duration) error {
	return UnsupportedError
}

func (*System) Stop() {
}

func (*System) Print() {
}

func (*System) Write(fmt string, output io.Writer) error {
	return nil
}
