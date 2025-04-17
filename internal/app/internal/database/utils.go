package database

import (
	"fmt"
)

func Abbr(s string, maxWidth int) string {
	n := (maxWidth / 2)
	return fmt.Sprintf("%s...%s", s[:n], s[len(s)-n:])
}

func FirstN(s string, n int) string {
	return fmt.Sprintf("%s...", s[:n])
}

func LastN(s string, n int) string {
	return fmt.Sprintf("...%s", s[len(s)-n:])
}
