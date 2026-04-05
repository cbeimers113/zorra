// Package sharecode_test implements unit tests for the sharecode package
package sharecode_test

import (
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cbeimers113/zorra/internal/channel"
	"github.com/cbeimers113/zorra/internal/sharecode"
	"github.com/cbeimers113/zorra/internal/testdata"
)

func TestMain(m *testing.M) {
	if err := channel.LoadChannels(); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func Test_addressing_CreateShareCode(t *testing.T) {
	tests := map[string]struct {
		ephem      net.IP
		want       string
		wantErrMsg string
	}{
		"Happy path - channel exists": {
			ephem: testdata.EphemAddrHasChannel,
			want:  testdata.ShareCodeHasChannel,
		},

		"Happy path - no channel": {
			ephem: testdata.EphemAddrNoChannel,
			want:  testdata.ShareCodeNoChannel,
		},

		"Happy path - channel ends on byte boundary": {
			ephem: testdata.EphemAddrByteBoundary,
			want:  testdata.ShareCodeByteBoundary,
		},

		"Sad path - invalid IPv6 address": {
			ephem:      net.ParseIP("not a valid address"),
			wantErrMsg: "invalid ephemeral IPv6 address",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := sharecode.CreateShareCode(tt.ephem)
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
			shareCode: testdata.ShareCodeHasChannel,
			want:      testdata.EphemAddrHasChannel,
		},

		"Happy path - no channel": {
			shareCode: testdata.ShareCodeNoChannel,
			want:      testdata.EphemAddrNoChannel,
		},

		"Happy path - channel ends on byte boundary": {
			shareCode: testdata.ShareCodeByteBoundary,
			want:      testdata.EphemAddrByteBoundary,
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
			got, err := sharecode.ReadShareCode(tt.shareCode)
			assert.Equal(t, tt.want, got)

			if tt.wantErrMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.wantErrMsg)
			}
		})
	}
}
