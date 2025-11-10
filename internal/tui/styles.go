package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Base styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("63")).
			PaddingBottom(1)

	QuestionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39"))

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	SelectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("42"))

	UnselectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	// Success screen styles
	SuccessIcon = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("42"))

	SuccessText = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39"))

	Label = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	Value = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("255"))

	Section = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("213")).
		PaddingBottom(1)

	Code = lipgloss.NewStyle().
		Foreground(lipgloss.Color("156")).
		Background(lipgloss.Color("236")).
		Padding(0, 1)

	SuccessStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("42")).
		Padding(1, 0)

	ErrorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("196")).
		Padding(1, 0)

	Tip = lipgloss.NewStyle().
		Foreground(lipgloss.Color("255"))
)

// DefaultStyle возвращает стили для non-interactive режима
type AppStyle struct {
	SuccessIcon *lipgloss.Style
	SuccessText *lipgloss.Style
	Label       *lipgloss.Style
	Value       *lipgloss.Style
	Section     *lipgloss.Style
	Code        *lipgloss.Style
	Tip         *lipgloss.Style
}

func DefaultStyle() AppStyle {
	return AppStyle{
		SuccessIcon: &SuccessIcon,
		SuccessText: &SuccessText,
		Label:       &Label,
		Value:       &Value,
		Section:     &Section,
		Code:        &Code,
		Tip:         &Tip,
	}
}
