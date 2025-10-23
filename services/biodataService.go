package services

import (
	"api-rect-go/db"
	models "api-rect-go/modals"
)


func GetAllBiodatas() ([]models.Biodata , error) {
	var biodatas []models.Biodata
	result := db.DB.Find(&biodatas)
	return biodatas , result.Error
}