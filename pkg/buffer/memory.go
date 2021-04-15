package buffer

import (
	"runtime"
	"syscall"
)

type Memory struct {
	TotalRam 	uint
	FreeRam 	uint
	AllocRam 	uint
}

func (m *Memory) Allocate() {

}

func MemStats() *Memory {
	in := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(in)
	if err != nil {
		panic(err)
	}

	m := new(Memory)
	m.TotalRam = uint(in.Totalram) * uint(in.Unit)
	m.FreeRam = uint(in.Freeram) * uint(in.Unit)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	m.AllocRam = uint(memStats.Alloc)

	return m
}