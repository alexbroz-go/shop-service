package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"product/config"
	"product/models"
	"regexp"
	"sort"

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

	if os.Getenv("MIGRATE_ON_STARTUP") == "1" {
		if err := applySQLMigrations("migrations"); err != nil {
			return fmt.Errorf("apply migrations: %w", err)
		}
	}
	return nil
}

func applySQLMigrations(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	re := regexp.MustCompile(`^\d+_.*\.up\.sql$`)
	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if re.MatchString(name) {
			files = append(files, name)
		}
	}
	sort.Strings(files)
	for _, fname := range files {
		path := filepath.Join(dir, fname)
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if len(content) == 0 {
			continue
		}
		if _, err := DB.Exec(string(content)); err != nil {
			return fmt.Errorf("migration %s: %w", fname, err)
		}
	}
	return nil
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
		"INSERT INTO products (title, description, price) VALUES ($1, $2, $3) RETURNING id",
		p.Title, p.Description, p.Price,
	).Scan(&id)
	return id, err
}

func GetProductByID(id int) (*models.Product, error) {
	p := &models.Product{}
	err := DB.QueryRow("SELECT id, title, description, price FROM products WHERE id=$1", id).
		Scan(&p.ID, &p.Title, &p.Description, &p.Price)
	return p, err
}

func GetAllProducts() ([]*models.Product, error) {
	rows, err := DB.Query("SELECT id, title, description, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		p := &models.Product{}
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func SearchProductsByTitle(keyword string, limit, offset int) ([]*models.Product, int, error) {
	// Общее количество
	var totalCount int
	err := DB.QueryRow("SELECT COUNT(*) FROM products WHERE title ILIKE '%' || $1 || '%'", keyword).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	rows, err := DB.Query(`
		SELECT id, title, description, price
		FROM products
		WHERE title ILIKE '%' || $1 || '%'
		LIMIT $2 OFFSET $3
	`, keyword, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		p := &models.Product{}
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Price); err != nil {
			return nil, 0, err
		}
		products = append(products, p)
	}
	return products, totalCount, nil
}
