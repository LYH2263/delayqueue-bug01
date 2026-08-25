package wal

import (
	"crypto/sha256"
	"encoding/hex"
)

func Checksum(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:8])
}

func Merge(a, b []byte) []byte {
	out := make([]byte, 0, len(a)+len(b))
	out = append(out, a...)
	out = append(out, b...)
	return out
}
