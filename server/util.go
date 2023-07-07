package main

import (
	"errors"
	"strconv"
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

func record_exists[K int64 | string](table string, id K) bool {
	var record_id int64
	var err error

	switch any(id).(type) {
	case string:
		record_id, err = strconv.ParseInt(string(id), 16, 64)
		if err != nil {
			return false
		}
	case int64:
		record_id = any(id).(int64)
	}

	var count int64
	database.Table(table).Where("id = ?", record_id).Count(&count)

	if count == 0 {
		return false
	} else {
		return true
	}

}
