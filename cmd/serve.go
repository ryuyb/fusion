package cmd

import (
	"fmt"
	"os"

	"github.com/ryuyb/fusion/internal/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",
	Long:  `Start the HTTP API server with all dependencies`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := app.RunServeApp(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to run serve: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().String("host", "0.0.0.0", "The host to bind the API server")
	serveCmd.Flags().IntP("port", "p", 8080, "Server port")

	_ = viper.BindPFlag("server.host", serveCmd.Flags().Lookup("host"))
	_ = viper.BindPFlag("server.port", serveCmd.Flags().Lookup("port"))
}
