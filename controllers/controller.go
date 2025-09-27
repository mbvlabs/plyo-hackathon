package controllers

import (
	"context"
	"net/http"

	"github.com/mbvlabs/plyo-hackathon/agents"
	"github.com/mbvlabs/plyo-hackathon/database"
	"github.com/mbvlabs/plyo-hackathon/router/cookies"
	"github.com/starfederation/datastar-go/datastar"
	"maragu.dev/goqite"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/maypok86/otter"
)

type Controllers struct {
	Assets         Assets
	API            API
	Pages          Pages
	ResearchBriefs ResearchBriefs
	Reports        Reports
}

func New(
	prelimAgent agents.PreliminaryResearch,
	db database.SQLite,
	q *goqite.Queue,
) (Controllers, error) {
	cacheBuilder, err := otter.NewBuilder[string, templ.Component](20)
	if err != nil {
		return Controllers{}, err
	}

	pageCacher, err := cacheBuilder.WithVariableTTL().Build()
	if err != nil {
		return Controllers{}, err
	}

	assets := newAssets()
	pages := newPages(db, q, pageCacher)
	api := newAPI(db)
	researchbriefs := newResearchBriefs(prelimAgent, db)
	reports := newReports(db, q)

	return Controllers{
		assets,
		api,
		pages,
		researchbriefs,
		reports,
	}, nil
}

func render(ctx echo.Context, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	appCtx := ctx.Get(string(cookies.AppKey))
	withAppCtx := context.WithValue(
		ctx.Request().Context(),
		cookies.AppKey,
		appCtx,
	)

	flashCtx := ctx.Get(string(cookies.FlashKey))
	withFlashCtx := context.WithValue(
		withAppCtx,
		cookies.FlashKey,
		flashCtx,
	)

	if err := t.Render(withFlashCtx, buf); err != nil {
		return err
	}

	return ctx.HTML(http.StatusOK, buf.String())
}

func getSSE(c echo.Context) *datastar.ServerSentEventGenerator {
	return datastar.NewSSE(c.Response(), c.Request())
}
