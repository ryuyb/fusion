package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Print the version, build time, git commit, and Go version information.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Fusion %s\n", version)
		fmt.Printf("  Build Time: %s\n", buildTime)
		fmt.Printf("  Git Commit: %s\n", gitCommit)
		fmt.Printf("  Go Version: %s\n", goVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
