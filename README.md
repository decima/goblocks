# Goblocks

A lightweight REST API for managing hierarchical content blocks with flexible storage backends.

## Features

- **RESTful API** for block management (CRUD operations)
- **Flexible Storage**: File system or in-memory storage
- **Path Validation**: Protection against path traversal attacks
- **Content-Type Validation**: MIME type validation for uploaded content
- **Upload Size Limits**: Configurable maximum upload size (default: 10MB)
- **Structured Logging**: Pretty-printed logs in development, JSON logs in production
- **Dependency Injection**: Clean architecture using Uber FX

## Architecture

The project follows a modular architecture:

```
app/
├── config/          # Configuration management
├── services/        # Business logic and storage implementations
│   └── blocks/      # Block management services
├── web/             # HTTP layer
│   ├── controllers/ # HTTP handlers
│   ├── router.go    # Route registration
│   └── server.go    # HTTP server
└── App.go           # Application metadata
```

## Installation

```bash
# Using Make
make install

# Or manually
go mod download
```

## Configuration

Configuration is managed via `goblocks.yaml` and environment variables:

```yaml
http:
  host: 0.0.0.0
  port: 8000
  max_upload_size: 10485760  # 10MB in bytes

blocks:
  storage:
    type: fs              # "fs" or "inMemory"
    path: ./data/         # Only for fs storage
```

Environment variables override config file (use `_` separator):
```bash
HTTP_PORT=9000 go run .
BLOCKS_STORAGE_TYPE=inMemory go run .
```

## Usage

### Quick Start with Make

```bash
# Install dependencies
make install

# Run in development mode
make dev

# Build the application
make build

# Run tests
make test

# Build and run
make run
```

### Development

```shell
# Using Make
make dev

# Or manually
DEBUG=1 go run .
```

### Production

```shell
# Using Make
make build
./bin/goblocks

# Or manually
go build -o goblocks
./goblocks
```

### Available Make Commands

**Build:**
- `make build` - Build the application
- `make build-linux` - Build for Linux
- `make build-windows` - Build for Windows
- `make build-mac` - Build for macOS
- `make build-all` - Build for all platforms

**Development:**
- `make dev` - Run in development mode with DEBUG=1
- `make run` - Build and run the application
- `make fmt` - Format Go code
- `make vet` - Run go vet
- `make lint` - Run golangci-lint

**Testing:**
- `make test` - Run all tests
- `make test-verbose` - Run tests in verbose mode
- `make test-coverage` - Generate coverage report (HTML)
- `make test-race` - Run tests with race detector
- `make test-bench` - Run benchmarks

**Utilities:**
- `make check` - Run format, vet, and tests
- `make ci` - Run full CI pipeline
- `make clean` - Remove build artifacts
- `make info` - Display project information
- `make help` - Show all available commands

## API Reference

### Create/Update Block

```http
PUT /blocks/{path}
Content-Type: <mime-type>

<content>
```

**Example:**
```http
PUT /blocks/my/document
Content-Type: text/plain

Hello, World!
```

### Get Block (Metadata)

```http
GET /blocks/{path}
```

Returns block metadata including children and content type.

**Response:**
```json
{
  "path": "my/document",
  "type": "text/plain",
  "size": 13,
  "children": []
}
```

### Get Block (Raw Content)

```http
GET /blocks/{path}?raw
```

Returns the raw content with original Content-Type header.

### Delete Block

```http
DELETE /blocks/{path}
```

Deletes the block and all its children.

## Security Features

- **Path Traversal Protection**: Paths are validated and sanitized
- **Maximum Path Depth**: Limited to 10 levels
- **Content-Type Validation**: MIME types must be valid
- **Upload Size Limits**: Configurable maximum file size
- **File Permissions**: Content files created with 0644 permissions

## Storage Backends

### File System (`fs`)

Stores blocks as directories with `.content` files containing JSON metadata:

```
data/
└── my/
    └── document/
        └── .content
```

### In-Memory (`inMemory`)

Stores blocks in memory using Go's `sync.Map`. Useful for testing or ephemeral data.

## Development

### Project Structure

- `main.go` - Application entry point
- `app/` - Core application code
- `libraries/` - Reusable utilities (logger, etc.)
- `tests/` - HTTP test files

### Running Tests

```bash
# Using Make (recommended)
make test                 # Run all tests
make test-verbose         # Run with verbose output
make test-coverage        # Generate HTML coverage report
make test-race            # Run with race detector

# Or manually
go test ./app/...
go test -cover ./app/...
go test -v ./app/...
```

**Test Coverage:**
- Block managers: 76.5%
- HTTP controllers: 80.8%

The test suite includes:
- Path validation tests (security)
- Content-Type validation tests
- In-memory storage tests
- File system storage tests
- HTTP integration tests
- Error handling tests

### Logging

The application uses structured logging with different formats:

- **Development** (`DEBUG=1`): Colorized, pretty-printed logs
- **Production**: JSON logs for machine parsing

## Error Handling

The API returns appropriate HTTP status codes:

- `200 OK` - Successful GET
- `202 Accepted` - Successful PUT
- `204 No Content` - Successful DELETE
- `403 Forbidden` - Invalid path, content type, or permissions
- `404 Not Found` - Block doesn't exist
- `422 Unprocessable Entity` - Other errors

## Dependencies

- [Uber FX](https://github.com/uber-go/fx) - Dependency injection
- [Viper](https://github.com/spf13/viper) - Configuration management
- [slog](https://pkg.go.dev/log/slog) - Structured logging

## License

See LICENSE file for details.