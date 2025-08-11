# rest-compre

`rest-compare` helps you verify that configuration data is consistent across different environments or deployments. It fetches JSON data from two REST API endpoints, optionally extracts specific sections using JSONPath expressions, and compares them while ignoring specified keys.

## Features

- Compare JSON data from two different API endpoints
- Extract specific sections of JSON using standard JSONPath expressions
- Ignore specified keys during comparison
- Detailed output showing exact differences
- Configurable via YAML configuration file
- Support for authentication headers
- Customizable HTTP timeout

## Installation

Building from source:

```bash
$ go build
```

## Usage

```bash
$ rest-compare config.yaml
```

### Exit Codes

- `0`: Endpoints contain identical configuration
- `1`: Endpoints contain different configuration
- `2`: Error occurred (invalid configuration, connection error, etc.)

## Configuration

Create a `config.yaml` file with the following structure:

```yaml
endpoints:
  - name: "Production"
    url: "https://api1.example.com/config"
    auth: "Basic dXNlcjpwYXNz"  # Optional
  - name: "Staging"
    url: "https://api2.example.com/config"
    auth: "Bearer token123"     # Optional

settings:
  timeout: 60                   # HTTP timeout in seconds
  ignoredKeys:                  # Keys to ignore during comparison
    - "id"
    - "timestamp"
    - "updated_at"
  jsonPath: "$.frontend.config" # Optional JSONPath expression
```

### Configuration Options

#### Endpoints

- `name`: Descriptive name for the endpoint
- `url`: Full URL of the API endpoint
- `auth`: Optional authentication header value

#### Settings

- `timeout`: HTTP request timeout in seconds (default: 30)
- `ignoredKeys`: List of keys to exclude from comparison
- `jsonPath`: JSONPath expression to extract from JSON (e.g., `$.frontend.config`)

## JSONPath

This tool supports standard JSONPath expressions as defined in RFC 9535. The JSONPath expression must return a single result for comparison. If multiple results are returned, an error will be raised.

## Examples

### Example 1: Compare Entire Response

```yaml
endpoints:
  - name: "Production"
    url: "https://api.example.com/v1/config"
  - name: "Staging"
    url: "https://staging-api.example.com/v1/config"

settings:
  ignoredKeys:
    - "timestamp"
    - "version"
```

### Example 2: Compare Specific Configuration Section

```yaml
endpoints:
  - name: "Service A"
    url: "https://service-a.example.com/config"
  - name: "Service B"
    url: "https://service-b.example.com/config"

settings:
  jsonPath: "$.database.connections.primary"
  ignoredKeys:
    - "connectionId"
```

### Example 3: Compare Feature Flags

```yaml
endpoints:
  - name: "US Region"
    url: "https://us.example.com/features"
  - name: "EU Region"
    url: "https://eu.example.com/features"

settings:
  jsonPath: "$.features[?@.enabled == true]"
  timeout: 30
```

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
