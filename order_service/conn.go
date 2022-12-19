package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"time"
)

func mustMapEnv(target *string, key string) {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("environment variable %s not set", key))
	}

	*target = v
}

func mustConnGRPC(ctx context.Context, conn **grpc.ClientConn, addr string) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	var err error
	*conn, err = grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Sprintf("failed to connect to grpc server with address %s, %v", addr, err))
	}
}

func mapConnections(s *server) {
	ctx := context.Background()

	mustMapEnv(&s.conns.catalogSvcAddr, "CATALOG_SVC_ADDR")
	mustMapEnv(&s.conns.cartSvcAddr, "CART_SVC_ADDR")
	mustMapEnv(&s.conns.shippingSvcAddr, "SHIPPING_SVC_ADDR")
	mustMapEnv(&s.conns.emailSvcAddr, "EMAIL_SVC_ADDR")
	mustMapEnv(&s.conns.paymentSvcAddr, "PAYMENT_SVC_ADDR")

	mustConnGRPC(ctx, &s.conns.catalogSvcConn, s.conns.catalogSvcAddr)
	mustConnGRPC(ctx, &s.conns.cartSvcConn, s.conns.cartSvcAddr)
	mustConnGRPC(ctx, &s.conns.shippingSvcConn, s.conns.shippingSvcAddr)
	mustConnGRPC(ctx, &s.conns.emailSvcConn, s.conns.emailSvcAddr)
	mustConnGRPC(ctx, &s.conns.paymentSvcConn, s.conns.paymentSvcAddr)
}
