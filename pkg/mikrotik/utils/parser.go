package utils

import "strconv"

func ParseInt(s string) int64 {
	if s == "" {
		return 0
	}
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func ParseBool(s string) bool {
	return s == "true" || s == "yes"
}

func FormatInt(n int64) string {
	return strconv.FormatInt(n, 10)
}
