# ginit - Go Project Generator

[English](README.md) | [Русский](README_RU.md)

**ginit** is a command-line tool for quickly creating structured Go projects with pre-configured templates for various application types.

## 🚀 Features

- **Interactive TUI interface** - user-friendly text-based user interface
- **CLI mode** - project creation via command line
- **Supported project types**:
  - **CLI** - command-line applications
  - **Web** - web applications with HTTP server
  - **Library** - libraries and packages
- **Automatic Git repository initialization**
- **Pre-configured project structure**
- **Ready-to-use configuration and logging templates**

## 📦 Installation

### Prerequisites

- Go 1.21+
- Git (for repository initialization)

### Install from source

```bash
git clone https://github.com/cardinalnsk/ginit.git
cd ginit
go build -o ginit ./cmd/ginit
sudo mv ginit /usr/local/bin/
```

### Install via go install

```bash
go install github.com/cardinalnsk/ginit/cmd/ginit@latest
```

## 🎯 Usage

### Interactive mode (TUI)

```bash
ginit
```

The interactive interface will guide you through the project creation process:

1. **Project name** - your project's name
2. **Module name** - Go module name (e.g.: github.com/user/project)
3. **Directory** - path for project creation
4. **Project type** - CLI, Web, or Library
5. **Git initialization** - create Git repository

### Command line (CLI)

```bash
# Create CLI project
ginit -name my-cli-app -module github.com/user/my-cli-app -dir ./my-cli-app -type cli

# Create Web project
ginit -name my-web-app -module github.com/user/my-web-app -dir ./my-web-app -type web

# Create Library project
ginit -name my-lib -module github.com/user/my-lib -dir ./my-lib -type library
```

#### Command line parameters

- `-name` - project name (required)
- `-module` - Go module name (required)
- `-dir` - directory for project creation (required)
- `-type` - project type: cli, web, library (required)
- `-vcs` - initialize Git repository (true/false, default: true)

## 🏗️ Project Structure

### CLI project

```
my-cli-app/
├── cmd/
│   └── my-cli-app/
│       └── main.go
├── internal/
│   ├── app/
│   │   └── app.go
│   ├── cli/
│   │   ├── cli.go
│   │   ├── commands.go
│   │   └── config.go
│   └── config/
│       └── config.go
├── pkg/
│   └── logger/
│       └── logger.go
├── go.mod
├── go.sum
└── README.md
```

### Web project

```
my-web-app/
├── cmd/
│   └── web/
│       └── main.go
├── internal/
│   ├── app/
│   │   └── app.go
│   ├── config/
│   │   └── config.go
│   ├── handlers/
│   │   └── handlers.go
│   └── server/
│       └── server.go
├── pkg/
│   └── logger/
│       └── logger.go
├── static/
├── templates/
├── go.mod
├── go.sum
└── README.md
```

### Library project

```
my-lib/
├── internal/
│   └── config/
│       └── config.go
├── pkg/
│   └── mylib/
│       └── mylib.go
├── go.mod
├── go.sum
└── README.md
```

## 🔧 Generated Files

### Core files

- **main.go** - application entry point
- **go.mod** - Go module file
- **README.md** - project documentation
- **.gitignore** - Git ignored files

### Configuration

- **config.go** - application configuration with env variables support
- Singleton system for configuration access
- Support for different configuration types for different project types

### Logging

- **logger.go** - logging utilities based on slog
- Support for text and JSON formats
- Configurable log levels

## 🎨 TUI Interface

The interactive interface provides:

- **Step-by-step wizard** for project creation
- **Input validation** - data correctness checking
- **Project structure preview**
- **Navigation keys**:
  - `Enter` - next step
  - `Tab` / `Shift+Tab` - switch between fields
  - `Ctrl+C` - exit
  - `h`/`l` or arrow keys - select project type
  - `y`/`n` - choose Git initialization

## 🛠️ Development

### Building the project

```bash
# Build ginit
go build -o ginit ./cmd/ginit

# Testing
go test ./...
```

### ginit code structure

```
ginit/
├── cmd/
│   └── ginit/
│       └── main.go          # Entry point
├── internal/
│   ├── generator/
│   │   ├── generator.go     # Project generation logic
│   │   ├── templates.go     # File templates
│   │   └── utils.go         # Utility functions
│   └── tui/
│       ├── model.go         # TUI model
│       └── styles.go        # Interface styles
├── go.mod
└── README.md
```

## 🐛 Troubleshooting

### Project creation issues

**Problem**: Project is not created completely
**Solution**: Make sure you don't press Enter multiple times during project creation


### Dependency issues

**Problem**: Import errors after project creation
**Solution**: Run `go mod tidy` in the created project directory

## 🤝 Contributing

We welcome contributions to the project! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

Distributed under the MIT License. See `LICENSE` file for details.

## 📞 Contact

- GitHub: [cardinalnsk](https://github.com/cardinalnsk)

---

⭐ If this project was helpful, please give it a star on GitHub!
