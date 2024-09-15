package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/gofrs/uuid"
	"github.com/kanowfy/ecom/cart_service/internal/repository"
	"github.com/kanowfy/ecom/cart_service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	logger *slog.Logger
	repo   repository.Repository
	*pb.UnimplementedCartServer
}

func New(logger *slog.Logger, repo repository.Repository) *service {
	return &service{
		logger: logger,
		repo:   repo,
	}
}

func (s *service) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.GetCartResponse, error) {
	id := req.GetUserId()
	if err := validateUUID(id); err != nil {
		return nil, fmt.Errorf("failed to validate user id: %w", err)
	}

	logger := s.logger.With(slog.String("user_id", id))

	logger.Info("received GetCart request")

	cart, err := s.repo.GetCart(req.GetUserId())
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNoRecord):
			return nil, status.Error(codes.NotFound, "no cart with matching user id found")
		default:
			logger.Error("failed to fetch cart", "error", err.Error())
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

	logger.Info("GetCart successful")
	return resp, nil
}

func (s *service) AddItem(ctx context.Context, req *pb.AddItemRequest) (*pb.None, error) {
	id := req.GetUserId()
	logger := s.logger.With(slog.String("user_id", id), slog.String("product_id", req.Item.GetProductId()), slog.Uint64("quantity", uint64(req.Item.GetQuantity())))

	if err := validateUUID(id); err != nil {
		return nil, err
	}

	logger.Info("received an AddItem request")

	item := &repository.CartItem{
		ProductId: req.Item.GetProductId(),
		Quantity:  req.Item.GetQuantity(),
	}

	if err := s.repo.AddItem(req.GetUserId(), item); err != nil {
		switch {
		case errors.Is(err, repository.ErrInCart):
			return nil, status.Error(codes.AlreadyExists, "item already in cart")
		default:
			logger.Error("failed to add item", "error", err.Error())
			return nil, status.Errorf(codes.Internal, "unable to add item: %v", err)
		}
	}

	logger.Info("AddItem successful")
	return &pb.None{}, nil
}

func (s *service) RemoveItem(ctx context.Context, req *pb.RemoveItemRequest) (*pb.None, error) {
	id := req.GetUserId()
	logger := s.logger.With(slog.String("user_id", id), slog.String("product_id", req.GetProductId()))
	if err := validateUUID(id, req.GetProductId()); err != nil {
		return nil, fmt.Errorf("failed to validate user id: %w", err)
	}

	logger.Info("received a RemoveItem request")

	if err := s.repo.RemoveItem(req.GetUserId(), req.GetProductId()); err != nil {
		logger.Error("failed to remove item", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "unable to remove item: %v", err)
	}

	logger.Info("RemoveItem successful")
	return &pb.None{}, nil
}

func validateUUID(ids ...string) error {
	for _, id := range ids {
		_, err := uuid.FromString(id)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
		}
	}

	return nil
}
