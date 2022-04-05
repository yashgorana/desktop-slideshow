package win32

import (
	"syscall"
	"unsafe"

	"github.com/gonutz/w32/v2"

	log "github.com/sirupsen/logrus"
)

type Display struct {
	WidthPx     uint32
	HeightPx    uint32
	BitsPerPx   uint32
	RefreshRate uint32
}

type ConnectedDisplays struct {
	count    uint8
	displays []Display
}

// GetLargestDisplay returns the Display with the largest resolution (by pixelss) in all ConnectedDisplays
func GetLargestDisplay() Display {
	largestDisp := Display{}

	displays := GetConnectedDisplays()
	if displays == nil {
		log.Fatal("No displays found!")
	}

	for _, disp := range displays.displays {
		if disp.WidthPx > largestDisp.WidthPx && disp.HeightPx > largestDisp.HeightPx {
			largestDisp = disp
		}
	}

	return largestDisp
}

// GetConnectedDisplays fetches the connected displays and their current settings (currently limited to screen resolution)
func GetConnectedDisplays() *ConnectedDisplays {
	monitors := &ConnectedDisplays{}

	result := w32.EnumDisplayMonitors(w32.HDC(0), nil, syscall.NewCallback(enumDisplayMonitorsCb), uintptr(unsafe.Pointer(monitors)))
	log.Debug("EnumDisplayMonitors: ", result)

	return monitors
}

// MONITORENUMPROC - Must return uintptr
func enumDisplayMonitorsCb(hMonitor w32.HMONITOR, hdc w32.HDC, lpRect *w32.RECT, dwData w32.LPARAM) uintptr {
	// Unwrap dwData as pointer to ConnectedMonitors
	monitors := (*ConnectedDisplays)(unsafe.Pointer(dwData))
	if result := getMonitorSettings(hMonitor); result != nil {
		log.Debug("EnumDisplayMonitors Callback: Monitor settings - ", result)

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
	log.Debug("getMonitorInfoW: ", info)

	devMode := w32.DEVMODE{}
	devMode.DmSize = uint16(unsafe.Sizeof(devMode))

	if result := w32.EnumDisplaySettingsEx(&info.SzDevice[0], w32.ENUM_CURRENT_SETTINGS, &devMode, 0); !result {
		log.Error("Failed to enumerate display settings")
		return nil
	}
	log.Debug("EnumDisplaySettingsEx: ", devMode)

	return &Display{
		WidthPx:     devMode.DmPelsWidth,
		HeightPx:    devMode.DmPelsHeight,
		BitsPerPx:   devMode.DmBitsPerPel,
		RefreshRate: devMode.DmDisplayFrequency,
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
