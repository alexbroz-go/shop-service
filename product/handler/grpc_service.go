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
	}, nil
}
func (s *server) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {
	products, totalCount, err := s.svc.SearchProducts(req.Query, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	var pbProducts []*pb.Product
	for _, p := range products {
		pbProducts = append(pbProducts, &pb.Product{
			Id:          int32(p.ID),
			Title:       p.Title,
			Description: p.Description,
			Price:       p.Price,
		})
	}

	return &pb.SearchProductsResponse{
		Products:   pbProducts,
		TotalCount: int32(totalCount),
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
	}
	id, err := s.svc.AddProduct(product)
	if err != nil {
		return nil, err
	}
	return &pb.AddProductResponse{ProductId: int32(id)}, nil
}
