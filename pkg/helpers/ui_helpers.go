package helpers

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pterm/pterm"
)

var (
	promptStyle = lipgloss.NewStyle().Bold(true)
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type textInputModel struct {
	prompt  string
	input   textinput.Model
	done    bool
	aborted bool
}

func newTextInputModel(prompt string, defaultValue string) textInputModel {
	ti := textinput.New()
	ti.SetValue(defaultValue)
	ti.CursorEnd()
	ti.Focus()

	return textInputModel{prompt: prompt, input: ti}
}

func (m textInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m textInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.done = true
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.aborted = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m textInputModel) View() string {
	if m.done || m.aborted {
		return ""
	}
	return fmt.Sprintf("%s\n%s\n%s\n", promptStyle.Render(m.prompt), m.input.View(), helpStyle.Render("(enter to confirm, esc to cancel)"))
}

// GetInputAsString prompts the user for a line of text using a bubbletea
// text input, pre-populated with defaultValue.
func GetInputAsString(prompt string, defaultValue string) string {
	m, err := tea.NewProgram(newTextInputModel(prompt, defaultValue)).Run()
	if err != nil {
		return defaultValue
	}

	result := m.(textInputModel)
	if result.aborted {
		return defaultValue
	}
	return result.input.Value()
}

type selectModel struct {
	prompt  string
	options []string
	cursor  int
	done    bool
	aborted bool
}

func newSelectModel(prompt string, options []string) selectModel {
	return selectModel{prompt: prompt, options: options}
}

func (m selectModel) Init() tea.Cmd {
	return nil
}

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "enter":
			m.done = true
			return m, tea.Quit
		case "ctrl+c", "esc":
			m.aborted = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m selectModel) View() string {
	if m.done || m.aborted {
		return ""
	}

	var b strings.Builder
	b.WriteString(promptStyle.Render(m.prompt))
	b.WriteString("\n")
	for i, option := range m.options {
		if i == m.cursor {
			b.WriteString(cursorStyle.Render("> " + option))
		} else {
			b.WriteString("  ")
			b.WriteString(option)
		}
		b.WriteString("\n")
	}
	b.WriteString(helpStyle.Render("(use arrows to navigate, enter to select, esc to cancel)"))
	return b.String()
}

// SelectOptionAsString prompts the user to choose one of options using a
// bubbletea interactive list.
func SelectOptionAsString(prompt string, options []string) string {
	if len(options) == 0 {
		return ""
	}

	m, err := tea.NewProgram(newSelectModel(prompt, options)).Run()
	if err != nil {
		return ""
	}

	result := m.(selectModel)
	if result.aborted {
		return ""
	}
	return result.options[result.cursor]
}

func RenderTable(data pterm.TableData) error {
	return pterm.DefaultTable.WithHasHeader().WithData(data).Render()
}
