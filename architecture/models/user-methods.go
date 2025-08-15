package models

import (
	"fmt"
	"net/mail"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func (u *User) ValidateNickname() error {
	if lng := len(u.Nickname); lng < 1 || 32 < lng {
		return fmt.Errorf("nickname: invalid lenght (%d)", lng)
	}
	for _, c := range u.Nickname {
		if !(unicode.IsLetter(c) || unicode.IsDigit(c)) {
			return fmt.Errorf("nickname: invalid character '%c'", c)
		}
	}
	return nil
}

func (u *User) ValidateEmail() error {
	if lng := len(u.Email); lng < 1 || 320 < lng {
		return fmt.Errorf("email: invalid lenght (%d)", lng)
	}
	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return err
	}
	return nil
}

// HashPassword - Creates hash and sets it for user password
func (u *User) HashPassword() error {
	pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}
	u.Password = string(pass)
	return nil
}

// CompareHashAndPassword - Compares user hashed password with field password
//
// returns: error if has error on encrypting field password;
// - true - if passowrds equal;
// - false - if passwords not equal;
func (u *User) CompareHashAndPassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	switch err {
	case bcrypt.ErrMismatchedHashAndPassword:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, fmt.Errorf("bcrypt.CompareHashAndPassword: %w", err)
	}
}
