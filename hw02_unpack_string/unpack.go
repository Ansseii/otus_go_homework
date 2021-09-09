package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(sample string) (string, error) {
	var prev rune
	var builder strings.Builder

	for _, current := range sample {
		if count, err := strconv.Atoi(string(current)); err == nil {
			if prev == 0 || unicode.IsDigit(prev) {
				return "", ErrInvalidString
			}
			repeat := strings.Repeat(string(prev), count)
			trimLast(&builder)
			builder.WriteString(repeat)
		} else {
			builder.WriteRune(current)
		}
		prev = current
	}

	return builder.String(), nil
}

func trimLast(builder *strings.Builder) {
	if str := builder.String(); str != "" {
		builder.Reset()
		builder.WriteString(str[:len(str)-1])
	}
}
