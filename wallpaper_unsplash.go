package main

import (
	"fmt"
	"net/url"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type UnsplashProvider struct {
	Source     string
	SearchTags string
}

func (p UnsplashProvider) GetApiInstance() WallpaperApi {
	return unsplashWallpaperApi(p)
}

type unsplashWallpaperApi struct {
	Source     string
	SearchTags string
}

func (api unsplashWallpaperApi) DownloadWallpaper(res *Resolution, toPath string) error {

	path, err := filepath.Abs(toPath)
	if err != nil {
		return err
	}
	url := api.unsplashUrlWithSize(res.Width, res.Height)

	log.Debug("UnsplashApi: Fetching image from ", url, " to ", path)

	return DownloadFile(url, toPath)
}

func (api unsplashWallpaperApi) unsplashUrlWithSize(w uint32, h uint32) string {
	return fmt.Sprintf("https://source.unsplash.com/%s/%dx%d/?%s", api.Source, w, h, url.QueryEscape(api.SearchTags))
}
