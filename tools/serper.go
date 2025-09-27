package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/openai/openai-go/v2"
)

const SerperToolName = "serper_search"

type Serper struct {
	apiKey string
}

func NewSerper(apiKey string) Serper {
	return Serper{apiKey}
}

type SerperRequest struct {
	Query       string `json:"q"`
	Autocorrect bool   `json:"autocorrect,omitempty"`
	Tbs         string `json:"tbs,omitempty"`
}

func (s *Serper) GetName() string {
	return SerperToolName
}

func (s *Serper) Query(query string) ([]byte, error) {
	req := SerperRequest{
		Query:       query,
		Autocorrect: false,
		Tbs:         "qdr:y",
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest(
		"POST",
		"https://google.serper.dev/search",
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

func (s *Serper) GetFunctionStructure() openai.ChatCompletionToolUnionParam {
	param := openai.ChatCompletionToolUnionParam{
		OfFunction: &openai.ChatCompletionFunctionToolParam{
			Function: openai.FunctionDefinitionParam{
				Name: SerperToolName,
				Description: openai.String(
					"Search the web for current information, news, facts, and data. Use this when you need recent information, current events, or want to verify information from the web.",
				),
				Parameters: openai.FunctionParameters{
					"type": "object",
					"properties": map[string]any{
						"query": map[string]string{
							"type":        "string",
							"description": "The search query to execute",
						},
					},
					"required": []string{"query"},
				},
			},
		},
	}

	return param
}

func (s *Serper) Execute(input json.RawMessage) (string, error) {
	var req struct {
		Query string `json:"query"`
	}

	if err := json.Unmarshal(input, &req); err != nil {
		return "", fmt.Errorf("failed to parse input: %w", err)
	}

	if req.Query == "" {
		return "", fmt.Errorf("query parameter is required")
	}

	result, err := s.Query(req.Query)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
