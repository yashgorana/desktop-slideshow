package main

import (
	"syscall"
	"unsafe"
)

const (
	SPI_SETDESKWALLPAPER = 0x0014
	SPIF_UPDATEINIFILE   = 0x01
	SPIF_SENDCHANGE      = 0x02
)

var (
	libUser32            = syscall.NewLazyDLL("user32.dll")
	systemParametersInfo = libUser32.NewProc("SystemParametersInfoW")
)

// SetFromFile sets the wallpaper for the current user.
func SetWallpaperFromFile(filename string) error {
	filenamePtr, err := syscall.UTF16PtrFromString(filename)
	if err != nil {
		return err
	}

	systemParametersInfo.Call(
		uintptr(SPI_SETDESKWALLPAPER),
		uintptr(0),
		uintptr(unsafe.Pointer(filenamePtr)),
		uintptr(SPIF_UPDATEINIFILE|SPIF_SENDCHANGE),
	)
	return nil
}

