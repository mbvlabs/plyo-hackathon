package main

import (
	"context"
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
	"github.com/mbvlabs/plyo-hackathon/providers"
	"github.com/mbvlabs/plyo-hackathon/router"
	"github.com/mbvlabs/plyo-hackathon/tools"

	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
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

func setupControllers(sqlite database.SQLite) (controllers.Controllers, error) {
	serper := tools.NewSerper(config.App.SerperAPIkey)
	openai := providers.NewClient(config.App.OpenAPIKey)

	// Create prelim agent
	prelimAgent := agents.NewPreliminaryResearch(
		openai,
		map[string]tools.Tooler{serper.GetName(): &serper},
	)
	ctrl, err := controllers.New(
		prelimAgent,
		sqlite,
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

	controllers, err := setupControllers(sqlite)
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
