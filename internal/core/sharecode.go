package core

import (
	"encoding/ascii85"
	"fmt"
	"net"
)

// IPv6ToShareCode returns the share code for an IPv6 address
func IPv6ToShareCode(addr net.IP) string {
	src := addr
	dst := make([]byte, ascii85.MaxEncodedLen(len(src)))
	n := ascii85.Encode(dst, src)

	return string(dst[:n])
}

// ShareCodeToIPv6 returns the IPv6 address for a share code
func ShareCodeToIPv6(shareCode string) (net.IP, error) {
	src := []byte(shareCode)
	dst := make([]byte, 39)
	n, _, err := ascii85.Decode(dst, src, true)
  if err != nil {
		return nil, fmt.Errorf("unable to decode share code %q: %w", shareCode, err)
	}

	return net.IP(dst[:n]), nil
}

