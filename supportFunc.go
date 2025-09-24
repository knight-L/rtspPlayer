package main

import (
	"crypto/rand"
	"fmt"
)

// Default streams signals
const (
	SignalStreamRestart = iota ///< Y   Restart
	SignalStreamStop
	SignalStreamClient
)

// generateUUID function make random uuid for clients and stream
func generateUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
