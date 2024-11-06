package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"sync/atomic"
	"time"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/grpc"
	"golang.org/x/sync/errgroup"
)

func (app Application) benchmarkWeaviate() error {

	if app.cfg.port == 0 {
		app.cfg.port = 8080
	}

	if app.cfg.grpcPort == 0 {
		app.cfg.grpcPort = 50051
	}

	app.logger.Info("Connecting to Weaviate")

	client, err := weaviate.NewClient(weaviate.Config{
		Host:   fmt.Sprintf("%s:%d", app.cfg.host, app.cfg.port),
		Scheme: "http",
		GrpcConfig: &grpc.Config{
			Host:    fmt.Sprintf("%s:%d", app.cfg.host, app.cfg.grpcPort),
			Secured: false,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create weaviate client -> %s", err.Error())
	}

	err = client.WaitForWeavaite(1 * time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to weaviate -> %s", err.Error())
	}

	className := "Pages"
	idCol := "idizzle"
	urlCol := "url"

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	exists, err := client.Schema().ClassExistenceChecker().WithClassName(className).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to check class existence -> %s", err.Error())
	}

	if !exists {
		return fmt.Errorf("pages class does not exist")
	}

	app.logger.Info(fmt.Sprintf("Starting benchmark of %v", className))

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

			vector := make([]float32, app.cfg.dimensions)

			for i := 0; i < app.cfg.dimensions; i++ {
				vector[i] = rand.Float32()
			}

			ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)

			_, err = client.GraphQL().Get().
				WithClassName(className).
				WithFields(
					graphql.Field{Name: idCol},
					graphql.Field{Name: urlCol},
					graphql.Field{
						Name: "_additional",
						Fields: []graphql.Field{
							{Name: "distance"},
						},
					},
				).
				WithNearVector(client.GraphQL().NearVectorArgBuilder().
					WithVector(vector)).
				WithLimit(app.cfg.topK).
				Do(ctx)

			if err != nil {
				return fmt.Errorf("failed to query weaviate -> %s", err.Error())
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
