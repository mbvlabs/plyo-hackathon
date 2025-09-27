package routes

import (
	"net/http"
)

const (
	reportsRoutePrefix = "/reports"
	reportsNamePrefix  = "reports"
)

var ReportRoutes = []Route{
	ReportCreate,
	ReportShow,
	ReportStreamProgress,
	ReportStreamGeneration,
}

var ReportCreate = Route{
	Name:         reportsNamePrefix + ".create",
	Path:         reportsRoutePrefix,
	Method:       http.MethodPost,
	Handler:      "Reports",
	HandleMethod: "Create",
}

var ReportShow = Route{
	Name:         reportsNamePrefix + ".show",
	Path:         reportsRoutePrefix + "/:id",
	Method:       http.MethodGet,
	Handler:      "Reports",
	HandleMethod: "Show",
}

var ReportStreamProgress = Route{
	Name:         reportsNamePrefix + ".stream",
	Path:         reportsRoutePrefix + "/:id/stream",
	Method:       http.MethodGet,
	Handler:      "Reports",
	HandleMethod: "TrackReportProgress",
}

var ReportStreamGeneration = Route{
	Name:         reportsNamePrefix + ".stream-generation",
	Path:         reportsRoutePrefix + "/:id/track-generation",
	Method:       http.MethodGet,
	Handler:      "Reports",
	HandleMethod: "TrackReportGeneration",
}

