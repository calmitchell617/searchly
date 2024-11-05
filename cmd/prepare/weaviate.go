package main

import (
	"fmt"
	"time"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

func (app Application) prepareWeaviate() error {

	if app.cfg.port == 0 {
		app.cfg.port = 8080
	}

	client, err := weaviate.NewClient(weaviate.Config{
		Host:   fmt.Sprintf("%s:%d", app.cfg.host, app.cfg.port),
		Scheme: "http",
	})
	if err != nil {
		return fmt.Errorf("failed to create weaviate client -> %s", err.Error())
	}

	err = client.WaitForWeavaite(1 * time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to weaviate -> %s", err.Error())
	}

	app.logger.Info("Connected to Weaviate")

	return nil
}
