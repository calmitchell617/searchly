package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"sync/atomic"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"golang.org/x/sync/errgroup"

	milvus "github.com/milvus-io/milvus-sdk-go/v2/client"
)

func (app Application) benchmarkMilvus() (results [][]string, err error) {

	if app.cfg.port == 0 {
		app.cfg.port = 19530
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	app.logger.Info("Connecting to Milvus")

	client, err := milvus.NewClient(ctx, milvus.Config{
		Address: fmt.Sprintf("%s:%d", app.cfg.host, app.cfg.port),
	})
	if err != nil {
		return results, fmt.Errorf("failed to create milvus client -> %s", err.Error())
	}
	defer client.Close()

	collectionName := "pages"
	urlCol := "url"
	embeddingCol := "page_meaning"

	collExists, err := client.HasCollection(ctx, collectionName)
	if err != nil {
		return results, fmt.Errorf("failed to check collection existence -> %s", err.Error())
	}
	if !collExists {
		return results, fmt.Errorf("pages collection does not exist")
	}

	app.logger.Info("Loading pages collection into RAM")

	err = client.LoadCollection(context.Background(), collectionName, false)
	if err != nil {
		return results, fmt.Errorf("failed to load collection -> %s", err.Error())
	}

	app.logger.Info(fmt.Sprintf("Starting benchmark of %v", collectionName))

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

			embeddingList := make([]float32, app.cfg.dimensions)

			for i := 0; i < app.cfg.dimensions; i++ {
				embeddingList[i] = rand.Float32()
			}

			vec2search := []entity.Vector{
				entity.FloatVector(embeddingList),
			}

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			sp, err := entity.NewIndexIvfFlatSearchParam(16)
			if err != nil {
				return fmt.Errorf("failed to create search param -> %s", err.Error())
			}
			_, err = client.Search(ctx, collectionName, nil, "", []string{urlCol}, vec2search,
				embeddingCol, entity.COSINE, app.cfg.topK, sp)
			if err != nil {
				return fmt.Errorf("failed to search -> %s", err.Error())
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
