package email

import (
	"errors"
	"net/mail"
	"regexp"
	"strings"
)

var (
	ErrEmailEmpty   = errors.New("email não pode ser vazio")
	ErrEmailInvalid = errors.New("email inválido")
	ErrEmailTooLong = errors.New("email muito longo")
)

// Regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$`)

const maxEmailLength = 254

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return ErrEmailEmpty
	}
	if len(email) > maxEmailLength {
		return ErrEmailTooLong
	}

	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return ErrEmailInvalid
	}

	if !emailRegex.MatchString(addr.Address) {
		return ErrEmailInvalid
	}

	return nil
}
