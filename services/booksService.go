package services

import (
	"api-rect-go/db"
	models "api-rect-go/modals"
)

func GetAllBooks() ([]models.Book , error){
	var books []models.Book
	result := db.DB.Find(&books)
	return books , result.Error
}

func CreateBook(books *models.Book) error{
	result := db.DB.Create(books)
	return result.Error	
}