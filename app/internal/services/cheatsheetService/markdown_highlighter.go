package cheatsheetservice

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Heading styles
	h1Style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	h2Style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))
	h3Style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("11"))
	h4Style = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))

	// Code styles
	codeBlockStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	inlineCodeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))

	// List styles
	listItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))

	// Link styles
	linkStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))

	// Bold/Italic
	emphasisStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))

	// Regular patterns
	inlineCodeRegex = regexp.MustCompile("`([^`]+)`")
	boldRegex       = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	italicRegex     = regexp.MustCompile(`\*([^*]+)\*`)
	linkRegex       = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
)

// HighlightMarkdown applies syntax highlighting to markdown text
func HighlightMarkdown(markdown string) string {
	lines := strings.Split(markdown, "\n")
	var highlighted []string
	inCodeBlock := false

	for _, line := range lines {
		// Check for code block delimiters
		if strings.HasPrefix(line, "```") {
			if inCodeBlock {
				// End of code block
				highlighted = append(highlighted, codeBlockStyle.Render(line))
				inCodeBlock = false
			} else {
				// Start of code block
				inCodeBlock = true
				highlighted = append(highlighted, codeBlockStyle.Render(line))
			}
			continue
		}

		// If in code block, style everything as code
		if inCodeBlock {
			highlighted = append(highlighted, codeBlockStyle.Render(line))
			continue
		}

		// Style headings
		if strings.HasPrefix(line, "# ") {
			highlighted = append(highlighted, h1Style.Render(line))
			continue
		}
		if strings.HasPrefix(line, "## ") {
			highlighted = append(highlighted, h2Style.Render(line))
			continue
		}
		if strings.HasPrefix(line, "### ") {
			highlighted = append(highlighted, h3Style.Render(line))
			continue
		}
		if strings.HasPrefix(line, "#### ") {
			highlighted = append(highlighted, h4Style.Render(line))
			continue
		}

		// Style list items
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") || strings.HasPrefix(trimmed, "+ ") {
			// Get indentation
			indent := line[:len(line)-len(trimmed)]
			highlighted = append(highlighted, indent+listItemStyle.Render(trimmed))
			continue
		}

		// For numbered lists
		if len(trimmed) > 0 && trimmed[0] >= '0' && trimmed[0] <= '9' {
			if idx := strings.Index(trimmed, ". "); idx > 0 {
				indent := line[:len(line)-len(trimmed)]
				highlighted = append(highlighted, indent+listItemStyle.Render(trimmed))
				continue
			}
		}

		// Apply inline highlighting
		line = highlightInline(line)
		highlighted = append(highlighted, line)
	}

	return strings.Join(highlighted, "\n")
}

// highlightInline applies inline markdown highlighting (code, bold, links, etc.)
func highlightInline(line string) string {
	// Highlight inline code (do this first to avoid conflicts)
	line = inlineCodeRegex.ReplaceAllStringFunc(line, func(match string) string {
		return inlineCodeStyle.Render(match)
	})

	// Highlight links
	line = linkRegex.ReplaceAllStringFunc(line, func(match string) string {
		return linkStyle.Render(match)
	})

	// Highlight bold
	line = boldRegex.ReplaceAllStringFunc(line, func(match string) string {
		return emphasisStyle.Bold(true).Render(match)
	})

	// Highlight italic
	line = italicRegex.ReplaceAllStringFunc(line, func(match string) string {
		return emphasisStyle.Italic(true).Render(match)
	})

	return line
}
