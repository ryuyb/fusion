package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	envFile string

	// Version information
	version   string
	buildTime string
	gitCommit string
	goVersion string
)

var rootCmd = &cobra.Command{
	Use:   "fusion",
	Short: "Fusion - A lightweight and efficient streaming platform, powered by Go for reliability and Solid.js for smooth interactions.",
	Long:  `The Fusion platform enables users to effortlessly start and watch live streams with minimal delay and high-quality video.`,
}

func SetVersionInfo(ver, build, commit, goVer string) {
	version = ver
	buildTime = build
	gitCommit = commit
	goVersion = goVer
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to execute command: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ./configs/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&envFile, "env", "e", "", "environment (dev, prod)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("./configs")
		viper.SetConfigType("yaml")
		if envFile != "" {
			viper.SetConfigName(fmt.Sprintf("config.%s", envFile))
		} else {
			viper.SetConfigName("config")
		}
	}

	viper.SetEnvPrefix("FUSION")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
