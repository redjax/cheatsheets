package cheatsheetservice

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type viewerModel struct {
	viewport viewport.Model
	ready    bool
	content  string
}

func (m viewerModel) Init() tea.Cmd {
	return nil
}

func (m viewerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := 2
		footerHeight := 2
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m viewerModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		PaddingLeft(2)

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		PaddingLeft(2)

	header := headerStyle.Render("Cheatsheet Viewer")
	footer := footerStyle.Render(fmt.Sprintf(
		"↑/↓: scroll • k/j: scroll • PgUp/PgDn: page • g/G: top/bottom • q: quit | %3.f%%",
		m.viewport.ScrollPercent()*100,
	))

	return fmt.Sprintf("%s\n%s\n%s", header, m.viewport.View(), footer)
}

// ShowInViewer displays content using glamour rendering (glow's underlying library)
func ShowInViewer(rawMarkdown string) error {
	// Get terminal width for proper rendering
	width, _, err := getTerminalSize()
	if err != nil {
		width = 80 // fallback width
	}

	// Create glamour renderer with dark style
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width-4), // Leave some margin
	)
	if err != nil {
		return fmt.Errorf("error creating markdown renderer: %w", err)
	}

	// Render the markdown
	rendered, err := r.Render(rawMarkdown)
	if err != nil {
		return fmt.Errorf("error rendering markdown: %w", err)
	}

	// Display in the TUI viewer
	m := viewerModel{
		content: rendered,
	}

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running viewer: %w", err)
	}

	return nil
}

// getTerminalSize returns the width and height of the terminal
func getTerminalSize() (int, int, error) {
	width, height, err := getSize(os.Stdout.Fd())
	if err != nil {
		return 80, 24, nil // fallback
	}
	return width, height, nil
}
