// Package address_test implements unit tests for the address package
package address_test

import (
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vishvananda/netlink"

	"github.com/cbeimers113/zorra/internal/address"
	"github.com/cbeimers113/zorra/internal/test"
)

func Test_address_Ephemeral(t *testing.T) {
	tests := map[string]struct {
		linkList   func() ([]netlink.Link, error)
		addrList   func(netlink.Link, int) ([]netlink.Addr, error)
		want       net.IP
		wantErrMsg string
	}{
		"Happy path": {
			linkList: func() ([]netlink.Link, error) {
				return []netlink.Link{new(test.MockLink)}, nil
			},
			addrList: func(netlink.Link, int) ([]netlink.Addr, error) {
				return []netlink.Addr{test.MockAddr}, nil
			},
			want: test.EphemAddr,
		},

		"Sad path - could not find network interfaces": {
			linkList: func() ([]netlink.Link, error) {
				return nil, errors.New("mock error")
			},
			wantErrMsg: "could not find network interfaces",
		},

		"Sad path - could not list addresses on interface": {
			linkList: func() ([]netlink.Link, error) {
				return []netlink.Link{new(test.MockLink)}, nil
			},
			addrList: func(netlink.Link, int) ([]netlink.Addr, error) {
				return nil, errors.New("mock error")
			},
			wantErrMsg: "could not list addresses on interface",
		},

		"Sad path - no globally addressable IPv6 addresses": {
			linkList: func() ([]netlink.Link, error) {
				return nil, nil
			},
			addrList: func(netlink.Link, int) ([]netlink.Addr, error) {
				return nil, nil
			},
			wantErrMsg: "no globally addressable IPv6 addresses available",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Backup and restore mocked functions
			linkList := address.LinkList
			addrList := address.AddrList
			defer func() {
				address.LinkList = linkList
				address.AddrList = addrList
			}()

			// Mock out the netlink functions
			address.LinkList = tt.linkList
			address.AddrList = tt.addrList

			addr, err := address.Ephemeral()
			if tt.wantErrMsg == "" {
				assert.NoError(t, err)

				// Must get back a valid, globally addressable IPv6 address
				assert.Equal(t, tt.want, addr)
				assert.True(t, addr.IsGlobalUnicast())
				assert.False(t, addr.IsPrivate())
			} else {
				assert.ErrorContains(t, err, tt.wantErrMsg)
			}
		})
	}
}
