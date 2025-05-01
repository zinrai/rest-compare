package main

import (
	"io"
	"net/http"
	"time"
)

// Fetches JSON data from the specified URL and converts it to a map
func fetchJSON(url, auth string, timeout time.Duration) (map[string]interface{}, error) {
	// Set up HTTP client
	client := &http.Client{
		Timeout: timeout,
	}

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set authentication header if provided
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}

	// Set Accept header
	req.Header.Set("Accept", "application/json")

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	return parseJSON(body)
}
