package agents

import (
	"context"
	"fmt"

	"github.com/mbvlabs/plyo-hackathon/providers"
	"github.com/mbvlabs/plyo-hackathon/tools"
)

const researchOrchestratorSystemPrompt = `
You are the Research Orchestrator Agent managing the entire research workflow. Your responsibilities include:

- Coordinate task distribution among specialized agents
- Manage research priorities and time allocation
- Synthesize findings into coherent, actionable reports
- Identify information gaps and direct follow-up research
- Ensure consistency across different research streams
- Format final deliverables according to user specifications

Focus on delivering comprehensive, well-structured insights that directly address user research objectives. Maintain clear traceability of sources and methodology.
`

type ResearchOrchestrator struct {
	client providers.Client
	tools  map[string]tools.Tooler
}

func NewResearchOrchestrator(
	client providers.Client,
	tools map[string]tools.Tooler,
) ResearchOrchestrator {
	return ResearchOrchestrator{
		client: client,
		tools:  tools,
	}
}

func (r ResearchOrchestrator) Research(
	ctx context.Context,
	companyName string,
	companyURL string,
) (string, error) {
	userPrompt := fmt.Sprintf(
		`
Synthesize and coordinate research findings for %s (%s). Responsibilities:

- Integrate insights from all specialized research agents
- Identify cross-cutting themes and strategic implications
- Resolve conflicting information and data inconsistencies
- Create executive summary with key findings and recommendations
- Highlight critical gaps requiring additional research
- Prioritize insights by business impact and confidence level
- Format comprehensive research report
- Provide actionable strategic recommendations

Reference the official company information from %s and deliver cohesive analysis that addresses the original research objectives.
		`,
		companyName,
		companyURL,
		companyURL,
	)

	response, err := r.client.Prompt(
		ctx,
		providers.GPT41Mini,
		researchOrchestratorSystemPrompt,
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
