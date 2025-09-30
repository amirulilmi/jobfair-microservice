package utils

import (
	"regexp"
	"strings"
)

func GenerateSlug(text string) string {
	text = strings.ToLower(text)

	text = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == ' ' || r == '-' {
			return r
		}
		return -1
	}, text)

	text = strings.TrimSpace(text)

	reg := regexp.MustCompile(`\s+`)
	text = reg.ReplaceAllString(text, "-")

	reg = regexp.MustCompile(`-+`)
	text = reg.ReplaceAllString(text, "-")

	text = strings.Trim(text, "-")

	return text
}