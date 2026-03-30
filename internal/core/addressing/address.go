// Package addressing implements functionality for working with Zorra IPv6 addresses
package addressing

import (
	"errors"
	"fmt"
	"hash/fnv"
	"net"

	"github.com/vishvananda/netlink"

	"github.com/cbeimers113/zorra/internal/log"
	"github.com/cbeimers113/zorra/internal/state"
)

// ephemeralIPv6 returns this peer's ephemeral IPv6 address
func ephemeralIPv6() (net.IP, error) {
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

			return deriveEphemeral(addr.IP.Mask(net.CIDRMask(64, 128)), iface)
		}
	}

	return nil, errors.New("no globally addressable IPv6 addresses available")
}

// deriveEphemeral derives an ephemeral, Zorra-specific IPv6 address for this peer to use in this session.
// Format:
// | ISP prefix, customer ID, subnet | ID hash | Zorra hash |
// |            64 bits              | 16 bits |   48 bits  |
func deriveEphemeral(addr net.IP, iface string) (net.IP, error) {
	// Hash this peer's identity
	identity := state.Identity()
	idHash, err := hashString(identity, 2)
	if err != nil {
		return nil, fmt.Errorf("could not hash user identity %q :%w", identity, err)
	}

	// Hash "Zorra"
	zorraHash, err := hashString("Zorra", 6)
	if err != nil {
		return nil, fmt.Errorf("could not hash Zorra: %w", err)
	}

	// Fill in IP mask
	for i, b := range append(idHash, zorraHash...) {
		addr[8+i] = b
	}

	// TODO: address registration
	log.Debugf("Creating ephemeral IPv6 address on interface %q", iface)
	log.Debugf("Ephemeral IPv6 address is %s", addr.String())
	return addr, nil
}

// hashString hashes an input string with 64-bit FNV-1a and returns the lower n bytes in reverse
func hashString(str string, size int) ([]byte, error) {
	hash := fnv.New64a()

	if _, err := hash.Write([]byte(str)); err != nil {
		return nil, fmt.Errorf("could not hash input string %q: %w", str, err)
	}

	// Convert sum to bytes
	sum := hash.Sum64()
	var bytes []byte
	for i := range size {
		shift := uint64(8 * i)
		b := 0xff & (sum >> shift)
		bytes = append(bytes, byte(b))
	}

	return bytes, nil
}
