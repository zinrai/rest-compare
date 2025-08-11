package main

import (
	"strings"
	"testing"
)

func TestCompareJSON(t *testing.T) {
	testCases := []struct {
		name       string
		a          interface{}
		b          interface{}
		ignoreKeys []string
		wantEqual  bool
	}{
		{
			name: "identical objects",
			a: map[string]interface{}{
				"name":  "test",
				"value": 123,
			},
			b: map[string]interface{}{
				"name":  "test",
				"value": 123,
			},
			ignoreKeys: nil,
			wantEqual:  true,
		},
		{
			name: "different values",
			a: map[string]interface{}{
				"name":  "test",
				"value": 123,
			},
			b: map[string]interface{}{
				"name":  "test",
				"value": 456,
			},
			ignoreKeys: nil,
			wantEqual:  false,
		},
		{
			name: "ignored key differences",
			a: map[string]interface{}{
				"name":  "test",
				"value": 123,
				"id":    "A001",
			},
			b: map[string]interface{}{
				"name":  "test",
				"value": 123,
				"id":    "B002", // different value
			},
			ignoreKeys: []string{"id"},
			wantEqual:  true,
		},
		{
			name: "nested objects",
			a: map[string]interface{}{
				"nested": map[string]interface{}{
					"inner": "value1",
				},
			},
			b: map[string]interface{}{
				"nested": map[string]interface{}{
					"inner": "value2", // different value
				},
			},
			ignoreKeys: nil,
			wantEqual:  false,
		},
		{
			name: "nested objects with ignored key",
			a: map[string]interface{}{
				"nested": map[string]interface{}{
					"inner": "value1",
					"id":    "A001",
				},
			},
			b: map[string]interface{}{
				"nested": map[string]interface{}{
					"inner": "value1",
					"id":    "B002", // different value
				},
			},
			ignoreKeys: []string{"id"},
			wantEqual:  true,
		},
		{
			name: "arrays with different lengths",
			a: map[string]interface{}{
				"items": []interface{}{"a", "b", "c"},
			},
			b: map[string]interface{}{
				"items": []interface{}{"a", "b"},
			},
			ignoreKeys: nil,
			wantEqual:  false,
		},
		{
			name: "arrays with same length but different values",
			a: map[string]interface{}{
				"items": []interface{}{"a", "b", "c"},
			},
			b: map[string]interface{}{
				"items": []interface{}{"a", "b", "d"},
			},
			ignoreKeys: nil,
			wantEqual:  false,
		},
		{
			name: "missing keys",
			a: map[string]interface{}{
				"name":  "test",
				"extra": "only in A",
			},
			b: map[string]interface{}{
				"name": "test",
			},
			ignoreKeys: nil,
			wantEqual:  false,
		},
		{
			name: "missing keys but ignored",
			a: map[string]interface{}{
				"name":  "test",
				"extra": "only in A",
			},
			b: map[string]interface{}{
				"name": "test",
			},
			ignoreKeys: []string{"extra"},
			wantEqual:  true,
		},
		{
			name: "different types",
			a: map[string]interface{}{
				"value": 123,
			},
			b: map[string]interface{}{
				"value": "123", // number vs string
			},
			ignoreKeys: nil,
			wantEqual:  false,
		},
		{
			name: "complex nested structure",
			a: map[string]interface{}{
				"service": map[string]interface{}{
					"name":    "api",
					"version": "1.0",
					"config": map[string]interface{}{
						"timeout": 30,
						"retries": 3,
					},
					"endpoints": []interface{}{
						map[string]interface{}{
							"path":   "/users",
							"method": "GET",
						},
						map[string]interface{}{
							"path":   "/items",
							"method": "POST",
						},
					},
				},
			},
			b: map[string]interface{}{
				"service": map[string]interface{}{
					"name":    "api",
					"version": "1.0",
					"config": map[string]interface{}{
						"timeout": 60, // different value
						"retries": 3,
					},
					"endpoints": []interface{}{
						map[string]interface{}{
							"path":   "/users",
							"method": "GET",
						},
						map[string]interface{}{
							"path":   "/items",
							"method": "POST",
						},
					},
				},
			},
			ignoreKeys: nil,
			wantEqual:  false,
		},
		{
			name: "complex nested structure with ignored paths",
			a: map[string]interface{}{
				"service": map[string]interface{}{
					"name":    "api",
					"version": "1.0",
					"config": map[string]interface{}{
						"timeout": 30, // different but ignored
						"retries": 3,
					},
					"endpoints": []interface{}{
						map[string]interface{}{
							"path":   "/users",
							"method": "GET",
						},
						map[string]interface{}{
							"path":   "/items",
							"method": "POST",
						},
					},
				},
			},
			b: map[string]interface{}{
				"service": map[string]interface{}{
					"name":    "api",
					"version": "1.0",
					"config": map[string]interface{}{
						"timeout": 60, // different but ignored
						"retries": 3,
					},
					"endpoints": []interface{}{
						map[string]interface{}{
							"path":   "/users",
							"method": "GET",
						},
						map[string]interface{}{
							"path":   "/items",
							"method": "POST",
						},
					},
				},
			},
			ignoreKeys: []string{"timeout"},
			wantEqual:  true,
		},
		{
			name: "null vs empty object",
			a: map[string]interface{}{
				"data": nil,
			},
			b: map[string]interface{}{
				"data": map[string]interface{}{},
			},
			ignoreKeys: nil,
			wantEqual:  false,
		},
		{
			name: "deeply nested arrays",
			a: map[string]interface{}{
				"matrix": []interface{}{
					[]interface{}{1, 2, 3},
					[]interface{}{4, 5, 6},
				},
			},
			b: map[string]interface{}{
				"matrix": []interface{}{
					[]interface{}{1, 2, 3},
					[]interface{}{4, 5, 7}, // different value
				},
			},
			ignoreKeys: nil,
			wantEqual:  false,
		},
	}

	// Run each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			equal, _ := compareJSON(tc.a, tc.b, tc.ignoreKeys)
			if equal != tc.wantEqual {
				t.Errorf("compareJSON() got = %v, want %v", equal, tc.wantEqual)
			}
		})
	}
}

