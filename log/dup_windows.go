// +build windows

package log

import (
	"log"
	"os"
	"syscall"
)

var (
	kernel32         = syscall.MustLoadDLL("kernel32.dll")
	procSetStdHandle = kernel32.MustFindProc("SetStdHandle")
)

func Dup2File(file *os.File, fd int) error {
	stdHandle := syscall.STD_ERROR_HANDLE
	if fd == int(os.Stdout.Fd()) {
		stdHandle = syscall.STD_OUTPUT_HANDLE
	}
	err := setStdHandle(stdHandle, syscall.Handle(file.Fd()))
	if err != nil {
		return err
	}
	if fd == int(os.Stdout.Fd()) {
		os.Stdout = file
	} else {
		os.Stderr = file
		log.SetOutput(os.Stderr)
		logger.SetOutput(os.Stderr)
	}
	return nil
}

func setStdHandle(stdhandle int, handle syscall.Handle) error {
	r0, _, e1 := syscall.Syscall(procSetStdHandle.Addr(), 2, uintptr(stdhandle), uintptr(handle), 0)
	if r0 == 0 {
		if e1 != 0 {
			return error(e1)
		}
		return syscall.EINVAL
	}
	return nil
}
