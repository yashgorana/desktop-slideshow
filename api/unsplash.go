package api

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/yashgorana/desktop-slideshow/utils"
)

// Implements IWallpaperApi
type UnsplashApi struct {
	Source     string
	SearchTags []string
}

func (api UnsplashApi) DownloadWallpaper(size *WallpaperSize, toPath string) error {
	path, err := filepath.Abs(toPath)
	if err != nil {
		return err
	}

	url := api.unsplashUrlWithSize(size.Width, size.Height)
	log.Debug("UnsplashApi: Fetching image from ", url, " to ", path)

	return utils.DownloadFile(url, toPath)
}

func (api UnsplashApi) unsplashUrlWithSize(w uint32, h uint32) string {
	searchTags := ""
	if len(api.SearchTags) > 0 {
		searchTags = fmt.Sprintf("?%s", url.QueryEscape(strings.Join(api.SearchTags, ",")))
	}
	return fmt.Sprintf("https://source.unsplash.com/%s/%dx%d/%s", api.Source, w, h, searchTags)
}
