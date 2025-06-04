package dao_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"user-cart-order/database"
	"user-cart-order/models"
)

func TestMain(m *testing.M) {
	// Инициализация базы данных
	err := database.Init()
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

// Тест для создания и получения пользователя
func TestCreateAndGetUser(t *testing.T) {
	// Генерируем уникальный логин и email
	uniqueSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
	login := "testuser_" + uniqueSuffix
	email := login + "@example.com"

	user := models.User{
		Login:    login,
		Email:    email,
		Password: "password123",
	}

	id, err := database.CreateUser(user)
	if err != nil {
		t.Fatalf("Ошибка при создании пользователя: %v", err)
	}

	fetchedUser, err := database.GetUserByID(id)
	if err != nil {
		t.Fatalf("Ошибка при получении пользователя по ID: %v", err)
	}

	if fetchedUser.Login != user.Login {
		t.Errorf("Ожидали login=%s, получили=%s", user.Login, fetchedUser.Login)
	}
	if fetchedUser.Email != user.Email {
		t.Errorf("Ожидали email=%s, получили=%s", user.Email, fetchedUser.Email)
	}
	if fetchedUser.Password != user.Password {
		t.Errorf("Ожидали password=%s, получили=%s", user.Password, fetchedUser.Password)
	}
}

// Тест для создания и получения корзины
func TestCreateAndGetCart(t *testing.T) {
	uniqueSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
	login := "cartuser_" + uniqueSuffix
	email := login + "@example.com"

	user := models.User{
		Login:    login,
		Email:    email,
		Password: "pass",
	}
	userID, err := database.CreateUser(user)
	if err != nil {
		t.Fatalf("Ошибка при создании пользователя: %v", err)
	}

	cartID, err := database.CreateCart(userID)
	if err != nil {
		t.Fatalf("Ошибка при создании корзины: %v", err)
	}

	cart, err := database.GetCartByUserID(userID)
	if err != nil {
		t.Fatalf("Ошибка при получении корзины: %v", err)
	}

	if cart.UserID != userID {
		t.Errorf("Ожидали userID=%d, получили=%d", userID, cart.UserID)
	}
	if cart.ID != cartID {
		t.Errorf("Ожидали id=%d, получили=%d", cartID, cart.ID)
	}
}

// Тест для создания и получения заказа
func TestCreateAndGetOrder(t *testing.T) {
	uniqueSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
	login := "orderuser_" + uniqueSuffix
	email := login + "@example.com"

	user := models.User{
		Login:    login,
		Email:    email,
		Password: "pass",
	}
	userID, err := database.CreateUser(user)
	if err != nil {
		t.Fatalf("Ошибка при создании пользователя: %v", err)
	}

	productID := 1

	order := models.Order{
		UserID:    userID,
		ProductID: productID,
		Status:    "pending",
	}
	orderID, err := database.CreateOrder(order)
	if err != nil {
		t.Fatalf("Ошибка при создании заказа: %v", err)
	}

	fetchedOrder, err := database.GetOrderByID(orderID)
	if err != nil {
		t.Fatalf("Ошибка при получении заказа: %v", err)
	}

	if fetchedOrder.UserID != order.UserID {
		t.Errorf("Ожидали UserID=%d, получили=%d", order.UserID, fetchedOrder.UserID)
	}
	if fetchedOrder.ProductID != order.ProductID {
		t.Errorf("Ожидали ProductID=%d, получили=%d", order.ProductID, fetchedOrder.ProductID)
	}
	if fetchedOrder.Status != order.Status {
		t.Errorf("Ожидали Status=%s, получили=%s", order.Status, fetchedOrder.Status)
	}

	orders, err := database.GetOrderHistory(userID)
	if err != nil {
		t.Fatalf("Ошибка при получении истории заказов: %v", err)
	}
	found := false
	for _, o := range orders {
		if o.ID == orderID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Заказ с ID=%d не найден в истории", orderID)
	}
}
