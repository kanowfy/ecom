package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/kanowfy/ecom/order_service/internal/config"
	"github.com/kanowfy/ecom/order_service/internal/grpcconn"
	"github.com/kanowfy/ecom/order_service/internal/log"
	"github.com/kanowfy/ecom/order_service/internal/service"
	"github.com/kanowfy/ecom/order_service/pb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

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
	flag.StringVar(&cfg.UpstreamAddr.Catalog, "addr.catalog", cfg.UpstreamAddr.Catalog, "catalog service address")
	flag.StringVar(&cfg.UpstreamAddr.Cart, "addr.cart", cfg.UpstreamAddr.Cart, "cart service address")
	flag.StringVar(&cfg.UpstreamAddr.Shipping, "addr.shipping", cfg.UpstreamAddr.Shipping, "shipping service address")
	flag.StringVar(&cfg.UpstreamAddr.Payment, "addr.payment", cfg.UpstreamAddr.Payment, "payment service address")
	flag.StringVar(&cfg.UpstreamAddr.Email, "addr.email", cfg.UpstreamAddr.Email, "email service address")
	flag.StringVar(&cfg.Otel.GrpcEndpoint, "otel.grpcendpoint", cfg.Otel.GrpcEndpoint, "grpc collector endpoint")

	var level = slog.LevelDebug

	flag.Func("loglevel", "minimum log level", func(s string) error {
		if err := level.UnmarshalText([]byte(s)); err != nil {
			return err
		}

		return nil
	})

	flag.Parse()

	logger := log.New(os.Stdout, level, true)
	ctx := context.Background()
	tp, err := initTracer(ctx, cfg.Otel.GrpcEndpoint, "order service")
	if err != nil {
		logger.Error("failed to initialize tracer", "error", err)
		os.Exit(1)
	}
	defer tp.Shutdown(ctx)

	conns := new(grpcconn.Connection)
	if err := conns.Map(context.Background(), cfg.UpstreamAddr); err != nil {
		fmt.Printf("failed to establish upstream connections: %v", err)
		os.Exit(1)
	}
	service := service.New(logger, conns)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
	if err != nil {
		logger.Error("failed to announce address", "error", err.Error())
		os.Exit(1)
	}

	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	pb.RegisterOrderServer(s, service)
	logger.Info(fmt.Sprintf("gRPC server listening on %s:%d", cfg.Server.Host, cfg.Server.Port))
	if err := s.Serve(lis); err != nil {
		logger.Error("failed to serve", "error", err.Error())
		os.Exit(1)
	}
}
