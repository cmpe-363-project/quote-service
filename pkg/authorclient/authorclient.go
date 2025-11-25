package authorclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type NewClientConfig struct {
	BaseURL string
	Timeout time.Duration
}

type Author struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type VersionResponse struct {
	Version string `json:"version"`
}

type AuthorsResponse struct {
	Items []Author `json:"items"`
}

func NewClient(config NewClientConfig) *Client {
	return &Client{
		baseURL: config.BaseURL,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

func (c *Client) GetVersion() (string, error) {
	url := fmt.Sprintf("%s/api/version", c.baseURL)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to get version: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var versionResp VersionResponse
	if err := json.NewDecoder(resp.Body).Decode(&versionResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return versionResp.Version, nil
}

func (c *Client) GetAuthorsByIDs(ids []int) ([]Author, error) {
	if len(ids) == 0 {
		return []Author{}, nil
	}

	idsStr := make([]string, len(ids))
	for i, id := range ids {
		idsStr[i] = strconv.Itoa(id)
	}

	url := fmt.Sprintf("%s/api/authors/by-id?id=%s", c.baseURL, strings.Join(idsStr, ","))
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get authors: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var authorsResp AuthorsResponse
	if err := json.NewDecoder(resp.Body).Decode(&authorsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return authorsResp.Items, nil
}
