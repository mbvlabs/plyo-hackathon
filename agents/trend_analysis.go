package agents

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mbvlabs/plyo-hackathon/providers"
	"github.com/mbvlabs/plyo-hackathon/tools"
)

const TrendAnalysisJobName = "trend_analysis_job"

type TrendAnalysisJobParams struct {
	ReportID      uuid.UUID `json:"report_id"`
	CandidateName string    `json:"candidate_name"`
	CompanyURL    string    `json:"company_url"`
}

const trendAnalysisSystemPrompt = `
You are a Trend Analysis Agent focused on identifying and analyzing industry trends. Your responsibilities include:

- Track emerging technologies, business models, and industry innovations
- Analyze consumer behavior shifts and demographic trends
- Monitor regulatory changes and policy developments
- Identify cyclical patterns and seasonal variations
- Forecast future market directions and disruption risks
- Correlate macro-economic factors with industry performance

Distinguish between short-term fluctuations and long-term trends. Provide probability-weighted scenarios and timeline estimates.
`

type TrendAnalysis struct {
	client providers.Client
	tools  map[string]tools.Tooler
}

func NewTrendAnalysis(
	client providers.Client,
	tools map[string]tools.Tooler,
) TrendAnalysis {
	return TrendAnalysis{
		client: client,
		tools:  tools,
	}
}

func (r TrendAnalysis) Research(
	ctx context.Context,
	companyName string,
	companyURL string,
) (string, error) {
	userPrompt := fmt.Sprintf(
		`
Identify and analyze industry trends affecting %s (%s). Focus on:

- Emerging technologies and innovation trends
- Shifting consumer behaviors and preferences
- Regulatory changes and policy developments
- Economic factors and market cyclicality
- Demographic and social trends
- Competitive landscape evolution
- Future market directions and disruption risks
- Technology adoption patterns and digital transformation

Review %s to understand their current approach, then distinguish between short-term fluctuations and long-term structural trends. Provide timeline estimates and probability assessments.
		`,
		companyName,
		companyURL,
		companyURL,
	)

	response, err := r.client.Prompt(
		ctx,
		providers.GPT41Mini,
		trendAnalysisSystemPrompt,
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
			providers.GPT41Mini,
			trendAnalysisSystemPrompt,
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
