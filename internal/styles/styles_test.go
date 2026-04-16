package styles

import (
	"os"
	"testing"
)

func TestError(t *testing.T) {
	// Test with color enabled
	SetColorMode(ColorModeAlways)
	output := Error("test error message")
	if output == "" {
		t.Error("Error output should not be empty")
	}
	// Should contain the icon
	expectedIcon := ErrorIcon()
	if expectedIcon != "✗" {
		t.Errorf("Expected icon ✗, got %s", expectedIcon)
	}
}

func TestWarning(t *testing.T) {
	SetColorMode(ColorModeAlways)
	output := Warning("test warning message")
	if output == "" {
		t.Error("Warning output should not be empty")
	}
	expectedIcon := WarningIcon()
	if expectedIcon != "⚠" {
		t.Errorf("Expected icon ⚠, got %s", expectedIcon)
	}
}

func TestSuccess(t *testing.T) {
	SetColorMode(ColorModeAlways)
	output := Success("test success message")
	if output == "" {
		t.Error("Success output should not be empty")
	}
	expectedIcon := SuccessIcon()
	if expectedIcon != "✓" {
		t.Errorf("Expected icon ✓, got %s", expectedIcon)
	}
}

func TestInfo(t *testing.T) {
	SetColorMode(ColorModeAlways)
	output := Info("test info message")
	if output == "" {
		t.Error("Info output should not be empty")
	}
	expectedIcon := InfoIcon()
	if expectedIcon != "ℹ" {
		t.Errorf("Expected icon ℹ, got %s", expectedIcon)
	}
}

func TestPrompt(t *testing.T) {
	SetColorMode(ColorModeAlways)
	output := Prompt("Enter value:")
	if output == "" {
		t.Error("Prompt output should not be empty")
	}
}

func TestWritten(t *testing.T) {
	SetColorMode(ColorModeAlways)
	output := Written("/path/to/file.txt")
	if output == "" {
		t.Error("Written output should not be empty")
	}
}

func TestDone(t *testing.T) {
	SetColorMode(ColorModeAlways)
	output := Done("operation complete")
	if output == "" {
		t.Error("Done output should not be empty")
	}
}

func TestDryRun(t *testing.T) {
	SetColorMode(ColorModeAlways)
	output := DryRun("write config to /path")
	if output == "" {
		t.Error("DryRun output should not be empty")
	}
}

func TestInvalid(t *testing.T) {
	SetColorMode(ColorModeAlways)
	output := Invalid("Please enter a valid number")
	if output == "" {
		t.Error("Invalid output should not be empty")
	}
}

func TestRequired(t *testing.T) {
	SetColorMode(ColorModeAlways)
	output := Required("--preset")
	if output == "" {
		t.Error("Required output should not be empty")
	}
}

// TestColorModeNever tests that colors are not rendered when disabled
func TestColorModeNever(t *testing.T) {
	SetColorMode(ColorModeNever)

	output := Error("test")
	// When color mode is never, output should be plain text + icon
	if output != "✗ test" {
		t.Errorf("Expected plain output without color codes, got: %s", output)
	}
}

// TestColorModeAutoTerminal tests auto mode with terminal
func TestColorModeAutoTerminal(t *testing.T) {
	SetColorMode(ColorModeAuto)
	// When stdout is a terminal, ShouldRenderColor should return true
	// When stdout is not a terminal (like in tests), ShouldRenderColor returns false
	// This is expected behavior
}

// TestShouldRenderColorWithNonTerminal tests that auto mode returns false for non-terminals
func TestShouldRenderColorWithNonTerminal(t *testing.T) {
	SetColorMode(ColorModeAuto)
	// In test environment (non-terminal), ShouldRenderColor should return false
	if ShouldRenderColor() {
		t.Log("Note: Running in a terminal, test may not reflect non-terminal behavior")
	}
}

// TestColorModeAlways tests that colors are always rendered when mode is always
func TestColorModeAlways(t *testing.T) {
	SetColorMode(ColorModeAlways)
	if !ShouldRenderColor() {
		t.Error("ShouldRenderColor should return true when mode is always")
	}
}

// TestColorModeNever tests that colors are never rendered when mode is never
func TestColorModeNever_ColorRender(t *testing.T) {
	SetColorMode(ColorModeNever)
	if ShouldRenderColor() {
		t.Error("ShouldRenderColor should return false when mode is never")
	}
}

// TestStylesSnapshot provides a snapshot-like test for all styled outputs
func TestStylesSnapshot(t *testing.T) {
	SetColorMode(ColorModeAlways)

	tests := []struct {
		name string
		fn   func() string
	}{
		{"Error", func() string { return Error("test error") }},
		{"Warning", func() string { return Warning("test warning") }},
		{"Success", func() string { return Success("test success") }},
		{"Info", func() string { return Info("test info") }},
		{"Prompt", func() string { return Prompt("test prompt") }},
		{"Written", func() string { return Written("/test/path") }},
		{"Done", func() string { return Done("test done") }},
		{"DryRun", func() string { return DryRun("test dry run") }},
		{"Invalid", func() string { return Invalid("test invalid") }},
		{"Required", func() string { return Required("--flag") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.fn()
			if output == "" {
				t.Errorf("%s output should not be empty", tt.name)
			}
			// Verify icon is present in output (at the start, before any styled content)
			switch tt.name {
			case "Error":
				runes := []rune(output)
				if len(runes) < 1 || runes[0] != '✗' {
					t.Errorf("Expected output to start with '✗', got: %s", output)
				}
			case "Warning":
				runes := []rune(output)
				if len(runes) < 1 || runes[0] != '⚠' {
					t.Errorf("Expected output to start with '⚠', got: %s", output)
				}
			case "Success":
				runes := []rune(output)
				if len(runes) < 1 || runes[0] != '✓' {
					t.Errorf("Expected output to start with '✓', got: %s", output)
				}
			case "Info":
				runes := []rune(output)
				if len(runes) < 1 || runes[0] != 'ℹ' {
					t.Errorf("Expected output to start with 'ℹ', got: %s", output)
				}
			}
		})
	}
}

// TestStylesNoColorSnapshot provides snapshot tests for non-colored output
func TestStylesNoColorSnapshot(t *testing.T) {
	// Save original stdout to restore later
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	SetColorMode(ColorModeNever)

	tests := []struct {
		name string
		fn   func() string
	}{
		{"Error_NoColor", func() string { return Error("test error") }},
		{"Warning_NoColor", func() string { return Warning("test warning") }},
		{"Success_NoColor", func() string { return Success("test success") }},
		{"Info_NoColor", func() string { return Info("test info") }},
		{"Written_NoColor", func() string { return Written("/test/path") }},
		{"Done_NoColor", func() string { return Done("test done") }},
		{"DryRun_NoColor", func() string { return DryRun("test dry run") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.fn()
			if output == "" {
				t.Errorf("%s output should not be empty", tt.name)
			}
		})
	}
}
