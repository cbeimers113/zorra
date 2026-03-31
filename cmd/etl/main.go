package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/cbeimers113/zorra/internal/core/addressing"
	"github.com/cbeimers113/zorra/internal/log"
)

// ETL CI process to update the channel map
// data from the public RIR delegation files.

// RIR delegation file locations
var sources = map[addressing.RIR]string{
	addressing.ARIN:    "https://ftp.arin.net/pub/stats/arin/delegated-arin-extended-latest",
	addressing.RIPE:    "https://ftp.ripe.net/ripe/stats/delegated-ripencc-extended-latest",
	addressing.APNIC:   "https://ftp.apnic.net/stats/apnic/delegated-apnic-extended-latest",
	addressing.LACNIC:  "https://ftp.lacnic.net/pub/stats/lacnic/delegated-lacnic-extended-latest",
	addressing.AFRINIC: "https://ftp.afrinic.net/pub/stats/afrinic/delegated-afrinic-extended-latest",
}

// readSource reads an RIR delegation file and returns a slice of IPv6 prefixes
func readSource(rir addressing.RIR, source string) ([]string, error) {
	var prefixes []string
	log.Infof("Reading RIR delegation file of %s...", addressing.RIRName(rir))

	rsp, err := http.Get(source)
	if err != nil {
		return nil, fmt.Errorf("could not request file over HTTP: %w", err)
	}
	defer rsp.Body.Close()

	// Read the response into memory
	data, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read HTTP response body: %w", err)
	}

	// Parse file
	for line := range strings.SplitSeq(string(data), "\n") {
		cols := strings.Split(line, "|")
		if len(cols) < 7 || cols[2] != "ipv6" || (cols[6] != "allocated" && cols[6] != "assigned") {
			continue
		}

		// Skip prefixes already in the map
		prefix := cols[3]
		if _, ok := addressing.ChannelOf(rir, prefix); ok {
			continue
		}

		prefixes = append(prefixes, prefix)
	}

	return prefixes, nil
}

// updateChannels updates the channel map with any new prefixes from a delegation file
func updateChannels(rir addressing.RIR, prefixes []string) {
	log.Infof("Creating %d new channels for %s", len(prefixes), addressing.RIRName(rir))
	for _, prefix := range prefixes {
		addressing.AddChannel(rir, prefix)
	}
}

func main() {
	log.Info("Updating channel map from RIR delegation files")
	if err := addressing.LoadChannelMap(); err != nil {
		log.Warnf("Unable to read existing channel map: %s, rebuilding...", err.Error())
	}

	// Update each RIR's channels from their delegation file
	fmt.Println()
	for rir, source := range sources {
		prefixes, err := readSource(rir, source)
		if err != nil {
			log.Errorf("Unable to read RIR delegaiton file of %s: %s", addressing.RIRName(rir), err.Error())
			continue
		}

		if len(prefixes) == 0 {
			log.Infof("No new prefixes for %s", addressing.RIRName(rir))
			continue
		}

		updateChannels(rir, prefixes)
		log.Infof("Done updating %s\n", addressing.RIRName(rir))
	}

	// Save the channel map to disk
	if err := addressing.SaveChannels(); err != nil {
		log.Errorf("Unable to save channel map file: %s", err.Error())
		os.Exit(1)
	}

	log.Info("Channel map update complete")
}
