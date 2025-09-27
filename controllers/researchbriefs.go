package controllers

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mbvlabs/plyo-hackathon/agents"
	"github.com/mbvlabs/plyo-hackathon/database"
	"github.com/mbvlabs/plyo-hackathon/models"
	"github.com/mbvlabs/plyo-hackathon/router/cookies"
	"github.com/mbvlabs/plyo-hackathon/router/routes"
	"github.com/mbvlabs/plyo-hackathon/views"
)

type ResearchBriefs struct {
	agent agents.PreliminaryResearch
	db    database.SQLite
}

func newResearchBriefs(agent agents.PreliminaryResearch, db database.SQLite) ResearchBriefs {
	return ResearchBriefs{agent, db}
}

//
// func (r ResearchBriefs) Index(c echo.Context) error {
// 	page := int64(1)
// 	if p := c.QueryParam("page"); p != "" {
// 		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
// 			page = int64(parsed)
// 		}
// 	}
//
// 	perPage := int64(25)
// 	if pp := c.QueryParam("per_page"); pp != "" {
// 		if parsed, err := strconv.Atoi(pp); err == nil && parsed > 0 &&
// 			parsed <= 100 {
// 			perPage = int64(parsed)
// 		}
// 	}
//
// 	researchbriefsList, err := models.PaginateResearchBriefs(
// 		c.Request().Context(),
// 		r.db.Conn(),
// 		page,
// 		perPage,
// 	)
// 	if err != nil {
// 		return render(c, views.InternalError())
// 	}
//
// 	return c.HTML(http.StatusOK, "researchbriefs index - no views implemented")
// }
//
// func (r ResearchBriefs) Show(c echo.Context) error {
// 	researchbriefID, err := uuid.Parse(c.Param("id"))
// 	if err != nil {
// 		return render(c, views.BadRequest())
// 	}
//
// 	researchbrief, err := models.FindResearchBrief(c.Request().Context(), r.db.Conn(), researchbriefID)
// 	if err != nil {
// 		return render(c, views.NotFound())
// 	}
//
// 	return c.HTML(http.StatusOK, "researchbrief show - no views implemented")
// }
//
// func (r ResearchBriefs) New(c echo.Context) error {
// 	return c.HTML(http.StatusOK, "researchbrief new - no views implemented")
// }

type CreateResearchBriefFormPayload struct {
	Query string `json:"query"`
}

