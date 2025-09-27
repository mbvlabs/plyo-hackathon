package agents

import (
	"context"
	"fmt"

	"github.com/mbvlabs/plyo-hackathon/providers"
	"github.com/mbvlabs/plyo-hackathon/tools"
)

const reportGeneratorSystemPrompt = `
You are a Report Generator Agent specialized in synthesizing business intelligence research into comprehensive executive reports.

Your responsibilities:
- Integrate findings from multiple research agents into a cohesive narrative
- Identify cross-cutting themes and strategic implications
- Resolve conflicting information using confidence scores and source quality
- Prioritize insights by business impact and reliability
- Generate actionable strategic recommendations
- Present findings in executive-friendly format
- Highlight areas requiring additional research
- Identify five topics for further discussion when contacting the company

Focus on delivering clear, actionable insights that enable strategic decision-making. 

Your output must be a well-structured markdown report that executives can use for strategic planning.

Always include a section at the end that has all the sources used to generate the report.
`

type ReportGenerator struct {
	client providers.Client
	tools  map[string]tools.Tooler
}

func NewReportGenerator(
	client providers.Client,
	tools map[string]tools.Tooler,
) ReportGenerator {
	return ReportGenerator{
		client: client,
		tools:  tools,
	}
}

func (r ReportGenerator) Generate(
	ctx context.Context,
	companyName string,
	companyURL string,
	companyIntelligenceFindings string,
	competitiveLandscapeAnalysis string,
	marketDynamicsAssessment string,
	industryTrendAnalysis string,
) (string, error) {
	userPrompt := fmt.Sprintf(`
Generate a comprehensive business intelligence report for %s (%s).

COMPANY INTELLIGENCE FINDINGS:
%s

COMPETITIVE LANDSCAPE ANALYSIS:
%s

MARKET DYNAMICS ASSESSMENT:
%s

INDUSTRY TRENDS ANALYSIS:
%s

Synthesize all findings into a cohesive executive report with:
- Executive Summary
- Key Strategic Insights  
- Market Opportunities & Threats
- Competitive Positioning
- Strategic Recommendations
- Risk Assessment
- Confidence Scores for major conclusions

Present the analysis in a structured format that enables strategic decision-making.`,
		companyName,
		companyURL,
		companyIntelligenceFindings,
		competitiveLandscapeAnalysis,
		marketDynamicsAssessment,
		industryTrendAnalysis,
	)

	response, err := r.client.Prompt(
		ctx,
		providers.GPT41Mini,
		reportGeneratorSystemPrompt,
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
