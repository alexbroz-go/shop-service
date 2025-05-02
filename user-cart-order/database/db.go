package database

import (
	"database/sql"
	"fmt"
	"shop-service/user-cart-order/models"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() error {
	dsn := "host=localhost port=5433 user=user password=password dbname=user_cart_order sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("ошибка подключения: %w", err)
	}
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("ошибка ping: %w", err)
	}
	return createTables()
}

func createTables() error {
	// Таблицы (если не существуют)
	_, err := DB.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            login VARCHAR(50) UNIQUE NOT NULL,
            email VARCHAR(100) UNIQUE NOT NULL,
            password VARCHAR(255) NOT NULL
        );
        CREATE TABLE IF NOT EXISTS carts (
            id SERIAL PRIMARY KEY,
            user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
            created_at TIMESTAMP DEFAULT NOW()
        );
        CREATE TABLE IF NOT EXISTS orders (
            id SERIAL PRIMARY KEY,
            user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
            product_id INTEGER NOT NULL,
            created_at TIMESTAMP DEFAULT NOW(),
            status VARCHAR(20)
        );
    `)
	return err
}

//  Методы Таблицы users

func CreateUser(user models.User) (int, error) {
	var id int
	err := DB.QueryRow(
		"INSERT INTO users (login, email, password) VALUES ($1, $2, $3) RETURNING id",
		user.Login, user.Email, user.Password,
	).Scan(&id)
	return id, err
}

func GetUserByID(id int) (*models.User, error) {
	u := &models.User{}
	err := DB.QueryRow("SELECT id, login, email, password FROM users WHERE id=$1", id).
		Scan(&u.ID, &u.Login, &u.Email, &u.Password)
	return u, err
}

func GetUserByLogin(login string) (*models.User, error) {
	u := &models.User{}
	err := DB.QueryRow("SELECT id, login, email, password FROM users WHERE login=$1", login).
		Scan(&u.ID, &u.Login, &u.Email, &u.Password)
	return u, err
}

//  Методы таблицы Cart

func CreateCart(userID int) (int, error) {
	var id int
	err := DB.QueryRow(
		"INSERT INTO carts (user_id) VALUES ($1) RETURNING id",
		userID,
	).Scan(&id)
	return id, err
}

func GetCartByUserID(userID int) (*models.Cart, error) {
	c := &models.Cart{}
	err := DB.QueryRow("SELECT id, user_id, created_at FROM carts WHERE user_id=$1", userID).
		Scan(&c.ID, &c.UserID, &c.CreatedAt)
	return c, err
}

// Методы таблицы Order

func CreateOrder(order models.Order) (int, error) {
	var id int
	err := DB.QueryRow(
		"INSERT INTO orders (user_id, product_id, status) VALUES ($1, $2, $3) RETURNING id",
		order.UserID, order.ProductID, order.Status,
	).Scan(&id)
	return id, err
}

func GetOrderByID(id int) (*models.Order, error) {
	o := &models.Order{}
	err := DB.QueryRow("SELECT id, user_id, product_id, created_at, status FROM orders WHERE id=$1", id).
		Scan(&o.ID, &o.UserID, &o.ProductID, &o.CreatedAt, &o.Status)
	return o, err
}

func GetOrderHistory(userID int) ([]*models.Order, error) {
	rows, err := DB.Query("SELECT id, user_id, product_id, created_at, status FROM orders WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		o := &models.Order{}
		if err := rows.Scan(&o.ID, &o.UserID, &o.ProductID, &o.CreatedAt, &o.Status); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
