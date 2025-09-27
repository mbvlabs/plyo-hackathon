package agents

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mbvlabs/plyo-hackathon/providers"
	"github.com/mbvlabs/plyo-hackathon/tools"
)

const CompetitiveIntelligenceJobName = "competitive_intel_job"

type CompetitiveIntelligenceJobParams struct {
	ReportID      uuid.UUID `json:"report_id"`
	CandidateName string    `json:"candidate_name"`
	CompanyURL    string    `json:"company_url"`
}

const competitiveIntelligenceSystemPrompt = `
You are a Competitive Intelligence Agent focused on mapping competitive landscapes. Your responsibilities include:

- Identify direct, indirect, and emerging competitors using multiple discovery methods
- Analyze competitor positioning, pricing strategies, and market share
- Compare product features, customer reviews, and competitive advantages
- Track competitor funding, partnerships, and strategic moves
- Create competitor matrices and positioning maps
- Monitor competitor marketing strategies and messaging

Use the tools available to your disposal to search and gather competitive intelligence. For each analysis:
1. First search for "[company name] competitors"
2. Search for "[company name] vs [competitor]" comparisons
3. Search for "[industry] market landscape [current year]"
4. Search for recent funding/partnership news

Always provide sources and explain your search methodology.

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
		providers.GPT41,
		competitiveIntelligenceSystemPrompt,
		userPrompt,
		r.tools,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate research summary: %w", err)
	}

	if response == "" {
		responseTwo, err := r.client.Prompt(
			ctx,
			providers.GPT41,
			competitiveIntelligenceSystemPrompt,
			userPrompt,
			r.tools,
			nil,
		)
		if err != nil {
			return "", fmt.Errorf("failed to generate research summary: %w", err)
		}

		return responseTwo, nil
	}

	return response, nil
}
