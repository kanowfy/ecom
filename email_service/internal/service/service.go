package service

import (
	"context"
	"log/slog"

	"github.com/kanowfy/ecom/email_service/internal/mailer"
	"github.com/kanowfy/ecom/email_service/pb"
	"go.opentelemetry.io/otel"
)

const tracerName = "github.com/kanowfy/ecom/email_service/service"

var tracer = otel.Tracer(tracerName)

type orderItem struct {
	Item struct {
		Product_id string
		Quantity   uint32
	}
	Cost int64
}

type service struct {
	logger *slog.Logger
	mailer mailer.Mailer
	*pb.UnimplementedEmailServer
}

func New(logger *slog.Logger, mailer mailer.Mailer) *service {
	return &service{
		logger: logger,
		mailer: mailer,
	}
}

func (s *service) SendConfirmation(ctx context.Context, req *pb.SendConfirmationRequest) (*pb.None, error) {
	ctx, span := tracer.Start(ctx, "send_confirmation_email")
	defer span.End()

	s.logger.Info("received a SendConfirmation request")

	var items []*orderItem
	for _, it := range req.GetOrder().GetItems() {
		var ordItem orderItem
		ordItem.Item = struct {
			Product_id string
			Quantity   uint32
		}{Product_id: it.GetItem().ProductId, Quantity: it.GetItem().Quantity}
		ordItem.Cost = it.GetCost()
		items = append(items, &ordItem)
	}

	data := map[string]interface{}{
		"order_id":             req.GetOrder().OrderId,
		"shipping_tracking_id": req.GetOrder().ShippingTrackingId,
		"shipping_cost":        req.GetOrder().ShippingCost,
		"shipping_address": map[string]interface{}{
			"street_address": req.GetOrder().ShippingAddress.StreetAddress,
			"city":           req.GetOrder().ShippingAddress.City,
			"country":        req.GetOrder().ShippingAddress.Country,
			"zip_code":       req.GetOrder().ShippingAddress.ZipCode,
		},
		"items": items,
	}

	go func() {
		err := s.mailer.Send(req.GetEmail(), "confirmation.tmpl", data)
		if err != nil {
			s.logger.Error("failed to send email", "error", err.Error())
		}
	}()

	return &pb.None{}, nil
}
