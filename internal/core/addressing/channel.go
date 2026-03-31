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
	// The map is itself an array of maps where at each index
	// is a mapping of IPv6 prefixes from the RIR delegation file
	// of the RIR at that index to their Zorra channel ID
	channelMap []map[string]int
)

func init() {
	// Initialize the map to handle missing or corrupt data file
	channelMap = make([]map[string]int, 5)
	channelMap[APNIC] = make(map[string]int)
	channelMap[AFRINIC] = make(map[string]int)
	channelMap[ARIN] = make(map[string]int)
	channelMap[LACNIC] = make(map[string]int)
	channelMap[RIPE] = make(map[string]int)
}

// LoadChannelMap parses the embedded channel map data
func LoadChannelMap() error {
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
func ChannelOf(rir RIR, prefix string) (int, bool) {
	if rir == UnknownRIR {
		return 0, false
	}

	ch, ok := channelMap[rir][prefix]
	return ch, ok
}

// AddChannel adds a new channel to the channel map
func AddChannel(rir RIR, prefix string) {
	if _, ok := ChannelOf(rir, prefix); ok {
		return
	}

	channelMap[rir][prefix] = len(channelMap[rir])
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
