package email

import (
	"context"
	"log"
	"strings"

	language "cloud.google.com/go/language/apiv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

var screeningMailSubj = []string{
	"選考結果",
}
var oinoriWordsBody = []string{
	"残念",
	"貴殿の今後のご活躍を",
	"期待に添",
	"誠に申し訳ございません",
}

var acceptanceWordsBody = []string{
	"是非参加",
	"ぜひ参加",
	"ご参加いただきたく",
	"お会いできますこと",
	"心よりお待ちしております",
}

// Classify if it's a "oinori" email using its subject
func ClassifyScreeningMailBySubj(subj string) bool {
	for _, word := range screeningMailSubj {
		if strings.Contains(subj, word) {
			return true
		}
	}
	return false
}

// Classify if it's a "oinori" email using its body
func ClassifyOinoriMailByBody(body string) bool {
	for _, word := range oinoriWordsBody {
		if strings.Contains(body, word) {
			return true
		}
	}
	return false
}

func ClassifyOinoriMailBySentiment(s string) bool {
	ctx := context.Background()
	client, err := language.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := analyzeSentiment(ctx, client, s)
	if resp != nil && resp.DocumentSentiment.Score < 0 {
		return true
	} else {
		return false
	}
}

func ClassifyAcceptanceMailByBody(body string) bool {
	for _, word := range acceptanceWordsBody {
		if strings.Contains(body, word) {
			return true
		}
	}
	return false
}

func analyzeSentiment(ctx context.Context, client *language.Client, text string) (*languagepb.AnalyzeSentimentResponse, error) {
	return client.AnalyzeSentiment(ctx, &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
	})
}
