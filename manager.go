package main

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const (
	CacheDirName = "desktop-slideshow"
)

type WallpaperManager struct {
	Provider   WallpaperProvider
	Resolution *Resolution
}

func (mgr WallpaperManager) UpdateWallpaper() error {
	dlDir, err := wallpaperDownloadDir()
	if err != nil {
		return err
	}

	dlPath := filepath.Join(dlDir, "wallpaper.wpr")

	api := mgr.Provider.GetApiInstance()
	if err := api.DownloadWallpaper(mgr.Resolution, dlPath); err != nil {
		return err
	}

	log.Debug("Setting wallpaper from path: ", dlPath)

	if err := SetWallpaperFromFile(dlPath); err != nil {
		return err
	}

	return nil
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
