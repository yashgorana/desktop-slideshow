package api

type WallpaperSize struct {
	Width  uint32
	Height uint32
}

type IWallpaperApi interface {
	DownloadWallpaper(size *WallpaperSize, toPath string) error
}
