package gocopy

import (
	"fmt"
	"io"
	"net/http"
)

type sourceFetcher interface {
	Fetch() ([]byte, error)
}

type HTTPFetcher struct {
	URL string
}

func (H HTTPFetcher) Fetch() ([]byte, error) {
	resp, err := http.Get(H.URL)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("error downloading file: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
