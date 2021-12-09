package torrent

import (
	"crypto/rand"
	"os"
)

var emptyBytes [20]byte

func GetOrCreatePeerID(p string) ([20]byte, error) {
	idb, err := os.ReadFile(p)
	if err == nil {
		var out [20]byte
		copy(out[:], idb)

		return out, nil
	}

	if !os.IsNotExist(err) {
		return emptyBytes, err
	}

	var out [20]byte
	_, err = rand.Read(out[:])
	if err != nil {
		return emptyBytes, err
	}

	return out, os.WriteFile(p, out[:], 0755)
}
