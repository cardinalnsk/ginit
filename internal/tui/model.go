package tui

import (
	"fmt"
	"strings"

	"github.com/cardinalnsk/ginit/internal/generator"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	step            int
	projectName     textinput.Model
	moduleName      textinput.Model
	directory       textinput.Model
	projectType     string
	initVCS         bool
	quitting        bool
	success         bool
	error           error
	config          generator.Config
	creatingProject bool
}

func NewModel() Model {
	// Инициализируем поля ввода
	project := textinput.New()
	project.Placeholder = "my-awesome-app"
	project.Focus()
	project.CharLimit = 50
	project.Width = 50
	project.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))

	module := textinput.New()
	module.Placeholder = "github.com/username/my-awesome-app"
	module.CharLimit = 100
	module.Width = 50
	module.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))

	dir := textinput.New()
	dir.Placeholder = "my-awesome-app"
	dir.CharLimit = 100
	dir.Width = 50
	dir.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))

	return Model{
		step:        0,
		projectName: project,
		moduleName:  module,
		directory:   dir,
		initVCS:     true,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.quitting {
		return m, tea.Quit
	}

	// Handle final step (project creation result)
	if m.step == 5 {
		switch msg.(type) {
		case tea.KeyMsg:
			// Any key to quit after seeing result
			m.quitting = true
			return m, tea.Quit
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			if m.creatingProject {
				// Игнорируем Enter во время создания проекта
				return m, nil
			}
			if m.step < 4 {
				m.step++
				// Передаем фокус следующему полю
				if m.step == 1 {
					m.moduleName.Focus()
				} else if m.step == 2 {
					m.directory.Focus()
				}
				return m, nil
			} else {
				// Создаем проект
				m.step = 5 // Переходим к шагу создания проекта
				m.creatingProject = true
				return m, m.createProject()
			}

		case "backspace":
			if m.step > 0 {
				m.step--
				// Передаем фокус предыдущему полю
				if m.step == 0 {
					m.projectName.Focus()
				} else if m.step == 1 {
					m.moduleName.Focus()
				} else if m.step == 2 {
					m.directory.Focus()
				}
				return m, nil
			}

		case "up", "down":
			if m.step == 3 {
				if m.projectType == "" {
					m.projectType = "cli"
				}
				if msg.String() == "up" {
					switch m.projectType {
					case "cli":
						m.projectType = "library"
					case "web":
						m.projectType = "cli"
					case "library":
						m.projectType = "web"
					}
				} else if msg.String() == "down" {
					switch m.projectType {
					case "cli":
						m.projectType = "web"
					case "web":
						m.projectType = "library"
					case "library":
						m.projectType = "cli"
					}
				}
				return m, nil
			}

		case "left", "right", "h", "l":
			if m.step == 4 {
				m.initVCS = !m.initVCS
				return m, nil
			}

		case "y", "Y":
			if m.step == 4 {
				m.initVCS = true
				return m, nil
			}

		case "n", "N":
			if m.step == 4 {
				m.initVCS = false
				return m, nil
			}
		}
	}

	// Обновляем соответствующий input
	var cmd tea.Cmd
	switch m.step {
	case 0:
		m.projectName, cmd = m.projectName.Update(msg)
	case 1:
		m.moduleName, cmd = m.moduleName.Update(msg)
	case 2:
		m.directory, cmd = m.directory.Update(msg)
	}

	return m, cmd
}

func (m *Model) createProject() tea.Cmd {
	return func() tea.Msg {
		projectName := strings.TrimSpace(m.projectName.Value())
		if projectName == "" {
			projectName = m.projectName.Placeholder
		}

		moduleName := strings.TrimSpace(m.moduleName.Value())
		if moduleName == "" {
			moduleName = projectName
		}

		directory := strings.TrimSpace(m.directory.Value())
		if directory == "" {
			directory = projectName
		}

		m.config = generator.Config{
			ProjectName: projectName,
			ModuleName:  moduleName,
			Directory:   directory,
			ProjectType: m.projectType,
			InitVCS:     m.initVCS,
		}

		// Debug: print config values
		fmt.Printf("TUI DEBUG: ProjectName=%s, ModuleName=%s, Directory=%s\n", projectName, moduleName, directory)

		err := generator.InitProject(m.config)
		if err != nil {
			m.error = err
		} else {
			m.success = true
		}
		return nil
	}
}

func (m Model) View() string {
	if m.quitting {
		return "" // Пустая строка, чтобы не мешать выводу success message
	}

	if m.step == 5 {
		if m.error != nil {
			return ErrorStyle.Render("❌ Error creating project: " + m.error.Error())
		}
		return SuccessStyle.Render("✅ Project created successfully!")
	}

	var b strings.Builder

	b.WriteString(TitleStyle.Render("🚀 Go Project Initializer"))
	b.WriteString("\n\n")

	switch m.step {
	case 0:
		b.WriteString(QuestionStyle.Render("What's your project name?"))
		b.WriteString("\n\n")
		b.WriteString(m.projectName.View())
		b.WriteString(HelpStyle.Render("\n\nPress Enter to continue, Ctrl+C to quit"))

	case 1:
		b.WriteString(QuestionStyle.Render("What's your Go module name?"))
		b.WriteString("\n\n")
		b.WriteString(m.moduleName.View())
		b.WriteString(HelpStyle.Render("\n\nPress Enter to continue, Backspace to go back"))

	case 2:
		b.WriteString(QuestionStyle.Render("Where should we create the project?"))
		b.WriteString("\n\n")
		b.WriteString(m.directory.View())
		b.WriteString(HelpStyle.Render("\n\nPress Enter to continue, Backspace to go back"))

	case 3:
		b.WriteString(QuestionStyle.Render("What type of project do you want to create?"))
		b.WriteString("\n\n")
		if m.projectType == "" {
			m.projectType = "cli"
		}
		if m.projectType == "cli" {
			b.WriteString(SelectedStyle.Render("• CLI Application") + "\n")
		} else {
			b.WriteString(UnselectedStyle.Render("  CLI Application") + "\n")
		}
		if m.projectType == "web" {
			b.WriteString(SelectedStyle.Render("• Web Application") + "\n")
		} else {
			b.WriteString(UnselectedStyle.Render("  Web Application") + "\n")
		}
		if m.projectType == "library" {
			b.WriteString(SelectedStyle.Render("• Library") + "\n")
		} else {
			b.WriteString(UnselectedStyle.Render("  Library") + "\n")
		}
		b.WriteString(HelpStyle.Render("\n\nUse ↑/↓ to select, Enter to continue, Backspace to go back"))

	case 4:
		b.WriteString(QuestionStyle.Render("Initialize Git repository?"))
		b.WriteString("\n\n")
		if m.initVCS {
			b.WriteString(SelectedStyle.Render("✓ Yes") + "   " + UnselectedStyle.Render("No"))
		} else {
			b.WriteString(UnselectedStyle.Render("Yes") + "   " + SelectedStyle.Render("✓ No"))
		}
		b.WriteString(HelpStyle.Render("\n\nUse ←/→ or Y/N to toggle, Enter to create project"))
	}

	return b.String()
}
