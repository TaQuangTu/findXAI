package search

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	apiKey     string
	engineID   string
	httpClient *http.Client
}

func NewClient(apiKey, engineID string) *Client {
	return &Client{
		apiKey:     apiKey,
		engineID:   engineID,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type GoogleSearchResult struct {
	Items []struct {
		Title   string `json:"title"`
		Link    string `json:"link"`
		Snippet string `json:"snippet"`
	} `json:"items"`
}

func (c *Client) Search(ctx context.Context, query string, params map[string]string) ([]SearchResult, error) {
	baseURL := "https://www.googleapis.com/customsearch/v1"

	queryParams := url.Values{}
	queryParams.Add("key", c.apiKey)
	queryParams.Add("cx", c.engineID)
	queryParams.Add("q", query)

	// Add optional parameters
	for key, value := range params {
		if value != "" {
			queryParams.Add(key, value)
		}
	}

	fullURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result GoogleSearchResult

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	results := make([]SearchResult, 0, len(result.Items))
	for _, item := range result.Items {
		results = append(results, SearchResult{
			Title:   item.Title,
			Link:    item.Link,
			Snippet: item.Snippet,
		})
	}

	return results, nil
}

type SearchResult struct {
	Title   string
	Link    string
	Snippet string
}
