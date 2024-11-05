package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
)

type config struct {
	system      string
	host        string
	port        int
	dimensions  int
	numEntities int
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
	flag.IntVar(&cfg.dimensions, "dimensions", 1024, "Vector dimensions")
	flag.IntVar(&cfg.numEntities, "num-entities", 100000, "Number of entities")

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

	logger.Info(fmt.Sprintf("Done preparing %s", cfg.system))
}
