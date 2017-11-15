package utils

import (
	"strconv"
	"strings"
)

func StringToInt(s string) (int, error) {
	return strconv.Atoi(strings.Join(strings.Split(s, ","), ""))
}
