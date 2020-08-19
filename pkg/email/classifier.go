package email

import "strings"

var oinoriWords = []string{
	"残念ながら",
	"貴殿の今後のご活躍を",
	"期待に添えない",
	"慎重に選考",
}

// Classify if it's a "oinori" email
func ClassifyMail(body string) bool {
	for _, word := range oinoriWords {
		if strings.Contains(body, word) {
			return true
		}
	}
	return false
}
