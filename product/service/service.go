package service

import (
	"product/database"
	"product/models"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) GetProductByID(id int) (*models.Product, error) {
	return database.GetProductByID(id)
}

func (s *Service) AddProduct(product models.Product) (int, error) {
	return database.CreateProduct(product)
}
func (s *Service) SearchProducts(query string, limit, offset int) ([]*models.Product, int, error) {
	return database.SearchProductsByTitle(query, limit, offset)
}

func (s *Service) DeleteProduct(id int) (bool, error) {
	return database.DeleteProduct(id)
}
