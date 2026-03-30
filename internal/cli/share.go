package cli

import (
	"github.com/spf13/cobra"

	"github.com/cbeimers113/zorra/internal/core/addressing"
	"github.com/cbeimers113/zorra/internal/log"
)

var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "Determine your Zorra share code",
	Run: func(*cobra.Command, []string) {
		shareCode, err := addressing.CreateShareCode()
		if err != nil {
			log.Errorf("Unable to determine share code: %s", err.Error())
			return
		}

		log.Infof("Zorra share code: %s", shareCode)
	},
}

func init() {
	rootCmd.AddCommand(shareCmd)
}
