// Package test contains data used throughout the Zorra test suite
package test

import (
	"net"

	"github.com/vishvananda/netlink"
)

var (
	// hash of the test identity + Zorra
	addrSuffix = "d071:7122:3c5f:626e"

	// Identity used by "this" peer in tests
	Identity = "zorra-test"

	// Mock base IPv6 address and network interface for address derivation in tests
	BaseAddr      = net.ParseIP("2605:59c0:dead:beef:0000:1111:2222:3333")
	InterfaceName = "test0"

	// Mock derived ephemeral IPv6 address from the above BaseAddr
	EphemAddr = net.ParseIP("2605:59c0:dead:beef:" + addrSuffix)

	// Mock ephemeral addresses and their share codes for testing different cases
	EphemAddrHasChannel = net.ParseIP("2605:59c0:1710:beef:" + addrSuffix)
	ShareCodeHasChannel = Identity + ".wBcQvu8.14ef"

	EphemAddrNoChannel = net.ParseIP("dead:beef:1234:1010:" + addrSuffix)
	ShareCodeNoChannel = Identity + ".3q2-7xI0EBA._"

	EphemAddrByteBoundary = net.ParseIP("2630:1111:2222:3333:" + addrSuffix)
	ShareCodeByteBoundary = Identity + ".EREiIjMz.2b6c"
)

// MockLink implements the netlink.Link interface
type MockLink struct{}

// Attrs returns the an empty *netlink.LinkAttrs
func (m *MockLink) Attrs() *netlink.LinkAttrs {
	return new(netlink.LinkAttrs)
}

// Type returns the type of the MockLink
func (m *MockLink) Type() string {
	return "mock"
}

// MockAddr is a mock representation of a netlink.Addr
// that encapsulates the BaseAddr
var MockAddr = netlink.Addr{
	IPNet: &net.IPNet{
		IP:   BaseAddr,
		Mask: net.CIDRMask(64, 128),
	},
}
