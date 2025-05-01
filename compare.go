package main

import (
	"fmt"
	"reflect"

	"github.com/google/go-cmp/cmp"
)

// Compares two JSON objects and returns whether they are equal and difference information
func compareJSON(a, b interface{}, ignoreKeys []string) (bool, []string) {
	// Set up go-cmp options
	opts := []cmp.Option{
		cmp.FilterPath(func(p cmp.Path) bool {
			if len(p) == 0 {
				return false
			}
			lastStep := p[len(p)-1]

			// Handle MapIndex case
			if mapIdx, ok := lastStep.(cmp.MapIndex); ok {
				key := fmt.Sprintf("%v", mapIdx.Key())
				// Check if key is in ignoreKeys list
				for _, ignoreKey := range ignoreKeys {
					if key == ignoreKey {
						return true // Filter (ignore) this path
					}
				}
			}
			return false
		}, cmp.Ignore()),
	}

	// Execute comparison
	diff := cmp.Diff(a, b, opts...)

	// If no difference, they are equal
	if diff == "" {
		return true, nil
	}

	// Format difference information
	diffs := formatDifferences(a, b, ignoreKeys)
	return false, diffs
}

// Formats difference information in a readable format
func formatDifferences(a, b interface{}, ignoreKeys []string) []string {
	var result []string

	// Detect differences through deep comparison
	diffPaths := findDiffPaths(a, b, "", ignoreKeys)

	for _, diffInfo := range diffPaths {
		result = append(result, fmt.Sprintf("- Path: %s\n  A: %v\n  B: %v",
			diffInfo.path, formatValue(diffInfo.valueA), formatValue(diffInfo.valueB)))
	}

	return result
}

// Represents difference information
type diffInfo struct {
	path   string
	valueA interface{}
	valueB interface{}
}

// Recursively finds difference paths between two objects
func findDiffPaths(a, b interface{}, currentPath string, ignoreKeys []string) []diffInfo {
	// Check for type differences
	typeA, typeB := reflect.TypeOf(a), reflect.TypeOf(b)
	if typeA != typeB {
		return []diffInfo{{currentPath, a, b}}
	}

	// Check for value equality
	if reflect.DeepEqual(a, b) {
		return nil
	}

	// Handle maps
	mapA, okA := a.(map[string]interface{})
	mapB, okB := b.(map[string]interface{})
	if okA && okB {
		return findMapDifferences(mapA, mapB, currentPath, ignoreKeys)
	}

	// Handle slices
	sliceA, okA := a.([]interface{})
	sliceB, okB := b.([]interface{})
	if okA && okB {
		return findSliceDifferences(sliceA, sliceB, currentPath, ignoreKeys)
	}

	// Handle other cases (primitive values, etc.)
	return []diffInfo{{currentPath, a, b}}
}

// Finds differences between two maps
func findMapDifferences(mapA, mapB map[string]interface{}, currentPath string, ignoreKeys []string) []diffInfo {
	var results []diffInfo

	// Collect all keys from both maps
	allKeys := make(map[string]bool)
	for k := range mapA {
		allKeys[k] = true
	}
	for k := range mapB {
		allKeys[k] = true
	}

	// Check each key for differences
	for k := range allKeys {
		// Skip keys in ignore list
		if isIgnoredKey(k, ignoreKeys) {
			continue
		}

		valueA, existsA := mapA[k]
		valueB, existsB := mapB[k]

		// Create path for this key
		newPath := createPath(currentPath, k)

		// Handle keys that exist in only one map
		if !existsA {
			results = append(results, diffInfo{newPath, "[missing]", valueB})
			continue
		}
		if !existsB {
			results = append(results, diffInfo{newPath, valueA, "[missing]"})
			continue
		}

		// Recursively compare values that exist in both maps
		results = append(results, findDiffPaths(valueA, valueB, newPath, ignoreKeys)...)
	}

	return results
}

// Finds differences between two slices
func findSliceDifferences(sliceA, sliceB []interface{}, currentPath string, ignoreKeys []string) []diffInfo {
	// Handle different lengths
	if len(sliceA) != len(sliceB) {
		return []diffInfo{{currentPath,
			fmt.Sprintf("array[%d]", len(sliceA)),
			fmt.Sprintf("array[%d]", len(sliceB))}}
	}

	var results []diffInfo

	// Compare each element
	for i := 0; i < len(sliceA); i++ {
		newPath := fmt.Sprintf("%s[%d]", currentPath, i)
		results = append(results, findDiffPaths(sliceA[i], sliceB[i], newPath, ignoreKeys)...)
	}

	return results
}

// Checks if a key is in the ignore list
func isIgnoredKey(key string, ignoreKeys []string) bool {
	for _, ignoreKey := range ignoreKeys {
		if key == ignoreKey {
			return true
		}
	}
	return false
}

// Creates a new path by appending a key to the current path
func createPath(currentPath, key string) string {
	if currentPath == "" {
		return key
	}
	return currentPath + "." + key
}

// Converts a value to a string representation
func formatValue(v interface{}) string {
	if v == nil {
		return "null"
	}

	switch val := v.(type) {
	case string:
		if val == "[missing]" {
			return val
		}
		return fmt.Sprintf("\"%s\"", val)
	default:
		return fmt.Sprintf("%v", v)
	}
}
