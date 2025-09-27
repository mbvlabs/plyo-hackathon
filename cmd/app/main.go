package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/mbvlabs/plyo-hackathon/agents"
	"github.com/mbvlabs/plyo-hackathon/config"
	"github.com/mbvlabs/plyo-hackathon/controllers"
	"github.com/mbvlabs/plyo-hackathon/database"
	"github.com/mbvlabs/plyo-hackathon/models"
	"github.com/mbvlabs/plyo-hackathon/providers"
	"github.com/mbvlabs/plyo-hackathon/router"
	"github.com/mbvlabs/plyo-hackathon/tools"

	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
	"maragu.dev/goqite"
	"maragu.dev/goqite/jobs"
)

var appVersion string

func startServer(ctx context.Context, srv *http.Server, env string) error {
	if env == config.ProdEnvironment {
		eg, egCtx := errgroup.WithContext(ctx)

		eg.Go(func() error {
			if err := srv.ListenAndServe(); err != nil &&
				err != http.ErrServerClosed {
				return fmt.Errorf("server error: %w", err)
			}
			return nil
		})

		eg.Go(func() error {
			<-egCtx.Done()
			slog.InfoContext(ctx, "initiating graceful shutdown")
			shutdownCtx, cancel := context.WithTimeout(
				ctx,
				10*time.Second,
			)
			defer cancel()
			if err := srv.Shutdown(shutdownCtx); err != nil {
				return fmt.Errorf("shutdown error: %w", err)
			}
			return nil
		})

		if err := eg.Wait(); err != nil {
			slog.InfoContext(ctx, "wait error", "e", err)
			return err
		}

		return nil
	}

	return srv.ListenAndServe()
}

func setupControllers(
	sqlite database.SQLite,
	q *goqite.Queue,
	prelimAgent agents.PreliminaryResearch,
) (controllers.Controllers, error) {
	ctrl, err := controllers.New(
		prelimAgent,
		sqlite,
		q,
	)
	if err != nil {
		return controllers.Controllers{}, err
	}

	return ctrl, nil
}

