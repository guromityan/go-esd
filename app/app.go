package app

import (
	"fmt"
)

func Hello() {
	fmt.Println("hello")
}

func SetHeader(header []string, f func(x, y int, val interface{}) error) {
	for i, h := range header {
		f(i+1, 3, h)
	}
}
