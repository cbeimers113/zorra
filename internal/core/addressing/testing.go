package addressing

import (
	"net"

	"github.com/cbeimers113/zorra/internal/state"
)

var (
	// hash of the test identity + Zorra
	testAddrSuffix = "d071:7122:3c5f:626e"

	TestEphemAddrHasChannel = net.ParseIP("2605:59c0:1710:beef:" + testAddrSuffix)
	TestShareCodeHasChannel = state.TestIdentity + ".wBcQvu8.14ef"

	TestEphemAddrNoChannel = net.ParseIP("dead:beef:1234:1010:" + testAddrSuffix)
	TestShareCodeNoChannel = state.TestIdentity + ".3q2-7xI0EBA._"

	TestEphemAddrByteBoundary = net.ParseIP("2630:1111:2222:3333:" + testAddrSuffix)
	TestShareCodeByteBoundary = state.TestIdentity + ".EREiIjMz.2b6c"
)
