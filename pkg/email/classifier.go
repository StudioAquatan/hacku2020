package email

import (
	"context"
	"log"
	"os"
	"strings"

	language "cloud.google.com/go/language/apiv1"
	"github.com/golang/protobuf/proto"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

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

func ClassifyMailBySentiment(s string) bool {
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

func printResp(v proto.Message, err error) {
	if err != nil {
		log.Fatal(err)
	}
	proto.MarshalText(os.Stdout, v)
}
