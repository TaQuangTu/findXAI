package search

import (
	"context"
	"findx/internal/libhttp"
	"fmt"
	"net/url"
	"time"
)

type Client struct {
	httpClient libhttp.IHttpRequest
}

func NewClient() *Client {
	return &Client{
		httpClient: libhttp.NewHttpClient(),
	}
}

type GoogleSearchResult struct {
	Items []struct {
		Title   string `json:"title"`
		Link    string `json:"link"`
		Snippet string `json:"snippet"`
	} `json:"items"`
}

func (c *Client) Search(ctx context.Context, apiKey, engineID string, params map[string]string) (_ []SearchResult, status int, err error) {
	baseURL := "https://www.googleapis.com/customsearch/v1"

	queryParams := url.Values{}
	queryParams.Add("key", apiKey)
	queryParams.Add("cx", engineID)

	// Add optional parameters
	for key, value := range params {
		if value != "" {
			queryParams.Add(key, value)
		}
	}

	var (
		fullURL = fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())
		result  GoogleSearchResult
	)
	status, err = c.httpClient.Request(
		ctx,
		libhttp.RequestOption{
			RequestTimeout: 10 * time.Second,
			Method:         libhttp.GET,
			Url:            fullURL,
		},
		&result)

	results := make([]SearchResult, 0, len(result.Items))
	for _, item := range result.Items {
		results = append(results, SearchResult{
			Title:   item.Title,
			Link:    item.Link,
			Snippet: item.Snippet,
		})
	}
	return results, status, err
}

type SearchResult struct {
	Title   string
	Link    string
	Snippet string
}
