package main

import (
	"context"
	"log"

	"github.com/kanowfy/ecom/email_service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type orderItem struct {
	Item struct {
		Product_id string
		Quantity   uint32
	}
	Cost int64
}

func (s *server) SendConfirmation(ctx context.Context, req *pb.SendConfirmationRequest) (*pb.None, error) {
	log.Printf("received a SendConfirmation request with data: %+v", req)

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

	err := s.mailer.Send(req.GetEmail(), "confirmation.tmpl", data)
	if err != nil {
		log.Printf("send email error: %v", err)
		return nil, status.Errorf(codes.Internal, "send email error: %v", err)
	}

	log.Println("SendConfirmation successful")
	return &pb.None{}, status.New(codes.OK, "").Err()
}
