package main

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

type UnsplashProvider struct {
	Source     string
	SearchTags []string
}

func (p UnsplashProvider) GetApiInstance() IWallpaperApi {
	return unsplashWallpaperApi(p)
}

type unsplashWallpaperApi struct {
	Source     string
	SearchTags []string
}

func (api unsplashWallpaperApi) DownloadWallpaper(size *WallpaperSize, toPath string) error {
	path, err := filepath.Abs(toPath)
	if err != nil {
		return err
	}

	url := api.unsplashUrlWithSize(size.Width, size.Height)
	log.Debug("UnsplashApi: Fetching image from ", url, " to ", path)

	return DownloadFile(url, toPath)
}

func (api unsplashWallpaperApi) unsplashUrlWithSize(w uint32, h uint32) string {
	searchTags := ""
	if len(api.SearchTags) > 0 {
		searchTags = fmt.Sprintf("?%s", url.QueryEscape(strings.Join(api.SearchTags, ",")))
	}
	return fmt.Sprintf("https://source.unsplash.com/%s/%dx%d/%s", api.Source, w, h, searchTags)
}
