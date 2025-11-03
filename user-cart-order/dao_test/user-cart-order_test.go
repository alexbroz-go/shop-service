package dao_test

import (
	"os"
	"testing"

	"user-cart-order/database"
	"user-cart-order/models"
)

func TestMain(m *testing.M) {
	if err := database.Init(); err != nil {
		panic(err)
	}
	// Clean DB state for stable tests
	if _, err := database.DB.Exec("TRUNCATE TABLE orders, carts, users RESTART IDENTITY CASCADE;"); err != nil {
		panic(err)
	}
	code := m.Run()
	// Optional: clean after run as well
	_, _ = database.DB.Exec("TRUNCATE TABLE orders, carts, users RESTART IDENTITY CASCADE;")
	os.Exit(code)
}

func createTestUser(t *testing.T) models.User {
	t.Helper()
	user := models.User{
		Login:    "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	id, err := database.CreateUser(user)
	if err != nil {
		t.Fatalf("Ошибка при создании пользователя: %v", err)
	}
	user.ID = id
	return user
}
func deleteTestUser(t *testing.T, userID int) {
	t.Helper()
	success, err := database.DeleteUserByID(userID)
	if err != nil {
		t.Fatalf("Ошибка при удалении пользователя: %v", err)
	}
	if !success {
		t.Fatalf("Удаление пользователя не удалось, id=%d", userID)
	}
}

func TestCreateAndGetUser(t *testing.T) {
	user := createTestUser(t)
	defer deleteTestUser(t, user.ID)

	fetchedUser, err := database.GetUserByID(user.ID)
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

func TestCreateAndGetCart(t *testing.T) {
	user := createTestUser(t)
	defer deleteTestUser(t, user.ID)

	cartID, err := database.CreateCart(user.ID)
	if err != nil {
		t.Fatalf("Ошибка при создании корзины: %v", err)
	}

	cart, err := database.GetCartByUserID(user.ID)
	if err != nil {
		t.Fatalf("Ошибка при получении корзины: %v", err)
	}

	if cart.UserID != user.ID {
		t.Errorf("Ожидали userID=%d, получили=%d", user.ID, cart.UserID)
	}
	if cart.ID != cartID {
		t.Errorf("Ожидали id=%d, получили=%d", cartID, cart.ID)
	}
}

func TestCreateAndGetOrder(t *testing.T) {
	user := createTestUser(t)
	defer deleteTestUser(t, user.ID)

	order := models.Order{
		UserID:    user.ID,
		ProductID: 1,
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

	orders, err := database.GetOrderHistory(user.ID)
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