func (r ResearchBriefs) Create(c echo.Context) error {
	var payload CreateResearchBriefFormPayload
	if err := c.Bind(&payload); err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"could not parse CreateResearchBriefFormPayload",
			"error",
			err,
		)

		return render(c, views.NotFound())
	}

	result, err := r.agent.Research(c.Request().Context(), payload.Query)
	if err != nil {
		return err
	}

	data := models.CreateResearchBriefData{
		IdentificationStatus: result.IdentificationStatus,
		CompanyName:          result.CompanyName,
		OfficialDomain:       result.OfficialDomain,
		Headquarters:         result.Headquarters,
		Industry:             result.Industry,
		CompanyType:          result.CompanyType,
		Status:               result.Status,
		GeographicScope:      result.GeographicScope,
		ResearchDepth:        result.ResearchDepth,
		ConfidenceScore:      result.ConfidenceScore,
		LastUpdated:          time.Now(),
	}

	researchbrief, err := models.CreateResearchBrief(
		c.Request().Context(),
		r.db.Conn(),
		data,
	)
	if err != nil {
		if flashErr := cookies.AddFlash(c, cookies.FlashError, fmt.Sprintf("Failed to create researchbrief: %v", err)); flashErr != nil {
			return flashErr
		}
		return c.Redirect(http.StatusSeeOther, routes.ResearchBriefNew.Path)
	}

	// Store CompanyCandidates
	for _, candidate := range result.CompanyCandidates {
		_, err := models.CreateCompanyCandidates(
			c.Request().Context(),
			r.db.Conn(),
			models.CreateCompanyCandidatesData{
				ResearchBriefID: researchbrief.ID.String(),
				Name:            candidate.Name,
				Domain:          candidate.Domain,
				Description:     candidate.Description,
				Industry:        candidate.Industry,
				Location:        candidate.Location,
			},
		)
		if err != nil {
			slog.ErrorContext(
				c.Request().Context(),
				"failed to create company candidate",
				"error", err,
				"research_brief_id", researchbrief.ID,
			)
		}
	}

	// Store SpecialConsiderations
	for _, consideration := range result.SpecialConsiderations {
		_, err := models.CreateSpecialConsiderations(
			c.Request().Context(),
			r.db.Conn(),
			models.CreateSpecialConsiderationsData{
				ResearchBriefID: researchbrief.ID.String(),
				Consideration:   consideration,
			},
		)
		if err != nil {
			slog.ErrorContext(
				c.Request().Context(),
				"failed to create special consideration",
				"error", err,
				"research_brief_id", researchbrief.ID,
			)
		}
	}

	// Store Sources
	for _, source := range result.Sources {
		_, err := models.CreateSources(
			c.Request().Context(),
			r.db.Conn(),
			models.CreateSourcesData{
				ResearchBriefID: researchbrief.ID.String(),
				SourceUrl:       source,
			},
		)
		if err != nil {
			slog.ErrorContext(
				c.Request().Context(),
				"failed to create source",
				"error", err,
				"research_brief_id", researchbrief.ID,
			)
		}
	}

	// Store AgentGuidance
	for key, value := range result.AgentGuidance {
		_, err := models.CreateAgentGuidance(
			c.Request().Context(),
			r.db.Conn(),
			models.CreateAgentGuidanceData{
				ResearchBriefID: researchbrief.ID.String(),
				GuidanceKey:     key,
				GuidanceValue:   value,
			},
		)
		if err != nil {
			slog.ErrorContext(
				c.Request().Context(),
				"failed to create agent guidance",
				"error", err,
				"research_brief_id", researchbrief.ID,
			)
		}
	}

	// if flashErr := cookies.AddFlash(c, cookies.FlashSuccess, "ResearchBrief created successfully"); flashErr != nil {
	// 	return render(c, views.InternalError())
	// }

	// Fetch related data
	companyCandidates, err := models.FindCompanyCandidatesByResearchBriefID(
		c.Request().Context(),
		r.db.Conn(),
		researchbrief.ID.String(),
	)
	if err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"failed to fetch company candidates",
			"error", err,
			"research_brief_id", researchbrief.ID,
		)
		companyCandidates = []models.CompanyCandidates{}
	}

	specialConsiderations, err := models.FindSpecialConsiderationsByResearchBriefID(
		c.Request().Context(),
		r.db.Conn(),
		researchbrief.ID.String(),
	)
	if err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"failed to fetch special considerations",
			"error", err,
			"research_brief_id", researchbrief.ID,
		)
		specialConsiderations = []models.SpecialConsiderations{}
	}

	sources, err := models.FindSourcesByResearchBriefID(
		c.Request().Context(),
		r.db.Conn(),
		researchbrief.ID.String(),
	)
	if err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"failed to fetch sources",
			"error", err,
			"research_brief_id", researchbrief.ID,
		)
		sources = []models.Sources{}
	}

	agentGuidances, err := models.FindAgentGuidancesByResearchBriefID(
		c.Request().Context(),
		r.db.Conn(),
		researchbrief.ID.String(),
	)
	if err != nil {
		slog.ErrorContext(
			c.Request().Context(),
			"failed to fetch agent guidances",
			"error", err,
			"research_brief_id", researchbrief.ID,
		)
		agentGuidances = []models.AgentGuidance{}
	}

	sse := getSSE(c)
	return sse.PatchElementTempl(
		views.PreliminaryResearchResults(
			researchbrief,
			companyCandidates,
			specialConsiderations,
			sources,
			agentGuidances,
		),
	)
}

