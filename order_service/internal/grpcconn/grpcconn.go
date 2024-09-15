package grpcconn

import (
	"context"
	"fmt"
	"time"

	"github.com/kanowfy/ecom/order_service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Connection struct {
	CatalogSvc  *grpc.ClientConn
	CartSvc     *grpc.ClientConn
	ShippingSvc *grpc.ClientConn
	EmailSvc    *grpc.ClientConn
	PaymentSvc  *grpc.ClientConn
}

func dialService(ctx context.Context, conn **grpc.ClientConn, addr string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	var err error
	*conn, err = grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to grpc server with address %s: %w", addr, err)
	}

	return nil
}

func (c *Connection) Map(ctx context.Context, addrs config.UpstreamAddr) error {
	pairs := []struct {
		conn *grpc.ClientConn
		addr string
	}{
		{
			conn: c.CatalogSvc,
			addr: addrs.Catalog,
		},
		{
			conn: c.CartSvc,
			addr: addrs.Cart,
		},
		{
			conn: c.ShippingSvc,
			addr: addrs.Shipping,
		},
		{
			conn: c.EmailSvc,
			addr: addrs.Email,
		},
		{
			conn: c.PaymentSvc,
			addr: addrs.Payment,
		},
	}

	for _, pair := range pairs {
		if err := dialService(ctx, &pair.conn, pair.addr); err != nil {
			return fmt.Errorf("failed to map connection: %w", err)
		}
	}
	return nil
}
