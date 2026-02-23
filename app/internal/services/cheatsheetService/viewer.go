package cheatsheetservice

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
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

// ShowInViewer displays content in a custom TUI viewer
func ShowInViewer(rawMarkdown string) error {
	// Apply syntax highlighting to the markdown
	highlighted := HighlightMarkdown(rawMarkdown)

	m := viewerModel{
		content: highlighted,
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
