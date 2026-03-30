package addressing 

import (
  "encoding/ascii85"
	"net"
)

// CreateShareCode determines this peer's ephemeral IPv6 address and encodes it into a share code
func CreateShareCode() (string, error) {
	src, err := ephemeralIPv6()
	if err != nil {
		return "", err
	}

	// Encode the address as ASCII85
	dst := make([]byte, ascii85.MaxEncodedLen(len(src)))
	n := ascii85.Encode(dst, src)

	return string(dst[:n]), nil
}

// ReadShareCode parses a share code into an IPv6 address
func ReadShareCode(shareCode string) (net.IP, error) {
	src := []byte(shareCode)
	dst := make([]byte, 39)
	n, _, err := ascii85.Decode(dst, src, true)
  if err != nil {
		return nil, err 
	}

	return net.IP(dst[:n]), nil
}
