package agents

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mbvlabs/plyo-hackathon/providers"
	"github.com/mbvlabs/plyo-hackathon/tools"
)

const MarketDynamicsJobName = "market_dynamics_job"

type MarketDynamicsJobParams struct {
	ReportID      uuid.UUID `json:"report_id"`
	CandidateName string    `json:"candidate_name"`
	CompanyURL    string    `json:"company_url"`
}

const marketDynamicsSystemPrompt = `
You are a Market Dynamics Agent specializing in market analysis and sizing. Your role encompasses:

- Calculate total addressable market (TAM), serviceable addressable market (SAM), and serviceable obtainable market (SOM)
- Analyze market structure, key players, and barriers to entry
- Identify market drivers, restraints, and growth opportunities
- Map customer segments, buying behaviors, and decision-making processes
- Assess regulatory environment and policy impacts
- Evaluate market maturity and lifecycle stage

Use multiple data sources and methodologies. Provide quantified insights with clear assumptions and limitations.
`

type MarketDynamics struct {
	client providers.Client
	tools  map[string]tools.Tooler
}

func NewMarketDynamics(
	client providers.Client,
	tools map[string]tools.Tooler,
) MarketDynamics {
	return MarketDynamics{
		client: client,
		tools:  tools,
	}
}

func (r MarketDynamics) Research(
	ctx context.Context,
	companyName string,
	companyURL string,
) (string, error) {
	userPrompt := fmt.Sprintf(
		`
Analyze market dynamics and opportunities for %s (%s). Examine:

- Total addressable market (TAM), serviceable addressable market (SAM), and serviceable obtainable market (SOM)
- Market size, growth rates, and revenue projections
- Market structure, key players, and concentration levels
- Customer segments, buying behaviors, and decision factors
- Market drivers, restraints, and growth opportunities
- Barriers to entry and competitive dynamics
- Regulatory environment and policy impacts
- Market maturity and lifecycle stage

Reference %s to understand their market positioning and provide quantified market insights with methodology and assumptions.
		`,
		companyName,
		companyURL,
		companyURL,
	)

	response, err := r.client.Prompt(
		ctx,
		providers.GPT41Mini,
		marketDynamicsSystemPrompt,
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
