package win32

import (
	"syscall"
)

var (
	libUser32            = syscall.NewLazyDLL("user32.dll")
	systemParametersInfo = libUser32.NewProc("SystemParametersInfoW")
	libGetMonitorInfoW   = libUser32.NewProc("GetMonitorInfoW")
)
