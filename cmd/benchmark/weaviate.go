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

func (app Application) benchmarkWeaviate() (results [][]string, err error) {

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
		return results, fmt.Errorf("failed to create weaviate client -> %s", err.Error())
	}

	err = client.WaitForWeavaite(1 * time.Second)
	if err != nil {
		return results, fmt.Errorf("failed to connect to weaviate -> %s", err.Error())
	}

	className := "Pages"
	idCol := "idizzle"
	urlCol := "url"

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	exists, err := client.Schema().ClassExistenceChecker().WithClassName(className).Do(ctx)
	if err != nil {
		return results, fmt.Errorf("failed to check class existence -> %s", err.Error())
	}

	if !exists {
		return results, fmt.Errorf("pages class does not exist")
	}

	app.logger.Info(fmt.Sprintf("Starting benchmark of %v", className))

	eg := errgroup.Group{}
	eg.SetLimit(app.cfg.concurrency)

	var queryCounter int32

	start := time.Now()

	lastPrintlinesCheckTime := time.Now()
	var lastPrintlinesNumQueries int32 = 0

	lastResultsCheckTime := time.Now()
	var lastResultsNumQueries int32 = 0

	for time.Since(start) < app.cfg.duration {

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

		// printlines for development
		if time.Since(lastPrintlinesCheckTime) > app.cfg.logFrequency {
			currentQueries := atomic.LoadInt32(&queryCounter)
			app.logger.Info(fmt.Sprintf("%v completing %.0f actions per second", app.cfg.system, float64(currentQueries-lastPrintlinesNumQueries)/time.Since(lastPrintlinesCheckTime).Seconds()))
			lastPrintlinesCheckTime = time.Now()
			lastPrintlinesNumQueries = queryCounter
		}

		// log results to app.results
		if time.Since(lastResultsCheckTime) > 1*time.Minute {
			currentQueries := atomic.LoadInt32(&queryCounter)
			numQueriesSinceLastCheck := currentQueries - lastResultsNumQueries

			roundedMinutes := int(time.Since(start).Minutes())

			results = append(results, []string{fmt.Sprint(roundedMinutes), fmt.Sprint(numQueriesSinceLastCheck)})

			lastResultsCheckTime = time.Now()
			lastResultsNumQueries = queryCounter
		}
	}

	err = eg.Wait()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}

	// log results to app.results
	if time.Since(lastResultsCheckTime) > 1*time.Minute {
		currentQueries := atomic.LoadInt32(&queryCounter)
		numQueriesSinceLastCheck := currentQueries - lastResultsNumQueries

		roundedMinutes := int(time.Since(start).Minutes())

		results = append(results, []string{fmt.Sprint(roundedMinutes), fmt.Sprint(numQueriesSinceLastCheck)})
	}

	app.logger.Info(fmt.Sprintf("%v completed %v actions in %v, rate of %.0f per second", app.cfg.system, queryCounter, app.cfg.duration, float64(queryCounter)/time.Since(start).Seconds()))

	return results, nil
}
