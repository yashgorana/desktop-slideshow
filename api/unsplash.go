package api

import (
	"fmt"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/yashgorana/desktop-slideshow/utils"
)

// Implements IWallpaperApi
type UnsplashApi struct {
	Source     string
	SearchTags []string
}

func (api UnsplashApi) GetWallpaperUrl(size *WallpaperSize) (string, error) {
	return api.unsplashUrlWithSize(size.Width, size.Height), nil
}

func (api UnsplashApi) DownloadWallpaper(size *WallpaperSize, toPath string) error {
	imgUrl, _ := api.GetWallpaperUrl(size)
	log.Info("UnsplashApi: URL ", imgUrl)

	return utils.DownloadFile(imgUrl, toPath)
}

func (api UnsplashApi) unsplashUrlWithSize(w uint32, h uint32) string {
	searchTags := ""

	if len(api.SearchTags) > 0 {
		searchTags = url.QueryEscape(strings.Join(api.SearchTags, ","))
	}

	url := url.URL{
		Scheme:   "https",
		Host:     "source.unsplash.com",
		Path:     fmt.Sprintf("/%s/%dx%d", api.Source, w, h),
		RawQuery: searchTags,
	}
	return url.String()
}
