package main

import "strings"

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
