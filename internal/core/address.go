// Package core implements the core functionality of Zorra
package core

import (
	"errors"
	"fmt"
	"net"

	"github.com/vishvananda/netlink"

	"github.com/cbeimers113/zorra/internal/log"
)

// Address returns this peer's ephemeral IPv6 address
func Address() (net.IP, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return nil, fmt.Errorf("unable to list network interfaces: %w", err)
	}

	for _, link := range links {
		iface := link.Attrs().Name
		addrs, err := netlink.AddrList(link, netlink.FAMILY_V6)
		if err != nil {
			return nil, fmt.Errorf("unable to find IPv6 addresses for interface %q: %w", iface, err)
		}

		for _, addr := range addrs {
			if !addr.IP.IsGlobalUnicast() || addr.IP.IsPrivate() {
				continue
			}

			return deriveIPv6(addr.IP.Mask(net.CIDRMask(64, 128)), iface)
		}
	}

	return nil, errors.New("no globally addressable IPv6 addresses available")
}

// deriveIPv6 derives an ephemeral, Zorra-specific IPv6 address for this peer to use in this session
func deriveIPv6(mask net.IP, iface string) (net.IP, error) {
	// TODO: For testing; replace with secure pseudorandom derivation
	for i, b := range []byte{
		0xab,
		0xcd,
		0xef,
		0xfe,
		0xdc,
		0xba,
		0x11,
		0x99,
	} {
		mask[8+i] = b
	}

	log.Debugf("Creating ephemeral IPv6 address on interface %q", iface)
	// TODO: address registration
	return mask, nil
}
