package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/kanowfy/ecom/frontend/internal/config"
	"github.com/kanowfy/ecom/frontend/internal/grpcconn"
	"github.com/kanowfy/ecom/frontend/internal/handlers"
	"github.com/kanowfy/ecom/frontend/internal/log"
	"github.com/kanowfy/ecom/frontend/internal/router"
	"github.com/kanowfy/ecom/frontend/internal/templatecache"
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
	flag.StringVar(&cfg.UpstreamAddr.Order, "addr.order", cfg.UpstreamAddr.Order, "order service address")
	flag.StringVar(&cfg.Cookie.SID, "cookie.sid", cfg.Cookie.SID, "cookie session id")
	flag.IntVar(&cfg.Cookie.MaxAge, "cookie.maxage", cfg.Cookie.MaxAge, "cookie max age")
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
	tp, err := initTracer(ctx, cfg.Otel.GrpcEndpoint, "frontend")
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

	templateCache, err := templatecache.New("./ui/templates/")
	if err != nil {
		logger.Error("failed to create template cache", "error", err)
		os.Exit(1)
	}

	h := handlers.New(logger, templateCache, conns)
	router := router.New(h, cfg)

	logger.Info(fmt.Sprintf("listening on port %d", cfg.Server.Port), "config", fmt.Sprintf("%+v", cfg), "conns", conns)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port), router); err != nil {
		logger.Error("failed to server", "error", err)
		os.Exit(1)
	}
}
