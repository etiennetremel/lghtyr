package cmd

import (
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/etiennetremel/lghtyr/pkg/builder"
)

type runCmdFlagsStruct struct {
	config    string
}

var runCmdFlags runCmdFlagsStruct

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [JOB NAME]",
	Short: "Run a given job from a builder.yaml config file",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("missing argument job name is missing")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		jobName := args[0]

		log.Info().Msgf("Executing job %s", jobName)

		// display execution time
		start := time.Now()
		defer func() {
			elapsed := time.Since(start)
			log.Info().Msgf("Job %s took %s to execute", jobName, elapsed)
		}()

		j, err := builder.NewBuilder(runCmdFlags.config)
		if err != nil {
			log.Error().Err(err).Msg("Failed instanciating build")
			return
		}

		err = j.RunJob(jobName)
		if err != nil {
			log.Error().Err(err).Msg("failed")
			return
		}
		return
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	flags := runCmd.Flags()
	flags.StringVarP(&runCmdFlags.config, "config", "c", "", "builder file (default is builder.yaml)")
}
