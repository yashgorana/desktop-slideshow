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

	mgr := WallpaperManager{
		Provider: ProviderUnsplash,
		Resolution: &Resolution{
			Width:  3840,
			Height: 2160,
		},
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
