package main

import (
	"time"

	"github.com/aopoltorzhicky/bcdhub/internal/helpers"
	"github.com/aopoltorzhicky/bcdhub/internal/jsonload"
	"github.com/aopoltorzhicky/bcdhub/internal/logger"
)

func main() {
	var cfg config
	if err := jsonload.StructFromFile("config.json", &cfg); err != nil {
		logger.Fatal(err)
	}
	cfg.print()

	helpers.InitSentry(cfg.Sentry.Env, cfg.Sentry.DSN, cfg.Sentry.Debug)
	helpers.SetTagSentry("project", cfg.Sentry.Project)
	defer helpers.CatchPanicSentry()

	ctx, err := newContext(cfg)
	if err != nil {
		logger.Error(err)
		helpers.CatchErrorSentry(err)
		return
	}

	// Initial syncronization
	if err = process(ctx); err != nil {
		logger.Error(err)
		helpers.CatchErrorSentry(err)
	}

	// Update state by ticker
	ticker := time.NewTicker(time.Duration(cfg.UpdateTimer) * time.Second)
	for range ticker.C {
		if err = process(ctx); err != nil {
			logger.Error(err)
			helpers.CatchErrorSentry(err)
		}
	}
}