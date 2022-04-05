package main

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	CacheDirName     = "desktop-slideshow"
	ProviderBing     = "bing"
	ProviderUnsplash = "unsplash"
)

func NewWallpaperManager(config *Configuration) WallpaperManager {
	display := GetLargestDisplay()
	log.Infof("Detected display with resolution %dx%d", display.WidthPx, display.HeightPx)

	downloadDir, err := wallpaperDownloadDir()
	if err != nil {
		log.Fatal(err)
	}

	return WallpaperManager{
		DownloadDir: downloadDir,
		Provider:    getProviderByConfig(config),
		WallpaperSize: &WallpaperSize{
			Width:  display.WidthPx,
			Height: display.HeightPx,
		},
	}
}

type WallpaperManager struct {
	DownloadDir   string
	Provider      IWallpaperProvider
	WallpaperSize *WallpaperSize
}

func (mgr WallpaperManager) UpdateWallpaper() error {
	dlPath := filepath.Join(mgr.DownloadDir, "wallpaper.wpr")

	api := mgr.Provider.GetApiInstance()

	log.Info("Downloading wallpaper")
	if err := api.DownloadWallpaper(mgr.WallpaperSize, dlPath); err != nil {
		return err
	}

	log.Info("Setting wallpaper from path: ", dlPath)
	if err := SetWallpaperFromFile(dlPath); err != nil {
		return err
	}

	return nil
}

func getProviderByConfig(config *Configuration) IWallpaperProvider {

	if strings.ToLower(config.Provider) == ProviderBing {
		log.Info("Initialized Bing Provider")
		return &BingProvider{
			Market: config.Bing.Market,
		}
	} else if strings.ToLower(config.Provider) == ProviderUnsplash {
		log.Info("Initialized Unplash Provider")
		return &UnsplashProvider{
			Source:     randomSource(config.Unsplash.Sources),
			SearchTags: config.Unsplash.SearchTags,
		}
	} else {
		log.Fatal("Incorrect provider %s", config.Provider)
	}

	return nil
}

func randomSource(slice []string) string {
	return slice[rand.Intn(len(slice))]
}

func wallpaperDownloadDir() (string, error) {
	localAppData, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	downloadPath := filepath.Join(localAppData, CacheDirName)

	if _, err := os.Stat(downloadPath); os.IsNotExist(err) {
		if err := os.MkdirAll(downloadPath, os.ModePerm); err != nil {
			return "", err
		}
	}

	return downloadPath, nil
}
