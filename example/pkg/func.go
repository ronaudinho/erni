package pkg

import (
	"errors"
)

var err = errors.New("pkg error")

func Fa(s string) (string, error) {
	ret, _ := Fb(s)
	return ret, err
}

func Fb(s string) (string, error) {
	ret, _ := Fc(s)
	return ret, err
}

func Fc(s string) (string, error) {
	return s, err
}
