# restcompre

`rest-compare` helps you verify that configuration data is consistent across different environments or deployments. It fetches JSON data from two REST API endpoints, optionally extracts specific sections using a JSON path, and compares them while ignoring specified keys.

## Features

- Compare JSON data from two different API endpoints
- Extract specific sections of JSON using dot notation path
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
  jsonPath: "frontend.config"   # Optional JSON path to compare
```

### Configuration Options

#### Endpoints

- `name`: Descriptive name for the endpoint
- `url`: Full URL of the API endpoint
- `auth`: Optional authentication header value

#### Settings

- `timeout`: HTTP request timeout in seconds (default: 30)
- `ignoredKeys`: List of keys to exclude from comparison
- `jsonPath`: Dot notation path to extract from JSON (e.g., `frontend.config`)

## JSON Path

The JSON path uses dot notation to navigate through nested objects. For example:

```
frontend.webserver.routes
```

This will extract the `routes` property from within the `webserver` property of the `frontend` object.

Array notation is also supported:

```
items[0].name
```

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
