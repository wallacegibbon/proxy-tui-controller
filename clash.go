package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	defaultClashURL = "http://127.0.0.1:9090"
	proxiesPath     = "/proxies"
)

var (
	mockMode  = os.Getenv("MOCK_CLASH") == "1"
	apiSecret = os.Getenv("MIHOMO_SECRET")
)

type ClashClient struct {
	baseURL    string
	secret     string
	httpClient *http.Client
}

type Proxy struct {
	Name    string                 `json:"name"`
	Type    string                 `json:"type"`
	Now     string                 `json:"now"`
	All     []string               `json:"all"`
	History []ProxyHistory         `json:"history"`
	Uptime  string                 `json:"uptime"`
	Extra   map[string]interface{} `json:"extra"`
}

type ProxyHistory struct {
	Time  string `json:"time"`
	Delay int    `json:"delay"`
}

type ProxiesResponse struct {
	Proxies map[string]Proxy `json:"proxies"`
}

func NewClashClient(baseURL string) *ClashClient {
	if baseURL == "" {
		baseURL = defaultClashURL
	}
	return &ClashClient{
		baseURL:    baseURL,
		secret:     apiSecret,
		httpClient: &http.Client{},
	}
}

func (c *ClashClient) addAuthHeader(req *http.Request) {
	if c.secret != "" {
		req.Header.Set("Authorization", "Bearer "+c.secret)
	}
}

func (c *ClashClient) GetProxies() (*ProxiesResponse, error) {
	if mockMode {
		// Return mock data for testing
		proxies := make(map[string]Proxy)
		proxies["Proxy"] = Proxy{
			Name: "Proxy",
			Type: "Selector",
			Now:  "Proxy-1",
			All:  []string{"Proxy-1", "Proxy-2", "Proxy-3"},
		}
		proxies["Auto"] = Proxy{
			Name: "Auto",
			Type: "URLTest",
			Now:  "Auto-2",
			All:  []string{"Auto-1", "Auto-2", "Auto-3", "Auto-4"},
		}
		return &ProxiesResponse{Proxies: proxies}, nil
	}

	url := c.baseURL + proxiesPath
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	c.addAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get proxies: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var result ProxiesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (c *ClashClient) SelectProxy(groupName, proxyName string) error {
	url := c.baseURL + proxiesPath + "/" + groupName

	payload := map[string]string{
		"name": proxyName,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	c.addAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to select proxy: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *ClashClient) TestDelay(groupName, proxyName string, testURL string) (int, error) {
	url := c.baseURL + proxiesPath + "/" + groupName + "/delay"

	if testURL == "" {
		testURL = "http://www.gstatic.com/generate_204"
	}

	payload := map[string]string{
		"url":     testURL,
		"timeout": "5000",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(body))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	c.addAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to test delay: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return result["delay"], nil
}
