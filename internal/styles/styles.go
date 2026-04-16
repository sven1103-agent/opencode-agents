// Package styles provides centralized styling for the CLI using lipgloss.
// It defines consistent visual presentation with accessible color + icon schemes.
package styles

import (
	"os"

	"github.com/charmbracelet/lipgloss"
)

// ColorMode represents the color output mode
type ColorMode string

const (
	ColorModeAuto   ColorMode = "auto"
	ColorModeAlways ColorMode = "always"
	ColorModeNever  ColorMode = "never"
)

// Global color mode setting
var colorMode ColorMode = ColorModeAlways

// SetColorMode sets the global color mode
func SetColorMode(mode ColorMode) {
	colorMode = mode
}

// GetColorMode returns the current color mode
func GetColorMode() ColorMode {
	return colorMode
}

// ShouldRenderColor determines whether to render colors based on the current mode
func ShouldRenderColor() bool {
	switch colorMode {
	case ColorModeAlways:
		return true
	case ColorModeNever:
		return false
	case ColorModeAuto:
		return isTerminal(os.Stdout)
	}
	return true
}

// isTerminal checks if the output is a terminal
func isTerminal(f *os.File) bool {
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

// Define color palette
var (
	errorColor   = lipgloss.Color("#FF5555")
	warningColor = lipgloss.Color("#F1FA8C")
	successColor = lipgloss.Color("#50FA7B")
	infoColor    = lipgloss.Color("#8BE9FD")
	promptColor  = lipgloss.Color("#BD93F9")
)

// Define style definitions

// ErrorStyle for error messages (red + ✗)
var ErrorStyle = lipgloss.Style{}.
	Foreground(errorColor).
	Bold(true)

// ErrorIcon returns the error icon
func ErrorIcon() string {
	return "✗"
}

// WarningStyle for warning messages (yellow + ⚠)
var WarningStyle = lipgloss.Style{}.
	Foreground(warningColor).
	Bold(true)

// WarningIcon returns the warning icon
func WarningIcon() string {
	return "⚠"
}

// SuccessStyle for success messages (green + ✓)
var SuccessStyle = lipgloss.Style{}.
	Foreground(successColor).
	Bold(true)

// SuccessIcon returns the success icon
func SuccessIcon() string {
	return "✓"
}

// InfoStyle for info messages (blue + ℹ)
var InfoStyle = lipgloss.Style{}.
	Foreground(infoColor).
	Bold(true)

// InfoIcon returns the info icon
func InfoIcon() string {
	return "ℹ"
}

// PromptStyle for interactive prompts (purple)
var PromptStyle = lipgloss.Style{}.
	Foreground(promptColor).
	Bold(true)

// HelpHeaderStyle for help text headers
var HelpHeaderStyle = lipgloss.Style{}.
	Foreground(infoColor).
	Underline(true)

// HelpCommandStyle for command names in help
var HelpCommandStyle = lipgloss.Style{}.
	Foreground(promptColor)

// Output styles that respect color mode

// Error outputs an error message with icon
func Error(msg string) string {
	if ShouldRenderColor() {
		return ErrorStyle.Render(ErrorIcon() + " " + msg)
	}
	return ErrorIcon() + " " + msg
}

// Warning outputs a warning message with icon
func Warning(msg string) string {
	if ShouldRenderColor() {
		return WarningStyle.Render(WarningIcon() + " " + msg)
	}
	return WarningIcon() + " " + msg
}

// Success outputs a success message with icon
func Success(msg string) string {
	if ShouldRenderColor() {
		return SuccessStyle.Render(SuccessIcon() + " " + msg)
	}
	return SuccessIcon() + " " + msg
}

// Info outputs an info message with icon
func Info(msg string) string {
	if ShouldRenderColor() {
		return InfoStyle.Render(InfoIcon() + " " + msg)
	}
	return InfoIcon() + " " + msg
}

// Prompt outputs a prompt message
func Prompt(msg string) string {
	if ShouldRenderColor() {
		return PromptStyle.Render(msg)
	}
	return msg
}

// DryRun outputs a dry-run message
func DryRun(msg string) string {
	if ShouldRenderColor() {
		return InfoStyle.Render("dry-run: " + msg)
	}
	return "dry-run: " + msg
}

// Written outputs a written file message
func Written(path string) string {
	if ShouldRenderColor() {
		return SuccessStyle.Render("written: " + path)
	}
	return "written: " + path
}

// Done outputs a done message
func Done(msg string) string {
	if ShouldRenderColor() {
		return SuccessStyle.Render("done: " + msg)
	}
	return "done: " + msg
}

// Invalid outputs an invalid selection message
func Invalid(msg string) string {
	if ShouldRenderColor() {
		return ErrorStyle.Render("Invalid selection. " + msg)
	}
	return "Invalid selection. " + msg
}

// Required outputs a required flag message
func Required(flag string) string {
	if ShouldRenderColor() {
		return ErrorStyle.Render(flag + " is required")
	}
	return flag + " is required"
}

// HelpHeader outputs a help header
func HelpHeader(msg string) string {
	if ShouldRenderColor() {
		return HelpHeaderStyle.Render(msg)
	}
	return msg
}

// HelpCommand outputs a help command name
func HelpCommand(cmd string) string {
	if ShouldRenderColor() {
		return HelpCommandStyle.Render(cmd)
	}
	return cmd
}
