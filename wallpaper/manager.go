package wallpaper

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/yashgorana/desktop-slideshow/api"
	"github.com/yashgorana/desktop-slideshow/win32"

	log "github.com/sirupsen/logrus"
)

const (
	CacheDirName     = "desktop-slideshow"
	ProviderBing     = "bing"
	ProviderUnsplash = "unsplash"
)

type WallpaperManager struct {
	DownloadDir   string
	Provider      IWallpaperProvider
	WallpaperSize *api.WallpaperSize
}

func NewManager(config *WallpaperManagerConfig) WallpaperManager {
	display := win32.GetLargestDisplay()
	log.Infof("Detected display with resolution %dx%d", display.WidthPx, display.HeightPx)

	downloadDir, err := wallpaperDownloadDir()
	if err != nil {
		log.Fatal(err)
	}

	return WallpaperManager{
		DownloadDir: downloadDir,
		Provider:    getProviderByConfig(config),
		WallpaperSize: &api.WallpaperSize{
			Width:  display.WidthPx,
			Height: display.HeightPx,
		},
	}
}

func (mgr WallpaperManager) UpdateWallpaper() error {
	dlPath := filepath.Join(mgr.DownloadDir, "wallpaper.wpr")

	api := mgr.Provider.GetApiInstance()

	log.Info("Downloading wallpaper")
	if err := api.DownloadWallpaper(mgr.WallpaperSize, dlPath); err != nil {
		return err
	}

	log.Info("Setting wallpaper from path: ", dlPath)
	if err := win32.SetWallpaperFromFile(dlPath); err != nil {
		return err
	}

	return nil
}

func getProviderByConfig(config *WallpaperManagerConfig) IWallpaperProvider {

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
