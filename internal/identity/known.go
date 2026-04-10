package identity

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cbeimers113/zorra/internal/log"
)

// IsKnown returns whether we know the public key of a peer
func IsKnown(peer string) bool {
	_, ok := knownPeers[peer]
	return ok
}

// PublicKey returns the public key of a known peer
func PublicKey(peer string) (string, error) {
	publicKey, ok := knownPeers[peer]
	if !ok {
		return "", fmt.Errorf("unknown peer: %q", peer)
	}

	return publicKey, nil
}

// Remember adds a peer's share code and public key to the known peers map
func Remember(peer, publicKey string) {
	knownPeers[peer] = publicKey
	log.Infof("Remembering peer %q", peer)

	data, err := json.Marshal(knownPeers)
	if err != nil {
		log.Warnf("Unable to marshal known peers map: %s", err.Error())
		return
	}

	if err = os.WriteFile(filepath.Join(zorraDirPath, knownPeersFile), data, 0o644); err != nil {
		log.Warnf("Unable to write known peers map to disk: %s", err.Error())
	}
}
