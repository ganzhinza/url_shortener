package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"url_shortener/pkg/db/memdb"
)

var api *API

const urlsize = 10

func TestMain(m *testing.M) {
	api = New(memdb.NewDB(urlsize))

	m.Run()
}

func TestAPI_saveURL(t *testing.T) {
	formData := url.Values{}
	formData.Set("url", "http://example.com")

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	api.Router().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Error, got: %d, want: %d", rr.Code, http.StatusOK)
	}
	if len(rr.Body.String()) != urlsize {
		t.Errorf("Wrong ans len got: %d, want %d", len(rr.Body.String()), urlsize)
	}

}

func TestAPI_getLongURL(t *testing.T) {
	//make short url
	formData := url.Values{}
	originalURL := "http://example.com"
	formData.Set("url", originalURL)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	api.Router().ServeHTTP(rr, req)

	//get short url
	short := rr.Body.String()

	queryParams := url.Values{}
	queryParams.Set("shortURL", short)

	req = httptest.NewRequest(http.MethodGet, "/?"+queryParams.Encode(), nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	api.Router().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Error, got: %d, want: %d", rr.Code, http.StatusOK)
	}

	url := rr.Body.String()
	if url != originalURL {
		t.Errorf("Unexpected url, got %s, want %s", url, originalURL)
	}
}
