package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	// Parse command line arguments
	flag.Parse()

	// Check positional arguments
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s config.yaml\n", os.Args[0])
		os.Exit(2)
	}

	configPath := args[0]

	// Load configuration file
	config, err := LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config file: %v\n", err)
		os.Exit(2)
	}

	// Fetch endpoints and timeout
	endpointA, endpointB := config.GetDefaultEndpoints()
	timeout := config.GetTimeout()

	// Display endpoint information
	fmt.Sprintf("Comparing:\n  A: %s (%s)\n  B: %s (%s)\n\n",
		endpointA.Name, endpointA.URL, endpointB.Name, endpointB.URL)

	// Fetch JSON from both endpoints
	jsonA, jsonB, err := fetchEndpointData(endpointA, endpointB, timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}

	// Process JSON data based on path
	dataA, dataB, err := processJSONData(jsonA, jsonB, config.GetJSONPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}

	// Compare JSON
	compareAndReportResults(dataA, dataB, config.GetIgnoredKeys())
}

// Fetches JSON from both endpoints
func fetchEndpointData(endpointA, endpointB Endpoint, timeout int) (map[string]interface{}, map[string]interface{}, error) {
	// Fetch from endpoint A
	jsonA, err := fetchJSON(endpointA.URL, endpointA.Auth, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, nil, fmt.Errorf("Error fetching from endpoint A: %v", err)
	}

	// Fetch from endpoint B
	jsonB, err := fetchJSON(endpointB.URL, endpointB.Auth, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, nil, fmt.Errorf("Error fetching from endpoint B: %v", err)
	}

	return jsonA, jsonB, nil
}

// Extracts data based on JSON path
func processJSONData(jsonA, jsonB map[string]interface{}, jsonPath string) (interface{}, interface{}, error) {
	// Use entire JSON if no path specified
	if jsonPath == "" {
		return jsonA, jsonB, nil
	}

	// Extract from endpoint A
	extractedA, err := extractPath(jsonA, jsonPath)
	if err != nil {
		return nil, nil, fmt.Errorf("Error extracting JSON path from endpoint A: %v", err)
	}

	// Extract from endpoint B
	extractedB, err := extractPath(jsonB, jsonPath)
	if err != nil {
		return nil, nil, fmt.Errorf("Error extracting JSON path from endpoint B: %v", err)
	}

	return extractedA, extractedB, nil
}

// Compares data and reports results
func compareAndReportResults(dataA, dataB interface{}, ignoredKeys []string) {
	equal, diffs := compareJSON(dataA, dataB, ignoredKeys)

	if !equal {
		fmt.Println("Difference found:")
		for _, diff := range diffs {
			fmt.Println(diff)
		}
		fmt.Println("\nEndpoints contain different configuration.")
		os.Exit(1)
	}

	fmt.Println("Endpoints contain identical configuration.")
	os.Exit(0)
}
