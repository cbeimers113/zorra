// Package channel_test implements unit tests for the channel package
package channel_test

import (
	"net/netip"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cbeimers113/zorra/internal/channel"
)

func TestMain(m *testing.M) {
	if err := channel.LoadChannels(); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func Test_channel_ChannelOf(t *testing.T) {
	tests := map[string]struct {
		prefix      netip.Prefix
		wantChannel int
		wantOk      bool
	}{
		"Channel exists": {
			prefix:      netip.MustParsePrefix("2605:59c0::/28"),
			wantChannel: 5359,
			wantOk:      true,
		},

		"Channel doesn't exist": {
			prefix: netip.MustParsePrefix("dead:beef::/12"),
			wantOk: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, ok := channel.ChannelOf(tt.prefix)
			assert.Equal(t, tt.wantChannel, got)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}

func Test_channel_PrefixOf(t *testing.T) {
	tests := map[string]struct {
		channel    int
		wantPrefix netip.Prefix
		wantOk     bool
	}{
		"Prefix exists": {
			channel:    5359,
			wantPrefix: netip.MustParsePrefix("2605:59c0::/28"),
			wantOk:     true,
		},

		"Prefix doesn't exist": {
			channel: -1,
			wantOk:  false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, ok := channel.PrefixOf(tt.channel)
			assert.Equal(t, tt.wantPrefix, got)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}

func Test_channel_AddChannel(t *testing.T) {
	tests := map[string]struct {
		prefix  netip.Prefix
		wantAdd bool
	}{
		"New channel": {
			prefix:  netip.MustParsePrefix("dead:beef::/14"),
			wantAdd: true,
		},

		"Channel already exists": {
			prefix:  netip.MustParsePrefix("2605:59c0::/28"),
			wantAdd: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			current, ok := channel.ChannelOf(tt.prefix)
			assert.Equal(t, tt.wantAdd, !ok)

			channel.AddChannel(tt.prefix)
			added, ok := channel.ChannelOf(tt.prefix)
			assert.True(t, ok)
			assert.Equal(t, tt.wantAdd, added != current)
		})
	}
}
