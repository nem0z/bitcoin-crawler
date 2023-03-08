package message

import "errors"

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
