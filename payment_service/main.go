package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/kanowfy/ecom/payment_service/internal/config"
	"github.com/kanowfy/ecom/payment_service/internal/log"
	"github.com/kanowfy/ecom/payment_service/internal/service"
	"github.com/kanowfy/ecom/payment_service/pb"
	"google.golang.org/grpc"
)

func main() {
	var cfg config.Config

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("failed to parse environment variable: %v", err)
		os.Exit(1)
	}

	flag.StringVar(&cfg.Host, "srv.host", cfg.Host, "server host")
	flag.IntVar(&cfg.Port, "srv.port", cfg.Port, "server port")

	var level = slog.LevelDebug

	flag.Func("loglevel", "minimum log level", func(s string) error {
		if err := level.UnmarshalText([]byte(s)); err != nil {
			return err
		}

		return nil
	})
	flag.Parse()

	logger := log.New(os.Stdout, level, true)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		logger.Error("failed to announce address", "error", err.Error())
		os.Exit(1)
	}

	s := grpc.NewServer()
	pb.RegisterPaymentServer(s, service.New(logger))
	logger.Info(fmt.Sprintf("gRPC server listening on %s:%d", cfg.Host, cfg.Port))
	if err := s.Serve(lis); err != nil {
		logger.Error("failed to serve", "error", err.Error())
		os.Exit(1)
	}
}
