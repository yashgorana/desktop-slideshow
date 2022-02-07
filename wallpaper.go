package main

type Resolution struct {
	Width  uint32
	Height uint32
}

type WallpaperAPI interface {
	DownloadWallpaper(resolution *Resolution, toPath string) error
}
