package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cardinalnsk/ginit/internal/generator"
	"github.com/cardinalnsk/ginit/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Флаги для non-interactive режима
	name := flag.String("name", "", "Project name")
	module := flag.String("module", "", "Go module name")
	dir := flag.String("dir", "", "Custom directory name")
	projectType := flag.String("type", "cli", "Project type: cli, web, or library")
	noVCS := flag.Bool("no-vcs", false, "Skip VCS initialization")
	nonInteractive := flag.Bool("non-interactive", false, "Disable interactive mode")

	flag.Parse()

	// Non-interactive режим
	if *nonInteractive || (*name != "" && len(flag.Args()) > 0) {
		runNonInteractive(name, module, dir, projectType, noVCS, flag.Args())
		return
	}

	// Interactive режим с BubbleTea
	runInteractive()
}

func runNonInteractive(name *string, module *string, dir *string, projectType *string, noVCS *bool, args []string) {
	// Логика как раньше
	if *name == "" && len(args) > 0 {
		*name = args[0]
	}

	if *name == "" {
		printUsage()
		os.Exit(1)
	}

	if *module == "" {
		*module = *name
	}

	if *dir == "" {
		*dir = *name
	}

	config := generator.Config{
		ProjectName: *name,
		ModuleName:  *module,
		Directory:   *dir,
		ProjectType: *projectType,
		InitVCS:     !*noVCS,
	}

	err := generator.InitProject(config)
	if err != nil {
		log.Fatalf("Error initializing project: %v", err)
	}

	printSuccessMessage(config)
}

func runInteractive() {
	// Запускаем TUI
	p := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running interactive mode: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("🚀 Go Project Initializer")
	fmt.Println("")
	fmt.Println("Usage: ginit [flags] <project-name>")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  -name string          Project name")
	fmt.Println("  -module string        Go module name (default: project name)")
	fmt.Println("  -dir string           Custom directory name (default: project name)")
	fmt.Println("  -type string          Project type: cli, web, or library (default: cli)")
	fmt.Println("  -no-vcs               Skip VCS initialization")
	fmt.Println("  -non-interactive      Disable interactive mode")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  ginit                             # Interactive mode")
	fmt.Println("  ginit my-project                  # Quick start")
	fmt.Println("  ginit -name=myapp -module=github.com/user/myapp")
	fmt.Println("  ginit my-project -no-vcs -non-interactive")
}

func printSuccessMessage(config generator.Config) {
	absPath, _ := filepath.Abs(config.Directory)

	style := tui.DefaultStyle()

	fmt.Println("")
	fmt.Println(style.SuccessIcon.Render("🎉 ") + style.SuccessText.Render("Project initialized successfully!"))
	fmt.Println("")
	fmt.Println(style.Label.Render("📁 Project:   ") + style.Value.Render(config.ProjectName))
	fmt.Println(style.Label.Render("📦 Module:    ") + style.Value.Render(config.ModuleName))
	fmt.Println(style.Label.Render("📂 Directory: ") + style.Value.Render(absPath))
	fmt.Println("")
	fmt.Println(style.Section.Render("🚀 Next steps:"))
	fmt.Println("")

	steps := []string{
		"cd " + config.Directory,
		"go mod tidy",
		"go build -o bin/" + config.ProjectName + " ./cmd/" + config.ProjectName,
		"./bin/" + config.ProjectName,
	}

	if config.InitVCS {
		steps = append(steps, "git add .", "git commit -m \"Initial commit\"")
	}

	for i, step := range steps {
		emoji := "📂"
		switch i {
		case 1:
			emoji = "🔧"
		case 2:
			emoji = "🏗️"
		case 3:
			emoji = "▶️"
		case 4:
			emoji = "📝"
		case 5:
			emoji = "📝"
		}
		fmt.Println(style.Label.Render("  "+emoji+" ") + style.Code.Render(step))
	}

	fmt.Println("")
	fmt.Println(style.Section.Render("💡 Tips:"))
	tips := []string{
		"Use 'go run ./cmd/" + config.ProjectName + "' for quick testing",
		"Check out the README.md for more details",
		"Modify internal/config for your needs",
	}

	for _, tip := range tips {
		fmt.Println(style.Label.Render("  • ") + style.Tip.Render(tip))
	}

	fmt.Println("")
	fmt.Println(style.SuccessText.Render("Happy coding!") + " 👨‍💻👩‍💻")
	fmt.Println("")
}
