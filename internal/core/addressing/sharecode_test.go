package addressing_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cbeimers113/zorra/internal/core/addressing"
)

func Test_addressing_CreateShareCode(t *testing.T) {
	tests := map[string]struct {
		ephem      net.IP
		want       string
		wantErrMsg string
	}{
		"Happy path - channel exists": {
			ephem: addressing.TestEphemAddrHasChannel,
			want:  addressing.TestShareCodeHasChannel,
		},

		"Happy path - no channel": {
			ephem: addressing.TestEphemAddrNoChannel,
			want:  addressing.TestShareCodeNoChannel,
		},

		"Happy path - channel ends on byte boundary": {
			ephem: addressing.TestEphemAddrByteBoundary,
			want:  addressing.TestShareCodeByteBoundary,
		},

		"Sad path - invalid IPv6 address": {
			ephem:      net.ParseIP("not a valid address"),
			wantErrMsg: "invalid ephemeral IPv6 address",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := addressing.CreateShareCode(tt.ephem)
			assert.Equal(t, tt.want, got)

			if tt.wantErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.wantErrMsg)
			}
		})
	}
}

func Test_addressing_ReadShareCode(t *testing.T) {
	tests := map[string]struct {
		shareCode  string
		want       net.IP
		wantErrMsg string
	}{
		"Happy path - channel exists": {
			shareCode: addressing.TestShareCodeHasChannel,
			want:      addressing.TestEphemAddrHasChannel,
		},

		"Happy path - no channel": {
			shareCode: addressing.TestShareCodeNoChannel,
			want:      addressing.TestEphemAddrNoChannel,
		},

		"Happy path - channel ends on byte boundary": {
			shareCode: addressing.TestShareCodeByteBoundary,
			want:      addressing.TestEphemAddrByteBoundary,
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
			got, err := addressing.ReadShareCode(tt.shareCode)
			assert.Equal(t, tt.want, got)

			if tt.wantErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.wantErrMsg)
			}
		})
	}
}
