package account

import "strings"

// Email is a value object for account email.
type Email string

// ParseEmail validates and returns Email value object.
func ParseEmail(raw string) (Email, error) {
	e := Email(strings.TrimSpace(raw))
	if !strings.Contains(string(e), "@") {
		return "", ErrInvalidEmail
	}
	return e, nil
}

func (e Email) String() string {
	return string(e)
}
