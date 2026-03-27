# Zorra Peer Discovery

A decentralized, infrastructure-free peer discovery mechanism for Zorra, exploiting the structure of IPv6 global unicast addresses and deterministic interface ID derivation.

---

## Overview

Zorra's default connection model requires peers to exchange full IPv6 addresses out of band. This document describes an optional discovery mechanism that reduces the information a peer needs to share to a short **channel identifier**, enabling Zorra to reconstruct a small, scannable set of candidate addresses guaranteed to contain the real peer's address.

Zorra's fully-decentralized philosophy is respected; no servers, no DHT, and no infrastructure. Discovery is pure math and public data.

---

## Background: IPv6 Global Unicast Address Structure

A globally routable IPv6 address is 128 bits with the following structure:

```
|      32 bits      |   16 bits    |   16 bits  |       64 bits        |
|   ISP prefix /32  |   customer   |   subnet   |    interface ID      |
|   (public, RIR)   |   (opaque)   | (heuristic)|  (Zorra: derived)    |
```

- **ISP prefix (/32):** Publicly registered in RIR databases (ARIN, RIPE, APNIC, LACNIC, AFRINIC). One per ISP. Fixed and known.
- **Customer bits (16):** Assigned by the ISP to the customer. Opaque; not publicly registered at the customer level. Clustered sequentially or regionally in practice.
- **Subnet ID (16):** Assigned by the customer's router within their /48. Home routers almost universally use `0x0000` or `0x0001` for the primary subnet.
- **Interface ID (64):** Identifies the device within the subnet. In Zorra's model, this is deterministically derived from the peer's public key.

---

## The Channel

A **channel** is a short integer that maps to an ISP's /32 IPv6 prefix via a lookup table derived from public RIR data.

```
channel 42  →  2607:fea8::/32   (example ISP prefix)
```

Channels are:
- Small integers: easily shared verbally or in a short message
- Publicly derivable: the mapping is built from RIR WHOIS/routing data, no proprietary source
- Stable: ISP prefix allocations rarely change
- Regional: channel ranges can optionally be segmented by RIR region

Zorra can determine a peer's channel by looking up their ISP's prefix in the channel table.

---

## Structured Interface ID

For discovery to work, the interface ID must be deterministic and known to both peers ahead of time.

Zorra derives the interface ID from known public data, such as the peer's **public key**:

```
interface_id = SHA256(zorra_public_key)[0:8]  // first 64 bits of hash
```

This means:
- The interface ID is stable as long as the public key doesn't rotate
- Any peer who knows your public key can compute your interface ID
- The interface ID has no exploitable structure; it appears random, which is good for privacy

### Configuring the Zorra Address

Zorra combines the ISP-assigned /64 prefix (which includes the customer and subnet bits) with the derived interface ID to form a Zorra-specific IPv6 address, configured on the network interface at runtime:

```
[ISP /32] + [customer 16 bits] + [subnet 16 bits] + [derived interface ID 64 bits]
```

Configuring a custom address on a network interface requires elevated privileges (`sudo` or `CAP_NET_ADMIN` on Linux). This is consistent with other networking tools (WireGuard, Tailscale, OpenVPN).

Zorra removes the address from the interface on clean exit.

---

## Discovery Algorithm

Given a peer's **channel** and **public key**, Zorra reconstructs a candidate address set and scans it.

### Known at scan time

| Component | Source | Bits |
|---|---|---|
| ISP prefix | Channel lookup table | 32 |
| Customer bits | Unknown | 16 |
| Subnet ID | Heuristic | 16 |
| Interface ID | Derived from public key | 64 |

### Subnet heuristic

Home routers assign subnet IDs sequentially starting from `0x0000`. Zorra scans a small list of common subnet values:

```
0x0000, 0x0001, 0x0002, 0x0003
```

This covers the vast majority of home networks. The heuristic can be extended cheaply.

### Scan space

With the subnet heuristic, the unknown space is the 16-bit customer allocation:

```
2^16 = 65,536 candidate addresses per subnet value
× 4 subnet values = 262,144 total candidates
```

This is scannable in seconds on a modern connection.

### Confirmation

A candidate address is confirmed by attempting a Zorra key exchange handshake. The handshake only succeeds with the real peer; wrong addresses either don't respond or fail cryptographic verification.

### Algorithm

```
function discover(channel, peer_public_key):
    isp_prefix = channel_table[channel]          // known /32
    interface_id = SHA256(peer_public_key)[0:8]  // derived 64 bits

    for subnet in [0x0000, 0x0001, 0x0002, 0x0003]:
        for customer in range(0, 65536):
            candidate = isp_prefix
                      + customer           // 16 bits
                      + subnet             // 16 bits
                      + interface_id       // 64 bits
            async_probe(candidate)         // fire and forget

    wait for successful handshake response
    return confirmed_address
```

Probes are sent asynchronously in parallel, rate-limited to avoid triggering abuse detection.

---

## What a Peer Shares

Under this model, a peer shares only:

```
channel    // A short identifier                     
public_key // If known from a previous connection, or derived from a known password
```

The full IPv6 address is never shared. Address changes due to ISP prefix rotation are handled automatically as long as the peer's ISP (channel) and public key remain stable.

---

## Failure Modes

| Condition | Result |
|---|---|
| Peer's subnet ID not in heuristic list | Discovery fails, fall back to manual address sharing |
| ISP prefix not in channel table | No channel available, fall back to manual |
| Peer's public key has rotated | Interface ID mismatch, discovery fails |
| Peer offline | No handshake response, discovery times out |

---

## Privacy Considerations

- Scanning 262,144 addresses on an ISP's prefix may be visible to network monitoring. Scans should be rate-limited and spread over time.
- The derived interface ID is stable and linkable to a public key. Peers who rotate their keys will get a new interface ID and must reshare their discovery information.
- The channel reveals the peer's ISP to anyone who receives the discovery string.

---

## Future Directions

- **Regional channel namespacing:** Channel integers namespaced by RIR region to reduce collisions and improve scan locality
- **Subnet heuristic expansion:** Collect data on real-world subnet ID distributions to improve heuristic coverage
- **Key rotation protocol:** A mechanism for peers to announce key rotation without losing discoverability

---

## Relationship to v1

Discovery is a **v2+ feature**. Zorra v1 uses ASCII85-encoded full IPv6 addresses shared out of band. The structured interface ID (derived from public key) is backwards-compatible with v1, so discovery can be layered on top without a breaking change to the address format.

