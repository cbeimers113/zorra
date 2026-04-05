// Package address implements logic for working with Zorra IPv6 addresses
package address

import (
	"errors"
	"fmt"
	"net"

	"github.com/vishvananda/netlink"

	"github.com/cbeimers113/zorra/internal/hash"
	"github.com/cbeimers113/zorra/internal/log"
	"github.com/cbeimers113/zorra/internal/state"
)

// EphemeralIPv6 returns this peer's ephemeral IPv6 address:
// a Zorra-specific IPv6 address for this peer to use in this session.
// Format:
// | ISP prefix, customer ID, subnet | ID hash | Zorra hash |
// |            64 bits              | 16 bits |   48 bits  |
func EphemeralIPv6() (net.IP, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return nil, fmt.Errorf("could not find network interfaces: %w", err)
	}

	// Find an interface with IPv6 addresses
	for _, link := range links {
		iface := link.Attrs().Name
		addrs, err := netlink.AddrList(link, netlink.FAMILY_V6)
		if err != nil {
			return nil, fmt.Errorf("could not list addresses on interface %q: %w", iface, err)
		}

		// Find a globally addressable address and derive a new address from its prefix
		for _, addr := range addrs {
			if !addr.IP.IsGlobalUnicast() || addr.IP.IsPrivate() {
				continue
			}

			// Hash this peer's identity and fill in the second half of the address
			addr := addr.IP.Mask(net.CIDRMask(64, 128))
			idHash := hash.String(state.Identity(), 2)
			for i, b := range append(idHash, hash.Zorra...) {
				addr[8+i] = b
			}

			log.Debugf("Ephemeral IPv6 address is %s", addr.String())
			return addr, nil
		}
	}

	return nil, errors.New("no globally addressable IPv6 addresses available")
}

// RegisterEphemeralIPv6 registers the ephemeral IPv6 address on the given network interface
func RegisterEphemeralIPv6(addr net.IP, iface string) {
	// TODO: address registration
	log.Debugf("Registering address %q on interface %q", addr.String(), iface)
}

// DeregisterEphemeralIPv6 removes the ephemeral IPv6 address from the given network interface
func DeregisterEphemeralIPv6(addr net.IP, iface string) {
	// TODO: address deregistration
	log.Debugf("Deregistering address %q from interface %q", addr.String(), iface)
}
