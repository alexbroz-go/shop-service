package dao_test

import (
	"os"
	"product/database"
	"product/models"
	"testing"
)

func TestMain(m *testing.M) {
	err := database.Init()
	if err != nil {
		panic(err)
	}
	// Clean DB for stable tests
	if _, err := database.DB.Exec("TRUNCATE TABLE products RESTART IDENTITY CASCADE;"); err != nil {
		panic(err)
	}
	code := m.Run()
	// Optional cleanup after tests
	_, _ = database.DB.Exec("TRUNCATE TABLE products RESTART IDENTITY CASCADE;")
	os.Exit(code)
}

// Тест для функции CreateProduct и GetProductByID
func TestCreateAndGetProduct(t *testing.T) {
	// Создаем тестовый продукт
	testProduct := models.Product{
		Title:       "Test Product",
		Description: "Test Description",
		Price:       99.99,
	}

	// Создаем продукт в базе
	id, err := database.CreateProduct(testProduct)
	if err != nil {
		t.Fatalf("Ошибка при создании продукта: %v", err)
	}

	// Получаем продукт по ID
	fetchedProduct, err := database.GetProductByID(id)
	if err != nil {
		t.Fatalf("Ошибка при получении продукта: %v", err)
	}

	// Проверяем поля
	if fetchedProduct.Title != testProduct.Title {
		t.Errorf("Ожидали Title=%s, получили=%s", testProduct.Title, fetchedProduct.Title)
	}
	if fetchedProduct.Description != testProduct.Description {
		t.Errorf("Ожидали Description=%s, получили=%s", testProduct.Description, fetchedProduct.Description)
	}
	if fetchedProduct.Price != testProduct.Price {
		t.Errorf("Ожидали Price=%f, получили=%f", testProduct.Price, fetchedProduct.Price)
	}

	// Удаляем тестовый продукт
	success, err := database.DeleteProduct(id)
	if err != nil {
		t.Errorf("Ошибка при удалении продукта: %v", err)
	}
	if !success {
		t.Errorf("Удаление продукта не удалось, id=%d", id)
	}
}

// Тест для функции GetAllProducts
func TestGetAllProducts(t *testing.T) {
	// Проверяем наличие продуктов
	products, err := database.GetAllProducts()
	if err != nil {
		t.Fatalf("Ошибка при получении всех продуктов: %v", err)
	}

	if len(products) == 0 {
		// Если продуктов нет, добавим один
		testProduct := models.Product{
			Title:       "Test Product",
			Description: "Test Description",
			Price:       99.99,
		}
		_, err := database.CreateProduct(testProduct)
		if err != nil {
			t.Fatalf("Ошибка при добавлении тестового продукта: %v", err)
		}
		t.Log("Добавлен тестовый продукт")
	}

	// Повторно получаем все продукты
	products, err = database.GetAllProducts()
	if err != nil {
		t.Fatalf("Ошибка при получении всех продуктов: %v", err)
	}
	t.Logf("Всего продуктов: %d", len(products))
}

// Тест для функции SearchProductsByTitle
func TestSearchProductsByTitle(t *testing.T) {
	// Предварительно добавим тестовый продукт для поиска
	testProduct := models.Product{
		Title:       "UniqueTestTitle123",
		Description: "Desc",
		Price:       10.0,
	}
	id, err := database.CreateProduct(testProduct)
	if err != nil {
		t.Fatalf("Ошибка при создании тестового продукта: %v", err)
	}
	defer database.DeleteProduct(id) // удалим после теста

	// Выполняем поиск по уникальному названию
	products, total, err := database.SearchProductsByTitle("UniqueTestTitle123", 10, 0)
	if err != nil {
		t.Fatalf("Ошибка при поиске: %v", err)
	}
	if total == 0 {
		t.Errorf("Ожидали найти хотя бы один продукт, найдено: %d", total)
	}
	found := false
	for _, p := range products {
		if p.ID == id {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Не нашли добавленный тестовый продукт по названию")
	}
}
