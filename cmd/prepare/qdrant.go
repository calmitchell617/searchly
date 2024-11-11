package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/qdrant/go-client/qdrant"
)

func (app Application) prepareQdrant() error {

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

	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if collectionExists {
		app.logger.Info("Pages collection already exists, dropping")
		err = client.DeleteCollection(ctx, collectionName)
		if err != nil {
			return fmt.Errorf("failed to delete collection -> %s", err.Error())
		}
	}

	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	app.logger.Info("Creating pages collection")

	efConstruct := uint64(app.cfg.efConstruction)
	m := uint64(app.cfg.m)
	onDisk := false

	client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     uint64(app.cfg.dimensions),
			Distance: qdrant.Distance_Cosine,
			HnswConfig: &qdrant.HnswConfigDiff{
				EfConstruct: &efConstruct,
				M:           &m,
				OnDisk:      &onDisk,
			},
		}),
	})

	app.logger.Info(fmt.Sprintf("Inserting %d random entities", app.cfg.numEntities))

	chunkSize := 10000
	start := time.Now()

	for counter := 0; counter < app.cfg.numEntities; counter += chunkSize {

		embeddingList := make([]*qdrant.Vectors, chunkSize)
		points := make([]*qdrant.PointStruct, chunkSize)

		for i := 0; i < chunkSize; i++ {
			id := int64(counter + i)
			url := fmt.Sprintf("%spage-%d", app.exampleSite, counter+i)
			vec := make([]float32, app.cfg.dimensions)

			for j := 0; j < app.cfg.dimensions; j++ {
				vec[j] = rand.Float32()
			}

			embeddingList[i] = qdrant.NewVectors(vec...)

			points[i] = &qdrant.PointStruct{
				Id:      qdrant.NewIDNum(uint64(id)),
				Vectors: embeddingList[i],
				Payload: qdrant.NewValueMap(map[string]any{"url": url}),
			}
		}

		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		_, err = client.Upsert(ctx, &qdrant.UpsertPoints{
			CollectionName: collectionName,
			Points:         points,
		})
		if err != nil {
			return fmt.Errorf("failed to upsert points -> %s", err.Error())
		}

		app.logger.Info(fmt.Sprintf("%.2f%% done inserting data", float64(counter+chunkSize)/float64(app.cfg.numEntities)*100))
	}

	app.logger.Info(fmt.Sprintf("Inserted %d entities in %v", app.cfg.numEntities, time.Since(start)))

	return nil
}
