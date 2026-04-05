// Package testdata contains data used throughout the Zorra test suite
package testdata

import "net"

var (
	Identity = "zorra-test"

	// hash of the test identity + Zorra
	addrSuffix = "d071:7122:3c5f:626e"

	EphemAddrHasChannel = net.ParseIP("2605:59c0:1710:beef:" + addrSuffix)
	ShareCodeHasChannel = Identity + ".wBcQvu8.14ef"

	EphemAddrNoChannel = net.ParseIP("dead:beef:1234:1010:" + addrSuffix)
	ShareCodeNoChannel = Identity + ".3q2-7xI0EBA._"

	EphemAddrByteBoundary = net.ParseIP("2630:1111:2222:3333:" + addrSuffix)
	ShareCodeByteBoundary = Identity + ".EREiIjMz.2b6c"
)
