package cli

import (
	"github.com/spf13/cobra"

	"github.com/cbeimers113/zorra/internal/log"
)

var rootCmd = &cobra.Command{
	Use:   "zorra",
	Short: "Serverless peer-to-peer encrypted data transfer over IPv6",
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&log.Verbose, "verbose", "v", false, "enable verbose logging")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
