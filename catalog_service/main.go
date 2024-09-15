package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kanowfy/ecom/catalog_service/internal/config"
	"github.com/kanowfy/ecom/catalog_service/internal/log"
	"github.com/kanowfy/ecom/catalog_service/internal/repository"
	"github.com/kanowfy/ecom/catalog_service/internal/service"
	"github.com/kanowfy/ecom/catalog_service/pb"
	"google.golang.org/grpc"
)

func main() {
	var cfg config.Config

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("failed to parse environment variable: %v", err)
		os.Exit(1)
	}

	flag.StringVar(&cfg.Server.Host, "srv.host", cfg.Server.Host, "server host")
	flag.IntVar(&cfg.Server.Port, "srv.port", cfg.Server.Port, "server port")
	flag.StringVar(&cfg.DB.Url, "db.url", cfg.DB.Url, "database connection string")

	var level = slog.LevelDebug

	flag.Func("loglevel", "minimum log level", func(s string) error {
		if err := level.UnmarshalText([]byte(s)); err != nil {
			return err
		}

		return nil
	})

	flag.Parse()

	logger := log.New(os.Stdout, level, true)
	pool, err := pgxpool.New(context.Background(), cfg.DB.Url)
	if err != nil {
		logger.Error("failed to obtain connection pool", "error", err)
		os.Exit(1)
	}

	if err := pool.Ping(context.Background()); err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	repo := repository.New(pool)
	service := service.New(logger, repo)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
	if err != nil {
		logger.Error("failed to announce address", "error", err.Error())
		os.Exit(1)
	}

	s := grpc.NewServer()
	pb.RegisterCatalogServer(s, service)
	logger.Info(fmt.Sprintf("gRPC server listening on %s:%d", cfg.Server.Host, cfg.Server.Port))
	if err := s.Serve(lis); err != nil {
		logger.Error("failed to serve", "error", err.Error())
		os.Exit(1)
	}
}
