// Package sharecode_test implements unit tests for the sharecode package
package sharecode_test

import (
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cbeimers113/zorra/internal/channel"
	"github.com/cbeimers113/zorra/internal/sharecode"
	"github.com/cbeimers113/zorra/internal/test"
)

func TestMain(m *testing.M) {
	if err := channel.LoadChannels(); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func Test_sharecode_Create(t *testing.T) {
	tests := map[string]struct {
		ephem      net.IP
		want       string
		wantErrMsg string
	}{
		"Happy path - channel exists": {
			ephem: test.EphemAddrHasChannel,
			want:  test.ShareCodeHasChannel,
		},

		"Happy path - no channel": {
			ephem: test.EphemAddrNoChannel,
			want:  test.ShareCodeNoChannel,
		},

		"Happy path - channel ends on byte boundary": {
			ephem: test.EphemAddrByteBoundary,
			want:  test.ShareCodeByteBoundary,
		},

		"Sad path - invalid IPv6 address": {
			ephem:      net.ParseIP("not a valid address"),
			wantErrMsg: "invalid ephemeral IPv6 address",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := sharecode.Create(tt.ephem)
			assert.Equal(t, tt.want, got)

			if tt.wantErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.wantErrMsg)
			}
		})
	}
}

func Test_sharecode_Read(t *testing.T) {
	tests := map[string]struct {
		shareCode  string
		want       net.IP
		wantErrMsg string
	}{
		"Happy path - channel exists": {
			shareCode: test.ShareCodeHasChannel,
			want:      test.EphemAddrHasChannel,
		},

		"Happy path - no channel": {
			shareCode: test.ShareCodeNoChannel,
			want:      test.EphemAddrNoChannel,
		},

		"Happy path - channel ends on byte boundary": {
			shareCode: test.ShareCodeByteBoundary,
			want:      test.EphemAddrByteBoundary,
		},

		"Sad path - invalid share code": {
			shareCode:  "a.b", // Not enough sections
			wantErrMsg: "invalid share code",
		},

		"Sad path - invalid channel": {
			shareCode:  "test.code.invalid", // Channel not hex
			wantErrMsg: "invalid channel",
		},

		"Sad path - unknown channel": {
			shareCode:  "test.code.ffffffff", // Channel too large
			wantErrMsg: "unknown channel",
		},

		"Sad path - could not decode opaque bits": {
			shareCode:  "test.!@#$.ff", // Opaque not base64
			wantErrMsg: "could not decode opaque bits",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := sharecode.Read(tt.shareCode)
			assert.Equal(t, tt.want, got)

			if tt.wantErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.wantErrMsg)
			}
		})
	}
}
