// Package transport implements the sending and recieving of UDP packets
package transport

import (
	"context"
	"fmt"
	"net"

	"github.com/cbeimers113/zorra/internal/log"
)

const (
	// Buffer size for the inbound channel; too small will
	// block the listener and drop packets, too large will
	// buffer stale packets
	inboundBufferSize = 256

	// Packet data size; this is the theoretical max
	// size of a UDP packet in bytes
	packetDataSize = 65535
)

// Packet represents a UDP packet; it contains
// the raw data bytes and the sender's address
type Packet struct {
	Data []byte
	Addr *net.UDPAddr
}

// Transport represents a UDP transport session.
// It maintains a UDP listener and a channel
// to write incoming packets to
type Transport struct {
	conn *net.UDPConn

	// Inbound UDP packet channel;
	// unbuffered and must be read
	// on caller's goroutine
	Inbound chan Packet
}

// New initializes a new Transport session
func New(addr *net.UDPAddr) (*Transport, error) {
	conn, err := net.ListenUDP("udp6", addr)
	if err != nil {
		return nil, fmt.Errorf("could not create UDP listener: %w", err)
	}

	return &Transport{
		conn:    conn,
		Inbound: make(chan Packet, inboundBufferSize),
	}, nil
}

// Listen waits for incoming UDP packets and
// writes them onto the inbound channel.
// This method should be called asynchronously
func (t *Transport) Listen(ctx context.Context) {
	buf := make([]byte, packetDataSize)

listening:
	for {
		select {
		case <-ctx.Done():
			break listening
		default:
			n, addr, err := t.conn.ReadFromUDP(buf)
			if err != nil {
				log.Warnf("UDP read error: %s", err.Error())
				return
			}

			pkt := Packet{
				Addr: addr,
				Data: make([]byte, n),
			}

			// Copy the buffered bytes into the packet
			// so that we can reuse the buffer slice
			copy(pkt.Data, buf[:n])
			t.Inbound <- pkt
		}
	}
}

// Send sends a UDP packet with the given data to the specified address
func (t *Transport) Send(data []byte, addr *net.UDPAddr) error {
	_, err := t.conn.WriteToUDP(data, addr)
	return err
}

// Close closes the UDP listener connection
func (t *Transport) Close() {
	t.conn.Close()
}
