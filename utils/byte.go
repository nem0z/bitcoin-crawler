package utils

import (
	"crypto/rand"
	"errors"
)

func ToByteFixedSize(data []byte, size int) error {
	if len(data) > size {
		return errors.New("data lenght is greater than expected length")
	}

	if len(data) == size {
		return nil
	}

	bytesToAdd := make([]byte, size-len(data))
	data = append(data, bytesToAdd...)

	return nil
}

func RemoveTrailingZeros(data []byte) []byte {
	for i := len(data) - 1; i >= 0; i-- {
		if data[i] != 0 {
			return data[:i+1]
		}
	}
	return []byte{}
}

func CreateNonce(size int) ([]byte, error) {
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	return bytes, err
}
