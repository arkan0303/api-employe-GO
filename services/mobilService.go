package services

import (
	"api-rect-go/db"
	models "api-rect-go/modals"
)

func GetAllMobils() ([]models.Mobil, error) {
	var mobil []models.Mobil
	result := db.DB.Find(&mobil)
	return mobil, result.Error
}

func CreateMobil(mobil *models.Mobil) error{
	result := db.DB.Create(mobil)
	return result.Error
}