package main

import (
	"syscall"
	"unsafe"

	"github.com/gonutz/w32/v2"
	log "github.com/sirupsen/logrus"
)

const (
	SPI_SETDESKWALLPAPER = 0x0014
	SPIF_UPDATEINIFILE   = 0x01
	SPIF_SENDCHANGE      = 0x02
)

var (
	libUser32            = syscall.NewLazyDLL("user32.dll")
	systemParametersInfo = libUser32.NewProc("SystemParametersInfoW")
	libGetMonitorInfoW   = libUser32.NewProc("GetMonitorInfoW")
)

type Display struct {
	WidthPx  uint32
	HeightPx uint32
}

type ConnectedDisplays struct {
	count    uint8
	displays []Display
}

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

// GetConnectedDisplays fetches the connected displays and their current settings (currently limited to screen resolution)
func GetConnectedDisplays() *ConnectedDisplays {
	monitors := &ConnectedDisplays{}

	result := w32.EnumDisplayMonitors(w32.HDC(0), nil, syscall.NewCallback(enumDisplayMonitorsCb), uintptr(unsafe.Pointer(monitors)))
	log.Debug("EnumDisplayMonitors result: ", result)

	return monitors
}

// MONITORENUMPROC - Must return uintptr
func enumDisplayMonitorsCb(hMonitor w32.HMONITOR, hdc w32.HDC, lpRect *w32.RECT, dwData w32.LPARAM) uintptr {
	// Unwrap dwData as pointer to ConnectedMonitors
	monitors := (*ConnectedDisplays)(unsafe.Pointer(dwData))

	if result := getMonitorSettings(hMonitor); result != nil {
		monitors.count = monitors.count + 1
		monitors.displays = append(monitors.displays, *result)
	}

	return uintptr(w32.TRUE)
}

func getMonitorSettings(hMonitor w32.HMONITOR) *Display {
	info := &w32.MONITORINFOEX{}
	if result := getMonitorInfoW(hMonitor, info); !result {
		log.Error("Failed to get monitor settings")
		return nil
	}
	log.Debug("GetMonitorInfoW: Success")

	devMode := w32.DEVMODE{}
	devMode.DmSize = uint16(unsafe.Sizeof(devMode))

	if result := w32.EnumDisplaySettingsEx(&info.SzDevice[0], w32.ENUM_CURRENT_SETTINGS, &devMode, 0); !result {
		log.Error("Failed to enumerate display settings")
		return nil
	}
	log.Debug("EnumDisplaySettingsEx: Success")

	return &Display{
		WidthPx:  devMode.DmPelsWidth,
		HeightPx: devMode.DmPelsHeight,
		// BitsPerPx: devMode.DmBitsPerPel,
		// RefreshRate: devMode.DmDisplayFrequency,
	}
}

// w32 expects MONITORINFO as argument, but we need MONITORINFOEX
func getMonitorInfoW(hMonitor w32.HMONITOR, lmpi *w32.MONITORINFOEX) bool {
	if lmpi != nil {
		lmpi.CbSize = uint32(unsafe.Sizeof(*lmpi))
	}
	ret, _, _ := libGetMonitorInfoW.Call(
		uintptr(hMonitor),
		uintptr(unsafe.Pointer(lmpi)),
	)
	return ret != 0
}
