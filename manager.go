package main

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const (
	ProviderUnsplash = "unsplash"
	ProviderBing     = "bing"
)

type BingArgs struct {
	Market string
}

type UnsplashArgs struct {
	SearchTag string
}

type WallpaperManager struct {
	Provider     string
	Resolution   *Resolution
	ProviderArgs interface{}
}

func (mgr *WallpaperManager) UpdateWallpaper() error {
	dlDir, err := wallpaperDownloadDir()
	if err != nil {
		return err
	}

	dlPath := filepath.Join(dlDir, "wallpaper.wpr")

	api := apiFromProvider(mgr.Provider, mgr.ProviderArgs)
	if err := api.DownloadWallpaper(mgr.Resolution, dlPath); err != nil {
		return err
	}

	log.Debug("Setting wallpaper from path: ", dlPath)

	if err := SetWallpaperFromFile(dlPath); err != nil {
		return err
	}

	return nil
}

func apiFromProvider(providerName string, providerArgs interface{}) WallpaperAPI {
	switch providerName {
	case ProviderUnsplash:
		return &UnsplashWallpaperApi{
			SearchTags: providerArgs.(UnsplashArgs).SearchTag,
		}
	case ProviderBing:
		return &BingWallpaperApi{
			Market: providerArgs.(BingArgs).Market,
		}
	default:
		log.Fatal("Yikes. No such provider ", providerName)
	}
	return nil
}

func wallpaperDownloadDir() (string, error) {
	localAppData, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	downloadPath := filepath.Join(localAppData, "daily-wallpaper")

	if _, err := os.Stat(downloadPath); os.IsNotExist(err) {
		if err := os.MkdirAll(downloadPath, os.ModePerm); err != nil {
			return "", err
		}
	}

	return downloadPath, nil
}
