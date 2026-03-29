// Package cli implements the Zorra command line interface
package cli

import (
	"context"
	"os"
)

// Execute calls the Zorra CLI with the user-provided arguments
func Execute(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
