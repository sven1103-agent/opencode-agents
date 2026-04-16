// Package styles provides centralized styling for the CLI using lipgloss.
// It defines consistent visual presentation with accessible color + icon schemes.
package styles

import (
	"os"
	"strings"

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

// Sensible color palette - muted, professional tones
var (
	errorColor   = lipgloss.Color("#E06C75") // Muted red
	warningColor = lipgloss.Color("#E5C07B") // Muted yellow
	successColor = lipgloss.Color("#98C379") // Muted green
	infoColor    = lipgloss.Color("#61AFEF") // Muted blue
	promptColor  = lipgloss.Color("#C678DD") // Muted purple
	keyColor     = lipgloss.Color("#D19A66") // Muted orange for keys
	valueColor   = lipgloss.Color("#ABB2BF") // Light gray for values
	mutedColor   = lipgloss.Color("#5C6370") // Gray for secondary text
)

// Style definitions

// ErrorStyle for error messages
var ErrorStyle = lipgloss.Style{}.
	Foreground(errorColor).
	Bold(true)

// ErrorIcon returns the error icon
func ErrorIcon() string {
	return "✗"
}

// WarningStyle for warning messages
var WarningStyle = lipgloss.Style{}.
	Foreground(warningColor).
	Bold(true)

// WarningIcon returns the warning icon
func WarningIcon() string {
	return "⚠"
}

// SuccessStyle for success messages
var SuccessStyle = lipgloss.Style{}.
	Foreground(successColor).
	Bold(true)

// SuccessIcon returns the success icon
func SuccessIcon() string {
	return "✓"
}

// InfoStyle for info messages
var InfoStyle = lipgloss.Style{}.
	Foreground(infoColor).
	Bold(true)

// InfoIcon returns the info icon
func InfoIcon() string {
	return "ℹ"
}

// PromptStyle for interactive prompts
var PromptStyle = lipgloss.Style{}.
	Foreground(promptColor)

// KeyStyle for key labels (subtle orange)
var KeyStyle = lipgloss.Style{}.
	Foreground(keyColor)

// ValueStyle for values (light gray)
var ValueStyle = lipgloss.Style{}.
	Foreground(valueColor)

// MutedStyle for secondary/muted text
var MutedStyle = lipgloss.Style{}.
	Foreground(mutedColor)

// HelpHeaderStyle for help text headers
var HelpHeaderStyle = lipgloss.Style{}.
	Foreground(infoColor).
	Underline(true)

// HelpCommandStyle for command names in help
var HelpCommandStyle = lipgloss.Style{}.
	Foreground(promptColor)

// Output functions that respect color mode

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

// Prompt outputs a prompt message (no icon, just color)
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

// SectionHeader outputs a prominent section header with separator
func SectionHeader(title string) string {
	sep := "─"
	titleWithColon := title + ":"
	if ShouldRenderColor() {
		// Make header more prominent with accent color and bold
		header := lipgloss.Style{}.Foreground(promptColor).Bold(true).Render(titleWithColon)
		separator := MutedStyle.Render(strings.Repeat(sep, 36))
		return header + "\n" + separator
	}
	return titleWithColon + "\n" + strings.Repeat(sep, 36)
}

// Highlight outputs highlighted/important text
func Highlight(msg string) string {
	if ShouldRenderColor() {
		return lipgloss.Style{}.Foreground(promptColor).Bold(true).Render(msg)
	}
	return msg
}

// SubHeader outputs a subsection header (subtle)
func SubHeader(title string) string {
	if ShouldRenderColor() {
		return KeyStyle.Render("▸ " + title)
	}
	return "▸ " + title
}

// KeyValue outputs a key-value pair
func KeyValue(key, value string) string {
	if ShouldRenderColor() {
		return KeyStyle.Render(key+":") + " " + ValueStyle.Render(value)
	}
	return key + ": " + value
}

// KeyValueMuted outputs a key-value pair with muted value
func KeyValueMuted(key, value string) string {
	if ShouldRenderColor() {
		return KeyStyle.Render(key+":") + " " + MutedStyle.Render(value)
	}
	return key + ": " + value
}

// Muted outputs muted/secondary text
func Muted(msg string) string {
	if ShouldRenderColor() {
		return MutedStyle.Render(msg)
	}
	return msg
}

// TableStyle renders a modern table with box-drawing characters
func TableStyle(headers []string, rows [][]string) string {
	if len(headers) == 0 {
		return ""
	}

	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Add padding
	for i := range widths {
		widths[i] += 2 // 1 space on each side
	}

	var lines []string

	// Build top border
	topBorder := "┌"
	for i, w := range widths {
		if i > 0 {
			topBorder += "┬"
		}
		topBorder += strings.Repeat("─", w)
	}
	topBorder += "┐"
	if ShouldRenderColor() {
		lines = append(lines, MutedStyle.Render(topBorder))
	} else {
		lines = append(lines, topBorder)
	}

	// Build header row
	headerRow := "│"
	for i, h := range headers {
		padding := widths[i] - len(h)
		leftPad := padding / 2
		rightPad := padding - leftPad
		if ShouldRenderColor() {
			headerRow += lipgloss.Style{}.Foreground(promptColor).Bold(true).Render(strings.Repeat(" ", leftPad)+h+strings.Repeat(" ", rightPad)) + "│"
		} else {
			headerRow += strings.Repeat(" ", leftPad) + h + strings.Repeat(" ", rightPad) + "│"
		}
	}
	lines = append(lines, headerRow)

	// Build separator
	sep := "├"
	for i, w := range widths {
		if i > 0 {
			sep += "┼"
		}
		sep += strings.Repeat("─", w)
	}
	sep += "┤"
	if ShouldRenderColor() {
		lines = append(lines, MutedStyle.Render(sep))
	} else {
		lines = append(lines, sep)
	}

	// Build data rows with optional zebra striping
	for rowIdx, row := range rows {
		rowStr := "│"
		for i, cell := range row {
			if i >= len(widths) {
				continue
			}
			padding := widths[i] - len(cell) - 1 // -1 for left padding
			content := " " + cell + strings.Repeat(" ", padding)

			if ShouldRenderColor() && rowIdx%2 == 1 {
				// Zebra striping - every other row gets slightly different background
				rowStr += lipgloss.Style{}.Foreground(valueColor).Render(content) + "│"
			} else {
				rowStr += content + "│"
			}
		}
		lines = append(lines, rowStr)
	}

	// Build bottom border
	bottomBorder := "└"
	for i, w := range widths {
		if i > 0 {
			bottomBorder += "┴"
		}
		bottomBorder += strings.Repeat("─", w)
	}
	bottomBorder += "┘"
	if ShouldRenderColor() {
		lines = append(lines, MutedStyle.Render(bottomBorder))
	} else {
		lines = append(lines, bottomBorder)
	}

	return strings.Join(lines, "\n")
}
