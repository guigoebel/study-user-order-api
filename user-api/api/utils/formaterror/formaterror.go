package formaterror

import (
	"errors"
	"strings"
)

// TODO - MELHORAR RESPOSTAS
func FormatError(err string) error {

	if strings.Contains(err, "name") {
		return errors.New("Name Already exists")
	}

	if strings.Contains(err, "email") {
		return errors.New("Email Already exists")
	}

	if strings.Contains(err, "description") {
		return errors.New("Description Already exists")
	}
	if strings.Contains(err, "hashedCpf") {
		return errors.New("Incorrect CPF")
	}
	return errors.New("Incorrect Details")
}
