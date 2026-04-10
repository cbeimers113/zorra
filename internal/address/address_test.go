// Package address_test implements unit tests for the address package
package address_test

import (
	"errors"
	"net"
	"syscall"
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
		wantAddr   net.IP
		wantLink   netlink.Link
		wantErrMsg string
	}{
		"Happy path": {
			linkList: func() ([]netlink.Link, error) {
				return []netlink.Link{new(test.MockLink)}, nil
			},
			addrList: func(netlink.Link, int) ([]netlink.Addr, error) {
				return []netlink.Addr{test.MockAddr}, nil
			},
			wantAddr: test.EphemAddr,
			wantLink: new(test.MockLink),
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

			addr, link, err := address.Ephemeral()
			if tt.wantErrMsg == "" {
				assert.NoError(t, err)

				// Must get back a valid, globally addressable IPv6 address
				assert.Equal(t, tt.wantAddr, addr)
				assert.True(t, addr.IsGlobalUnicast())
				assert.False(t, addr.IsPrivate())
				assert.Equal(t, tt.wantLink, link)
			} else {
				assert.ErrorContains(t, err, tt.wantErrMsg)
			}
		})
	}
}

func Test_address_Register(t *testing.T) {
	link := new(test.MockLink)

	tests := map[string]struct {
		addErr  error
		wantErr bool
	}{
		"Happy path": {
			addErr:  nil,
			wantErr: false,
		},

		"Happy path - already registered": {
			addErr:  syscall.EEXIST,
			wantErr: false,
		},

		"Sad path": {
			addErr:  errors.New("mock error"),
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Backup and restore mocked AddrAdd
			addrAdd := address.AddrAdd
			defer func() {
				address.AddrAdd = addrAdd
			}()

			// Mock out AddrAdd
			address.AddrAdd = func(l netlink.Link, a *netlink.Addr) error {
				assert.Equal(t, link.Attrs().Name, l.Attrs().Name)
				assert.Equal(t, test.BaseAddr, a.IP)

				ones, bits := a.Mask.Size()
				assert.Equal(t, 128, ones)
				assert.Equal(t, 128, bits)
				return tt.addErr
			}

			err := address.Register(test.BaseAddr, link)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_address_Deregister(t *testing.T) {
	link := new(test.MockLink)

	tests := map[string]struct {
		delErr  error
		wantErr bool
	}{}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Backup and restore mocked AddrDel
			addrDel := address.AddrDel
			defer func() {
				address.AddrDel = addrDel
			}()

			// Mock out AddrDel
			address.AddrDel = func(l netlink.Link, a *netlink.Addr) error {
				assert.Equal(t, link.Attrs().Name, l.Attrs().Name)
				assert.Equal(t, test.BaseAddr, a.IP)

				ones, bits := a.Mask.Size()
				assert.Equal(t, 128, ones)
				assert.Equal(t, 128, bits)
				return tt.delErr
			}

			err := address.Deregister(test.BaseAddr, link)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
