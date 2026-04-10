package cli

import (
	"github.com/spf13/cobra"

	"github.com/cbeimers113/zorra/internal/address"
	"github.com/cbeimers113/zorra/internal/identity"
	"github.com/cbeimers113/zorra/internal/log"
)

var connCmd = &cobra.Command{
	Use:   "conn <sharecode>",
	Short: "Request a connection with a peer",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		ephem, link, err := address.Ephemeral()
		if err != nil {
			log.Errorf("Unable to create ephemeral IPv6 address: %s", err.Error())
			return
		}

		if err = address.Register(ephem, link); err != nil {
			log.Errorf("Unable to register ephemeral IPv6 address on interface %q: %s", link.Attrs().Name, err.Error())
			return
		}

		defer func() {
			if err = address.Deregister(ephem, link); err != nil {
				log.Warnf("Unable to deregister ephemeral IPv6 address from interface %q: %s", link.Attrs().Name, err.Error())
			}
		}()

		peer := args[0]
		log.Infof("Requesting connection with %s...", peer)

		if identity.IsKnown(args[0]) {
			// If peer is known, initiate a Noise IK handshake
			// TODO
			log.Debugf("Initiating Noise IK handshake with %s...", peer)
		} else {
			// If unknown, initiate an NNpsk0 handshake and TOFU prompt
			log.Debugf("Initiating Noise NNpsk0 handshake with %s...", peer)
		}
	},
}

func init() {
	rootCmd.AddCommand(connCmd)
}
