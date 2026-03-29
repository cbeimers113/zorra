# Zorra

A fully decentralized peer-to-peer encrypted data transfer tool built on a direct IPv6 tunnel primitive.

Zorra lets you send messages and files directly to a friend; no cloud, no relay, and no accounts are required.

---

## Goals

- **Direct peer-to-peer**: data travels directly between peers with no intermediary server
- **Decentralized**: no coordination infrastructure, no accounts, no persistent services
- **Emergent mesh**: no formal mesh topology, meshes arise as a consequence of multilateral connections
- **Encrypted**: all data is encrypted end-to-end with a custom key exchange layer
- **Simple CLI**: the interface is minimal and obvious
- **Portable**: a single Go binary, no kernel modules, no third party dependencies

## Non-Goals

- **Not a VPN**: applications cannot route arbitrary traffic through Zorra
- **Not a mesh network**: there is no concept of a persistent named network or topology
- **Not a sync tool**: Zorra is for one-shot transfers, not continuous synchronization
- **Not production infrastructure**: Zorra is a personal tool, not designed for enterprise scale
- **No IPv4 support**: Zorra is IPv6 only

---

## Design

### Philosophy

**Decentralization + Simplicity = Elegance**

Zorra is built around a single primitive: a secure encrypted tunnel between exactly two peers. There is no session to join, no network to configure, and no server to register with. Two peers exchange addresses out of band, establish a direct connection, and communicate.

A mesh, if one ever emerges, is a consequence of multilateral connections rather than a first-class concept. This keeps the protocol simple and the trust model clear.

### Transport

Zorra uses **UDP over IPv6** exclusively.

IPv6 eliminates NAT as a concern for addressing and eliminates the need for extra infrastructure such as STUN and TURN servers.

UDP is used for the tunnel transport, giving full control over the reliability and ordering semantics at the application layer.

### Connectivity

Both peers are assumed to be behind home routers with stateful IPv6 firewalls. Connection is established via **UDP hole punching**:

1. Both peers learn each other's IPv6 address and chosen port out of band (e.g. pasted over Signal)
2. Both peers simultaneously send a UDP packet to each other's address
3. Each outbound packet opens a firewall rule for return traffic
4. The connection is established

This works reliably on typical consumer home routers. No STUN server, no relay, and no rendezvous infrastructure is required.

### Key Exchange (TBD)

Zorra implements a custom key exchange layer inspired by the [Noise Protocol Framework](https://noiseprotocol.org/). Peers perform an authenticated handshake before any data is transferred. Keys are ephemeral per session.

Initial trust is established via **TOFU (Trust On First Use)**; the first time you connect to a peer's key you accept and remember it. Key changes on subsequent connections produce a warning.

### Encryption (TBD)

All tunnel traffic is encrypted. The specific cipher suite is TBD during implementation, but will draw from the same primitives WireGuard uses: ChaCha20-Poly1305 for symmetric encryption, Curve25519 for key exchange.

### File Transfer Protocol

On top of the encrypted tunnel, Zorra implements a simple file transfer protocol:

- Files are chunked and streamed over the tunnel
- Transfer is resumable, so interrupted transfers can be continued
- Multiple files can be queued in a single session
- Progress is reported in the CLI

### Address Stability

IPv6 addresses on home networks are subject to change due to ISP prefix rotation and OS privacy extensions (RFC 4941). Zorra solves this automatically by registering a custom, deterministic IPv6 address for the duration of a session. Peers are expected to exchange connection details out of band, with a tentative solution for peer discovery in v2.

---

## Usage (TODO)

```bash
zorra
```

---

## Architecture

```
┌─────────────────────────────────┐
│           CLI (cobra)           │
├─────────────────────────────────┤
│       File Transfer Protocol    │
├─────────────────────────────────┤
│        Encrypted Tunnel         │
├─────────────────────────────────┤
│         Key Exchange            │
├─────────────────────────────────┤
│          UDP / IPv6             │
└─────────────────────────────────┘
```

---

## Comparison

| | Zorra | Magic Wormhole | Croc | Syncthing |
|---|---|---|---|---|
| Direct peer-to-peer | ✅ | ❌ relay | ❌ relay | ✅ |
| No server required | ✅ | ❌ | ❌ | ✅ |
| IPv6 native | ✅ | partial | partial | partial |
| Resumable transfers | ✅ | ❌ | ✅ | ✅ |
| One-shot transfer | ✅ | ✅ | ✅ | ❌ |
| No account required | ✅ | ✅ | ✅ | ✅ |

---

## Post-v1 Goals

- 🧭 **Peer Discovery**: challenging for a zero-infrastructure IPv6-based protocol
  - Proposed solution: [DISCOVERY](./docs/DISCOVERY.md)

---

## Prior Art & Inspiration

- [WireGuard](https://www.wireguard.com/): tunnel design and cryptographic primitives
- [Tailscale](https://tailscale.com/): NAT traversal and connectivity model
- [Magic Wormhole](https://magic-wormhole.readthedocs.io/): UX inspiration
- [Noise Protocol Framework](https://noiseprotocol.org/): key exchange design
- [Croc](https://github.com/schollz/croc): Go file transfer reference

