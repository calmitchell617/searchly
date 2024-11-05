package main

import (
	"context"
	"fmt"
	"time"

	"math/rand"

	milvus "github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

func (app Application) prepareMilvus() error {

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
		return fmt.Errorf("failed to create milvus client -> %s", err.Error())
	}
	defer client.Close()

	collectionName := "pages"
	idCol := "id"
	urlCol := "url"
	embeddingCol := "page_meaning"

	collExists, err := client.HasCollection(ctx, collectionName)
	if err != nil {
		return fmt.Errorf("failed to check collection existence -> %s", err.Error())
	}
	if collExists {
		app.logger.Info("Pages collection already exists, dropping")
		_ = client.DropCollection(ctx, collectionName)
	}

	schema := entity.NewSchema().WithName(collectionName).WithDescription("Web pages to semantically search").
		WithField(entity.NewField().WithName(idCol).WithDataType(entity.FieldTypeInt64).WithIsPrimaryKey(true).WithIsAutoID(false)).
		WithField(entity.NewField().WithName(urlCol).WithDataType(entity.FieldTypeVarChar).WithMaxLength(256)).
		WithField(entity.NewField().WithName(embeddingCol).WithDataType(entity.FieldTypeFloatVector).WithDim(int64(app.cfg.dimensions)))

	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	app.logger.Info("Creating pages collection")

	err = client.CreateCollection(ctx, schema, entity.DefaultShardNumber)
	if err != nil {
		return fmt.Errorf("failed to create collection -> %s", err.Error())
	}

	ctx = context.Background()

	app.logger.Info(fmt.Sprintf("Inserting %d random entities", app.cfg.numEntities))

	chunkSize := 10000
	mValue := 16
	efConstruction := 64
	start := time.Now()
	exampleSite := "http://example.com/"

	for counter := 0; counter < app.cfg.numEntities; counter += chunkSize {

		idList := make([]int64, chunkSize)
		urlList := make([]string, chunkSize)
		embeddingList := make([][]float32, chunkSize)

		for i := 0; i < chunkSize; i++ {
			idList[i] = int64(counter + i)
			urlList[i] = exampleSite + fmt.Sprintf("page-%d", counter+i)
			vec := make([]float32, app.cfg.dimensions)
			for j := 0; j < app.cfg.dimensions; j++ {
				vec[j] = rand.Float32()
			}
			embeddingList[i] = vec
		}

		idColData := entity.NewColumnInt64(idCol, idList)
		urlColData := entity.NewColumnVarChar(urlCol, urlList)
		embeddingColData := entity.NewColumnFloatVector(embeddingCol, app.cfg.dimensions, embeddingList)

		if _, err := client.Insert(ctx, collectionName, "", idColData, urlColData, embeddingColData); err != nil {
			return fmt.Errorf("failed to insert random data into `%v, err: %v", collectionName, err)
		}

		app.logger.Info(fmt.Sprintf("%.2f%% done inserting data", float64(counter+chunkSize)/float64(app.cfg.numEntities)*100))
	}

	app.logger.Info(fmt.Sprintf("Data inserted in %v", time.Since(start)))

	app.logger.Info("Flushing data")
	start = time.Now()

	if err := client.Flush(ctx, collectionName, false); err != nil {
		return fmt.Errorf("failed to flush data into %v, err: %v", collectionName, err)
	}

	app.logger.Info(fmt.Sprintf("Data flushed in %v", time.Since(start)))

	app.logger.Info(fmt.Sprintf("Creating index on %s. This could take a while.", embeddingCol))
	start = time.Now()

	idx, err := entity.NewIndexHNSW(entity.L2, mValue, efConstruction)
	if err != nil {
		return fmt.Errorf("failed to create index, err: %v", err)
	}
	if err := client.CreateIndex(ctx, collectionName, embeddingCol, idx, false); err != nil {
		return fmt.Errorf("failed to create index on %s, err: %v", embeddingCol, err)
	}

	app.logger.Info(fmt.Sprintf("Index created in %v", time.Since(start)))

	return nil
}
