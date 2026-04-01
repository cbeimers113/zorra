package addressing

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/netip"
	"strconv"
	"strings"

	"github.com/cbeimers113/zorra/internal/log"
	"github.com/cbeimers113/zorra/internal/state"
)

// CreateShareCode determines this peer's ephemeral IPv6 address,
// channel, and identity and encodes them into a share code
func CreateShareCode() (string, error) {
	ephem, err := ephemeralIPv6()
	if err != nil {
		return "", err
	}

	// Parse the ephemeral IPv6 address
	addr, ok := netip.AddrFromSlice(ephem)
	if !ok {
		return "", fmt.Errorf("invalid ephemeral IPv6 address: %q", ephem.String())
	}

	// Find the longest IPv6 prefix of this address that's in the channel map
	bits := 64
	channel := 0
	for bits > 0 {
		if channel, ok = ChannelOf(netip.PrefixFrom(addr, bits).Masked()); ok {
			break
		}

		bits--
	}

	// Find the remaining opaque customer bits
	opaque := make([]byte, 0, 8-bits/8)
	mod := bits % 8

	for i := bits / 8; i < 8; i++ {
		next := ephem[i]

		// If the prefix didn't end cleanly on a byte boundary, mask its lower bits
		if i == 0 && mod != 0 {
			next = ephem[bits/8] & (0xff >> mod)
		}

		opaque = append(opaque, next)
	}

	// Share code = ID.opaque.channel
	// ID: no encoding, easiest for human-to-human transfer
	// opaque bits: base64, best balance of compact and readable
	// channel: hex, small channel space and easy to work with
	shareCode := state.Identity() + "." + base64.StdEncoding.EncodeToString(opaque) + "."

	// If no channel was found, use a non-numerical "no channel" identifier
	if bits == 0 {
		shareCode += "_"
	} else {
		shareCode += fmt.Sprintf("%x", channel)
	}

	return shareCode, nil
}

// ReadShareCode parses a share code into an IPv6 address
func ReadShareCode(shareCode string) (net.IP, error) {
	parts := strings.Split(shareCode, ".")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid share code: %q", shareCode)
	}

	// Extract encoded parts
	encID := strings.Join(parts[:len(parts)-2], ".")
	encOpaque := parts[len(parts)-2]
	encChannel := parts[len(parts)-1]

	var (
		prefix netip.Prefix
		ok     bool
	)

	// Parse the channel into a prefix if applicable
	if encChannel != "_" {
		channel, err := strconv.ParseUint(encChannel, 16, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid channel: %q", encChannel)
		}

		prefix, ok = PrefixOf(int(channel))
		if !ok {
			return nil, fmt.Errorf("unknown channel: %q", encChannel)
		}
	}

	// Parse the opaque bits
	opaque, err := base64.StdEncoding.DecodeString(encOpaque)
	if err != nil {
		return nil, fmt.Errorf("could not decode opaque bits %q: %w", encOpaque, err)
	}

	// Combine the prefix and opaque bits into the first half of the IPv6 address
	var addr [16]byte
	// TODO: Pick up here
	return nil, nil
}
