package gocopy

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPFetcher_Fetch(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		fmt.Fprint(w, "Hello, World!")
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{
			name:    "Successful fetch",
			url:     server.URL + "/hello",
			want:    "Hello, World!",
			wantErr: false,
		},
		{
			name:    "Non-existent file",
			url:     server.URL + "/foobar",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := HTTPFetcher{
				URL: tt.url,
			}
			got, err := h.Fetch()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, string(got))
			}
		})
	}
}
