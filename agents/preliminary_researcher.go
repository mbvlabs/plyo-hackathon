package agents

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mbvlabs/plyo-hackathon/providers"
	"github.com/mbvlabs/plyo-hackathon/tools"
	"github.com/openai/openai-go/v2"
)

// RESEARCH BRIEF - [Company Name]
// ================================
//
// CONFIRMED IDENTITY:
// - Official Name: [Verified company name]
// - Primary Domain: [company.com]
// - Headquarters: [City, Country]
// - Industry: [Primary classification]
// - Company Type: [Public/Private/Subsidiary]
// - Status: [Active/Acquired/Defunct]
//
// RESEARCH PARAMETERS:
// - Geographic Scope: [Global/Regional/Local]
// - Market Focus: [B2B/B2C/B2B2C]
// - Research Depth: [Comprehensive/Standard/Basic]
// - Special Considerations: [Any unique factors]
//
// AGENT GUIDANCE:
// - Company Intelligence: Focus on [specific areas]
// - Competitive Intelligence: Competitor scope [direct/indirect/emerging]
// - Market Dynamics: Market definition [specific market boundaries]
// - Trend Analysis: Industry focus [specific industry segments]

type ResearchBrief struct {
	// Company Identification
	IdentificationStatus string             `json:"identification_status" jsonschema:"required" jsonschema_description:"Status of company identification"`
	CompanyCandidates    []CompanyCandidate `json:"company_candidates"    jsonschema:"required" jsonschema_description:"List of potential company matches when ambiguous"`

	// Core Company Data
	CompanyName           string   `json:"company_name"           jsonschema:"required" jsonschema_description:"The verified official name of the company"`
	OfficialDomain        string   `json:"official_domain"        jsonschema:"required" jsonschema_description:"The primary domain of the company"`
	Headquarters          string   `json:"headquarters"           jsonschema:"required" jsonschema_description:"The location of company headquarters"`
	Industry              string   `json:"industry"               jsonschema:"required" jsonschema_description:"The primary industry classification"`
	CompanyType           string   `json:"company_type"           jsonschema:"required" jsonschema_description:"The type of company organization"`
	Status                string   `json:"status"                 jsonschema:"required" jsonschema_description:"Current operational status of the company"`
	GeographicScope       string   `json:"geographic_scope"       jsonschema:"required" jsonschema_description:"The geographic scope of operations"`
	ResearchDepth         string   `json:"research_depth"         jsonschema:"required" jsonschema_description:"The depth of research conducted"`
	SpecialConsiderations []string `json:"special_considerations" jsonschema:"required" jsonschema_description:"Any unique factors or considerations about the company"`
	ConfidenceScore       float64  `json:"confidence_score"       jsonschema:"required" jsonschema_description:"Confidence level in the research findings (0.0-1.0)"`
	Sources               []string `json:"sources"                jsonschema:"required" jsonschema_description:"List of sources used for the research"`
	LastUpdated           string   `json:"last_updated"           jsonschema:"required" jsonschema_description:"Timestamp of when the research was last updated"`

	// Optional fields
	AgentGuidance map[string]string `json:"agent_guidance,omitempty" jsonschema_description:"Guidance for specialized research agents"`
}

type CompanyCandidate struct {
	Name        string `json:"name"        jsonschema:"required" jsonschema_description:"Candidate company name"`
	Domain      string `json:"domain"      jsonschema:"required" jsonschema_description:"Official website URL for this candidate"`
	Description string `json:"description" jsonschema:"required" jsonschema_description:"Brief description to help distinguish this candidate"`
	Industry    string `json:"industry"    jsonschema:"required" jsonschema_description:"Primary industry of this candidate"`
	Location    string `json:"location"    jsonschema:"required" jsonschema_description:"Headquarters or primary location"`
}

const preliminaryResearchSystemPrompt = `You are a research assistant tasked with gathering comprehensive information about companies.

When given a company name or URL, FIRST verify you have the correct company by:
- If given only a company name, search for and identify the official website URL
- If multiple companies share similar names, return ALL candidates with their official URLs and brief descriptions
- Ask the user to confirm which specific company they want researched
- Only proceed with detailed research once company identity is confirmed

Once the correct company is identified, use the available search tools to find current, accurate information about:
- Basic company information (what they do, when founded, location)
- Key products or services  
- Recent news or developments
- Financial information if available
- Notable achievements or milestones

Always include the official company URL in your findings. Provide a well-structured summary based on your findings.`

var researchBriefSchema = openai.ResponseFormatJSONSchemaJSONSchemaParam{
	Name: "research_brief",
	Description: openai.String(
		"Preliminary research brief used to identify correct company and show basic data",
	),
	Schema: GenerateSchema[ResearchBrief](),
	Strict: openai.Bool(true),
}

type PreliminaryResearch struct {
	client providers.Client
	tools  map[string]tools.Tooler
}

func NewPreliminaryResearch(
	client providers.Client,
	tools map[string]tools.Tooler,
) *PreliminaryResearch {
	return &PreliminaryResearch{
		client: client,
		tools:  tools,
	}
}

func (r *PreliminaryResearch) Research(
	ctx context.Context,
	companyName string,
) (ResearchBrief, error) {
	userPrompt := fmt.Sprintf(
		"Research and provide a comprehensive summary about the company: %s",
		companyName,
	)

	response, err := r.client.Prompt(
		ctx,
		providers.GPT41Mini,
		preliminaryResearchSystemPrompt,
		userPrompt,
		r.tools,
		nil,
	)
	if err != nil {
		return ResearchBrief{}, fmt.Errorf("failed to generate research summary: %w", err)
	}

	finalResponse, err := r.client.Prompt(
		ctx,
		providers.GPT41Mini,
		"Your job is to make sure that the final response adheres to the specific schema. You will receive a string as the user prompt, as well as a schema, return the user prompt in the specified user format.",
		response,
		r.tools,
		&researchBriefSchema,
	)
	if err != nil {
		return ResearchBrief{}, fmt.Errorf("failed to generate research summary: %w", err)
	}

	var brief ResearchBrief
	if err := json.Unmarshal([]byte(finalResponse), &brief); err != nil {
		return ResearchBrief{}, fmt.Errorf("failed to unmarshal research brief: %w", err)
	}

	return brief, nil
}
