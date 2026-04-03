// Package addressing implements functionality for working with Zorra IPv6 addresses
package addressing

import "hash/fnv"

// Compute the hash of "Zorra" on init for use in derived IPv6 addresses
var zorraHash = hashString("Zorra", 6)

// hashString hashes an input string with 64-bit FNV-1a and returns the lower n bytes in reverse
func hashString(str string, size uint8) (bytes []byte) {
	hash := fnv.New64a()
	hash.Write([]byte(str))
	sum := hash.Sum64()

	// Limit to 8 bytes
	for i := range min(size, 8) {
		shift := uint64(8 * i)
		b := 0xff & (sum >> shift)
		bytes = append(bytes, byte(b))
	}

	return bytes
}
