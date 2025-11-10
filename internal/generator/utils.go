package generator

import (
	"os"
	"path/filepath"
	"text/template"
)

func createFileFromTemplate(path, tmpl string, data any) error {
	// Создаем директорию если нужно
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	t := template.Must(template.New("").Parse(tmpl))
	return t.Execute(file, data)
}

func createGitignore() error {
	gitignore := `# Binaries
bin/
dist/

# Dependencies
vendor/

# Environment files
.env
.env.local

# IDE
.vscode/
.idea/

# Build artifacts
*.exe
*.dll
*.so
*.dylib

# Test output
coverage.txt
profile.out

# Go workspace
go.work
go.work.sum
`

	return os.WriteFile(".gitignore", []byte(gitignore), 0644)
}

func createReadme(projectName string) error {
	readme := `# ` + projectName + `

A Go project generated with ginit.

## Features

- Modern Go project structure
- Configuration management
- Structured logging with slog
- Ready for production

## Getting Started

### Prerequisites
- Go 1.21+ (for slog support)

### Installation

1. Build the project:
` + "```bash" + `
go build -o bin/` + projectName + ` ./cmd/` + projectName + `
` + "```" + `

2. Run:
` + "```bash" + `
./bin/` + projectName + `
` + "```" + `

### Development

Run with hot reload (if you have air/gin installed):
` + "```bash" + `
air
# or
gin -i run cmd/` + projectName + `/main.go
` + "```" + `

## Project Structure

` + "```" + `
` + projectName + `/
├── cmd/` + projectName + `/main.go
├── internal/
│   ├── config/     # Configuration management
│   └── app/        # Application logic
├── pkg/
│   ├── logger/     # slog-based logging
│   └── utils/      # Shared utilities
` + "```" + `

## Configuration

The application uses environment variables for configuration:

- ` + "`HTTP_PORT`" + `: Port for HTTP server (default: :8080)
- ` + "`LOG_LEVEL`" + `: Log level (debug, info, warn, error) (default: info)
- ` + "`DB_URL`" + `: Database connection string

## Logging

Uses Go's built-in slog package for structured logging.

Example:
` + "```go" + `
log.InfoContext(ctx, "user logged in", "user_id", userID, "ip", ipAddress)
` + "```" + `
`

	return os.WriteFile("README.md", []byte(readme), 0644)
}
