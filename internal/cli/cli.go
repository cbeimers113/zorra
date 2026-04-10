// Package cli implements the Zorra command line interface
package cli

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/term"

	"github.com/cbeimers113/zorra/internal/hash"
	"github.com/cbeimers113/zorra/internal/log"
)

const minPasswordLength = 10

// Execute calls the Zorra CLI with the user-provided arguments
func Execute(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil || log.ErrorStatus() {
		os.Exit(1)
	}
}

// getPasswordHash gets a password from the user on the command line
// and returns the argon2id hash of it
func getPasswordHash() ([]byte, error) {
	pw, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return nil, fmt.Errorf("could not read password: %w", err)
	}

	// Enfore a minimum password length
	if len(pw) < minPasswordLength {
		return nil, fmt.Errorf("password must be at least %d bytes long", minPasswordLength)
	}

	return hash.Password(pw), nil
}
