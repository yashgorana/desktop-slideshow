package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	client = http.Client{
		Timeout: 30 * time.Second,
	}
)

func DownloadFile(url string, path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	// Create a new file at path
	dlF, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dlF.Close()

	// Download file from URL
	response, err := client.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Check server response
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %d", response.StatusCode)
	}

	// Writer the body to file
	if _, err := io.Copy(dlF, response.Body); err != nil {
		return err
	}
	return dlF.Sync()
}
