package cmd

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	logLevel string
	logFormat string
)

// rootCmd represents the root command
var rootCmd = &cobra.Command{
	Use:   "builder",
	Short: "Builder execute a pipeline based on a builder.yaml config file",
}

func init() {
	cobra.OnInitialize(initConfig)

	flags := rootCmd.PersistentFlags()

	flags.StringVarP(&logLevel, "logLevel", "", "info", "Log level: trace, debug, info, warn,"+ "error, fatal or panic")
	flags.StringVarP(&logFormat, "logFormat", "", "pretty", "Log format: json or pretty")
}

func initConfig() {
	// logs
	logLevel, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		log.Error().Err(err)
		return
	}

	zerolog.SetGlobalLevel(logLevel)

	if logFormat == "pretty" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()
}

// Execute adds all child commands to the root command and sets flags
// appropriately. This is called by main.main(). It only needs to happen once
// to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

