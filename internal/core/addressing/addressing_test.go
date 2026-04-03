// Package addressing_test implements unit tests for the addressing package
package addressing_test

import (
	"os"
	"testing"

	"github.com/cbeimers113/zorra/internal/core/addressing"
)

func TestMain(m *testing.M) {
	if err := addressing.LoadChannels(); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
