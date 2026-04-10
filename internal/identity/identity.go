// Package identity implements logic for managing Zorra identities
package identity

import (
	"encoding/json"
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"regexp"

	"github.com/cbeimers113/zorra/internal/log"
)

const (
	// Where Zorra data is stored in the user home directory
	zorraDir = ".zorra"

	// File within the Zorra dir containing this peer's identity
	identityFile = "identity"

	// File within the Zorra dir containing this peer's known peers' public keys
	knownPeersFile = "known_peers"

	// The default identity to use when we can't determine the user's username
	defaultIdentity = "zorra-user"
)

var (
	// The absolute path to the Zorra directory
	zorraDirPath string

	// This peer's identity
	identity string = defaultIdentity

	// This peer's known peers' public keys
	knownPeers map[string]string
)

// init reads identity data when this package is first used at runtime
func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Warnf("Unable to determine user home directory; defaulting to CWD: %s", err.Error())
	}

	// Create Zorra dir if needed
	zorraDirPath = filepath.Join(home, zorraDir)
	if err = os.MkdirAll(zorraDirPath, 0o700); err != nil {
		log.Warnf("Unable to create Zorra home directory; defaulting to CWD: %s", err.Error())
	} else {
		os.Chmod(zorraDirPath, 0o700)
	}

	// Read or create identity
	if idBytes, err := os.ReadFile(filepath.Join(zorraDirPath, identityFile)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Debug("Creating new identity")
		} else {
			log.Warnf("Unable to read identity: %s", err.Error())
		}

		username := defaultIdentity
		if currentUser, err := user.Current(); err != nil {
			log.Warnf("Unable to determine username, falling back to %q: %s", defaultIdentity, err.Error())
		} else {
			username = currentUser.Username
		}

		// Save new identity to disk
		SetThis(username)
	} else {
		// Clean the read bytes and store the identity in memory
		re := regexp.MustCompile(`\s+`)
		identity = re.ReplaceAllString(string(idBytes), "")
	}

	// Initialize and read known peers
	knownPeers = make(map[string]string)
	if knownPeersBytes, err := os.ReadFile(filepath.Join(zorraDirPath, knownPeersFile)); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Warnf("Unable to read known peers: %s", err.Error())
	} else if err == nil {
		if err = json.Unmarshal(knownPeersBytes, &knownPeers); err != nil {
			log.Warnf("Unable to parse known peers: %s", err.Error())
		}
	}
}
