package buffer

import (
	"runtime"
	"syscall"
)

func loadSysTotalMemory() (uint, error) {
	in := &syscall.Sysinfo_t{}

	err := syscall.Sysinfo(in)
	if err != nil {
		return 0, err
	}

	return uint(in.Totalram) * uint(in.Unit), err
}

func loadSysMemoryAlloc() uint {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return uint(memStats.Alloc)
}