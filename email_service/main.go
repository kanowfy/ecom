package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/kanowfy/ecom/email_service/internal/config"
	"github.com/kanowfy/ecom/email_service/internal/log"
	"github.com/kanowfy/ecom/email_service/internal/mailer"
	"github.com/kanowfy/ecom/email_service/internal/service"
	"github.com/kanowfy/ecom/email_service/pb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"google.golang.org/grpc"
)

const SvcName = "email-svc"

func main() {
	var cfg config.Config

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("failed to parse environment variable: %v", err)
		os.Exit(1)
	}

	flag.StringVar(&cfg.Server.Host, "srv.host", cfg.Server.Host, "server host")
	flag.IntVar(&cfg.Server.Port, "srv.port", cfg.Server.Port, "server port")
	flag.StringVar(&cfg.Mail.Host, "mail.host", cfg.Mail.Host, "mail host")
	flag.IntVar(&cfg.Mail.Port, "mail.port", cfg.Mail.Port, "mail port")
	flag.StringVar(&cfg.Mail.Username, "mail.username", cfg.Mail.Username, "mail username")
	flag.StringVar(&cfg.Mail.Password, "mail.password", cfg.Mail.Password, "mail password")
	flag.StringVar(&cfg.Mail.Sender, "mail.sender", cfg.Mail.Sender, "mail sender")
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
	tp, err := initTracer(ctx, cfg.Otel.GrpcEndpoint, SvcName)
	if err != nil {
		logger.Error("failed to initialize tracer", "error", err)
		os.Exit(1)
	}
	defer tp.Shutdown(ctx)

	mp, err := initMetrics(ctx, cfg.Otel.GrpcEndpoint, SvcName)
	if err != nil {
		logger.Error("failed to initialize tracer", "error", err)
		os.Exit(1)
	}
	defer mp.Shutdown(ctx)

	mailer := mailer.New(cfg.Mail.Host, cfg.Mail.Port, cfg.Mail.Username, cfg.Mail.Password, cfg.Mail.Sender)
	service := service.New(logger, mailer)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
	if err != nil {
		logger.Error("failed to announce address", "error", err.Error())
		os.Exit(1)
	}

	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	pb.RegisterEmailServer(s, service)
	logger.Info(fmt.Sprintf("gRPC server listening on %s:%d", cfg.Server.Host, cfg.Server.Port))
	if err := s.Serve(lis); err != nil {
		logger.Error("failed to serve", "error", err.Error())
		os.Exit(1)
	}
}
