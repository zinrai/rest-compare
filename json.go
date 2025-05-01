package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
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

// Extracts a value from a JSON object at the specified path
// Path is specified in dot notation like "settings.timeout"
func extractPath(data map[string]interface{}, path string) (interface{}, error) {
	if path == "" {
		return data, nil
	}

	parts := strings.Split(path, ".")
	current := interface{}(data)

	for _, part := range parts {
		current, err := navigatePathElement(current, part)
		if err != nil {
			return nil, err
		}

		if current == nil {
			return nil, errors.New("path element not found: " + part)
		}
	}

	return current, nil
}

// Handles navigation to a single path element
func navigatePathElement(data interface{}, pathElement string) (interface{}, error) {
	switch currentType := data.(type) {
	case map[string]interface{}:
		value, exists := currentType[pathElement]
		if !exists {
			return nil, errors.New("path element not found: " + pathElement)
		}
		return value, nil

	case []interface{}:
		index, err := parseIndex(pathElement)
		if err != nil {
			return nil, err
		}

		if index < 0 || index >= len(currentType) {
			return nil, errors.New("array index out of range: " + pathElement)
		}
		return currentType[index], nil

	default:
		return nil, errors.New("cannot traverse path: non-object element")
	}
}

// Parses a string to an array index
func parseIndex(s string) (int, error) {
	// Check if string is wrapped in square brackets
	if len(s) >= 2 && s[0] == '[' && s[len(s)-1] == ']' {
		// Extract the number between brackets
		numStr := s[1 : len(s)-1]
		return strconv.Atoi(numStr)
	}

	// Try direct conversion
	return strconv.Atoi(s)
}
