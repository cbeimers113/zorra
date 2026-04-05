package cli

import (
	"github.com/spf13/cobra"

	"github.com/cbeimers113/zorra/internal/log"
)

var connCmd = &cobra.Command{
	Use:   "conn <peer's share code>",
	Short: "Request a connection with a peer",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		log.Infof("Requesting connection with %q...", args[0])
	},
}

func init() {
	rootCmd.AddCommand(connCmd)
}
