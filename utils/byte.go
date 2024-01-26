package utils

import (
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
