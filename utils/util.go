package utils

import "strings"

func GetBase62String() string {
	alphabetLower := "abcdefghijklmnopqrstuvwxyz"
	alphabetUpper := strings.ToUpper(alphabetLower)
	numeric := "0123456789"
	base62 := alphabetLower + alphabetUpper + numeric

	return base62
}
