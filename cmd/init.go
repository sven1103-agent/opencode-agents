package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/qbicsoftware/occo/internal/preset"
	"github.com/qbicsoftware/occo/internal/schema"
	"github.com/qbicsoftware/occo/internal/styles"
	"github.com/spf13/cobra"
)

// validateOutputPath ensures the output path stays within the project root
// to prevent path traversal attacks.
func validateOutputPath(projectRoot, outputPath string) error {
	absProjectRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		return fmt.Errorf("failed to resolve project root: %w", err)
	}
	absOutputPath, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("failed to resolve output path: %w", err)
	}

	absProjectRoot = filepath.Clean(absProjectRoot)
	absOutputPath = filepath.Clean(absOutputPath)

	if !strings.HasPrefix(absOutputPath, absProjectRoot+string(filepath.Separator)) {
		if absOutputPath != absProjectRoot {
			return fmt.Errorf("invalid output path: path traversal detected (output must be within project root)")
		}
	}

	return nil
}

var (
	initProjectRoot string
	initOutput      string
	initForce       bool
	initDryRun      bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a project with OpenCode configuration and install schemas",
	Long: `Initialize a project by:
1. Copying the default config to the project root
2. Installing schemas to .opencode/schemas/

The default config is bundled with the installation.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit()
	},
}

func init() {
	initCmd.Flags().StringVar(&initProjectRoot, "project-root", ".", "Project root directory")
	initCmd.Flags().StringVar(&initOutput, "output", "opencode.json", "Output file path")
	initCmd.Flags().BoolVar(&initForce, "force", false, "Overwrite existing files")
	initCmd.Flags().BoolVar(&initDryRun, "dry-run", false, "Show what would be done without doing it")

	rootCmd.AddCommand(initCmd)
}

func runInit() error {
	// Resolve project root
	projectRoot, err := filepath.Abs(initProjectRoot)
	if err != nil {
		return fmt.Errorf("invalid project root: %w", err)
	}

	// Check if project root exists
	if _, err := os.Stat(projectRoot); os.IsNotExist(err) {
		return fmt.Errorf("project root does not exist: %s", projectRoot)
	}

	// Resolve output path
	outputPath := filepath.Join(projectRoot, initOutput)

	// Validate output path to prevent path traversal
	if err := validateOutputPath(projectRoot, outputPath); err != nil {
		return err
	}

	// Get default config (bundled with installation)
	configData, err := preset.GetDefaultConfig()
	if err != nil {
		return fmt.Errorf("failed to get default config: %w", err)
	}

	// Dry run mode
	if initDryRun {
		fmt.Println(styles.SectionHeader("Dry Run"))
		fmt.Println(styles.SubHeader("Operations"))
		fmt.Println(styles.DryRun(fmt.Sprintf("write config to %s", outputPath)))
		fmt.Println(styles.DryRun(fmt.Sprintf("install schemas to %s/.opencode/schemas/", projectRoot)))
		fmt.Println()
		fmt.Println(styles.Done("dry-run complete"))
		return nil
	}

	// Write config file
	if err := preset.WriteConfig(outputPath, configData, initForce); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	fmt.Println(styles.Written(outputPath))

	// Install schemas
	opencodeDir := filepath.Join(projectRoot, ".opencode")
	if err := schema.InstallAll(opencodeDir, initForce); err != nil {
		return fmt.Errorf("failed to install schemas: %w", err)
	}
	fmt.Println(styles.Written(fmt.Sprintf("%s/.opencode/schemas/", projectRoot)))

	fmt.Println()
	fmt.Println(styles.Done("init complete"))

	return nil
}
