package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/kanowfy/ecom/catalog_service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr string
)

func init() {
	flag.StringVar(&addr, "addr", "localhost:3001", "server address")
	flag.Parse()
}

func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error dialing server: %v", err)
	}

	defer conn.Close()

	client := pb.NewCatalogClient(conn)

	product := &pb.Product{
		Name:        "Nike Air Jordan 1 Low",
		Description: "Inspired by the original that debuted in 1985, the Air Jordan 1 Low offers a clean, classic look that's familiar yet always fresh.",
		Category:    "Shoes",
		PriceVnd:    1500000,
		Image:       "https://static.nike.com/a/images/t_PDP_1280_v1/f_auto,q_auto:eco/1e463dee-799d-4fba-9b32-0f7e0bb9d5f5/air-jordan-1-low-shoes-6Q1tFM.png",
	}

	resp, err := client.AddProduct(context.Background(), &pb.AddProductRequest{Product: product})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("AddProduct: product added with id %s", resp.Id)

	fmt.Println("---------------------------------------")
	resp1, err := client.GetProductById(context.Background(), &pb.GetProductByIdRequest{Id: resp.Id})
	fmt.Printf("GetProductById: received product: %v", resp1.Product)
	fmt.Println("---------------------------------------")
	resp2, err := client.ListProducts(context.Background(), &pb.None{})
	for _, p := range resp2.Products {
		fmt.Printf("ListProducts: product: %v", p)
	}
	fmt.Println("---------------------------------------")
	resp3, err := client.SearchProducts(context.Background(), &pb.SearchProductsRequest{Query: "jordan"})
	for _, p := range resp3.Results {
		fmt.Printf("SearchProducts: product: %v", p)
	}
}
