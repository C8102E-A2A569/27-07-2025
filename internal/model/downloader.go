package model

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"zip-archive/internal/config"
)

type Downloader struct {
	client      *http.Client
	allowedExts map[string]bool
	maxSizeMB   int64
}

func New(cfg *config.FilesConfig) *Downloader {
	allowedExts := make(map[string]bool)
	for _, ext := range cfg.AllowedExt {
		allowedExts[ext] = true
	}

	return &Downloader{
		client: &http.Client{
			Timeout: cfg.DownloadTimeout,
		},
		allowedExts: allowedExts,
		maxSizeMB:   int64(cfg.MaxSizeMB) * 1024 * 1024,
	}
}

func (d *Downloader) Download(ctx context.Context, url string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", err
	}
	resp, err := d.client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	if resp.ContentLength > d.maxSizeMB {
		return nil, "", fmt.Errorf("file too large")
	}
	filename := d.extractFilename(url, resp)
	if !d.isAllowedExt(filename) {
		return nil, "", fmt.Errorf("unsupported file type")
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return data, filename, nil
}

func (d *Downloader) extractFilename(url string, resp *http.Response) string {
	if cd := resp.Header.Get("Content-Disposition"); cd != "" {
		if strings.Contains(cd, "filename=") {
			parts := strings.Split(cd, "filename=")
			if len(parts) > 1 {
				filename := strings.Trim(parts[1], `"`)
				if filename != "" {
					return filename
				}
			}
		}
	}
	return path.Base(url)
}

func (d *Downloader) isAllowedExt(filename string) bool {
	ext := strings.ToLower(path.Ext(filename))
	return d.allowedExts[ext]
}
