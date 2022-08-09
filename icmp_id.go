package main

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixMilli())
}

// requires write lock
func getIcmpId() uint16 {
	return uint16(rand.Uint32())
}

func putIcmpId(uint16) {
	// no-op
}

// TODO: try harder to avoid conflicting id
