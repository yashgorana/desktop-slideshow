package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/yashgorana/desktop-slideshow/utils"
)

// Implements IWallpaperApi
type BingApi struct {
	Market string
}

func (api BingApi) DownloadWallpaper(size *WallpaperSize, toPath string) error {
	path, err := filepath.Abs(toPath)
	if err != nil {
		return err
	}

	resp, err := getWallpaperMetadata(api.Market, 1)
	if err != nil {
		return err
	}

	// Get the latest image
	img := resp.Images[0]

	// always work with the UHD image
	imgUrlUHD := img.Urlbase + "_UHD.jpg"

	url, err := url.Parse(imgUrlUHD)
	if err != nil {
		return err
	}

	query := url.Query()
	query.Set("w", fmt.Sprintf("%d", size.Width))
	query.Set("h", fmt.Sprintf("%d", size.Height))

	url.Scheme = "https"
	url.Host = "bing.com"
	url.RawQuery = query.Encode()

	dlUrl := url.String()

	log.Debug("BingWallpaperApi: Fetching image from ", dlUrl, " to ", path)

	return utils.DownloadFile(dlUrl, path)
}

func getWallpaperMetadata(market string, count uint8) (*bingWallpaperApiResponse, error) {
	apiUrl := urlFromProps(&bingApiUrlProps{
		format: "js",
		mkt:    market,
		n:      count,
		idx:    0,
	})

	log.Debug("BingWallpaperApi: Metadata URL=", apiUrl)

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response bingWallpaperApiResponse
	json.NewDecoder(resp.Body).Decode(&response)

	return &response, nil
}

func urlFromProps(props *bingApiUrlProps) string {
	url := url.URL{
		Scheme: "https",
		Host:   "bing.com",
		Path:   "HPImageArchive.aspx",
	}

	q := url.Query()
	q.Set("format", props.format)
	q.Set("mkt", props.mkt)
	q.Set("idx", strconv.Itoa(int(props.idx)))
	q.Set("n", strconv.Itoa(int(props.n)))
	url.RawQuery = q.Encode()

	return url.String()
}
