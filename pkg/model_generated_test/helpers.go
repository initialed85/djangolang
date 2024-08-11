package model_generated_test

import (
	"io"
	"net/http"
)

type HTTPClient struct {
	httpClient *http.Client
}

func (h *HTTPClient) Get(url string) (resp *http.Response, err error) {
	return h.httpClient.Get(url)
}

func (h *HTTPClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	return h.httpClient.Post(url, contentType, body)
}

func (h *HTTPClient) Put(url, contentType string, body io.Reader) (*http.Response, error) {
	r, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")

	return h.httpClient.Do(r)
}

func (h *HTTPClient) Patch(url, contentType string, body io.Reader) (*http.Response, error) {
	r, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")

	return h.httpClient.Do(r)
}

func (h *HTTPClient) Delete(url string) (*http.Response, error) {
	r, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")

	return h.httpClient.Do(r)
}
