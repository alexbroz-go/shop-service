package handler

import (
	"context"
	"product/models"
	pb "product/proto"
	"product/service"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedProductServiceServer
	svc *service.Service
}

func RegisterProductServer(grpcServer *grpc.Server) {
	pb.RegisterProductServiceServer(grpcServer, &server{
		svc: service.New(),
	})
}

func (s *server) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	product, err := s.svc.GetProductByID(int(req.ProductId))
	if err != nil {
		return nil, err
	}
	return &pb.Product{
		Id:          int32(product.ID),
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
		Count:       int32(product.Count),
	}, nil
}
func (s *server) GetProductByTitle(ctx context.Context, req *pb.GetProductByTitleRequest) (*pb.Product, error) {
	product, err := s.svc.GetProductByTitle(req.Title)
	if err != nil {
		return nil, err
	}
	return &pb.Product{
		Id:          int32(product.ID),
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
		Count:       int32(product.Count),
	}, nil
}

func (s *server) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	success, err := s.svc.DeleteProduct(int(req.ProductId))
	if err != nil {
		return nil, err
	}
	return &pb.DeleteProductResponse{Success: success}, nil
}

func (s *server) AddProduct(ctx context.Context, req *pb.AddProductRequest) (*pb.AddProductResponse, error) {
	product := models.Product{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Count:       int(req.Count),
	}
	id, err := s.svc.AddProduct(product)
	if err != nil {
		return nil, err
	}
	return &pb.AddProductResponse{ProductId: int32(id)}, nil
}
