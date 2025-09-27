// Package openai should provide the functionality to intract with the openai API
package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	apiKey string
}

type PromptRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

type PromptResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type Model string

const (
	// Expensive/Premium models
	GPT4o     Model = "gpt-4o"
	GPT4oMini Model = "gpt-4o-mini"
	GPT4Turbo Model = "gpt-4-turbo"
	GPT4      Model = "gpt-4"

	// Cheap/Fast models
	GPT35Turbo Model = "gpt-3.5-turbo"
)

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

func (c *Client) Prompt(model Model, systemPrompt, userPrompt string) (string, error) {
	messages := []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}{}

	if systemPrompt != "" {
		messages = append(messages, struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{Role: "system", Content: systemPrompt})
	}

	messages = append(messages, struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}{Role: "user", Content: userPrompt})

	req := PromptRequest{
		Model:    string(model),
		Messages: messages,
	}

	var lastErr error
	for attempt := range 5 {
		resp, err := c.makeRequest(req)
		if err == nil && len(resp.Choices) > 0 {
			return resp.Choices[0].Message.Content, nil
		}
		lastErr = err

		if attempt < 4 {
			backoff := time.Duration(1<<attempt) * time.Second
			time.Sleep(backoff)
		}
	}

	return "", fmt.Errorf("failed after 5 attempts: %v", lastErr)
}

func (c *Client) makeRequest(req PromptRequest) (*PromptResponse, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(
		"POST",
		"https://api.openai.com/v1/chat/completions",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != 200 {
		body, _ := io.ReadAll(httpResp.Body)
		return nil, fmt.Errorf("API error %d: %s", httpResp.StatusCode, string(body))
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var resp PromptResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
