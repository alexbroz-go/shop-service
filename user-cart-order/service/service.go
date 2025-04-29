package service

import (
	"context"
	pb "shop-service/shop-service/proto"
	"shop-service/user-cart-order/database"
	"shop-service/user-cart-order/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service struct{}

func New() *Service { return &Service{} }

// Регистрация
func (s *Service) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := models.User{
		Login:    req.Login,
		Email:    req.Email,
		Password: string(hashedPassword),
	}
	_, err := database.CreateUser(user)
	if err != nil {
		return &pb.RegisterResponse{Success: false, Message: "Ошибка при регистрации"}, err
	}
	return &pb.RegisterResponse{Success: true, Message: "Успешная регистрация"}, nil
}

// Вход
func (s *Service) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := database.GetUserByLogin(req.Login)
	if err != nil || user.Password == "" {
		return &pb.LoginResponse{Success: false, Message: "Некорректный логин или пароль"}, nil
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return &pb.LoginResponse{Success: false, Message: "Некорректный логин или пароль"}, nil
	}
	return &pb.LoginResponse{Success: true, Message: "Успешный вход", UserId: int32(user.ID)}, nil
}

// Получить или создать юзера
func (s *Service) GetOrCreateUser(ctx context.Context, req *pb.UserRequest) (*pb.User, error) {
	user, err := database.GetUserByID(int(req.UserId))
	if err != nil {
		// По сути, тут можно создать, если нет — для демонстрации возьмем, что создаем нового по логину
		// Но лучше созадть явно
		return nil, err
	}
	return &pb.User{
		Id:    int32(user.ID),
		Login: user.Login,
		Email: user.Email,
	}, nil
}

// Получить корзину
func (s *Service) GetCart(ctx context.Context, req *pb.UserRequest) (*pb.Cart, error) {
	cart, err := database.GetCartByUserID(int(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.Cart{
		Id:         int32(cart.ID),
		UserId:     int32(cart.UserID),
		CreatedAt:  cart.CreatedAt.Format("2006-01-02 15:04:05"),
		ProductIds: []int32{},
	}, nil
}

// Добавить товар в корзину
func (s *Service) AddToCart(ctx context.Context, req *pb.AddToCartRequest) (*pb.Cart, error) {
	// В простом варианте — возвращаем текущую корзину
	cart, err := database.GetCartByUserID(int(req.UserId))
	if err != nil {
		// Если корзины нет — создадим
		cartID, err := database.CreateCart(int(req.UserId))
		if err != nil {
			return nil, err
		}
		cart, _ = database.GetCartByUserID(int(req.UserId))
		cart.ID = cartID
		cart.CreatedAt = time.Now()
	}
	// Логика добавления товара (расширение по необходимости)
	return &pb.Cart{
		Id:         int32(cart.ID),
		UserId:     int32(cart.UserID),
		CreatedAt:  cart.CreatedAt.Format("2006-01-02 15:04:05"),
		ProductIds: []int32{},
	}, nil
}

// Создать заказ
func (s *Service) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	order := models.Order{
		UserID:    int(req.UserId),
		ProductID: int(req.ProductId),
		Status:    req.Status,
	}
	id, err := database.CreateOrder(order)
	if err != nil {
		return nil, err
	}
	order.ID = id
	order.CreatedAt = time.Now()
	return &pb.Order{
		Id:        int32(order.ID),
		UserId:    int32(order.UserID),
		ProductId: int32(order.ProductID),
		CreatedAt: order.CreatedAt.Format("2006-01-02 15:04:05"),
		Status:    order.Status,
	}, nil
}

// Получить историю заказов
func (s *Service) GetOrderHistory(ctx context.Context, req *pb.UserRequest) (*pb.OrderList, error) {
	orders, err := database.GetOrderHistory(int(req.UserId))
	if err != nil {
		return nil, err
	}
	var pbOrders []*pb.Order
	for _, o := range orders {
		pbOrders = append(pbOrders, &pb.Order{
			Id:        int32(o.ID),
			UserId:    int32(o.UserID),
			ProductId: int32(o.ProductID),
			CreatedAt: o.CreatedAt.Format("2006-01-02 15:04:05"),
			Status:    o.Status,
		})
	}
	return &pb.OrderList{Orders: pbOrders}, nil
}

// Взять конкретный заказ
func (s *Service) GetOrder(ctx context.Context, req *pb.OrderRequest) (*pb.Order, error) {
	order, err := database.GetOrderByID(int(req.OrderId))
	if err != nil {
		return nil, err
	}
	return &pb.Order{
		Id:        int32(order.ID),
		UserId:    int32(order.UserID),
		ProductId: int32(order.ProductID),
		CreatedAt: order.CreatedAt.Format("2006-01-02 15:04:05"),
		Status:    order.Status,
	}, nil
}

func (s *Service) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.User, error) {
	// Получите пользователя из базы по req.UserId
	user, err := database.GetUserByID(int(req.UserId))
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id:    int32(user.ID),
		Login: user.Login,
		Email: user.Email,
	}, nil
}
