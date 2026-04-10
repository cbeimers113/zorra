// Package cli implements the Zorra command line interface
package cli

import (
	"context"
	"os"

	"github.com/cbeimers113/zorra/internal/log"
)

// Execute calls the Zorra CLI with the user-provided arguments
func Execute(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil || log.ErrorStatus() {
		os.Exit(1)
	}
}

