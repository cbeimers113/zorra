// Package address implements logic for working with Zorra IPv6 addresses
package address

import (
	"errors"
	"fmt"
	"net"

	"github.com/vishvananda/netlink"

	"github.com/cbeimers113/zorra/internal/hash"
	"github.com/cbeimers113/zorra/internal/identity"
	"github.com/cbeimers113/zorra/internal/log"
)

// LinkList is a late-bound function that returns a slice of netlink.Link
// representing the network interfaces found on this device. The production
// implmentation is provided by vishvananda/netlink and wraps syscalls to the
// netlink kernel module. It can be reassigned to mock its response in tests.
var LinkList func() ([]netlink.Link, error) = netlink.LinkList

// AddrList is a late-bound function that returns a list of addresses on
// a network interface, optionally filtered by family. The production
// implementation is provided by vishvananda/netlink and wraps syscalls to
// the netlink kernel module. It can be reassigned to mock its response in tests.
var AddrList func(netlink.Link, int) ([]netlink.Addr, error) = netlink.AddrList

// Ephemeral returns this peer's ephemeral IPv6 address:
// a Zorra-specific IPv6 address for this peer to use in this session.
// Requires the user's ISP to provide IPv6 and an existing globally
// addressable IPv6 address to use as a base.
// Format:
// | ISP prefix, customer ID, subnet | ID hash | Zorra hash |
// |            64 bits              | 16 bits |   48 bits  |
func Ephemeral() (net.IP, error) {
	links, err := LinkList()
	if err != nil {
		return nil, fmt.Errorf("could not find network interfaces: %w", err)
	}

	// Find an interface with IPv6 addresses
	for _, link := range links {
		iface := link.Attrs().Name
		addrs, err := AddrList(link, netlink.FAMILY_V6)
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
			idHash := hash.String(identity.This(), 2)
			for i, b := range append(idHash, hash.Zorra...) {
				addr[8+i] = b
			}

			log.Debugf("Ephemeral IPv6 address is %s", addr.String())
			return addr, nil
		}
	}

	return nil, errors.New("no globally addressable IPv6 addresses available")
}

// Register registers the given address on the given network interface
func Register(addr net.IP, iface string) {
	// TODO: address registration
	log.Debugf("Registering address %q on interface %q", addr.String(), iface)
}

// Deregister removes the given address from the given network interface
func Deregister(addr net.IP, iface string) {
	// TODO: address deregistration
	log.Debugf("Deregistering address %q from interface %q", addr.String(), iface)
}
