package addressing

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/json"
	"io"
	"os"
)

var (
	//go:embed data/channel-map.gz
	channelMapRaw []byte

	// The in-memory channel map:
	channelMap map[string]int
)

// LoadChannelMap parses the embedded channel map data
func LoadChannelMap() error {
	channelMap = make(map[string]int)
	gzr, err := gzip.NewReader(bytes.NewBuffer(channelMapRaw))
	if err != nil {
		return err
	}

	data, err := io.ReadAll(gzr)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &channelMap); err != nil {
		return err
	}

	return nil
}

// ChannelOf returns the channel ID of an IPv6 prefix, if it exists in the map
func ChannelOf(prefix string) (int, bool) {
	ch, ok := channelMap[prefix]
	return ch, ok
}

// AddChannel adds a new channel to the channel map
func AddChannel(prefix string) {
	if _, ok := ChannelOf(prefix); ok {
		return
	}

	channelMap[prefix] = len(channelMap)
}

// SaveChannels writes the channel map to the disk
func SaveChannels() error {
	data, err := json.Marshal(channelMap)
	if err != nil {
		return err
	}

	file, err := os.OpenFile("channel-map.gz", os.O_CREATE|os.O_WRONLY, 0o0644)
	if err != nil {
		return err
	}
	defer file.Close()

	gzw := gzip.NewWriter(file)
	defer gzw.Close()

	_, err = gzw.Write(data)
	return err
}
