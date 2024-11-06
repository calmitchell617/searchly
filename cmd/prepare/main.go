package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"
)

type config struct {
	system         string
	host           string
	port           int
	grpcPort       int
	dimensions     int
	numEntities    int
	efConstruction int
	m              int
}

type Application struct {
	logger      *slog.Logger
	cfg         config
	exampleSite string
}

func main() {

	cfg := config{}

	flag.StringVar(&cfg.system, "system", "", "System type")
	flag.StringVar(&cfg.host, "host", "localhost", "Hostname")
	flag.IntVar(&cfg.port, "port", 0, "Port")
	flag.IntVar(&cfg.grpcPort, "grpc-port", 0, "gRPC port")
	flag.IntVar(&cfg.dimensions, "dimensions", 1024, "Vector dimensions")
	flag.IntVar(&cfg.numEntities, "num-entities", 100000, "Number of entities")
	flag.IntVar(&cfg.efConstruction, "ef-construction", 64, "efConstruction")
	flag.IntVar(&cfg.m, "m", 16, "m")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if cfg.system == "" {
		logger.Error("System flag not provided")
		os.Exit(1)
	}

	app := Application{
		logger:      logger,
		cfg:         cfg,
		exampleSite: "http://example.com/",
	}

	var err error
	start := time.Now()

	switch cfg.system {
	case "milvus":
		err = app.prepareMilvus()
	case "weaviate":
		err = app.prepareWeaviate()
	case "qdrant":
		err = app.prepareQdrant()
	default:
		logger.Error("Invalid system type")
		os.Exit(1)
	}

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info(fmt.Sprintf("Done preparing %s in %v", cfg.system, time.Since(start)))
}
