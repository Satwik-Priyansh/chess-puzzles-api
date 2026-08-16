package validation

import (
	"fmt"
	"net/mail"
	"regexp"
)

var RegUsername = regexp.MustCompile("^[a-zA-Z0-9_]+$")
var RegUCIMoves = regexp.MustCompile("^[a-h][1-8][a-h][1-8][qrbn]?$")

func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return err
	}
	return nil
}
func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 30 {
		return fmt.Errorf("username must be between 3 and 30 characters")
	}
	if !RegUsername.MatchString(username) {
		return fmt.Errorf("username can only contain letters, numbers and underscores")
	}
	return nil
}
func ValidatePassword(password string) error {
	if len(password) >= 8 && len(password) <= 72 {
		return nil
	}
	return fmt.Errorf("password must be between 8 and 72 characters")
}
func ValidateUCIMoves(moves []string) error {
	if len(moves) == 0 {
		return fmt.Errorf("no moves to evaluate")
	}
	for _, move := range moves {
		ok := RegUCIMoves.MatchString(move)
		if !ok {
			return fmt.Errorf("invalid UCI move format: %s", move)
		}
	}
	return nil
}
