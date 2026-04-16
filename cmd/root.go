package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/qbicsoftware/occo/internal/styles"
	"github.com/qbicsoftware/occo/internal/version"
	"github.com/spf13/cobra"
)

var versionFlag bool
var colorMode string

// styledHelpTemplate returns a cobra help template with color styling
func styledHelpTemplate() string {
	return `{{with (or .Long .Short)}}{{. | wrapWithStyle "description"}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
}

// styledUsageTemplate returns a styled usage template
func styledUsageTemplate() string {
	return `{{styleHeader "Usage:"}}{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

{{styleHeader "Aliases:"}}
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

{{styleHeader "Examples:"}}
{{styleExamples .Example}}{{end}}{{if .HasAvailableSubCommands}}

{{styleHeader "Available Commands:"}}
{{styleCommands .Commands}}{{end}}{{if .HasAvailableLocalFlags}}

{{styleHeader "Flags:"}}
{{.LocalFlags.FlagUsages | styleFlags | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

{{styleHeader "Global Flags:"}}
{{.InheritedFlags.FlagUsages | styleFlags | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

{{styleHeader "Additional help topics:"}}{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}

// wrapWithStyle wraps text with the given style
func wrapWithStyle(style, text string) string {
	switch style {
	case "info":
		return styles.Info(text)
	case "muted":
		return styles.Muted(text)
	case "description":
		// Split into title (first line) and description (rest)
		lines := strings.Split(text, "\n")
		if len(lines) == 0 {
			return text
		}

		// First line is the title - make it bold
		title := lipgloss.Style{}.Foreground(lipgloss.Color("#ABB2BF")).Bold(true).Render(lines[0])

		// Rest is description in medium gray
		var desc string
		if len(lines) > 1 {
			desc = styles.ValueStyle.Render(strings.Join(lines[1:], "\n"))
		}

		if desc != "" {
			return title + "\n" + desc
		}
		return title
	default:
		return text
	}
}

// styleHeaderFunc returns a function that styles section headers
func styleHeaderFunc() func(string) string {
	return func(s string) string {
		return styles.SectionHeader(s)
	}
}

// styleCommandFunc returns a function that styles command names
func styleCommandFunc() func(string) string {
	return func(s string) string {
		return styles.Highlight(s)
	}
}

// styleCommandsFunc returns a function that styles and aligns the commands list
func styleCommandsFunc() func([]*cobra.Command) string {
	return func(commands []*cobra.Command) string {
		var lines []string

		// First pass: find max command name width
		maxWidth := 0
		for _, cmd := range commands {
			if cmd.IsAvailableCommand() || cmd.Name() == "help" {
				if len(cmd.Name()) > maxWidth {
					maxWidth = len(cmd.Name())
				}
			}
		}

		// Second pass: format each command with consistent alignment
		for _, cmd := range commands {
			if cmd.IsAvailableCommand() || cmd.Name() == "help" {
				padding := strings.Repeat(" ", maxWidth-len(cmd.Name())+2) // +2 for minimum spacing
				// Use plain text for description (not muted) for better readability
				lines = append(lines, "  "+styles.Highlight(cmd.Name())+padding+cmd.Short)
			}
		}

		return strings.Join(lines, "\n")
	}
}

// styleExamplesFunc styles example text
func styleExamplesFunc() func(string) string {
	return func(s string) string {
		// Indent and add subtle styling to examples
		lines := strings.Split(s, "\n")
		for i, line := range lines {
			if line != "" {
				// Just indent, keep examples in plain/white text for readability
				lines[i] = "  " + line
			}
		}
		return strings.Join(lines, "\n")
	}
}

// styleFlagsFunc styles flag usage text with consistent column alignment
func styleFlagsFunc() func(string) string {
	return func(s string) string {
		lines := strings.Split(s, "\n")

		// First pass: find the maximum width of "flag+type" column
		maxFlagWidth := 0
		for _, line := range lines {
			if line == "" {
				continue
			}
			trimmed := strings.TrimLeft(line, " \t")
			if strings.HasPrefix(trimmed, "-") {
				// Find where description starts (two or more spaces)
				for j := 0; j < len(trimmed)-1; j++ {
					if trimmed[j] == ' ' && trimmed[j+1] == ' ' {
						if j > maxFlagWidth {
							maxFlagWidth = j
						}
						break
					}
				}
			}
		}

		// Second pass: format each line with consistent alignment
		for i, line := range lines {
			if line == "" {
				continue
			}

			trimmed := strings.TrimLeft(line, " \t")

			// Check if this line starts with a flag
			if strings.HasPrefix(trimmed, "-") {
				// Find where description starts
				descStart := -1
				for j := 0; j < len(trimmed)-1; j++ {
					if trimmed[j] == ' ' && trimmed[j+1] == ' ' {
						descStart = j
						break
					}
				}

				if descStart > 0 {
					flagAndType := trimmed[:descStart]
					desc := strings.TrimLeft(trimmed[descStart:], " ")

					// Split flag name from type
					flagEnd := strings.Index(flagAndType, " ")
					if flagEnd == -1 {
						flagEnd = len(flagAndType)
					}

					flagName := flagAndType[:flagEnd]
					typePart := ""
					if flagEnd < len(flagAndType) {
						typePart = flagAndType[flagEnd:]
					}

					// Calculate padding needed to reach maxFlagWidth
					currentWidth := len(flagAndType)
					padding := strings.Repeat(" ", maxFlagWidth-currentWidth+2) // +2 for minimum spacing

					// Reconstruct with colored flag name, original spacing, and normal (not muted) description
					lines[i] = "  " + styles.KeyStyle.Render(flagName) + typePart + padding + desc
				} else {
					// No description found, just color the flag name
					parts := strings.SplitN(trimmed, " ", 2)
					if len(parts) >= 1 {
						lines[i] = "  " + styles.KeyStyle.Render(parts[0])
						if len(parts) == 2 {
							lines[i] += " " + parts[1]
						}
					}
				}
			} else {
				// Not a flag line, keep as-is
				lines[i] = line
			}
		}

		return strings.Join(lines, "\n")
	}
}

func init() {
	// Set up custom template functions
	cobra.AddTemplateFunc("styleHeader", styleHeaderFunc())
	cobra.AddTemplateFunc("styleCommand", styleCommandFunc())
	cobra.AddTemplateFunc("styleCommands", styleCommandsFunc())
	cobra.AddTemplateFunc("styleExamples", styleExamplesFunc())
	cobra.AddTemplateFunc("styleFlags", styleFlagsFunc())
	cobra.AddTemplateFunc("wrapWithStyle", func(style string, text string) string {
		return wrapWithStyle(style, text)
	})
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "occo",
	Short: "occo - OpenCode configuration manager",
	Long: `occo is the OpenCode configuration manager CLI.

Manage OpenCode configurations, including presets, sources, and bundle operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag {
			fmt.Printf("occo %s\n", version.Version)
			os.Exit(0)
		}
		_ = cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	// Set color mode from flag
	styles.SetColorMode(styles.ColorMode(colorMode))

	rootCmd.PersistentFlags().BoolVarP(&versionFlag, "version", "v", false, "Print version information")
	rootCmd.PersistentFlags().StringVar(&colorMode, "color", "always", "Color output mode: auto, always, never")

	// Set custom help template
	rootCmd.SetHelpTemplate(styledHelpTemplate())
	rootCmd.SetUsageTemplate(styledUsageTemplate())

	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
