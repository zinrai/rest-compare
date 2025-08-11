package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/theory/jsonpath"
)

// Converts a JSON string to a map
func parseJSON(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Extracts a value from a JSON object at the specified JSONPath
// Path should be a valid JSONPath expression like "$.settings.timeout" or "$..name"
func extractPath(data map[string]interface{}, path string) (interface{}, error) {
	if path == "" {
		return data, nil
	}

	// Parse the JSONPath expression
	p, err := jsonpath.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("invalid JSONPath expression: %w", err)
	}

	// Execute the JSONPath query
	results := p.Select(data)

	// Check if any results were found
	if len(results) == 0 {
		return nil, errors.New("no matching elements found for JSONPath: " + path)
	}

	// For comparison purposes, we expect a single result
	// If multiple results are found, return an error
	if len(results) > 1 {
		return nil, fmt.Errorf("JSONPath returned multiple results (%d), expected single result", len(results))
	}

	return results[0], nil
}
