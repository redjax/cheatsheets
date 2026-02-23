package cheatsheetservice

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type viewerModel struct {
	viewport        viewport.Model
	ready           bool
	content         string
	originalContent string // Content without search highlights
	rawContent      string // Original unrendered content for searching
	searchMode      bool
	searchInput     textinput.Model
	searchQuery     string
	searchLines     []int // Line numbers containing matches
	currentMatch    int   // Current match index
}

func (m viewerModel) Init() tea.Cmd {
	return nil
}

// performSearch searches for the query in the rendered content and updates search state
func (m *viewerModel) performSearch() {
	m.searchLines = []int{}
	m.currentMatch = 0

	if m.searchQuery == "" {
		// Restore original content without highlights
		m.content = m.originalContent
		m.viewport.SetContent(m.content)
		return
	}

	// Search in the RENDERED content, not raw markdown
	// This ensures line numbers match what's actually displayed
	lines := strings.Split(m.originalContent, "\n")
	query := strings.ToLower(m.searchQuery)

	// Count each individual occurrence of the search term
	for i, line := range lines {
		lowerLine := strings.ToLower(line)
		// Count how many times the query appears on this line
		lineContent := lowerLine
		pos := 0
		for {
			idx := strings.Index(lineContent[pos:], query)
			if idx == -1 {
				break
			}
			// Add this line number for each occurrence
			m.searchLines = append(m.searchLines, i)
			pos += idx + len(query) // Move past this match
		}
	}

	// searchLines is already in order since we process lines sequentially
	// and add matches in order within each line

	// Apply highlighting to the content
	m.highlightMatches()
}

// highlightMatches highlights all occurrences of the search query in the content
func (m *viewerModel) highlightMatches() {
	if m.searchQuery == "" {
		m.content = m.originalContent
		m.viewport.SetContent(m.content)
		return
	}

	// Create highlight styles
	// Yellow background for regular matches
	highlightStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("226")).
		Foreground(lipgloss.Color("0")).
		Bold(true)

	// Bright cyan background for current match
	currentHighlightStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("51")).
		Foreground(lipgloss.Color("0")).
		Bold(true)

	// Determine which line and which occurrence on that line is the current match
	var currentMatchLine int = -1
	var currentMatchOccurrence int = 0
	if len(m.searchLines) > 0 {
		currentMatchLine = m.searchLines[m.currentMatch]
		// Count how many times this line appears before current match
		for i := 0; i < m.currentMatch; i++ {
			if m.searchLines[i] == currentMatchLine {
				currentMatchOccurrence++
			}
		}
	}

	// Split content into lines and highlight matches
	lines := strings.Split(m.originalContent, "\n")
	query := strings.ToLower(m.searchQuery)

	for i, line := range lines {
		// Find all occurrences of the search query (case-insensitive)
		lowerLine := strings.ToLower(line)
		var highlighted strings.Builder
		lastPos := 0
		matchCount := 0

		for {
			pos := strings.Index(lowerLine[lastPos:], query)
			if pos == -1 {
				// No more matches, append the rest of the line
				highlighted.WriteString(line[lastPos:])
				break
			}

			// Adjust position to absolute position in line
			pos += lastPos

			// Append text before match
			highlighted.WriteString(line[lastPos:pos])

			// Determine which style to use
			var styleToUse lipgloss.Style
			if i == currentMatchLine && matchCount == currentMatchOccurrence {
				// This specific occurrence gets the current highlight
				styleToUse = currentHighlightStyle
			} else {
				// Other matches get regular highlight
				styleToUse = highlightStyle
			}

			// Append highlighted match
			matchText := line[pos : pos+len(m.searchQuery)]
			highlighted.WriteString(styleToUse.Render(matchText))

			// Move past this match
			lastPos = pos + len(m.searchQuery)
			matchCount++
		}

		lines[i] = highlighted.String()
	}

	// Update the content with highlights
	m.content = strings.Join(lines, "\n")
	m.viewport.SetContent(m.content)
}

