package cli

import (
	"github.com/spf13/cobra"

	"github.com/cbeimers113/zorra/internal/address"
	"github.com/cbeimers113/zorra/internal/log"
	"github.com/cbeimers113/zorra/internal/sharecode"
)

var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "Determine your Zorra share code",
	Run: func(*cobra.Command, []string) {
		ephem, err := address.EphemeralIPv6()
		if err != nil {
			log.Errorf("Unable to create ephemeral IPv6 address: %s", err.Error())
			return
		}

		shareCode, err := sharecode.CreateShareCode(ephem)
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
