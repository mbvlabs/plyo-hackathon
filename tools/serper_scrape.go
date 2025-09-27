package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/openai/openai-go/v2"
)

const SerperScrapeToolName = "serper_scrape"

type SerperScrape struct {
	apiKey string
}

func NewSerperScrape(apiKey string) SerperScrape {
	return SerperScrape{apiKey}
}

type SerperScrapeRequest struct {
	URL string `json:"url"`
}

func (s *SerperScrape) GetName() string {
	return SerperScrapeToolName
}

func (s *SerperScrape) Scrape(url string) ([]byte, error) {
	req := SerperScrapeRequest{
		URL: url,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest(
		"POST",
		"https://scrape.serper.dev",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("X-API-KEY", s.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"API request failed with status %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	return body, nil
}

func (s *SerperScrape) GetFunctionStructure() openai.ChatCompletionToolUnionParam {
	param := openai.ChatCompletionToolUnionParam{
		OfFunction: &openai.ChatCompletionFunctionToolParam{
			Function: openai.FunctionDefinitionParam{
				Name: SerperScrapeToolName,
				Description: openai.String(
					"Scrape content from any website URL. Use this to extract text, data, and information from web pages.",
				),
				Parameters: openai.FunctionParameters{
					"type": "object",
					"properties": map[string]any{
						"url": map[string]string{
							"type":        "string",
							"description": "The URL of the website to scrape",
						},
					},
					"required": []string{"url"},
				},
			},
		},
	}

	return param
}

func (s *SerperScrape) Execute(input json.RawMessage) (string, error) {
	var req struct {
		URL string `json:"url"`
	}

	if err := json.Unmarshal(input, &req); err != nil {
		return "", fmt.Errorf("failed to parse input: %w", err)
	}

	if req.URL == "" {
		return "", fmt.Errorf("url parameter is required")
	}

	result, err := s.Scrape(req.URL)
	if err != nil {
		return "", err
	}

	return string(result), nil
}