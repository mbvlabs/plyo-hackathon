// Package providers should provide the functionality to intract with the openai API
package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mbvlabs/plyo-hackathon/tools"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

type Client struct {
	client openai.Client
}

type Model string

const (
	// Expensive/Premium models
	GPT41     Model = "gpt-4.1"
	GPT41Mini Model = "gpt-4.1-mini"
	GPT4Turbo Model = "gpt-4-turbo"
	GPT4      Model = "gpt-4"

	// Cheap/Fast models
	GPT35Turbo Model = "gpt-3.5-turbo"
)

func NewClient(apiKey string) Client {
	return Client{openai.NewClient(option.WithAPIKey(apiKey))}
}

func (c *Client) retryWithBackoff(ctx context.Context, fn func() error, maxRetries int) error {
	for attempt := 0; attempt <= maxRetries; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		if attempt == maxRetries {
			return err
		}

		backoffDuration := time.Duration(1<<attempt) * time.Second
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoffDuration):
		}
	}
	return nil
}

func (c *Client) Prompt(
	ctx context.Context,
	model Model,
	systemPrompt, userPrompt string,
	tools map[string]tools.Tooler,
	responseFormat *openai.ResponseFormatJSONSchemaJSONSchemaParam,
) (string, error) {
	messages := []openai.ChatCompletionMessageParamUnion{}

	if systemPrompt != "" {
		messages = append(messages, openai.SystemMessage(systemPrompt))
	}

	messages = append(messages, openai.UserMessage(userPrompt))

	agentTools := make([]openai.ChatCompletionToolUnionParam, len(tools))
	i := 0
	for _, tool := range tools {
		agentTools[i] = tool.GetFunctionStructure()
		i++
	}

	params := openai.ChatCompletionNewParams{
		Model:    string(model),
		Messages: messages,
		Tools:    agentTools,
	}

	if responseFormat != nil {
		params.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{JSONSchema: *responseFormat},
		}
	}

	var resp *openai.ChatCompletion
	err := c.retryWithBackoff(ctx, func() error {
		var apiErr error
		resp, apiErr = c.client.Chat.Completions.New(ctx, params)
		return apiErr
	}, 3)
	if err != nil {
		return "", fmt.Errorf("failed to create completion after retries: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from API")
	}

	toolCalls := resp.Choices[0].Message.ToolCalls

	if len(toolCalls) == 0 {
		return resp.Choices[0].Message.Content, nil
	}

	params.Messages = append(params.Messages, resp.Choices[0].Message.ToParam())
	for _, toolCall := range toolCalls {
		if toolCall.Function.Name != "" {
			var args map[string]any
			err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
			if err != nil {
				return "", fmt.Errorf("failed to unmarshal tool arguments: %w", err)
			}

			result, err := tools[toolCall.Function.Name].Execute(
				json.RawMessage(toolCall.Function.Arguments),
			)
			if err != nil {
				return "", fmt.Errorf("tool execution failed: %w", err)
			}

			params.Messages = append(params.Messages, openai.ToolMessage(result, toolCall.ID))
		}
	}

	err = c.retryWithBackoff(ctx, func() error {
		var apiErr error
		resp, apiErr = c.client.Chat.Completions.New(ctx, params)
		return apiErr
	}, 3)
	if err != nil {
		return "", fmt.Errorf("failed to create completion after retries: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}
