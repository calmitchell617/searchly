package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"sync/atomic"
	"time"

	"github.com/qdrant/go-client/qdrant"
	"golang.org/x/sync/errgroup"
)

func (app Application) benchmarkQdrant() error {

	if app.cfg.port == 0 {
		app.cfg.port = 6334
	}

	app.logger.Info("Connecting to Qdrant")

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

	collectionName := "pages"

	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	collectionExists, err := client.CollectionExists(ctx, collectionName)
	if err != nil {
		return fmt.Errorf("failed to check collection existence -> %s", err.Error())
	}

	if !collectionExists {
		return fmt.Errorf("pages collection does not exist")
	}

	app.logger.Info(fmt.Sprintf("Starting benchmark of %v", collectionName))

	eg := errgroup.Group{}
	eg.SetLimit(app.cfg.concurrency)

	var queryCounter int32

	start := time.Now()

	lastCheckTime := time.Now()
	var lastQueries int32 = 0

	for time.Since(start) < app.cfg.duration {

		if time.Since(lastCheckTime) > 3*time.Second {
			currentQueries := atomic.LoadInt32(&queryCounter)
			app.logger.Info(fmt.Sprintf("%v completing %.0f actions per second", app.cfg.system, float64(currentQueries-lastQueries)/time.Since(lastCheckTime).Seconds()))
			lastCheckTime = time.Now()
			lastQueries = queryCounter
		}

		eg.Go(func() error {

			embeddingList := make([]float32, app.cfg.dimensions)

			for i := 0; i < app.cfg.dimensions; i++ {
				embeddingList[i] = rand.Float32()
			}

			ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			_, err = client.Query(context.Background(), &qdrant.QueryPoints{
				CollectionName: collectionName,
				Query:          qdrant.NewQuery(embeddingList...),
			})
			if err != nil {
				return fmt.Errorf("failed to query collection -> %s", err.Error())
			}

			atomic.AddInt32(&queryCounter, 1)
			return nil
		})
	}

	err = eg.Wait()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}

	app.logger.Info(fmt.Sprintf("%v completed %v actions in %v, rate of %.0f per second", app.cfg.system, queryCounter, app.cfg.duration, float64(queryCounter)/time.Since(start).Seconds()))

	return nil
}
