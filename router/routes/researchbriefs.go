package routes

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const (
	researchbriefsRoutePrefix = "/researchbriefs"
	researchbriefsNamePrefix  = "researchbriefs"
)

var ResearchBriefRoutes = []Route{
	// ResearchBriefIndex,
	// ResearchBriefShow.Route,
	// ResearchBriefNew,
	ResearchBriefCreate,
	// ResearchBriefEdit.Route,
	// ResearchBriefUpdate.Route,
	// ResearchBriefDestroy.Route,
}

var ResearchBriefIndex = Route{
	Name:         researchbriefsNamePrefix + ".index",
	Path:         researchbriefsRoutePrefix,
	Method:       http.MethodGet,
	Handler:      "ResearchBriefs",
	HandleMethod: "Index",
}

var ResearchBriefShow = researchbriefsShow{
	Route: Route{
		Name:         researchbriefsNamePrefix + ".show",
		Path:         researchbriefsRoutePrefix + "/:id",
		Method:       http.MethodGet,
		Handler:      "ResearchBriefs",
		HandleMethod: "Show",
	},
}

type researchbriefsShow struct {
	Route
}

func (r researchbriefsShow) GetPath(id uuid.UUID) string {
	return strings.Replace(r.Path, ":id", id.String(), 1)
}

var ResearchBriefNew = Route{
	Name:         researchbriefsNamePrefix + ".new",
	Path:         researchbriefsRoutePrefix + "/new",
	Method:       http.MethodGet,
	Handler:      "ResearchBriefs",
	HandleMethod: "New",
}

var ResearchBriefCreate = Route{
	Name:         researchbriefsNamePrefix + ".create",
	Path:         researchbriefsRoutePrefix,
	Method:       http.MethodPost,
	Handler:      "ResearchBriefs",
	HandleMethod: "Create",
}

var ResearchBriefEdit = researchbriefsEdit{
	Route: Route{
		Name:         researchbriefsNamePrefix + ".edit",
		Path:         researchbriefsRoutePrefix + "/:id/edit",
		Method:       http.MethodGet,
		Handler:      "ResearchBriefs",
		HandleMethod: "Edit",
	},
}

type researchbriefsEdit struct {
	Route
}

func (r researchbriefsEdit) GetPath(id uuid.UUID) string {
	return strings.Replace(r.Path, ":id", id.String(), 1)
}

var ResearchBriefUpdate = researchbriefsUpdate{
	Route: Route{
		Name:         researchbriefsNamePrefix + ".update",
		Path:         researchbriefsRoutePrefix + "/:id",
		Method:       http.MethodPut,
		Handler:      "ResearchBriefs",
		HandleMethod: "Update",
	},
}

type researchbriefsUpdate struct {
	Route
}

func (r researchbriefsUpdate) GetPath(id uuid.UUID) string {
	return strings.Replace(r.Path, ":id", id.String(), 1)
}

var ResearchBriefDestroy = researchbriefsDestroy{
	Route: Route{
		Name:         researchbriefsNamePrefix + ".destroy",
		Path:         researchbriefsRoutePrefix + "/:id",
		Method:       http.MethodDelete,
		Handler:      "ResearchBriefs",
		HandleMethod: "Destroy",
	},
}

type researchbriefsDestroy struct {
	Route
}

func (r researchbriefsDestroy) GetPath(id uuid.UUID) string {
	return strings.Replace(r.Path, ":id", id.String(), 1)
}
