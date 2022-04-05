package api

type WallpaperSize struct {
	Width  uint32
	Height uint32
}

type IWallpaperApi interface {
	GetWallpaperUrl(size *WallpaperSize) (string, error)
	DownloadWallpaper(size *WallpaperSize, toPath string) error
}