func TestCompareJSONDifferences(t *testing.T) {
	// Test cases for checking specific difference reporting
	testCases := []struct {
		name           string
		a              interface{}
		b              interface{}
		ignoreKeys     []string
		expectedPaths  []string // expected path fragments in difference output
		unexpectedPath string   // path that should NOT be in the difference output
	}{
		{
			name: "simple value difference",
			a: map[string]interface{}{
				"name":  "test",
				"value": 123,
			},
			b: map[string]interface{}{
				"name":  "test",
				"value": 456,
			},
			ignoreKeys:    nil,
			expectedPaths: []string{"value"},
		},
		{
			name: "nested value difference",
			a: map[string]interface{}{
				"config": map[string]interface{}{
					"timeout": 30,
					"retries": 3,
				},
			},
			b: map[string]interface{}{
				"config": map[string]interface{}{
					"timeout": 60,
					"retries": 3,
				},
			},
			ignoreKeys:    nil,
			expectedPaths: []string{"config.timeout"},
		},
		{
			name: "missing key",
			a: map[string]interface{}{
				"name":    "test",
				"value":   123,
				"missing": "only in A",
			},
			b: map[string]interface{}{
				"name":  "test",
				"value": 123,
			},
			ignoreKeys:    nil,
			expectedPaths: []string{"missing"},
		},
		{
			name: "array element difference",
			a: map[string]interface{}{
				"items": []interface{}{"a", "b", "c"},
			},
			b: map[string]interface{}{
				"items": []interface{}{"a", "b", "d"},
			},
			ignoreKeys:    nil,
			expectedPaths: []string{"items[2]"},
		},
		{
			name: "ignored differences",
			a: map[string]interface{}{
				"id":    "A001", // ignored
				"value": 123,    // different
				"name":  "test", // same
			},
			b: map[string]interface{}{
				"id":    "B002", // ignored
				"value": 456,    // different
				"name":  "test", // same
			},
			ignoreKeys:     []string{"id"},
			expectedPaths:  []string{"value"},
			unexpectedPath: "id",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			equal, diffs := compareJSON(tc.a, tc.b, tc.ignoreKeys)

			// Should not be equal if we expect differences
			if len(tc.expectedPaths) > 0 && equal {
				t.Error("Objects should not be equal")
			}

			// Verify differences include expected paths
			for _, expectedPath := range tc.expectedPaths {
				found := false
				for _, diff := range diffs {
					if strings.Contains(diff, expectedPath) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected to find difference with path containing '%s', but not found in: %v", expectedPath, diffs)
				}
			}

			// Verify differences do not include unexpected paths
			if tc.unexpectedPath != "" {
				for _, diff := range diffs {
					if strings.Contains(diff, tc.unexpectedPath) {
						t.Errorf("Found unexpected path '%s' in differences: %s", tc.unexpectedPath, diff)
						break
					}
				}
			}
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	t.Run("formatValue", func(t *testing.T) {
		cases := []struct {
			input    interface{}
			expected string
		}{
			{nil, "null"},
			{123, "123"},
			{"hello", "\"hello\""},
			{"[missing]", "[missing]"}, // Special case
			{true, "true"},
			{[]interface{}{1, 2, 3}, "[1 2 3]"},
		}

		for _, c := range cases {
			result := formatValue(c.input)
			if result != c.expected {
				t.Errorf("formatValue(%v) = %s, want %s", c.input, result, c.expected)
			}
		}
	})

	t.Run("isIgnoredKey", func(t *testing.T) {
		ignoreKeys := []string{"id", "timestamp", "created_at"}

		if !isIgnoredKey("id", ignoreKeys) {
			t.Error("'id' should be ignored")
		}
		if isIgnoredKey("name", ignoreKeys) {
			t.Error("'name' should not be ignored")
		}
	})
}

func TestMapDifferences(t *testing.T) {
	t.Run("findMapDifferences", func(t *testing.T) {
		mapA := map[string]interface{}{
			"shared": "value",
			"onlyA":  "a value",
			"diff":   123,
		}
		mapB := map[string]interface{}{
			"shared": "value",
			"onlyB":  "b value",
			"diff":   456,
		}

		// Test with no ignored keys
		diffs := findMapDifferences(mapA, mapB, "root", nil)

		// Should find 3 differences: onlyA, onlyB, and diff
		if len(diffs) != 3 {
			t.Errorf("Expected 3 differences, got %d", len(diffs))
		}

		// Test with ignored key
		diffs = findMapDifferences(mapA, mapB, "root", []string{"diff"})

		// Should find 2 differences: onlyA and onlyB (diff is ignored)
		if len(diffs) != 2 {
			t.Errorf("Expected 2 differences after ignoring 'diff', got %d", len(diffs))
		}

		// Verify diff paths are correct
		foundOnlyA := false
		foundOnlyB := false

		for _, diff := range diffs {
			if diff.path == "root.onlyA" && diff.valueB == "[missing]" {
				foundOnlyA = true
			}
			if diff.path == "root.onlyB" && diff.valueA == "[missing]" {
				foundOnlyB = true
			}
		}

		if !foundOnlyA {
			t.Error("Missing expected difference for onlyA key")
		}
		if !foundOnlyB {
			t.Error("Missing expected difference for onlyB key")
		}
	})
}

func TestSliceDifferences(t *testing.T) {
	t.Run("findSliceDifferences", func(t *testing.T) {
		// Test different lengths
		sliceA := []interface{}{1, 2, 3}
		sliceB := []interface{}{1, 2}

		diffs := findSliceDifferences(sliceA, sliceB, "items", nil)
		if len(diffs) != 1 {
			t.Errorf("Expected 1 difference for slices with different lengths, got %d", len(diffs))
		}

		// Test same length but different values
		sliceA = []interface{}{1, 2, 3}
		sliceB = []interface{}{1, 2, 4}

		diffs = findSliceDifferences(sliceA, sliceB, "items", nil)
		if len(diffs) < 1 {
			t.Error("Expected at least one difference for slices with different values")
		}

		foundDiffAtIndex2 := false
		for _, diff := range diffs {
			if diff.path == "items[2]" && diff.valueA == 3 && diff.valueB == 4 {
				foundDiffAtIndex2 = true
				break
			}
		}

		if !foundDiffAtIndex2 {
			t.Error("Missing expected difference at index 2")
		}
	})
}

func TestExtractPath(t *testing.T) {
	testCases := []struct {
		name      string
		data      map[string]interface{}
		path      string
		expected  interface{}
		wantError bool
	}{
		{
			name: "empty path returns full data",
			data: map[string]interface{}{
				"key": "value",
			},
			path: "",
			expected: map[string]interface{}{
				"key": "value",
			},
			wantError: false,
		},
		{
			name: "simple object path",
			data: map[string]interface{}{
				"settings": map[string]interface{}{
					"timeout": float64(30),
				},
			},
			path:      "$.settings.timeout",
			expected:  float64(30), // JSON numbers are float64
			wantError: false,
		},
		{
			name: "nested object path",
			data: map[string]interface{}{
				"frontend": map[string]interface{}{
					"webserver": map[string]interface{}{
						"port": float64(8080),
					},
				},
			},
			path:      "$.frontend.webserver.port",
			expected:  float64(8080),
			wantError: false,
		},
		{
			name: "array index access",
			data: map[string]interface{}{
				"items": []interface{}{"first", "second", "third"},
			},
			path:      "$.items[1]",
			expected:  "second",
			wantError: false,
		},
		{
			name: "non-existent path",
			data: map[string]interface{}{
				"settings": map[string]interface{}{
					"timeout": 30,
				},
			},
			path:      "$.settings.nonexistent",
			wantError: true,
		},
		{
			name: "invalid JSONPath syntax",
			data: map[string]interface{}{
				"key": "value",
			},
			path:      "$[invalid syntax",
			wantError: true,
		},
		{
			name: "path returning multiple results",
			data: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"name": "item1"},
					map[string]interface{}{"name": "item2"},
				},
			},
			path:      "$.items[*].name",
			wantError: true, // Our implementation expects single result
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := extractPath(tc.data, tc.path)

			if tc.wantError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				// For maps, we need to compare them properly
				if expectedMap, ok := tc.expected.(map[string]interface{}); ok {
					if resultMap, ok := result.(map[string]interface{}); ok {
						if !compareMaps(expectedMap, resultMap) {
							t.Errorf("Expected %v, got %v", tc.expected, result)
						}
					} else {
						t.Errorf("Expected map but got %T", result)
					}
				} else if result != tc.expected {
					t.Errorf("Expected %v, got %v", tc.expected, result)
				}
			}
		})
	}
}

// Helper function to compare maps
func compareMaps(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			return false
		}
	}
	return true
}
