package identity

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cbeimers113/zorra/internal/log"
	"github.com/cbeimers113/zorra/internal/test"
)

// This returns this peer's identity
func This() string {
	if testing.Testing() {
		return test.Identity
	}

	return identity
}

// SetThis sets this peer's identity
func SetThis(id string) {
	identity = id
	log.Infof("Set identity to %q", id)

	if err := os.WriteFile(filepath.Join(zorraDirPath, identityFile), []byte(id), 0o644); err != nil {
		log.Warnf("Unable to save identity to disk: %s", err.Error())
	}
}
