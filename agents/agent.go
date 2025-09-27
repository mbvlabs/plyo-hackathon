// Package agents exposes agents that can perform very specific tasks using theopenai client and tools exposed in the tools package
package agents

import (
	"github.com/invopop/jsonschema"
)

// five key topics to digg into during

func GenerateSchema[T any]() any {
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

// type ResearchAgent struct {
// 	client       providers.Client
// 	tools        map[string]tools.Tooler
// 	systemPrompt string
// }
//
// func NewResearchAgent(client providers.Client, tools map[string]tools.Tooler) *ResearchAgent {
// 	return &ResearchAgent{
// 		client: client,
// 		tools:  tools,
// 	}
// }
//
// func (r *ResearchAgent) Research(companyName string) (string, error) {
// 	systemPrompt := `You are a research assistant tasked with gathering comprehensive information about companies.
// 	Use the available search tools to find current, accurate information about the company the user asks about.
//
// 	When researching a company, search for:
// 	- Basic company information (what they do, when founded, location)
// 	- Key products or services
// 	- Recent news or developments
// 	- Financial information if available
// 	- Notable achievements or milestones
//
// 	Provide a well-structured summary based on your findings.`
//
// 	userPrompt := fmt.Sprintf(
// 		"Research and provide a comprehensive summary about the company: %s",
// 		companyName,
// 	)
//
// 	response, err := r.client.Prompt(
// 		context.Background(),
// 		providers.GPT4oMini,
// 		systemPrompt,
// 		userPrompt,
// 		r.tools,
// 	)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to generate research summary: %w", err)
// 	}
//
// 	return response, nil
// }
