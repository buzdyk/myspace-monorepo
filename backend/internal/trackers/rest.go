package trackers

import (
	"fmt"
	"net/http"
	"net/url"
)

type RestClient struct {
	client *http.Client
}

func NewRestClient() *RestClient {
	return &RestClient{
		client: &http.Client{},
	}
}

func (r *RestClient) Get(baseURI, path string, headers map[string]string, params map[string]string) (*http.Response, error) {
	fullURL := baseURI + path
	
	if len(params) > 0 {
		u, err := url.Parse(fullURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse URL: %w", err)
		}
		
		q := u.Query()
		for key, value := range params {
			q.Set(key, value)
		}
		u.RawQuery = q.Encode()
		fullURL = u.String()
	}
	
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	
	return resp, nil
}