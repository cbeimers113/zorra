// Package state implements logic to read and write Zorra state data on the disk
package state

import (
	"encoding/json"
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/cbeimers113/zorra/internal/log"
	"github.com/cbeimers113/zorra/internal/testdata"
)

const (
	// Where Zorra state data is stored in the user home directory
	zorraDir = ".zorra"

	// File within the Zorra dir containing this peer's identity
	identityFile = "identity"

	// File within the Zorra dir containing this peer's known hosts
	knownHostsFile = "known_hosts"

	// The default identity to use when we can't determine the user's username
	defaultIdentity = "zorra-user"
)

var (
	// The absolute path to the Zorra directory
	zorraDirPath string

	// This peer's identity
	identity string

	// This peer's known hosts
	knownHosts map[string]string
)

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

		SetIdentity(username)
	} else {
		re := regexp.MustCompile(`\s+`)
		identity = re.ReplaceAllString(string(idBytes), "")
	}

	// Initialize and read known hosts
	knownHosts = make(map[string]string)
	if knownHostBytes, err := os.ReadFile(filepath.Join(zorraDirPath, knownHostsFile)); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Warnf("Unable to read known hosts: %s", err.Error())
	} else if err == nil {
		if err = json.Unmarshal(knownHostBytes, &knownHosts); err != nil {
			log.Warnf("Unable to parse known hosts: %s", err.Error())
		}
	}
}

// Identity returns this peer's identity
func Identity() string {
	if testing.Testing() {
		return testdata.Identity
	}

	return identity
}

// SetIdentity sets this peer's identity
func SetIdentity(id string) {
	identity = id
	log.Infof("Set identity to %q", id)

	if err := os.WriteFile(filepath.Join(zorraDirPath, identityFile), []byte(id), 0o644); err != nil {
		log.Warnf("Unable to save identity to disk: %s", err.Error())
	}
}
