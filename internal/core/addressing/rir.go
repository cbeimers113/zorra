package addressing

import (
	"net"
)

type RIR int

const (
	APNIC RIR = iota
	AFRINIC
	ARIN
	LACNIC
	RIPE
	UnknownRIR
)

// RIRName returns the acronym of an RIR
func RIRName(rir RIR) string {
	switch rir {
	case APNIC:
		return "APNIC"
	case AFRINIC:
		return "AFRINIC"
	case ARIN:
		return "ARIN"
	case LACNIC:
		return "LACNIC"
	case RIPE:
		return "RIPE"
	}

	return "unknown"
}

// detectRIR determines a peer's RIR from its IPv6 address
func detectRIR(ip net.IP) RIR {
	b := ip[0]

	switch {
	case b >= 0x24 && b < 0x28: // 2400::/6
		return APNIC
	case b >= 0x28 && b < 0x2C: // 2800::/6
		return AFRINIC
	case b >= 0x2C && b < 0x30: // 2C00::/6
		return ARIN
	case b >= 0x30 && b < 0x34: // 3000::/6
		return LACNIC
	case b >= 0x34 && b < 0x38: // 3400::/6
		return RIPE
	}

	return UnknownRIR
}
