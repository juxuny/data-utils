package email

import (
	"github.com/pkg/errors"
	"strings"
)

func IsValidEmail(s string) error {
	l := strings.Split(s, "@")
	if len(l) != 2 {
		return errors.Errorf("invalid email address: %s", s)
	}
	tail := strings.Split(l[1], ".")
	if len(tail) != 2 {
		return errors.Errorf("invalid eamil address: %s", s)
	}
	return nil
}
