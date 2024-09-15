package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gofrs/uuid"
	"github.com/kanowfy/ecom/catalog_service/internal/repository"
	"github.com/kanowfy/ecom/catalog_service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	logger *slog.Logger
	repo   repository.Repository
	*pb.UnimplementedCatalogServer
}

func New(logger *slog.Logger, repo repository.Repository) *service {
	return &service{
		logger: logger,
		repo:   repo,
	}
}

func (s *service) GetProductById(ctx context.Context, req *pb.GetProductByIdRequest) (*pb.GetProductByIdResponse, error) {
	id := req.GetId()
	_, err := uuid.FromString(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}

	logger := s.logger.With(slog.String("product_id", id))

	logger.Info("received a GetProductById request")

	product, err := s.repo.GetById(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNoRecord):
			return nil, status.Error(codes.NotFound, "no product with matching id found")
		default:
			logger.Error("failed to fetch product", "error", err.Error())
			return nil, status.Errorf(codes.Internal, "failed to fetch product", err)
		}
	}

	resp := &pb.GetProductByIdResponse{
		Product: &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Category:    product.Category,
			PriceVnd:    product.Price,
			Image:       product.Image,
		},
	}

	logger.Info("GetProductById successful")
	return resp, nil
}

func (s *service) AddProduct(ctx context.Context, req *pb.AddProductRequest) (*pb.AddProductResponse, error) {
	s.logger.Info("received an AddProduct request")
	id, err := uuid.NewV4()
	if err != nil {
		s.logger.Error("failed to generate uuid: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to generate uuid: %v", err)
	}

	product := repository.Product{
		ID:          id.String(),
		Name:        req.Product.Name,
		Description: req.Product.Description,
		Category:    req.Product.Category,
		Price:       req.Product.PriceVnd,
		Image:       req.Product.Image,
	}

	if err = s.repo.Add(ctx, product); err != nil {
		s.logger.Error("can not add product: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to add product: %v", err)
	}

	s.logger.Info("AddProduct successful")

	resp := &pb.AddProductResponse{
		Id: id.String(),
	}

	return resp, nil
}
func (s *service) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {
	s.logger.Info("received a SearchProducts request")
	products, err := s.repo.List(ctx, req.GetQuery())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve products: %v", err)
	}

	var productsResponse []*pb.Product
	for _, p := range products {
		product := &pb.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Category:    p.Category,
			PriceVnd:    p.Price,
			Image:       p.Image,
		}
		productsResponse = append(productsResponse, product)
	}

	s.logger.Info("SearchProducts successful")
	resp := &pb.SearchProductsResponse{Results: productsResponse}
	return resp, nil
}
