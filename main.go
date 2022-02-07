package main

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

var (
	Environment string
	Version     string
)

func init() {
	if isProd() {
		logFilePath := filepath.Join(executablePath(), "app.log")
		logFp, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}
		log.SetOutput(logFp)
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetOutput(os.Stdout)
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	log.Info("App start. Production=", isProd())

	screenRes := getDisplayResolution()
	log.Infof("Detected screen with resolution %dx%d", screenRes.Width, screenRes.Height)

	mgr := WallpaperManager{
		Provider:   ProviderUnsplash,
		Resolution: screenRes,
		ProviderArgs: UnsplashArgs{
			SearchTag: "wallpaper",
		},
	}

	if err := mgr.UpdateWallpaper(); err != nil {
		log.Fatal(err)
	}
	log.Info("Done")
}

func isProd() bool {
	return Environment == "PROD"
}

func executablePath() string {
	ex, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(ex)
}

func getDisplayResolution() *Resolution {
	displays := GetConnectedDisplays()
	if displays == nil {
		log.Fatal("No displays found!")
	}

	return getLargestDisplayResolution(displays)
}

func getLargestDisplayResolution(displays *ConnectedDisplays) *Resolution {
	largestResolution := &Resolution{}

	for _, disp := range displays.displays {
		if disp.WidthPx > largestResolution.Width && disp.HeightPx > largestResolution.Height {
			largestResolution.Width = disp.WidthPx
			largestResolution.Height = disp.HeightPx
		}
	}

	return largestResolution
}