// func (r ResearchBriefs) Edit(c echo.Context) error {
// 	researchbriefID, err := uuid.Parse(c.Param("id"))
// 	if err != nil {
// 		return render(c, views.BadRequest())
// 	}
//
// 	researchbrief, err := models.FindResearchBrief(
// 		c.Request().Context(),
// 		r.db.Conn(),
// 		researchbriefID,
// 	)
// 	if err != nil {
// 		return render(c, views.NotFound())
// 	}
//
// 	return c.HTML(http.StatusOK, "researchbrief edit - no views implemented")
// }
//
// type UpdateResearchBriefFormPayload struct {
// 	IdentificationStatus string  `form:"identification_status"`
// 	CompanyName          string  `form:"company_name"`
// 	OfficialDomain       string  `form:"official_domain"`
// 	Headquarters         string  `form:"headquarters"`
// 	Industry             string  `form:"industry"`
// 	CompanyType          string  `form:"company_type"`
// 	Status               string  `form:"status"`
// 	GeographicScope      string  `form:"geographic_scope"`
// 	ResearchDepth        string  `form:"research_depth"`
// 	ConfidenceScore      float64 `form:"confidence_score"`
// 	LastUpdated          string  `form:"last_updated"`
// }
//
// func (r ResearchBriefs) Update(c echo.Context) error {
// 	researchbriefID, err := uuid.Parse(c.Param("id"))
// 	if err != nil {
// 		return render(c, views.BadRequest())
// 	}
//
// 	var payload UpdateResearchBriefFormPayload
// 	if err := c.Bind(&payload); err != nil {
// 		slog.ErrorContext(
// 			c.Request().Context(),
// 			"could not parse UpdateResearchBriefFormPayload",
// 			"error",
// 			err,
// 		)
//
// 		return render(c, views.NotFound())
// 	}
//
// 	data := models.UpdateResearchBriefData{
// 		ID:                   researchbriefID,
// 		IdentificationStatus: payload.IdentificationStatus,
// 		CompanyName:          payload.CompanyName,
// 		OfficialDomain:       payload.OfficialDomain,
// 		Headquarters:         payload.Headquarters,
// 		Industry:             payload.Industry,
// 		CompanyType:          payload.CompanyType,
// 		Status:               payload.Status,
// 		GeographicScope:      payload.GeographicScope,
// 		ResearchDepth:        payload.ResearchDepth,
// 		ConfidenceScore:      payload.ConfidenceScore,
// 		LastUpdated: func() time.Time {
// 			if payload.LastUpdated == "" {
// 				return time.Time{}
// 			}
// 			if t, err := time.Parse("2006-01-02", payload.LastUpdated); err == nil {
// 				return t
// 			}
// 			return time.Time{}
// 		}(),
// 	}
//
// 	researchbrief, err := models.UpdateResearchBrief(
// 		c.Request().Context(),
// 		r.db.Conn(),
// 		data,
// 	)
// 	if err != nil {
// 		if flashErr := cookies.AddFlash(c, cookies.FlashError, fmt.Sprintf("Failed to update researchbrief: %v", err)); flashErr != nil {
// 			return render(c, views.InternalError())
// 		}
// 		return c.Redirect(
// 			http.StatusSeeOther,
// 			routes.ResearchBriefEdit.GetPath(researchbriefID),
// 		)
// 	}
//
// 	if flashErr := cookies.AddFlash(c, cookies.FlashSuccess, "ResearchBrief updated successfully"); flashErr != nil {
// 		return render(c, views.InternalError())
// 	}
//
// 	return c.Redirect(http.StatusSeeOther, routes.ResearchBriefShow.GetPath(researchbrief.ID))
// }
//
// func (r ResearchBriefs) Destroy(c echo.Context) error {
// 	researchbriefID, err := uuid.Parse(c.Param("id"))
// 	if err != nil {
// 		return render(c, views.BadRequest())
// 	}
//
// 	err = models.DestroyResearchBrief(c.Request().Context(), r.db.Conn(), researchbriefID)
// 	if err != nil {
// 		if flashErr := cookies.AddFlash(c, cookies.FlashError, fmt.Sprintf("Failed to delete researchbrief: %v", err)); flashErr != nil {
// 			return render(c, views.InternalError())
// 		}
// 		return c.Redirect(http.StatusSeeOther, routes.ResearchBriefIndex.Path)
// 	}
//
// 	if flashErr := cookies.AddFlash(c, cookies.FlashSuccess, "ResearchBrief destroyed successfully"); flashErr != nil {
// 		return render(c, views.InternalError())
// 	}
//
// 	return c.Redirect(http.StatusSeeOther, routes.ResearchBriefIndex.Path)
// }
