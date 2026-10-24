// Package main is the entrypoint for the static SPA web server.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/nawaphonOHM/resume-web-server-golang/internal/config"
	"github.com/nawaphonOHM/resume-web-server-golang/internal/handler"
	mongorepo "github.com/nawaphonOHM/resume-web-server-golang/internal/repository/mongodb"
	"github.com/nawaphonOHM/resume-web-server-golang/internal/service"
	"github.com/nawaphonOHM/whatever/v3/pkg/logging"
	"github.com/nawaphonOHM/whatever/v3/pkg/mongodb"
	"github.com/nawaphonOHM/whatever/v3/pkg/rest"
	"github.com/nawaphonOHM/whatever/v3/pkg/storage"
)

var (
	errNilMongoClient   = errors.New("mongodb client is nil")
	errNilStorageClient = errors.New("storage client is nil")
)

func validateClients(mongoClient *mongodb.Client, storageClient *storage.Client) error {
	if mongoClient == nil {
		return errNilMongoClient
	}
	if storageClient == nil {
		return errNilStorageClient
	}
	return nil
}

func buildBluePrint(
	cfg *config.Config,
	mongoClient *mongodb.Client,
	storageClient *storage.Client,
) (*rest.BluePrint, error) {
	if err := validateClients(mongoClient, storageClient); err != nil {
		return nil, err
	}
	coll := mongoClient.Collection(cfg.MongoCollection)
	repo := mongorepo.NewMongoFileRepository(coll)
	svc, err := service.NewFileStorageService(repo, storageClient, cfg.StorageBucket, cfg.IndexFile)
	if err != nil {
		return nil, fmt.Errorf("init file service: %w", err)
	}
	h := handler.NewFileHandler(svc)
	return handler.RegisterRoutes(h), nil
}

func startServer(ctx context.Context, cfg *config.Config, mongoClient *mongodb.Client) error {
	storageClient, err := storage.Connect(ctx)
	if err != nil {
		return fmt.Errorf("connect storage: %w", err)
	}
	bp, err := buildBluePrint(cfg, mongoClient, storageClient)
	if err != nil {
		return err
	}
	return rest.StartREST(bp)
}

func disconnectClient(ctx context.Context, client *mongodb.Client) {
	if client == nil {
		return
	}
	if err := client.Disconnect(ctx); err != nil {
		logging.WarnContext(ctx, "failed to disconnect mongodb client", "error", err)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	mongoClient, err := mongodb.Connect(ctx)
	if err != nil {
		return fmt.Errorf("connect mongodb: %w", err)
	}
	defer disconnectClient(ctx, mongoClient)
	return startServer(ctx, cfg, mongoClient)
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		logging.ErrorContext(ctx, "server stopped with error", "error", err)
		os.Exit(1)
	}
}
