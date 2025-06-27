package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

// Unpack - Распаковка строки.
func Unpack(str string) (string, error) {
	var (
		sb                    strings.Builder
		prev                  rune
		escape, prevIsEscaped bool
	)

	runeStr := []rune(str)

	for i, v := range runeStr {
		switch {
		case escape:
			sb.WriteRune(v)
			prev = v
			escape = false
			prevIsEscaped = true
			continue
		case unicode.IsDigit(v):
			if i == 0 || unicode.IsDigit(prev) && !prevIsEscaped {
				return "", ErrInvalidString
			} else if v == '0' && !unicode.IsDigit(prev) {
				removeLastRune(&sb)
				prev = v
				prevIsEscaped = false
				continue
			}
			count := int(v - '0')
			sb.WriteString(strings.Repeat(string(prev), count-1))
			prev = v
			prevIsEscaped = false
			continue
		case v == '\\':
			escape = true
			prev = v
			continue
		default:
			sb.WriteRune(v)
			prev = v
		}
	}

	finalStr := sb.String()
	return finalStr, nil
}

// removeLastRune - удаление последнего символа в строке.
func removeLastRune(sb *strings.Builder) {
	str := sb.String()
	runeStr := []rune(str)
	sb.Reset()
	sb.WriteString(string(runeStr[:len(runeStr)-1]))
}
