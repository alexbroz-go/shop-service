package main

import (
	"log"
	"net"
	pb "shop-service/shop-service/proto"
	"shop-service/user-cart-order/database"
	"shop-service/user-cart-order/handler"

	"google.golang.org/grpc"
)

func main() {
	// Инициализация базы данных
	if err := database.Init(); err != nil {
		log.Fatalf("Ошибка инициализации базы: %v", err)
	}

	// Прослушка порта
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Невозможно запустить слушателя: %v", err)
	}

	// gRPC сервер
	grpcServer := grpc.NewServer()
	pb.RegisterUserCartOrderServiceServer(grpcServer, handler.NewServer())

	log.Println("Сервер запущен на порту :50052")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
