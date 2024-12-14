package rss

import (
	"crypto/sha256"
	"encoding/hex"
)

func NewDeterministicGUID(values ...string) string {
	digest := sha256.New()
	for _, value := range values {
		digest.Write([]byte(value))
	}

	return hex.EncodeToString(digest.Sum(nil))[0:16]
}
