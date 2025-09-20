package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the baseline by creating the target directory",
	Long: `Initialize the baseline by creating the target directory if it doesn't exist.

This command prepares the baseline directory structure for cloning repositories.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if verbose {
			fmt.Printf("Initializing baseline directory: %s\n", directory)
		}

		// Create the target directory if it doesn't exist
		if err := os.MkdirAll(directory, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", directory, err)
		}

		fmt.Printf("Successfully initialized baseline directory: %s\n", directory)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
