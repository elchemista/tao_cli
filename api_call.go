package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Prompt struct {
	Model            string   `json:"model"`
	Prompt           string   `json:"prompt"`
	Temperature      float64  `json:"temperature"`
	MaxTokens        int      `json:"max_tokens"`
	TopP             float64  `json:"top_p"`
	FrequencyPenalty float64  `json:"frequency_penalty"`
	PresencePenalty  float64  `json:"presence_penalty"`
	Stop             []string `json:"stop"`
}

type CompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string `json:"text"`
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

func getCommandFromAPI(prompt string, key string) string {
	url := "https://api.openai.com/v1/completions"

	promptData := Prompt{
		Model:            "text-davinci-003",
		Prompt:           fmt.Sprintf("Translate the following English instruction into a Unix terminal command: '%s'", prompt),
		Temperature:      1.0,
		MaxTokens:        256,
		TopP:             1.0,
		FrequencyPenalty: 0.0,
		PresencePenalty:  0.0,
		Stop:             []string{"\\n"},
	}

	jsonData, err := json.Marshal(promptData)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var response CompletionResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		panic(err)
	}

	// The generated command is in the 'text' property of the first choice
	command := response.Choices[0].Text
	return command
}
