package main

import (
	"context"
	"errors"
	"github.com/kanowfy/ecom/cart_service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (s *server) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.GetCartResponse, error) {
	err := s.validateUUID(req.GetUserId())
	if err != nil {
		return nil, err
	}
	log.Printf("received a GetCart request with user id: %s", req.GetUserId())

	cart, err := s.model.GetCart(req.GetUserId())
	if err != nil {
		switch {
		case errors.Is(err, ErrNoRecord):
			return nil, status.Error(codes.NotFound, "no cart with matching user id found")
		default:
			log.Println(err)
			return nil, status.Errorf(codes.Internal, "unable to fetch cart: %v", err)
		}
	}

	var items []*pb.CartItem

	for _, i := range cart.Items {
		item := &pb.CartItem{
			ProductId: i.ProductId,
			Quantity:  i.Quantity,
		}

		items = append(items, item)
	}

	resp := &pb.GetCartResponse{
		UserId: req.GetUserId(),
		Items:  items,
	}

	log.Println("GetCart successful")
	return resp, status.New(codes.OK, "").Err()
}

func (s *server) AddItem(ctx context.Context, req *pb.AddItemRequest) (*pb.None, error) {
	err := s.validateUUID(req.GetUserId())
	if err != nil {
		return nil, err
	}
	log.Printf("received an AddItem request with user id: %s. product id: %s, quantity: %d",
		req.GetUserId(), req.Item.GetProductId(), req.Item.GetQuantity())

	item := &CartItem{
		ProductId: req.Item.GetProductId(),
		Quantity:  req.Item.GetQuantity(),
	}

	err = s.model.AddItem(req.GetUserId(), item)
	if err != nil {
		switch {
		case errors.Is(err, ErrInCart):
			return nil, status.Error(codes.AlreadyExists, "item already in cart")
		default:
			log.Println(err)
			return nil, status.Errorf(codes.Internal, "unable to add item: %v", err)
		}
	}

	log.Println("AddItem successful")
	return &pb.None{}, status.New(codes.OK, "").Err()
}

func (s *server) RemoveItem(ctx context.Context, req *pb.RemoveItemRequest) (*pb.None, error) {
	err := s.validateUUID(req.GetUserId(), req.GetProductId())
	if err != nil {
		return nil, err
	}
	log.Printf("received a RemoveItem request with user id: %s. product id: %s",
		req.GetUserId(), req.GetProductId())

	err = s.model.RemoveItem(req.GetUserId(), req.GetProductId())
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "unable to remove item: %v", err)
	}

	log.Println("RemoveItem successful")
	return &pb.None{}, status.New(codes.OK, "").Err()
}
