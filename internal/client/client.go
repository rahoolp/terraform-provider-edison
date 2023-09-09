package edison

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-cleanhttp"
)

type Client struct {
	client  *http.Client
	baseURL *url.URL

	token string

	EAStores    *EAStoresService
	EHSClusters *EHSClustersService
	AWs         *AWsService
}

func NewClient(baseURL, token string) (*Client, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	c := &Client{
		client:  cleanhttp.DefaultPooledClient(),
		baseURL: base,
		token:   token,
	}
	c.EAStores = newEAStoreService("eastores", c)
	c.EHSClusters = newEHSClusterService("ehsclusters", c)
	c.AWs = newAWService("aws", c)
	return c, nil
}

func (c Client) NewRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("error parsing path: %w", err)
	}
	reqURL := c.baseURL.ResolveReference(u)
	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), body)
	if err != nil {
		return nil, err
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer: "+c.token)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (c Client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}