// jumpToMatch scrolls the viewport to the current match
func (m *viewerModel) jumpToMatch() {
	if len(m.searchLines) == 0 {
		return
	}

	// Get the line number of the current match
	lineNum := m.searchLines[m.currentMatch]

	// Get viewport dimensions
	viewportHeight := m.viewport.Height
	currentOffset := m.viewport.YOffset

	// Calculate the visible range (with some margin for safety)
	visibleStart := currentOffset + 1                // Add 1 line buffer at top
	visibleEnd := currentOffset + viewportHeight - 2 // Remove 2 lines buffer at bottom

	// Check if the match is comfortably visible (not at edges)
	if lineNum >= visibleStart && lineNum < visibleEnd {
		// Already visible with good context, no need to scroll
		return
	}

	// Position the match line with context
	// Try to center it, or at least show it with 5 lines of context above
	contextLines := 5
	targetOffset := lineNum - contextLines

	// Ensure we don't scroll past the beginning
	if targetOffset < 0 {
		targetOffset = 0
	}

	// Get total number of lines
	totalLines := len(strings.Split(m.content, "\n"))

	// Ensure we don't scroll past the end
	if targetOffset+viewportHeight > totalLines {
		targetOffset = totalLines - viewportHeight
		if targetOffset < 0 {
			targetOffset = 0
		}
	}

	// Set the viewport to show the match with context
	m.viewport.SetYOffset(targetOffset)
}

func (m viewerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle search mode input
		if m.searchMode {
			switch msg.String() {
			case "esc":
				// Exit search mode
				m.searchMode = false
				m.searchInput.Reset()
				return m, nil
			case "enter":
				// Execute search
				m.searchQuery = m.searchInput.Value()
				m.performSearch()
				m.searchMode = false
				if len(m.searchLines) > 0 {
					// Force scroll to first match on new search
					lineNum := m.searchLines[m.currentMatch]
					contextLines := 3
					targetOffset := lineNum - contextLines
					if targetOffset < 0 {
						targetOffset = 0
					}
					m.viewport.SetYOffset(targetOffset)
				}
				return m, nil
			default:
				// Update search input
				m.searchInput, cmd = m.searchInput.Update(msg)
				return m, cmd
			}
		}

		// Normal mode key handling
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			// Clear search and highlights
			if m.searchQuery != "" {
				m.searchQuery = ""
				m.searchLines = []int{}
				m.currentMatch = 0
				m.content = m.originalContent
				m.viewport.SetContent(m.content)
				return m, nil
			}
		case "/":
			// Enter search mode
			m.searchMode = true
			m.searchInput.Focus()
			return m, textinput.Blink
		case "n":
			// Next match
			if len(m.searchLines) > 0 {
				m.currentMatch = (m.currentMatch + 1) % len(m.searchLines)
				m.highlightMatches() // Re-highlight to update current match
				m.jumpToMatch()
			}
			return m, nil
		case "N":
			// Previous match
			if len(m.searchLines) > 0 {
				m.currentMatch = (m.currentMatch - 1 + len(m.searchLines)) % len(m.searchLines)
				m.highlightMatches() // Re-highlight to update current match
				m.jumpToMatch()
			}
			return m, nil
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

		// Update search input width
		m.searchInput.Width = msg.Width - 4
	}

	// Update viewport in normal mode
	if !m.searchMode {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

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

	var footer string
	if m.searchMode {
		// Show search input in footer
		footer = footerStyle.Render(fmt.Sprintf("Search: %s", m.searchInput.View()))
	} else if len(m.searchLines) > 0 {
		// Show search results info
		footer = footerStyle.Render(fmt.Sprintf(
			"Match %d/%d • n: next • N: prev • esc: clear • /: search • q: quit | %3.f%%",
			m.currentMatch+1,
			len(m.searchLines),
			m.viewport.ScrollPercent()*100,
		))
	} else if m.searchQuery != "" {
		// No matches found
		footer = footerStyle.Render(fmt.Sprintf(
			"No matches for '%s' • esc: clear • /: search • q: quit | %3.f%%",
			m.searchQuery,
			m.viewport.ScrollPercent()*100,
		))
	} else {
		// Normal mode
		footer = footerStyle.Render(fmt.Sprintf(
			"↑/↓: scroll • k/j: scroll • PgUp/PgDn: page • g/G: top/bottom • /: search • q: quit | %3.f%%",
			m.viewport.ScrollPercent()*100,
		))
	}

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

	// Initialize search input
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.CharLimit = 100

	// Display in the TUI viewer
	m := viewerModel{
		content:         rendered,
		originalContent: rendered,
		rawContent:      rawMarkdown,
		searchInput:     ti,
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
