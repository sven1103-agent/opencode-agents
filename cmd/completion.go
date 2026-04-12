package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completions",
	Long: `Generate shell completion scripts for supported shells.

Bash:
  # Source on the fly
  source <(oc completion bash)

  # Or install permanently
  oc completion bash | sudo tee /etc/bash_completion.d/oc > /dev/null

Zsh:
  # Option 1: Source on the fly (recommended)
  # Add to end of ~/.zshrc (after compinit):
  source <(oc completion zsh)
  # If you use alias (alias oc="opencode-config-cli"), also run:
  compdef _oc opencode-config-cli

  # Option 2: Save to completions dir
  oc completion zsh > ~/.zsh/completions/_oc
  # Ensure ~/.zshrc has: fpath=(~/.zsh/completions $fpath)
  # Then add: compdef _oc opencode-config-cli

  # Clear completion cache and restart:
  rm -f ~/.zcompdump && exec zsh

Fish:
  oc completion fish > ~/.config/fish/completions/oc.fish

PowerShell:
  oc completion powershell >> $PROFILE

For more details, see: https://github.com/sven1103-agent/opencode-config-cli#shell-completion
`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			return fmt.Errorf("unsupported shell: %s", args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