func setupRouter(ctrl controllers.Controllers) (*echo.Echo, error) {
	router, err := router.New(
		ctrl,
	)
	if err != nil {
		return nil, err
	}

	return router.SetupRoutes(), nil
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	sqlite, err := database.NewSQLite(ctx)
	if err != nil {
		return err
	}

	q := goqite.New(goqite.NewOpts{
		DB:   sqlite.Conn(),
		Name: "jobs",
	})

	r := jobs.NewRunner(jobs.NewRunnerOpts{
		Limit:        10,
		Log:          slog.Default(),
		PollInterval: 10 * time.Millisecond,
		Queue:        q,
	})

	serper := tools.NewSerper(config.App.SerperAPIkey)
	openai := providers.NewClient(config.App.OpenAPIKey)
	toolsMap := map[string]tools.Tooler{serper.GetName(): &serper}

	// Create agents
	companyIntel := agents.NewCompanyIntelligence(openai, toolsMap)
	competitiveIntel := agents.NewCompetitiveIntelligence(openai, toolsMap)
	marketDynamics := agents.NewMarketDynamics(openai, toolsMap)
	trendAnalysis := agents.NewTrendAnalysis(openai, toolsMap)
	reportGenerator := agents.NewReportGenerator(openai, nil)

	r.Register(agents.ReportGeneratorJobName, func(ctx context.Context, m []byte) error {
		var params agents.ReportGeneratorJobParams
		if err := json.Unmarshal(m, &params); err != nil {
			slog.ErrorContext(ctx, "failed to unmarshal job params", "error", err)
			return err
		}

		report, err := models.FindReport(ctx, sqlite.Conn(), params.ReportID)
		if err != nil {
			return err
		}

		slog.InfoContext(
			ctx,
			"####################### starting report generation",
			"report_id",
			params.ReportID,
		)
		result, err := reportGenerator.Generate(
			ctx,
			params.CandidateName,
			params.CompanyURL,
			report.CompanyIntelligenceData,
			report.CompetitiveIntelligenceData,
			report.MarketDynamicsData,
			report.TrendAnalysisData,
		)
		if err != nil {
			slog.ErrorContext(ctx, "research failed", "error", err)
			return err
		}

		if err := models.UpdateFinalReport(ctx, sqlite.Conn(), params.ReportID, result); err != nil {
			slog.ErrorContext(ctx, "failed to update final report", "error", err)
			return err
		}

		slog.InfoContext(ctx, "completed report", "report_id", params.ReportID)
		return nil
	})
	r.Register(agents.TrendAnalysisJobName, func(ctx context.Context, m []byte) error {
		var params agents.TrendAnalysisJobParams
		if err := json.Unmarshal(m, &params); err != nil {
			slog.ErrorContext(ctx, "failed to unmarshal job params", "error", err)
			return err
		}

		slog.InfoContext(
			ctx,
			"####################### starting trend analysis",
			"report_id",
			params.ReportID,
		)
		result, err := trendAnalysis.Research(ctx, params.CandidateName, params.CompanyURL)
		if err != nil {
			slog.ErrorContext(ctx, "research failed", "error", err)
			return err
		}

		slog.InfoContext(
			ctx,
			"####################### result of trend analysis",
			"result",
			result,
		)

		if err := models.UpdateTrendAnalysis(ctx, sqlite.Conn(), params.ReportID, result); err != nil {
			slog.ErrorContext(ctx, "failed to update trend analysis", "error", err)
			return err
		}

		if err := models.UpdateReportProgress(ctx, sqlite.Conn(), params.ReportID); err != nil {
			slog.ErrorContext(ctx, "failed to update report progress", "error", err)
			return err
		}
		slog.InfoContext(ctx, "completed update trend analysis", "report_id", params.ReportID)
		return nil
	})
	r.Register(agents.MarketDynamicsJobName, func(ctx context.Context, m []byte) error {
		var params agents.MarketDynamicsJobParams
		if err := json.Unmarshal(m, &params); err != nil {
			slog.ErrorContext(ctx, "failed to unmarshal job params", "error", err)
			return err
		}

		slog.InfoContext(
			ctx,
			"####################### starting market dynamics",
			"report_id",
			params.ReportID,
		)
		result, err := marketDynamics.Research(ctx, params.CandidateName, params.CompanyURL)
		if err != nil {
			slog.ErrorContext(ctx, "research failed", "error", err)
			return err
		}

		slog.InfoContext(
			ctx,
			"####################### result of market dynamics",
			"result",
			result,
		)

		if err := models.UpdateMarketDynamics(ctx, sqlite.Conn(), params.ReportID, result); err != nil {
			slog.ErrorContext(ctx, "failed to update market dynamics", "error", err)
			return err
		}

		if err := models.UpdateReportProgress(ctx, sqlite.Conn(), params.ReportID); err != nil {
			slog.ErrorContext(ctx, "failed to update report progress", "error", err)
			return err
		}
		slog.InfoContext(ctx, "completed update market dynamics", "report_id", params.ReportID)
		return nil
	})
	r.Register(agents.CompetitiveIntelligenceJobName, func(ctx context.Context, m []byte) error {
		var params agents.CompanyIntelligenceJobParams
		if err := json.Unmarshal(m, &params); err != nil {
			slog.ErrorContext(ctx, "failed to unmarshal job params", "error", err)
			return err
		}

		slog.InfoContext(
			ctx,
			"####################### starting competitive intelligence",
			"report_id",
			params.ReportID,
		)
		result, err := competitiveIntel.Research(ctx, params.CandidateName, params.CompanyURL)
		if err != nil {
			slog.ErrorContext(ctx, "research failed", "error", err)
			return err
		}

		slog.InfoContext(
			ctx,
			"####################### result of competitive intelligence",
			"result",
			result,
		)

		if err := models.UpdateCompetitiveIntelligence(ctx, sqlite.Conn(), params.ReportID, result); err != nil {
			slog.ErrorContext(ctx, "failed to update competitive intelligence", "error", err)
			return err
		}

		if err := models.UpdateReportProgress(ctx, sqlite.Conn(), params.ReportID); err != nil {
			slog.ErrorContext(ctx, "failed to update report progress", "error", err)
			return err
		}
		slog.InfoContext(ctx, "completed competitive intelligence", "report_id", params.ReportID)
		return nil
	})
	r.Register(agents.CompanyIntelligenceJobName, func(ctx context.Context, m []byte) error {
		var params agents.CompanyIntelligenceJobParams
		if err := json.Unmarshal(m, &params); err != nil {
			slog.ErrorContext(ctx, "failed to unmarshal job params", "error", err)
			return err
		}

		slog.InfoContext(
			ctx,
			"####################### starting company intelligence",
			"report_id",
			params.ReportID,
		)
		result, err := companyIntel.Research(ctx, params.CandidateName, params.CompanyURL)
		if err != nil {
			slog.ErrorContext(ctx, "research failed", "error", err)
			return err
		}

		slog.InfoContext(
			ctx,
			"####################### result of company intelligence",
			"result",
			result,
		)

		if err := models.UpdateCompanyIntelligence(ctx, sqlite.Conn(), params.ReportID, result); err != nil {
			slog.ErrorContext(ctx, "failed to update company intelligence", "error", err)
			return err
		}

		if err := models.UpdateReportProgress(ctx, sqlite.Conn(), params.ReportID); err != nil {
			slog.ErrorContext(ctx, "failed to update report progress", "error", err)
			return err
		}
		slog.InfoContext(ctx, "completed company intelligence", "report_id", params.ReportID)
		return nil
	})

	// r.Register("CompanyIntelligenceJobName", func(ctx context.Context, m []byte) error {
	// 	slog.InfoContext(ctx, "starting company intelligence", "report_id", report.ID)
	//
	// 	result, err := companyIntel.Research(ctx, candidate.Name, companyURL)
	//
	// 	if err := models.UpdateCompanyIntelligence(ctx, r.db.Conn(), report.ID, result); err != nil {
	// 		slog.ErrorContext(ctx, "failed to update company intelligence", "error", err)
	// 		companyDone <- err
	// 		return
	// 	}
	// 	models.UpdateReportProgress(ctx, r.db.Conn(), report.ID)
	// 	slog.InfoContext(ctx, "completed company intelligence", "report_id", report.ID)
	// 	return nil
	// })

	// Start the job runner and see the job run.
	go func() {
		r.Start(ctx)
	}()

	// Create prelim agent
	prelimAgent := agents.NewPreliminaryResearch(
		openai,
		map[string]tools.Tooler{serper.GetName(): &serper},
	)
	controllers, err := setupControllers(sqlite, q, prelimAgent)
	if err != nil {
		return err
	}

	handler, err := setupRouter(controllers)
	if err != nil {
		return err
	}

	port := config.App.ServerPort
	host := config.App.ServerHost

	srv := &http.Server{
		Addr:         fmt.Sprintf("%v:%v", host, port),
		Handler:      handler,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 10 * time.Minute,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
	}

	slog.InfoContext(ctx, "starting server", "host", host, "port", port)
	return startServer(ctx, srv, config.App.Env)
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
