package cli

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/cbeimers113/zorra/internal/core"
	"github.com/cbeimers113/zorra/internal/log"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Output IPv6 address, port, and share code",
	Run: func(*cobra.Command, []string) {
		addr, err := core.Address()
		if err != nil {
			log.Error(err.Error())
			os.Exit(1)
		}

		log.Infof("IPv6 address is %s", addr)
		log.Infof("Share code is %s", core.IPv6ToShareCode(addr))
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
