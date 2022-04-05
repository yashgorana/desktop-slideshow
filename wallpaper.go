package main

type WallpaperSize struct {
	Width  uint32
	Height uint32
}

type IWallpaperProvider interface {
	GetApiInstance() IWallpaperApi
}

type IWallpaperApi interface {
	DownloadWallpaper(size *WallpaperSize, toPath string) error
}
