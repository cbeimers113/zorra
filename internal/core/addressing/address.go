package addressing

import (
	"errors"
	"fmt"
	"net"

	"github.com/vishvananda/netlink"

	"github.com/cbeimers113/zorra/internal/log"
	"github.com/cbeimers113/zorra/internal/state"
)

// EphemeralIPv6 returns this peer's ephemeral IPv6 address
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

			return deriveEphemeral(addr.IP.Mask(net.CIDRMask(64, 128)), iface), nil
		}
	}

	return nil, errors.New("no globally addressable IPv6 addresses available")
}

// deriveEphemeral derives an ephemeral, Zorra-specific IPv6 address for this peer to use in this session.
// Format:
// | ISP prefix, customer ID, subnet | ID hash | Zorra hash |
// |            64 bits              | 16 bits |   48 bits  |
func deriveEphemeral(addr net.IP, iface string) net.IP {
	// Hash this peer's identity and fill in the second half of the address
	idHash := hashString(state.Identity(), 2)
	for i, b := range append(idHash, zorraHash...) {
		addr[8+i] = b
	}

	// TODO: address registration
	log.Debugf("Creating ephemeral IPv6 address on interface %q", iface)
	log.Debugf("Ephemeral IPv6 address is %s", addr.String())
	return addr
}
