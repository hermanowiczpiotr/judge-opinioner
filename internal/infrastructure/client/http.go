package client

import (
	"io"
	"net/http"
	"net/url"

	"github.com/kataras/iris/v12/x/errors"
)

type HttpClientInterface interface {
	GET(path string, queryParams map[string]string) ([]byte, error)
}

type HTTPClient struct {
}

func (client HTTPClient) GET(path string, queryParams map[string]string) ([]byte, error) {
	baseURL, err := url.Parse(path)
	if err != nil {
		return nil, errors.New("")
	}

	urlValues := url.Values{}
	for key, value := range queryParams {
		urlValues.Add(key, value)
	}

	baseURL.RawQuery = urlValues.Encode()

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "*/*")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
