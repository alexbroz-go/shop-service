package dao_test

import (
	"os"
	"testing"

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
	user := models.User{
		Login:    "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	// Создаем пользователя
	id, err := database.CreateUser(user)
	if err != nil {
		t.Fatalf("Ошибка при создании пользователя: %v", err)
	}

	// Получаем по ID
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

	// Удаляем пользователя
	// (если есть метод DeleteUser, его можно вызвать, иначе оставить)
}

// Тест для создания и получения корзины
func TestCreateAndGetCart(t *testing.T) {
	// Создаем пользователя для корзины
	user := models.User{
		Login:    "cartuser",
		Email:    "cartuser@example.com",
		Password: "pass",
	}
	userID, err := database.CreateUser(user)
	if err != nil {
		t.Fatalf("Ошибка при создании пользователя: %v", err)
	}
	defer func() {
		// Можно добавить метод удаления пользователя, если есть
	}()

	// Создаем корзину для пользователя
	cartID, err := database.CreateCart(userID)
	if err != nil {
		t.Fatalf("Ошибка при создании корзины: %v", err)
	}

	// Получаем корзину по userID
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

	// Можно добавить удаление корзины, если есть метод
}

// Тест для создания и получения заказа
func TestCreateAndGetOrder(t *testing.T) {
	// Создаем пользователя
	user := models.User{
		Login:    "orderuser",
		Email:    "orderuser@example.com",
		Password: "pass",
	}
	userID, err := database.CreateUser(user)
	if err != nil {
		t.Fatalf("Ошибка при создании пользователя: %v", err)
	}
	defer func() {
		// Можно добавить удаление пользователя
	}()

	// Создаем товар (предположим, что есть товар с id=1, или создайте его)
	// В вашем случае, возможно, нужно создать товар в базе, если есть таблица товаров.
	// Для примера возьмем ProductID=1
	productID := 1

	// Создаем заказ
	order := models.Order{
		UserID:    userID,
		ProductID: productID,
		Status:    "pending",
	}
	orderID, err := database.CreateOrder(order)
	if err != nil {
		t.Fatalf("Ошибка при создании заказа: %v", err)
	}

	// Получаем заказ по ID
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

	// Получение истории заказов пользователя
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
