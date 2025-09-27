package agents

import (
	"context"
	"fmt"

	"github.com/mbvlabs/plyo-hackathon/providers"
	"github.com/mbvlabs/plyo-hackathon/tools"
)

const competitiveIntelligenceSystemPrompt = `
You are a Competitive Intelligence Agent focused on mapping competitive landscapes. Your responsibilities include:

- Identify direct, indirect, and emerging competitors using multiple discovery methods
- Analyze competitor positioning, pricing strategies, and market share
- Compare product features, customer reviews, and competitive advantages
- Track competitor funding, partnerships, and strategic moves
- Create competitor matrices and positioning maps
- Monitor competitor marketing strategies and messaging

Prioritize current, publicly available information. Classify competitors by threat level and market overlap. Provide actionable competitive insights.`

type CompetitiveIntelligence struct {
	client providers.Client
	tools  map[string]tools.Tooler
}

func NewCompetitiveIntelligence(
	client providers.Client,
	tools map[string]tools.Tooler,
) CompetitiveIntelligence {
	return CompetitiveIntelligence{
		client: client,
		tools:  tools,
	}
}

func (r CompetitiveIntelligence) Research(
	ctx context.Context,
	companyName string,
	companyURL string,
) (string, error) {
	userPrompt := fmt.Sprintf(
		`
Perform competitive landscape analysis for %s (%s). Research and analyze:

- Direct competitors offering similar products/services
- Indirect competitors and alternative solutions
- Emerging threats and new market entrants
- Competitive positioning and market share estimates
- Competitor strengths, weaknesses, and strategies
- Pricing models and go-to-market approaches
- Recent competitive moves, funding, and partnerships

Use %s as your starting point to understand their positioning, then create competitor profiles with threat assessments and strategic implications.
		`,
		companyName,
		companyURL,
		companyURL,
	)

	response, err := r.client.Prompt(
		ctx,
		providers.GPT41Mini,
		competitiveIntelligenceSystemPrompt,
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
