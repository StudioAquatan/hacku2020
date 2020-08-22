package email

import "strings"

var oinoriWordsSubj = []string{
	"選考結果",
}
var oinoriWordsBody = []string{
	"残念",
	"貴殿の今後のご活躍を",
	"期待に添",
	"誠に申し訳ございません",
}

// Classify if it's a "oinori" email using its subject
func ClassifyMailBySubj(subj string) bool {
	for _, word := range oinoriWordsSubj {
		if strings.Contains(subj, word) {
			return true
		}
	}
	return false
}

// Classify if it's a "oinori" email using its body
func ClassifyMailByBody(body string) bool {
	for _, word := range oinoriWordsBody {
		if strings.Contains(body, word) {
			return true
		}
	}
	return false
}
