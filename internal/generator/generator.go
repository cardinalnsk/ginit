package generator

import (
	"fmt"
	"os"
	"os/exec"
)

type Config struct {
	ProjectName string
	ModuleName  string
	Directory   string
	ProjectType string
	InitVCS     bool
}

func InitProject(config Config) error {
	// Создаем директорию проекта
	if err := os.MkdirAll(config.Directory, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(config.Directory); err != nil {
		return fmt.Errorf("failed to change directory: %w", err)
	}

	// Debug: print current directory
	cwd, _ := os.Getwd()
	fmt.Printf("DEBUG: Current directory after chdir: %s\n", cwd)
	fmt.Printf("DEBUG: Config - Directory: %s, ProjectName: %s, ModuleName: %s\n", config.Directory, config.ProjectName, config.ModuleName)

	if err := initGoMod(config.ModuleName); err != nil {
		return err
	}

	if err := addDependencies(config.ProjectType); err != nil {
		return err
	}

	if err := createProjectStructure(config.ProjectName, config.ProjectType); err != nil {
		return err
	}

	if err := createMainFile(config.ProjectName, config.ProjectType); err != nil {
		return err
	}

	if err := createConfigFiles(config.ProjectType); err != nil {
		return err
	}

	if err := createLoggerFiles(config.ProjectType); err != nil {
		return err
	}

	if err := createAdditionalFiles(config.ProjectName, config.ProjectType); err != nil {
		return err
	}

	if err := createReadme(config.ProjectName); err != nil {
		return err
	}

	if config.InitVCS {
		if err := initVCS(); err != nil {
			return err
		}
	}

	return nil
}

func initGoMod(moduleName string) error {
	// Check if Go is available in PATH
	if _, err := exec.LookPath("go"); err != nil {
		// Go is not available, create go.mod file manually
		goModContent := fmt.Sprintf("module %s\n\ngo 1.21\n", moduleName)
		return os.WriteFile("go.mod", []byte(goModContent), 0644)
	}
	
	// Go is available, use go mod init
	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func addDependencies(projectType string) error {
	// Check if Go is available in PATH
	if _, err := exec.LookPath("go"); err != nil {
		fmt.Println("Go not found, skipping dependency installation")
		return nil
	}
	
	var dependencies []string
	
	switch projectType {
	case "cli":
		dependencies = []string{
			"github.com/caarlos0/env/v11",
		}
	case "web":
		dependencies = []string{
			"github.com/caarlos0/env/v11",
		}
	case "library":
		// Library projects typically don't need external dependencies
		dependencies = []string{}
	}
	
	for _, dep := range dependencies {
		cmd := exec.Command("go", "get", dep)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add dependency %s: %w", dep, err)
		}
	}
	
	return nil
}

func initVCS() error {
	if _, err := exec.LookPath("git"); err != nil {
		fmt.Println("Git not found, skipping VCS initialization")
		return nil
	}

	cmd := exec.Command("git", "init")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to init git: %w", err)
	}

	if err := createGitignore(); err != nil {
		return err
	}

	return nil
}
