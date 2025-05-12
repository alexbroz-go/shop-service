package main

import (
	"log"
	"net"
	"user-cart-order/database"
	"user-cart-order/handler"
	pb "user-cart-order/proto"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
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
