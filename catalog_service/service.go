package main

import (
	"context"
	"errors"
	"log"

	"github.com/gofrs/uuid"
	"github.com/kanowfy/ecom/catalog_service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) GetProductById(ctx context.Context, req *pb.GetProductByIdRequest) (*pb.GetProductByIdResponse, error) {
	_, err := uuid.FromString(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}
	log.Printf("received a GetProductById request with id %s", req.GetId())

	product, err := s.model.GetById(req.GetId())
	if err != nil {
		switch {
		case errors.Is(err, ErrNoRecord):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("error: %v", err)
			return nil, status.Errorf(codes.Internal, "GetProductById failed: %v", err)
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

	log.Printf("GetProductById successful")
	return resp, status.New(codes.OK, "").Err()
}

func (s *server) AddProduct(ctx context.Context, req *pb.AddProductRequest) (*pb.AddProductResponse, error) {
	log.Println("received an AddProduct request")
	id, err := uuid.NewV4()
	if err != nil {
		log.Printf("can not generate uuid: %v", err)
		return nil, status.Errorf(codes.Internal, "unable to generate uuid: %v", err)
	}

	product := &Product{
		ID:          id.String(),
		Name:        req.Product.Name,
		Description: req.Product.Description,
		Category:    req.Product.Category,
		Price:       req.Product.PriceVnd,
		Image:       req.Product.Image,
	}

	err = s.model.Add(product)
	if err != nil {
		log.Printf("can not add product: %v", err)
		return nil, status.Errorf(codes.Internal, "unable to add product: %v", err)
	}

	log.Printf("AddProduct successful")

	resp := &pb.AddProductResponse{
		Id: id.String(),
	}

	return resp, status.New(codes.OK, "").Err()
}

func (s *server) ListProducts(ctx context.Context, req *pb.None) (*pb.ListProductsResponse, error) {
	log.Println("received a ListProducts request")
	products, err := s.model.List("")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to retrieve products: %v", err)
	}

	productsResponse := []*pb.Product{}
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

	log.Println("ListProducts successful")
	resp := &pb.ListProductsResponse{Products: productsResponse}
	return resp, status.New(codes.OK, "").Err()
}

func (s *server) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {
	log.Println("received a SearchProducts request")
	products, err := s.model.List(req.GetQuery())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to retrieve products: %v", err)
	}

	productsResponse := []*pb.Product{}
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

	log.Println("SearchProducts successful")
	resp := &pb.SearchProductsResponse{Results: productsResponse}
	return resp, status.New(codes.OK, "").Err()
}
