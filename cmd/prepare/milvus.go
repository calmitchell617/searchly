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

	schema := entity.NewSchema().WithName(collectionName).
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
	start := time.Now()

	for counter := 0; counter < app.cfg.numEntities; counter += chunkSize {

		idList := make([]int64, chunkSize)
		urlList := make([]string, chunkSize)
		embeddingList := make([][]float32, chunkSize)

		for i := 0; i < chunkSize; i++ {
			idList[i] = int64(counter + i)
			urlList[i] = fmt.Sprintf("%spage-%d", app.exampleSite, counter+i)
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

	app.logger.Info(fmt.Sprintf("Creating index on %s. This takes 5X-10X as long as insertions", embeddingCol))
	start = time.Now()

	idx, err := entity.NewIndexDISKANN(entity.COSINE)
	if err != nil {
		return fmt.Errorf("failed to create index, err: %v", err)
	}
	if err := client.CreateIndex(ctx, collectionName, embeddingCol, idx, false); err != nil {
		return fmt.Errorf("failed to create index on %s, err: %v", embeddingCol, err)
	}

	app.logger.Info(fmt.Sprintf("Index created in %v", time.Since(start)))

	return nil
}
