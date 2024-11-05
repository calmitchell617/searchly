package main

import (
	"context"
	"fmt"
	"time"

	"github.com/qdrant/go-client/qdrant"
)

func (app Application) prepareQdrant() error {

	if app.cfg.port == 0 {
		app.cfg.port = 6334
	}

	client, err := qdrant.NewClient(&qdrant.Config{
		Host: app.cfg.host,
		Port: app.cfg.port,
	})
	if err != nil {
		return fmt.Errorf("failed to create qdrant client -> %s", err.Error())
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err = client.HealthCheck(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to qdrant -> %s", err.Error())
	}

	app.logger.Info("Connected to Qdrant")

	return nil
}
