package main

import (
	"context"
	"github.com/kanowfy/ecom/email_service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type orderItem struct {
	item struct {
		product_id string
		quantity uint32
	}
	price int64
}

func (s *server) SendConfirmation(ctx context.Context, req *pb.SendConfirmationRequest) (*pb.None, error) {
	log.Printf("received a SendConfirmation request with data: %+v", req)
	defer log.Println("SendConfirmation successful")

	var items []*orderItem
	for _, it := range req.GetOrder().GetItems() {
		var ordItem orderItem
		ordItem.item = struct {
			product_id string
			quantity   uint32
		}{product_id: it.GetItem().ProductId, quantity: it.GetItem().Quantity}
		ordItem.price = it.GetCost()
		items = append(items, &ordItem)
	}
	data := map[string]interface{}{
		"order_id": req.GetOrder().OrderId,
		"shipping_tracking_id": req.GetOrder().ShippingTrackingId,
		"shipping_address": map[string]interface{}{
			"street_address": req.GetOrder().ShippingAddress.StreetAddress,
			"city": req.GetOrder().ShippingAddress.City,
			"country": req.GetOrder().ShippingAddress.Country,
			"zip_code": req.GetOrder().ShippingAddress.ZipCode,
			},
			"items": items,
	}

	err := s.mailer.Send(req.GetEmail(), "confirmation.tmpl", data)
	if err != nil {
		log.Printf("send email error: %v", err)
		return nil, status.Errorf(codes.Internal, "send email error: %v", err)
	}

	return &pb.None{}, status.New(codes.OK, "").Err()
}
