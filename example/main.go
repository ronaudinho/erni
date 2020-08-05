package main

import (
	"errors"
	"fmt"

	"github.com/ronaudinho/erni/example/pkg"

	_ "github.com/ronaudinho/erni"
)

var err = errors.New("error")

func main() {
	a, _ := fa("a")
	fmt.Println(a)

	ia, _ := pkg.Fa("a")
	fmt.Println(ia)
}

func fa(s string) (string, error) {
	ret, _ := fb(s)
	return ret, err
}

func fb(s string) (string, error) {
	ret, _ := fc(s)
	return ret, err
}

func fc(s string) (string, error) {
	return s, err
}
