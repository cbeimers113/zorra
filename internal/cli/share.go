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
		ephem, err := addressing.EphemeralIPv6()
		if err != nil {
			log.Errorf("Unable to create ephemeral IPv6 address: %s", err.Error())
			return
		}

		shareCode, err := addressing.CreateShareCode(ephem)
		if err != nil {
			log.Errorf("Unable to determine share code: %s", err.Error())
			return
		}

		log.Infof("Zorra share code: %s", shareCode)
		addressing.ReadShareCode(shareCode)
	},
}

func init() {
	rootCmd.AddCommand(shareCmd)
}
