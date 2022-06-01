package models

import (
	"crypto/rand"
	"math/big"
)

func randomBase62(bytes int) string {
	b := make([]byte, bytes)
	_, err := rand.Read(b)
	if err != nil {
		panic("unable to read pseudorandom numbers")
	}
	n := new(big.Int)
	n.SetBytes(b)
	return n.Text(62)
}

// GenerateTag sets Tag field to a random base62 (base64 without / or +) with a length
// of 22 characters. Values are cryptographically pseudo-random and pseudo-unique with
// probability of getting same values the same as for random UUID.
func GenerateTag() string {
	// generate from more bytes because fixed-size base62 cannot be guaranteed at this size
	return randomBase62(20)[0:22]
}
