package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/grpc"
	"github.com/weaviate/weaviate/entities/models"
)

func (app Application) prepareWeaviate() error {

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
	embeddingCol := "page_meaning"

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	exists, err := client.Schema().ClassExistenceChecker().WithClassName(className).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to check class existence -> %s", err.Error())
	}

	if exists {
		app.logger.Info("Pages class already exists, dropping")

		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete class pages -> %s", err.Error())
		}
	}

	class := &models.Class{
		Class: className,
		Properties: []*models.Property{
			{
				Name:     idCol,
				DataType: []string{"int"},
			},
			{
				Name:     urlCol,
				DataType: []string{"text"},
			},
			{
				Name:     embeddingCol,
				DataType: []string{"number"},
			},
		},
		VectorIndexType: "hnsw",
		VectorIndexConfig: map[string]interface{}{
			"efConstruction": app.cfg.efConstruction,
			"m":              app.cfg.m,
		},
	}

	app.logger.Info("Creating pages class")

	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = client.Schema().ClassCreator().WithClass(class).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to create class pages -> %s", err.Error())
	}

	app.logger.Info(fmt.Sprintf("Inserting %d random entities", app.cfg.numEntities))

	chunkSize := 200
	start := time.Now()
	batcher := client.Batch().ObjectsBatcher()

	ctx = context.Background()

	for counter := 0; counter < app.cfg.numEntities; counter += chunkSize {

		dataObjs := make([]models.PropertySchema, chunkSize)

		for i := 0; i < chunkSize; i++ {
			vec := make([]float32, app.cfg.dimensions)
			for j := 0; j < app.cfg.dimensions; j++ {
				vec[j] = rand.Float32()
			}
			dataObjs[i] = map[string]interface{}{
				idCol:        counter + i,
				urlCol:       fmt.Sprintf("%spage-%d", app.exampleSite, counter+i),
				embeddingCol: vec,
			}
		}

		for _, dataObj := range dataObjs {
			batcher.WithObjects(&models.Object{
				Class:      className,
				Properties: dataObj,
			})
		}

		_, err = batcher.Do(ctx)
		if err != nil {
			return fmt.Errorf("failed to insert objects -> %s", err.Error())
		}

		app.logger.Info(fmt.Sprintf("%.2f%% done inserting data", float64(counter+chunkSize)/float64(app.cfg.numEntities)*100))
	}

	app.logger.Info(fmt.Sprintf("Inserted %d entities in %s", app.cfg.numEntities, time.Since(start)))

	return nil
}
