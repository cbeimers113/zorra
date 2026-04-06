// Package address_test implements unit tests for the address package
package address_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cbeimers113/zorra/internal/address"
)

func Test_address_Ephemeral(t *testing.T) {
	addr, err := address.Ephemeral()
	assert.NoError(t, err)

	// Must get back a valid, globally addressable IPv6 address
	assert.NotNil(t, addr.To16())
	assert.True(t, addr.IsGlobalUnicast())
	assert.False(t, addr.IsPrivate())
}

