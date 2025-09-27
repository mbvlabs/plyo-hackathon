package agents

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mbvlabs/plyo-hackathon/providers"
	"github.com/mbvlabs/plyo-hackathon/tools"
)

const CompanyIntelligenceJobName = "company_intel_job"

type CompanyIntelligenceJobParams struct {
	ReportID      uuid.UUID `json:"report_id"`
	CandidateName string    `json:"candidate_name"`
	CompanyURL    string    `json:"company_url"`
}

const companyIntelligenceSystemPrompt = `
You are a Company Intelligence Agent specialized in gathering comprehensive company information. Your role is to:

- Extract and analyze company websites, About pages, leadership bios, and product/service offerings
- Identify company size, revenue estimates, funding history, and business model
- Map organizational structure, key personnel, and recent company news
- Determine company's primary markets, customer segments, and value propositions
- Flag any compliance issues, controversies, or risk factors
- Output structured company profiles with confidence scores for each data point

You must use tools that directly scrape the company website provided.

Always verify information from multiple sources and note data freshness. Focus on factual, business-relevant intelligence.
	`

type CompanyIntelligence struct {
	client providers.Client
	tools  map[string]tools.Tooler
}

func NewCompanyIntelligence(
	client providers.Client,
	tools map[string]tools.Tooler,
) CompanyIntelligence {
	return CompanyIntelligence{
		client: client,
		tools:  tools,
	}
}

func (r CompanyIntelligence) Research(
	ctx context.Context,
	companyName string,
	companyURL string,
) (string, error) {
	userPrompt := fmt.Sprintf(`
		Conduct comprehensive company intelligence research for %s. 

		Company URL: %s

		Focus on:

		- Company overview, business model, and organizational structure
		- Financial performance, revenue estimates, and funding history
		- Leadership team, key personnel, and board members
		- Product/service portfolio and competitive positioning
		- Recent developments, news, partnerships, and strategic initiatives
		- Market presence, customer base, and geographic footprint
		- Any regulatory issues, controversies, or risk factors
		
		Provide detailed analysis with confidence scores and cite all sources.`,
		companyName,
		companyURL,
	)

	response, err := r.client.Prompt(
		ctx,
		providers.GPT41Mini,
		companyIntelligenceSystemPrompt,
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
			companyIntelligenceSystemPrompt,
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
