package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Adeithe/go-twitch/api"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type Upcoming struct {
	Dates       []string `json:"dates"`
	MaybeOnline bool     `json:"maybe_online"` // geaendert fon vilt_online
}

func classify(transcription string, video api.Video) (Upcoming, error) {
	client := openai.NewClient(os.Getenv("OPENAI_CHATGPT_TOKEN"))

	functionDefinitions := openai.FunctionDefinition{
		Name: "has_online_intend",
		Description: "Evaluate if the streamer will stream again. Set maybe_online=true if uncertain or no clear data.",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"dates": {
					Type: jsonschema.Array,
					Items: &jsonschema.Definition{
						Type: jsonschema.String,
					},
					Description: "Array of planned stream dates in RFC3339 format",
				},
				"maybe_online": { // geaendert von vilt_online
					Type: jsonschema.Boolean,
					Description: "True if streamer is uncertain about streaming plans",
				},
			},
		},
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     openai.GPT4oMini,
			Functions: []openai.FunctionDefinition{functionDefinitions},
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `You analyze streamer transcripts to predict future streams. Evaluate:
1. Return dates[] if streamer announced specific plans
2. Set maybe_online=true if:
   - Streamer is uncertain ("maybe", "not sure")
   - Mentions potential streams without dates
   - No clear streaming intentions
3. Return empty dates[] only when clearly stating no streams planned
Consider context, tone and streamer's usual patterns.`,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf(
						"Stream from %s. Analyze for streaming plans:\n%s",
						video.PublishedAt.Format(time.RFC3339),
						transcription,
					),
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return Upcoming{}, err
	}

	if len(resp.Choices) == 0 {
		fmt.Println("No choices found")
		return Upcoming{}, nil
	}

	// debug output
	if f, err := os.Create("chatgpt.json"); err == nil {
		defer f.Close()
		json.NewEncoder(f).Encode(resp)
	}

	var data Upcoming
	if err := json.Unmarshal(
		[]byte(resp.Choices[0].Message.FunctionCall.Arguments),
		&data,
	); err != nil {
		fmt.Printf("Unmarshal error: %v\nRaw: %s\n", 
			err,
			resp.Choices[0].Message.FunctionCall.Arguments,
		)
		return Upcoming{}, err
	}

	return data, nil
}