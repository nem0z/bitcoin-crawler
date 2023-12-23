package payload

import "crypto/rand"

func NewNonce() []byte {
	b := make([]byte, 8)
	rand.Read(b)
	return b
}
