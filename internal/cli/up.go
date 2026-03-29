package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up <peer address>",
	Short: "Attempt to connect to the given peer",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		peerAddr := args[0]
		fmt.Println("connect to", peerAddr)
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
