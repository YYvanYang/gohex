package vo

import (
	"regexp"
	"strings"
)

type Email struct {
	address string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func NewEmail(address string) (Email, error) {
	address = strings.TrimSpace(address)
	if address == "" {
		return Email{}, errors.ErrEmptyEmail
	}

	if !emailRegex.MatchString(address) {
		return Email{}, errors.ErrInvalidEmail
	}

	return Email{address: strings.ToLower(address)}, nil
}

func (e Email) String() string {
	return e.address
}

func (e Email) IsEmpty() bool {
	return e.address == ""
} 