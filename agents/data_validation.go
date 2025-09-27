package agents

import (
	"context"
	"fmt"

	"github.com/mbvlabs/plyo-hackathon/providers"
	"github.com/mbvlabs/plyo-hackathon/tools"
)

const dataValidationSystemPrompt = `
You are a Data Validation Agent responsible for ensuring research quality and accuracy. Your tasks include:

- Cross-reference information across multiple sources
- Identify and flag conflicting data points
- Assess source credibility and data freshness
- Standardize data formats and resolve inconsistencies
- Calculate confidence scores for research findings
- Highlight gaps in information and recommend additional research

Maintain strict quality standards. Clearly distinguish between verified facts, estimates, and assumptions.
`

type DataValidation struct {
	client providers.Client
	tools  map[string]tools.Tooler
}

func NewDataValidation(
	client providers.Client,
	tools map[string]tools.Tooler,
) DataValidation {
	return DataValidation{
		client: client,
		tools:  tools,
	}
}

func (r DataValidation) Research(
	ctx context.Context,
	companyName string,
	companyURL string,
) (string, error) {
	userPrompt := fmt.Sprintf(
		`
Validate and cross-reference research findings for %s (%s). Tasks include:

- Cross-check information across multiple sources including their official website
- Identify and flag conflicting data points
- Assess source credibility and information freshness
- Standardize data formats and resolve inconsistencies
- Calculate confidence scores for key findings
- Highlight information gaps and data quality issues
- Verify key metrics and claims against official company sources
- Flag any suspicious or unverified information

Use %s as the primary source for validation and provide data quality assessment with recommendations for additional verification.
		`,
		companyName,
		companyURL,
		companyURL,
	)

	response, err := r.client.Prompt(
		ctx,
		providers.GPT41Mini,
		dataValidationSystemPrompt,
		userPrompt,
		r.tools,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate research summary: %w", err)
	}

	// finalResponse, err := r.client.Prompt(
	// 	ctx,
	// 	providers.GPT41Mini,
	// 	"Your job is to make sure that the final response adheres to the specific schema. You will receive a string as the user prompt, as well as a schema, return the user prompt in the specified user format.",
	// 	response,
	// 	r.tools,
	// 	&researchBriefSchema,
	// )
	// if err != nil {
	// 	return ResearchBrief{}, fmt.Errorf("failed to generate research summary: %w", err)
	// }
	//
	// var brief ResearchBrief
	// if err := json.Unmarshal([]byte(finalResponse), &brief); err != nil {
	// 	return ResearchBrief{}, fmt.Errorf("failed to unmarshal research brief: %w", err)
	// }

	return response, nil
}
