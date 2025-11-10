package generator

import (
	"fmt"
	"os"
)

func createProjectStructure(projectName, projectType string) error {
	var dirs []string

	switch projectType {
	case "cli":
		dirs = []string{
			"cmd/" + projectName,
			"internal/config",
			"internal/cli",
			"internal/commands",
			"pkg/logger",
			"pkg/utils",
			"pkg/version",
		}
	case "web":
		dirs = []string{
			"cmd/" + projectName,
			"internal/config",
			"internal/app",
			"internal/handlers",
			"internal/middleware",
			"internal/models",
			"internal/repository",
			"internal/service",
			"pkg/logger",
			"pkg/utils",
			"pkg/database",
			"api/",
			"web/static/",
			"web/templates/",
		}
	case "library":
		dirs = []string{
			"internal/",
			"pkg/",
			"examples/",
			"docs/",
		}
	default:
		dirs = []string{
			"cmd/" + projectName,
			"internal/config",
			"internal/app",
			"pkg/logger",
			"pkg/utils",
		}
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

func createMainFile(projectName, projectType string) error {
	var mainTemplate string

	switch projectType {
	case "cli":
		mainTemplate = `package main

import (
	"context"
	"fmt"
	"os"

	"{{.Module}}/internal/config"
	"{{.Module}}/internal/cli"
	"{{.Module}}/pkg/logger"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Инициализация логгера
	log := logger.New(cfg.LogLevel)

	ctx := context.Background()
	log.InfoContext(ctx, "Starting {{.ProjectName}} CLI...")
	
	// Запуск CLI приложения
	if err := cli.Run(ctx, cfg); err != nil {
		log.ErrorContext(ctx, "CLI execution failed", "error", err)
		os.Exit(1)
	}
	
	log.InfoContext(ctx, "CLI application stopped")
}
`
	case "web":
		mainTemplate = `package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"{{.Module}}/internal/config"
	"{{.Module}}/internal/app"
	"{{.Module}}/pkg/logger"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Инициализация логгера
	log := logger.New(cfg.LogLevel)

	ctx := context.Background()
	log.InfoContext(ctx, "Starting {{.ProjectName}} web server...")
	
	// Создание и запуск приложения
	application := app.New(cfg, log)
	
	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		if err := application.Run(ctx); err != nil {
			log.ErrorContext(ctx, "Application failed", "error", err)
			os.Exit(1)
		}
	}()
	
	<-stop
	log.InfoContext(ctx, "Shutting down server...")
	
	if err := application.Shutdown(ctx); err != nil {
		log.ErrorContext(ctx, "Graceful shutdown failed", "error", err)
	}
	
	log.InfoContext(ctx, "Server stopped")
}
`
	case "library":
		mainTemplate = `package main

import (
	"fmt"
	"os"

	"{{.Module}}/pkg/version"
)

func main() {
	fmt.Printf("{{.ProjectName}} version %s\n", version.Version)
	fmt.Println("This is a library project. Run 'go test ./...' to run tests.")
	os.Exit(0)
}
`
	default:
		mainTemplate = `package main

import (
	"context"
	"fmt"
	"os"

	"{{.Module}}/internal/config"
	"{{.Module}}/pkg/logger"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Инициализация логгера
	log := logger.New(cfg.LogLevel)

	ctx := context.Background()
	log.InfoContext(ctx, "Starting {{.ProjectName}}...")
	
	fmt.Println("Hello, World!")
	
	log.InfoContext(ctx, "Application stopped")
}
`
	}

	data := struct {
		Module      string
		ProjectName string
	}{
		Module:      projectName,
		ProjectName: projectName,
	}

	return createFileFromTemplate("cmd/"+projectName+"/main.go", mainTemplate, data)
}

func createConfigFiles(projectType string) error {
	// Skip config files for library projects as they don't have internal/config directory
	if projectType == "library" {
		return nil
	}

	var configGo string

	switch projectType {
	case "cli":
		configGo = `package config

import (
	"sync"
)

type Config struct {
	LogLevel  string
	Verbose   bool
	ConfigPath string
}
}

var (
	instance *Config
	once     sync.Once
)

func Load() (*Config, error) {
	var err error
	once.Do(func() {
		instance = &Config{
			LogLevel:   "info",
			Verbose:    false,
			ConfigPath: "config.yaml",
		}
		// Здесь будет загрузка из env/config file
	})
	return instance, err
}
`
	case "web":
		configGo = `package config

import (
	"sync"
)

type Config struct {
	HTTPPort  string
	LogLevel  string
	DBURL     string
	RedisURL  string
	JWTSecret string
}

var (
	instance *Config
	once     sync.Once
)

func Load() (*Config, error) {
	var err error
	once.Do(func() {
		instance = &Config{
			HTTPPort:  ":8080",
			LogLevel:  "info",
			DBURL:     "postgres://user:pass@localhost:5432/db",
			RedisURL:  "redis://localhost:6379",
			JWTSecret: "your-secret-key",
		}
		// Здесь будет загрузка из env/config file
	})
	return instance, err
}
`
	case "library":
		configGo = `package config

import (
	"sync"
)

type Config struct {
	LogLevel string
}

var (
	instance *Config
	once     sync.Once
)

func Load() (*Config, error) {
	var err error
	once.Do(func() {
		instance = &Config{
			LogLevel: "info",
		}
		// Здесь будет загрузка из env/config file
	})
	return instance, err
}
`
	default:
		configGo = `package config

import (
	"sync"
)

type Config struct {
	HTTPPort  string
	LogLevel  string
	DBURL     string
}

var (
	instance *Config
	once     sync.Once
)

func Load() (*Config, error) {
	var err error
	once.Do(func() {
		instance = &Config{
			HTTPPort: ":8080",
			LogLevel: "info",
		}
		// Здесь будет загрузка из env/config file
	})
	return instance, err
}
`
	}

	return os.WriteFile("internal/config/config.go", []byte(configGo), 0644)
}

func createAdditionalFiles(projectName, projectType string) error {
	switch projectType {
	case "cli":
		if err := createCLIFiles(projectName); err != nil {
			return err
		}
	case "web":
		if err := createWebFiles(projectName); err != nil {
			return err
		}
	case "library":
		if err := createLibraryFiles(projectName); err != nil {
			return err
		}
	}
	return nil
}

func createLoggerFiles(projectType string) error {
	// Skip logger files for library projects as they don't have pkg/logger directory
	if projectType == "library" {
		return nil
	}

	loggerGo := `package logger

import (
	"log/slog"
	"os"
)

// New создает новый логгер slog с указанным уровнем
func New(level string) *slog.Logger {
	var logLevel slog.Level

	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}

// NewJSON создает логгер с JSON форматом
func NewJSON(level string) *slog.Logger {
	var logLevel slog.Level

	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}
`

	return createFileFromTemplate("pkg/logger/logger.go", loggerGo, nil)
}

func createCLIFiles(projectName string) error {
	// CLI package
	cliGo := `package cli

import (
	"context"
	"flag"
	"fmt"
	"os"

	"{{.Module}}/internal/config"
	"{{.Module}}/internal/commands"
	"{{.Module}}/pkg/logger"
)

func Run(ctx context.Context, cfg *config.Config) error {
	var (
		verbose bool
		version bool
	)

	flag.BoolVar(&verbose, "verbose", cfg.Verbose, "Enable verbose output")
	flag.BoolVar(&version, "version", false, "Show version information")
	flag.Parse()

	log := logger.New(cfg.LogLevel)

	if version {
		fmt.Println("{{.ProjectName}} v1.0.0")
		return nil
	}

	if verbose {
		log.InfoContext(ctx, "Verbose mode enabled")
	}

	// Execute command
	if flag.NArg() > 0 {
		cmd := flag.Arg(0)
		return commands.Execute(ctx, cmd, flag.Args()[1:], cfg, log)
	}

	// Default command
	return commands.Default(ctx, cfg, log)
}
`

	data := struct {
		Module      string
		ProjectName string
	}{
		Module:      projectName,
		ProjectName: projectName,
	}

	if err := createFileFromTemplate("internal/cli/cli.go", cliGo, data); err != nil {
		return err
	}

	// Commands package
	commandsGo := `package commands

import (
	"context"
	"fmt"

	"{{.Module}}/internal/config"
	"{{.Module}}/pkg/logger"
)

func Execute(ctx context.Context, cmd string, args []string, cfg *config.Config, log *slog.Logger) error {
	switch cmd {
	case "help":
		return Help(ctx, args, cfg, log)
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

func Default(ctx context.Context, cfg *config.Config, log *slog.Logger) error {
	log.InfoContext(ctx, "Running default command")
	fmt.Println("Welcome to {{.ProjectName}}!")
	fmt.Println("Use '{{.ProjectName}} help' for available commands.")
	return nil
}

func Help(ctx context.Context, args []string, cfg *config.Config, log *slog.Logger) error {
	fmt.Println("Available commands:")
	fmt.Println("  help    - Show this help message")
	fmt.Println("  version - Show version information")
	return nil
}
`

	return createFileFromTemplate("internal/commands/commands.go", commandsGo, data)
}

func createWebFiles(projectName string) error {
	// App package
	appGo := `package app

import (
	"context"
	"log/slog"
	"net/http"

	"{{.Module}}/internal/config"
	"{{.Module}}/internal/handlers"
)

type App struct {
	config *config.Config
	log    *slog.Logger
	server *http.Server
}

func New(cfg *config.Config, log *slog.Logger) *App {
	return &App{
		config: cfg,
		log:    log,
	}
}

func (a *App) Run(ctx context.Context) error {
	a.log.InfoContext(ctx, "Starting HTTP server", "port", a.config.HTTPPort)
	
	handler := handlers.New(a.config, a.log)
	
	a.server = &http.Server{
		Addr:    a.config.HTTPPort,
		Handler: handler,
	}
	
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	if a.server != nil {
		return a.server.Shutdown(ctx)
	}
	return nil
}
`

	data := struct {
		Module      string
		ProjectName string
	}{
		Module:      projectName,
		ProjectName: projectName,
	}

	if err := createFileFromTemplate("internal/app/app.go", appGo, data); err != nil {
		return err
	}

	// Handlers package
	handlersGo := `package handlers

import (
	"log/slog"
	"net/http"

	"{{.Module}}/internal/config"
)

type Handler struct {
	config *config.Config
	log    *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger) *Handler {
	return &Handler{
		config: cfg,
		log:    log,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.InfoContext(r.Context(), "HTTP request", "method", r.Method, "path", r.URL.Path)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(` + "`{\"message\": \"Hello from {{.ProjectName}}\"}`" + `))
}
`

	return createFileFromTemplate("internal/handlers/handlers.go", handlersGo, data)
}

func createLibraryFiles(projectName string) error {
	// Create version package directory
	if err := os.MkdirAll("pkg/version", 0755); err != nil {
		return fmt.Errorf("failed to create version directory: %w", err)
	}

	// Version package
	versionGo := `package version

// Version of the library
const Version = "v1.0.0"
`

	if err := os.WriteFile("pkg/version/version.go", []byte(versionGo), 0644); err != nil {
		return err
	}

	// Example usage
	exampleGo := `package main

import (
	"fmt"

	"{{.Module}}/pkg/version"
)

func main() {
	fmt.Printf("Using {{.ProjectName}} version: %s\n", version.Version)
	fmt.Println("This is an example of how to use the library.")
}
`

	data := struct {
		Module      string
		ProjectName string
	}{
		Module:      projectName,
		ProjectName: projectName,
	}

	return createFileFromTemplate("examples/example.go", exampleGo, data)
}
