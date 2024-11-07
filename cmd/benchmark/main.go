package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"
)

type config struct {
	system       string
	host         string
	port         int
	grpcPort     int
	duration     time.Duration
	dimensions   int
	concurrency  int
	topK         int
	logFrequency time.Duration
}

type Application struct {
	logger *slog.Logger
	cfg    config
}

func main() {

	cfg := config{}

	flag.StringVar(&cfg.system, "system", "", "System type")
	flag.StringVar(&cfg.host, "host", "localhost", "Hostname")
	flag.IntVar(&cfg.port, "port", 0, "Port")
	flag.IntVar(&cfg.grpcPort, "grpc-port", 0, "gRPC port")
	flag.IntVar(&cfg.concurrency, "concurrency", 64, "Concurrency")
	flag.IntVar(&cfg.dimensions, "dimensions", 1024, "Dimensions")
	flag.IntVar(&cfg.topK, "topK", 10, "TopK")
	flag.DurationVar(&cfg.duration, "duration", 1*time.Minute, "Duration")
	flag.DurationVar(&cfg.logFrequency, "log-frequency", 10*time.Second, "Log frequency")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if cfg.system == "" {
		logger.Error("System flag not provided")
		os.Exit(1)
	}

	app := Application{
		logger: logger,
		cfg:    cfg,
	}

	var err error
	var results [][]string

	switch cfg.system {
	case "milvus":
		results, err = app.benchmarkMilvus()
	case "weaviate":
		results, err = app.benchmarkWeaviate()
	case "qdrant":
		results, err = app.benchmarkQdrant()
	default:
		logger.Error("Invalid system type")
		os.Exit(1)
	}

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	file, err := os.Create(fmt.Sprintf("results/%v_%v.csv", cfg.system, time.Now().Unix()))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create CSV file: %v", err))
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"minute", "num_queries"})

	for _, row := range results {
		err = writer.Write(row)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to write row to CSV file: %v", err))
			os.Exit(1)
		}
	}

	logger.Info(fmt.Sprintf("Done benchmarking %v, results written to %v", cfg.system, file.Name()))
}
