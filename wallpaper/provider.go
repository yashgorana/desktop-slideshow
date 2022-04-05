package wallpaper

import "github.com/yashgorana/desktop-slideshow/api"

type IWallpaperProvider interface {
	GetApiInstance() api.IWallpaperApi
}

// Bing API Provider
type BingProvider struct {
	Market string
}

func (p BingProvider) GetApiInstance() api.IWallpaperApi {
	return api.BingApi(p)
}

// Unsplash API Provider
type UnsplashProvider struct {
	Source     string
	SearchTags []string
}

func (p UnsplashProvider) GetApiInstance() api.IWallpaperApi {
	return api.UnsplashApi(p)
}
