package openai

import (
	"context"
	"log"

	copenai "github.com/sashabaranov/go-openai"
)

func CreateImageClient(q string) string {
	var imageOpenAIClient *copenai.Client
	resp, err := imageOpenAIClient.CreateImage(
		context.Background(),
		copenai.ImageRequest{
			Prompt: q,
			N:      1,
			Size:   "512x512",
		},
	)
	if err != nil {
		log.Printf("openAIClient.CreateImage err=%+v\n", err)
		return ""
	}
	return resp.Data[0].URL
}
