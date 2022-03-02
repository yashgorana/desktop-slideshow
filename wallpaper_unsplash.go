package main

import (
	"fmt"
	"net/url"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type UnsplashProvider struct {
	SearchTags string
}

func (p UnsplashProvider) GetApiInstance() WallpaperApi {
	return unsplashWallpaperApi(p)
}

type unsplashWallpaperApi struct {
	SearchTags string
}

func (api unsplashWallpaperApi) DownloadWallpaper(res *Resolution, toPath string) error {

	path, err := filepath.Abs(toPath)
	if err != nil {
		return err
	}
	url := unsplashUrlWithArgs(res.Width, res.Height, api.SearchTags)

	log.Debug("UnsplashApi: Fetching image from ", url, " to ", path)

	return DownloadFile(url, toPath)
}

func unsplashUrlWithArgs(w uint32, h uint32, query string) string {
	return fmt.Sprintf("https://source.unsplash.com/random/%dx%d/?%s", w, h, url.QueryEscape(query))
}
