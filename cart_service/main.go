package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/go-redis/redis/v8"
	"github.com/kanowfy/ecom/cart_service/internal/config"
	"github.com/kanowfy/ecom/cart_service/internal/log"
	"github.com/kanowfy/ecom/cart_service/internal/repository"
	"github.com/kanowfy/ecom/cart_service/internal/service"
	"github.com/kanowfy/ecom/cart_service/pb"
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
	flag.StringVar(&cfg.DB.Addr, "db.addr", cfg.DB.Addr, "database address")
	flag.StringVar(&cfg.DB.Password, "db.password", cfg.DB.Password, "database password")
	flag.IntVar(&cfg.DB.Index, "db.index", cfg.DB.Index, "database index")

	var level = slog.LevelDebug

	flag.Func("loglevel", "minimum log level", func(s string) error {
		if err := level.UnmarshalText([]byte(s)); err != nil {
			return err
		}

		return nil
	})

	flag.Parse()

	logger := log.New(os.Stdout, level, true)
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.DB.Addr,
		Password: cfg.DB.Password,
		DB:       cfg.DB.Index,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logger.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}

	repo := repository.New(rdb)
	service := service.New(logger, repo)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
	if err != nil {
		logger.Error("failed to announce address", "error", err.Error())
		os.Exit(1)
	}

	s := grpc.NewServer()
	pb.RegisterCartServer(s, service)
	logger.Info(fmt.Sprintf("gRPC server listening on %s:%d", cfg.Server.Host, cfg.Server.Port))
	if err := s.Serve(lis); err != nil {
		logger.Error("failed to serve", "error", err.Error())
		os.Exit(1)
	}
}
