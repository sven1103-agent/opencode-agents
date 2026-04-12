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
  # Add to end of ~/.zshrc:
  source <(oc completion zsh)

  # Option 2: Save to completions dir
  oc completion zsh > ~/.zsh/completions/_oc
  # Ensure ~/.zshrc has: fpath=(~/.zsh/completions $fpath)

  # Clear completion cache and restart:
  rm -f ~/.zcompdump && exec zsh

  # Note: If you get "parse error" when pressing TAB, try:
  # 1. Ensure you're using the binary name (oc), not an alias
  # 2. Or add explicit completion binding after sourcing:
  #    compdef _oc oc

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
