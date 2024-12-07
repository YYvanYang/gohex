package vo

import (
	"golang.org/x/crypto/bcrypt"
	"unicode"
)

type Password struct {
	hash string
}

func NewPassword(plaintext string) (Password, error) {
	if err := validatePassword(plaintext); err != nil {
		return Password{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, err
	}

	return Password{hash: string(hash)}, nil
}

func NewPasswordFromHash(hash string) Password {
	return Password{hash: hash}
}

func (p Password) Compare(plaintext string) error {
	err := bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plaintext))
	if err != nil {
		return ErrInvalidPassword
	}
	return nil
}

func (p Password) Hash() string {
	return p.hash
}

func (p Password) IsEmpty() bool {
	return p.hash == ""
}

// validatePassword 验证密码复杂度
func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return ErrPasswordTooWeak
	}

	return nil
} 