package main

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/yashgorana/desktop-slideshow/app"
	"github.com/yashgorana/desktop-slideshow/utils"
	"github.com/yashgorana/desktop-slideshow/wallpaper"
)

var (
	Environment string
	Version     string
)

func isProd() bool {
	return Environment == "PROD"
}

func init() {
	if isProd() {
		logFilePath := filepath.Join(utils.ExecutableDir(), "app.log")
		logFp, err := os.Create(logFilePath)
		// logFp, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		log.SetOutput(logFp)
		log.SetLevel(log.InfoLevel)

		app.Init()
	} else {
		log.SetOutput(os.Stdout)
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	log.Info("App start. Production=", isProd())

	config := wallpaper.LoadConfig()
	mgr := wallpaper.NewManager(config)
	if err := mgr.UpdateWallpaper(); err != nil {
		log.Fatal(err)
	}
	log.Info("Done")
}
