// Package hash implements various hashes used by Zorra
package hash

import (
	"hash/fnv"
	"runtime"

	"golang.org/x/crypto/argon2"
)

// Zorra is the String hash of "Zorra", for use in derived IPv6 addresses
var Zorra = String("Zorra", 6)

// String hashes an input string with 64-bit FNV-1a and returns the lower n bytes in reverse.
// Not cryptographically secure, optimized for speed and high entropy in output bits
func String(str string, size uint8) (bytes []byte) {
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

// Password hashes an input password with argon2id, using a deterministic
// salt derived from the password byte slice
func Password(pw []byte) []byte {
	salt := make([]byte, 16)
	salt[0] = pw[len(pw)/2]

	for i := 1; i < len(salt); i++ {
		last := int(salt[i-1])
		next := last % len(pw)
		salt[i] = pw[next]
	}

	return argon2.IDKey(pw, salt, 1, 64*1024, uint8(runtime.NumCPU()), 32)
}
