//go:build !windows
// +build !windows

package cheatsheetservice

import (
	"golang.org/x/term"
)

// getSize returns the width and height of the terminal
func getSize(fd uintptr) (int, int, error) {
	width, height, err := term.GetSize(int(fd))
	return width, height, err
}
