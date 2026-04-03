package addressing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_addressing_hashString(t *testing.T) {
	tests := map[string]struct {
		str  string
		size uint8
		want []byte
	}{
		"Non-empty string, size in range": {
			str:  "hello, world",
			size: 3,
			want: []byte{0x3d, 0x63, 0xbe},
		},

		"Empty string, size in range": {
			str:  "",
			size: 6,
			want: []byte{0x25, 0x23, 0x22, 0x84, 0xe4, 0x9c},
		},

		"Non-empty string, size not in range": {
			str:  "zorra test",
			size: 30,
			want: []byte{0xf7, 0x0d, 0x80, 0xfa, 0xe4, 0x0f, 0x4f, 0xa1},
		},

		"Empty string, size not in range": {
			str:  "",
			size: 128,
			want: []byte{0x25, 0x23, 0x22, 0x84, 0xe4, 0x9c, 0xf2, 0xcb},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, hashString(tt.str, tt.size))
		})
	}
}
