package message

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
)

func formatCommandName(commandName string) ([]byte, error) {
	const maxLenght = 12

	if len(commandName) > maxLenght {
		return nil, errors.New("Command name can't exceed 12 bytes")
	}

	commandNameByte := []byte(commandName)

	if len(commandNameByte) == maxLenght {
		return commandNameByte, nil
	}

	bytesToAdd := make([]byte, 12-len(commandName))
	formatedCommandName := append(commandNameByte, bytesToAdd...)

	return formatedCommandName, nil
}

func RemoveTrailingZeros(data []byte) []byte {
	for i := len(data) - 1; i >= 0; i-- {
		if data[i] != 0 {
			return data[:i+1]
		}
	}
	return []byte{}
}

func ResolveCommandName(command []byte) string {
	return string(RemoveTrailingZeros(command))
}

func ReadVarInt(data []byte) (int64, []byte, error) {
	var value int64
	var shift uint

	for i, b := range data {
		value |= (int64(b&0x7F) << shift)

		if b&0x80 == 0 {
			return value, data[i+1:], nil
		}

		shift += 7

		if shift > 63 {
			return 0, nil, errors.New("ReadVarIntFromSlice: Varint is too long")
		}
	}

	return 0, nil, errors.New("ReadVarIntFromSlice: Incomplete varint in the slice")
}

func CreateNonce(size int) ([]byte, error) {
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	return bytes, err
}

func Checksum(data []byte) []byte {
	hash := sha256.Sum256(data)
	hash = sha256.Sum256(hash[:])
	return hash[:4]
}
