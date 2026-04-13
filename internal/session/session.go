// Package session implements the Zorra session lifecycle
package session

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/cbeimers113/zorra/internal/transport"
	"github.com/flynn/noise"
)

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

// How long to wait before a connection attempt times out
var handshakeTimeout = 30 * time.Second

// Session represents a time-bound communication session between
// two Zorra peers, from the perspective of "this" peer
type Session struct {
	addr      net.IP
	trns      *transport.Transport
	handshake *noise.HandshakeState
}

// New configures a new session with the given peer. The session is not
// opened until a handshake is initiated and completed by calling
// Initiate()
func New(addr net.IP, pwHash []byte) (*Session, error) {
	// Create the UDP transport session
	trns, err := transport.New(nil) // TODO: addr -> UDP addr
	if err != nil {
		return nil, fmt.Errorf("could not initialize transport session: %w", err)
	}

	// Configure handshake based on whether we have a session password;
	// no password: we know the peer and can do Noise IK
	// have password: we don't know the peer so will do Noise NNpsk0
	pattern := noise.HandshakeIK
	if len(pwHash) > 0 {
		pattern = noise.HandshakeNN
	}

	config := noise.Config{
		CipherSuite:  noise.NewCipherSuite(noise.DH25519, noise.CipherChaChaPoly, noise.HashSHA256),
		Pattern:      pattern,
		Initiator:    true,
		PresharedKey: pwHash,
		// StaticKeypair: TODO; for IK
		// PeerStatic: TODO: for IK
	}

	handshake, err := noise.NewHandshakeState(config)
	if err != nil {
		return nil, fmt.Errorf("could not configure handshake: %w", err)
	}

	return &Session{
		addr:      addr,
		trns:      trns,
		handshake: handshake,
	}, nil
}

// Initiate begins the handshake state machine to establish a secure connection to a peer
func (s *Session) Initiate(ctx context.Context) error {
	//
	return nil
}
