package main

import (
	"context"
	"github.com/kanowfy/ecom/shipping_service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (s *server) GetShippingCost(ctx context.Context, req *pb.GetShippingCostRequest) (*pb.GetShippingCostResponse, error) {
	log.Printf("received a GetShippingCost request with data: %+v", req)
	defer log.Println("GetShippingCost successful")
	return &pb.GetShippingCostResponse{
		Cost: 5,
	}, status.New(codes.OK, "").Err()
}

func (s *server) ShipOrder(ctx context.Context, req *pb.ShipOrderRequest) (*pb.ShipOrderResponse, error) {
	log.Printf("received a ShipOrder request with data: %+v", req)
	defer log.Println("ShipOrder successful")

	trackingId := CreateTrackingId()

	return &pb.ShipOrderResponse{
		TrackingId: trackingId,
	}, status.New(codes.OK, "").Err()
}
