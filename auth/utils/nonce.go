package utils

import "math/rand"

func GenToken(length int) string {
	chars := "0123456789"
	result := ""
	for range length {
		result += string(chars[rand.Intn(len(chars))])
	}
	return result
}
