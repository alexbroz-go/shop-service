package handler

import (
	"context"
	pb "shop-service/shop-service/proto"
	"shop-service/user-cart-order/service"
)

type server struct {
	pb.UnimplementedUserCartOrderServiceServer
	svc *service.Service
}

func NewServer() pb.UserCartOrderServiceServer {
	return &server{
		svc: service.New(),
	}
}

func (s *server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return s.svc.Register(ctx, req)
}

func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return s.svc.Login(ctx, req)
}

func (s *server) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.User, error) {
	return s.svc.GetUser(ctx, req)
}

func (s *server) GetCart(ctx context.Context, req *pb.UserRequest) (*pb.Cart, error) {
	return s.svc.GetCart(ctx, req)
}

func (s *server) AddToCart(ctx context.Context, req *pb.AddToCartRequest) (*pb.Cart, error) {
	return s.svc.AddToCart(ctx, req)
}

func (s *server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	return s.svc.CreateOrder(ctx, req)
}

func (s *server) GetOrderHistory(ctx context.Context, req *pb.UserRequest) (*pb.OrderList, error) {
	return s.svc.GetOrderHistory(ctx, req)
}

func (s *server) GetOrder(ctx context.Context, req *pb.OrderRequest) (*pb.Order, error) {
	return s.svc.GetOrder(ctx, req)
}
