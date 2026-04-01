package addressing

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/json"
	"io"
	"net/netip"
	"os"

	"github.com/cbeimers113/zorra/internal/log"
)

var (
	//go:embed data/channel-map.gz
	channelMapRaw []byte

	// The in-memory channel map: map IPv6 prefixes to channel IDs
	channelMap map[netip.Prefix]int

	// The in-memory inverse channel map: lookup IPv6 prefixes by channel
	prefixMap map[int]netip.Prefix
)

// LoadChannels parses the embedded channel map data
func LoadChannels() error {
	// Initialize the channel and prefix maps, and an intermediate mapping of CIDR prefixes to channels
	channelMap = make(map[netip.Prefix]int)
	prefixMap = make(map[int]netip.Prefix)
	cidrMap := make(map[string]int)

	// Read the compressed channel map payload
	gzr, err := gzip.NewReader(bytes.NewBuffer(channelMapRaw))
	if err != nil {
		return err
	}

	// Extract the gz bytes
	data, err := io.ReadAll(gzr)
	if err != nil {
		return err
	}

	// Unmarshal into the intermediate CIDR map
	if err := json.Unmarshal(data, &cidrMap); err != nil {
		return err
	}

	// Parse CIDR prefixes and populate the channel and prefix maps
	for cidr, channel := range cidrMap {
		prefix, err := netip.ParsePrefix(cidr)
		if err != nil {
			log.Warnf("Invalid CIDR-notation IPv6 prefix %q in channel map: %s", cidr, err.Error())
			continue
		}

	  channelMap[prefix.Masked()] = channel
		prefixMap[channel] = prefix
	}

	return nil
}

// ChannelOf returns the channel ID of an IPv6 prefix, if it exists in the map
func ChannelOf(prefix netip.Prefix) (int, bool) {
	ch, ok := channelMap[prefix]
	return ch, ok
}

// PrefixOf returns the IPv6 prefix of a channel ID, if it exists in the map
func PrefixOf(channel int) (netip.Prefix, bool) {
	prefix, ok := prefixMap[channel]
	return prefix, ok
}

// AddChannel adds a new channel to the channel map
func AddChannel(prefix netip.Prefix) {
	if _, ok := ChannelOf(prefix); ok {
		return
	}

	channelMap[prefix] = len(channelMap)
}

// SaveChannels writes the channel map to the disk
func SaveChannels() error {
	// Convert the channel map into an intermediate mapping of CIDR prefixes to channels
	cidrMap := make(map[string]int)
	for prefix, channel := range channelMap {
		cidrMap[prefix.String()] = channel
	}

	// Marshal the CIDR map
	data, err := json.Marshal(cidrMap)
	if err != nil {
		return err
	}

	// Create the output file
	file, err := os.OpenFile("channel-map.gz", os.O_CREATE|os.O_WRONLY, 0o0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Compress the map as gz and write to disk
	gzw := gzip.NewWriter(file)
	defer gzw.Close()

	_, err = gzw.Write(data)
	return err
}
