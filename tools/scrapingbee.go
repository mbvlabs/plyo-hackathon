package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/openai/openai-go/v2"
)

const ScrapingBeeToolName = "scrapingbee_scraper"

type ScrapingBee struct {
	apiKey string
}

func NewScrapingBee(apiKey string) ScrapingBee {
	return ScrapingBee{apiKey}
}

type ScrapingBeeRequest struct {
	URL      string `json:"url"`
	RenderJS bool   `json:"render_js,omitempty"`
}

func (s *ScrapingBee) GetName() string {
	return ScrapingBeeToolName
}

func (s *ScrapingBee) Scrape(targetURL string) ([]byte, error) {
	slog.Info("SCRAAAAAAAAAAAAAAAAAAAAAAPING BEEEEEEEEEEE")
	baseURL := "https://app.scrapingbee.com/api/v1"

	params := url.Values{}
	params.Add("api_key", s.apiKey)
	params.Add("url", targetURL)
	params.Add("render_js", "false") // Use cheapest option (1 credit vs 5)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	httpReq, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

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

func (s *ScrapingBee) GetFunctionStructure() openai.ChatCompletionToolUnionParam {
	param := openai.ChatCompletionToolUnionParam{
		OfFunction: &openai.ChatCompletionFunctionToolParam{
			Function: openai.FunctionDefinitionParam{
				Name: ScrapingBeeToolName,
				Description: openai.String(
					"Scrape web pages to extract HTML content. Use this when you need to get the content of a specific webpage. Uses the most cost-effective option (no JavaScript rendering).",
				),
				Parameters: openai.FunctionParameters{
					"type": "object",
					"properties": map[string]any{
						"url": map[string]string{
							"type":        "string",
							"description": "The full URL to scrape (must include http:// or https://)",
						},
					},
					"required": []string{"url"},
				},
			},
		},
	}

	return param
}

func (s *ScrapingBee) Execute(input json.RawMessage) (string, error) {
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
