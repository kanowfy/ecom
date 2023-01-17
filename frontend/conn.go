package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func mapConnections(app *application) {
	ctx := context.Background()

	mustMapEnv(&app.catalogSvcAddr, "CATALOG_SVC_ADDR")
	mustMapEnv(&app.cartSvcAddr, "CART_SVC_ADDR")
	mustMapEnv(&app.shippingSvcAddr, "SHIPPING_SVC_ADDR")
	mustMapEnv(&app.emailSvcAddr, "EMAIL_SVC_ADDR")
	mustMapEnv(&app.paymentSvcAddr, "PAYMENT_SVC_ADDR")
	mustMapEnv(&app.orderSvcAddr, "ORDER_SVC_ADDR")

	mustConnGRPC(ctx, &app.catalogSvcConn, app.catalogSvcAddr)
	mustConnGRPC(ctx, &app.cartSvcConn, app.cartSvcAddr)
	mustConnGRPC(ctx, &app.shippingSvcConn, app.shippingSvcAddr)
	mustConnGRPC(ctx, &app.emailSvcConn, app.emailSvcAddr)
	mustConnGRPC(ctx, &app.paymentSvcConn, app.paymentSvcAddr)
	mustConnGRPC(ctx, &app.orderSvcConn, app.orderSvcAddr)
}
