package watcher

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

var (
	modpsapi                     = windows.NewLazySystemDLL("psapi.dll")
	procGetProcessMemoryInfo     = modpsapi.NewProc("GetProcessMemoryInfo")
)
type SYSTEM_TIMES struct {
	CreateTime syscall.Filetime
	ExitTime   syscall.Filetime
	KernelTime syscall.Filetime
	UserTime   syscall.Filetime
}

type PROCESS_MEMORY_COUNTERS struct {
	CB                         uint32
	PageFaultCount             uint32
	PeakWorkingSetSize         uint64
	WorkingSetSize             uint64
	QuotaPeakPagedPoolUsage    uint64
	QuotaPagedPoolUsage        uint64
	QuotaPeakNonPagedPoolUsage uint64
	QuotaNonPagedPoolUsage     uint64
	PagefileUsage              uint64
	PeakPagefileUsage          uint64
}

func getProcessRow() (uint64, uint64, uint64, uint64) {
	h, err := syscall.GetCurrentProcess()
	if err != nil {
		return 0, 0, 0, 0
	}
	defer syscall.CloseHandle(h)

	userTime, kernelTime := getProcessCPURow(h)
	residentSize, totalSize := getProcessMemoryRow(h)
	return userTime, kernelTime, residentSize, totalSize
}

func getProcessCPURow(h syscall.Handle) (uint64, uint64) {
	var times SYSTEM_TIMES
	err := syscall.GetProcessTimes(
		h,
		&times.CreateTime,
		&times.ExitTime,
		&times.KernelTime,
		&times.UserTime,
	)
	if err != nil {
		return 0, 0
	}
	userTime := (uint64(times.UserTime.HighDateTime) << 32 | uint64(times.UserTime.LowDateTime)/10)/1000
	kernelTime := (uint64(times.KernelTime.HighDateTime) << 32 | uint64(times.KernelTime.LowDateTime)/10)/1000
	return userTime, kernelTime
}

func getProcessMemoryRow(h syscall.Handle) (uint64, uint64) {
	var mem PROCESS_MEMORY_COUNTERS
	r1, _, err := procGetProcessMemoryInfo.Call(uintptr(h), uintptr(unsafe.Pointer(&mem)), uintptr(unsafe.Pointer(&mem.CB)))
	if r1 == 0 {
		if err != nil {
			return 0, 0
		}
	}
	return mem.WorkingSetSize, mem.PagefileUsage
}