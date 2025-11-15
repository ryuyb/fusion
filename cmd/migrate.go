package cmd

import (
	"fmt"
	"os"

	"github.com/ryuyb/fusion/internal/app"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate [up|down|version]",
	Short: "Run database migrations",
	Long:  `Execute database schema migrations`,
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run migrations up",
	Run: func(cmd *cobra.Command, args []string) {
		err := app.RunMigrateApp(app.MigrateUp)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to run migrate up: %v\n", err)
		}
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback migrations",
	Run: func(cmd *cobra.Command, args []string) {
		err := app.RunMigrateApp(app.MigrateDown)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to run migrate down: %v\n", err)
		}
	},
}

var migrateVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show migration version",
	Run: func(cmd *cobra.Command, args []string) {
		err := app.RunMigrateVersionApp()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to run migrate version: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.AddCommand(migrateUpCmd, migrateDownCmd, migrateVersionCmd)
}
