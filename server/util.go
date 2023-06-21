package main

import (
	"errors"
	"strings"
)

func remove_from_array[T interface{}](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

func verify_string(s string, check string) bool {
	for _, v := range s {
		if !strings.Contains(check, string(v)) {
			return false
		}
	}
	return true
}

func random_string(chars string, length int) (string, error) {

	if length <= 0 {
		return "", errors.New("length must be greater than zero")
	} else if len(chars) == 0 {
		return "", errors.New("chars length must be greater than zero")
	}

	var rs string

	for i := 0; i < length; i++ {
		rs += string(chars[random.Intn(len(chars))])
	}

	return rs, nil

}
