package main

func remove_from_array[T interface{}](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}
