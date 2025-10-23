package services

import (
	"api-rect-go/db"
	models "api-rect-go/modals"
)

func GetAllProducts() ([]models.Product, error) {
	var product []models.Product
	result := db.DB.Find(&product)
	return product , result.Error
}

func CreateProduct(product *models.Product) error{
	result := db.DB.Create(product)
	return result.Error
}