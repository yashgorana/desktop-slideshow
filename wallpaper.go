package main

type Resolution struct {
	Width  uint32
	Height uint32
}

type WallpaperProvider interface {
	GetApiInstance() WallpaperApi
}

type WallpaperApi interface {
	DownloadWallpaper(resolution *Resolution, toPath string) error
}
