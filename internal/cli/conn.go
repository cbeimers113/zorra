package cli

import (
	"github.com/spf13/cobra"

	"github.com/cbeimers113/zorra/internal/address"
	"github.com/cbeimers113/zorra/internal/identity"
	"github.com/cbeimers113/zorra/internal/log"
	"github.com/cbeimers113/zorra/internal/session"
	"github.com/cbeimers113/zorra/internal/sharecode"
)

var connCmd = &cobra.Command{
	Use:   "conn <sharecode>",
	Short: "Request a connection with a peer",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Create and register the ephemeral address
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

		// Parse the share code
		addr, err := sharecode.Read(args[0])
		if err != nil {
			log.Errorf("Unable to read share code %q: %s", args[0], err.Error())
			return
		}

		// Check if we're connecting to a known peer, prompt for session password if not
		var pwHash []byte
		if !identity.IsKnown(args[0]) {
			pwHash, err = getPasswordHash()
			if err != nil {
				log.Errorf("Unable to get session password: %s", err.Error())
				return
			}
		}

		// Configure a session
		sess, err := session.New(addr, pwHash)
		if err != nil {
			log.Errorf("Unable to create session: %s", err.Error())
			return
		}

		// Connect to the peer
		log.Infof("Requesting connection with %s...", args[0])
		if err = sess.Initiate(cmd.Context()); err != nil {
			log.Errorf("Unable to connect to %s: %s", args[0], err.Error())
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(connCmd)
}
