package controllers

import (
	"github.com/mbvlabs/plyo-hackathon/database"
	"maragu.dev/goqite"

	"github.com/mbvlabs/plyo-hackathon/views"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/maypok86/otter"
)

type Pages struct {
	db    database.SQLite
	q     *goqite.Queue
	cache otter.CacheWithVariableTTL[string, templ.Component]
}

func newPages(
	db database.SQLite,
	q *goqite.Queue,
	cache otter.CacheWithVariableTTL[string, templ.Component],
) Pages {
	return Pages{db, q, cache}
}

func (p Pages) Home(c echo.Context) error {
	return render(c, views.Home())
}

func (p Pages) NotFound(c echo.Context) error {
	return render(c, views.NotFound())
}
