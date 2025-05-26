package database

import (
	"database/sql"
	"fmt"
	"product/config"
	"product/models"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() error {
	cfg := config.GetConfig()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	if err := DB.Ping(); err != nil {
		return err
	}

	return createTable()
}

func createTable() error {
	_, err := DB.Exec(`
        CREATE TABLE IF NOT EXISTS products (
            id SERIAL PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            description TEXT,
            price NUMERIC(10,2) NOT NULL,
            count INTEGER NOT NULL
        );
    `)
	return err
}

func GetProductByTitle(title string) (*models.Product, error) {
	p := &models.Product{}
	err := DB.QueryRow("SELECT id, title, description, price, count FROM products WHERE title=$1", title).
		Scan(&p.ID, &p.Title, &p.Description, &p.Price, &p.Count)
	return p, err
}

func DeleteProduct(id int) (bool, error) {
	result, err := DB.Exec("DELETE FROM products WHERE id=$1", id)
	if err != nil {
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected > 0, nil
}

func CreateProduct(p models.Product) (int, error) {
	var id int
	err := DB.QueryRow(
		"INSERT INTO products (title, description, price, count) VALUES ($1, $2, $3, $4) RETURNING id",
		p.Title, p.Description, p.Price, p.Count,
	).Scan(&id)
	return id, err
}

func GetProductByID(id int) (*models.Product, error) {
	p := &models.Product{}
	err := DB.QueryRow("SELECT id, title, description, price, count FROM products WHERE id=$1", id).
		Scan(&p.ID, &p.Title, &p.Description, &p.Price, &p.Count)
	return p, err
}

func GetAllProducts() ([]*models.Product, error) {
	rows, err := DB.Query("SELECT id, title, description, price, count FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		p := &models.Product{}
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Price, &p.Count); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
