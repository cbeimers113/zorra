// Package address implements logic for working with Zorra IPv6 addresses
package address

import (
	"errors"
	"fmt"
	"net"
	"syscall"

	"github.com/vishvananda/netlink"

	"github.com/cbeimers113/zorra/internal/hash"
	"github.com/cbeimers113/zorra/internal/identity"
	"github.com/cbeimers113/zorra/internal/log"
)

// LinkList is a late-bound function that returns a slice of netlink.Link
// representing the network interfaces found on this device. Can be reassigned
// in tests to simulate different scenarios
var LinkList func() ([]netlink.Link, error) = netlink.LinkList

// AddrList is a late-bound function that returns a list of addresses on
// a network interface, optionally filtered by family. Can be reassigned
// in tests to simulate different scenarios
var AddrList func(netlink.Link, int) ([]netlink.Addr, error) = netlink.AddrList

// AddrAdd is a late-bound function that registers an IP address to a network
// interface. Can be reassigned in tests to simulate different scenarios
var AddrAdd func(netlink.Link, *netlink.Addr) error = netlink.AddrAdd

// AddrDel is a late-bound function that deregisters an IP address from a network
// interface. Can be reassigned in tests to simulate different scenarios
var AddrDel func(netlink.Link, *netlink.Addr) error = netlink.AddrDel

// Ephemeral returns this peer's ephemeral IPv6 address and the network
// interface it can be registered on. The ephemeral address is a Zorra-specific
// address derived from an existing address on the network, used only for the
// duration of the Zorra session.
// Requires the user's ISP to provide IPv6 and an existing globally
// addressable IPv6 address to use as a base.
// Format:
// | ISP prefix, customer ID, subnet | ID hash | Zorra hash |
// |            64 bits              | 16 bits |   48 bits  |
func Ephemeral() (net.IP, netlink.Link, error) {
	links, err := LinkList()
	if err != nil {
		return nil, nil, fmt.Errorf("could not find network interfaces: %w", err)
	}

	// Find an interface with IPv6 addresses
	for _, link := range links {
		addrs, err := AddrList(link, netlink.FAMILY_V6)
		if err != nil {
			return nil, nil, fmt.Errorf("could not list addresses on interface %q: %w", link.Attrs().Name, err)
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
			return addr, link, nil
		}
	}

	return nil, nil, errors.New("no globally addressable IPv6 addresses available")
}

// Register registers the given address on the given network interface
func Register(addr net.IP, link netlink.Link) error {
	log.Debugf("Registering address %q on interface %q", addr.String(), link.Attrs().Name)
	err := AddrAdd(link, toNetlinkAddr(addr))

	// Zorra address already present
	if errors.Is(err, syscall.EEXIST) {
		log.Warnf("Ephemeral IPv6 address %q is already registered", addr.String())
		return nil
	}

	return err
}

// Deregister removes the given address from the given network interface
func Deregister(addr net.IP, link netlink.Link) error {
	log.Debugf("Deregistering address %q from interface %q", addr.String(), link.Attrs().Name)
	return AddrDel(link, toNetlinkAddr(addr))
}

// toNetlinkAddr creates a *netlink.Addr with 0 subnets from a net.IP
func toNetlinkAddr(addr net.IP) *netlink.Addr {
	return &netlink.Addr{
		IPNet: &net.IPNet{
			IP:   addr,
			Mask: net.CIDRMask(128, 128),
		},
	}
}
