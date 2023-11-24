package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/teng231/suggest/pkg/lm"
	"github.com/teng231/suggest/pkg/store"
)

func init() {
	rootCmd.AddCommand(buildLMCmd)
}

var buildLMCmd = &cobra.Command{
	Use:   "build-lm -c [config path]",
	Short: "builds ngram language model for the given config",
	Long:  `builds ngram language model for the given config and saves it in the binary format`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := lm.ReadConfig(configPath)

		if err != nil {
			return fmt.Errorf("couldn't read a config %w", err)
		}

		directory, err := store.NewFSDirectory(config.GetOutputPath())

		if err != nil {
			return fmt.Errorf("failed to create a fs directory: %w", err)
		}

		return lm.StoreBinaryLMFromGoogleFormat(directory, config)
	},
}
