package utils

import (
	"encoding/json"
	"os"

	"github.com/nem0z/bitcoin-crawler/peer"
)

func LoadAddr(path string) (addrs []*peer.Addr, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&addrs)
	if err != nil {
		return
	}
	return addrs, err
}
