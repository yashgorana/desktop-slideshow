package app

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/yashgorana/desktop-slideshow/utils"
	"github.com/yashgorana/desktop-slideshow/win32"
)

const (
	ExecutableName = "desktop-slideshow"
	RegistryPath   = "DesktopBackground\\shell\\Next Wallpaper\\command"
)

func Init() {
	utils.ResetWorkingDir()
	if isFirstRun() {
		log.Info("First time run. Setting up the application.")
		performFirstTimeSetup()
	}
}

func isFirstRun() bool {
	return !hasDesktopShellEntry()
}

func performFirstTimeSetup() {
	elevateProcess()
	createDesktopShellEntry()
}

func hasDesktopShellEntry() bool {
	val, err := win32.GetClassRootValue(RegistryPath, win32.RegistryDefaultName)

	if err != nil {
		log.Error("Failed to get desktop shell entry. error=", err)
		return false
	}

	return val == utils.ExecutablePath()
}

func createDesktopShellEntry() error {
	err := win32.SetClassRootValue(RegistryPath, win32.RegistryDefaultName, utils.ExecutablePath())
	if err != nil {
		log.Error("Failed to create desktop shell entry. error=", err)
		return err
	}

	return nil
}

func elevateProcess() {
	if win32.CheckElevated() {
		log.Info("Process is elevated")
		return
	}

	if err := win32.RunElevated(); err != nil {
		log.Error("Failed to elevate process. error=", err)
		os.Exit(1)
	}

	log.Info("New elevated process created. Quitting.")
	os.Exit(0)
}
