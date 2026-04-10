// Package session implements the Zorra session lifecycle
package session

// First connection (TOFU):

// 1. Both peers run `zorra conn <sharecode>`
// 2. Both prompted for session password
// 3. Both start sending NNpsk0 handshake initiation packets to peer's Zorra address
// 4. On receiving a valid initiation, respond with handshake response
// 5. Handshake completes, shared secret established, password authenticated
// 6. Over encrypted channel, exchange static public keys
// 7. Prompt: "Accept connection from alice@laptop? [y/n]"
// 8. On accept, store peer's static public key in known_peers
// 9. Connection established

// Subsequent connections (IK):

// 1. Both peers run `zorra conn <sharecode>`
// 2. No password prompt, static keys already known
// 3. Noise IK handshake, zero round trip, encrypted from first packet
// 4. Connection established

// Session represents a time-bound communication session between
// two Zorra peers, from the perspective of "this" peer
type Session struct{}
