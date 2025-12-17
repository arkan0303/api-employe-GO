package services

import (
	"api-rect-go/db"
	"api-rect-go/dto"
	// models "api-rect-go/modals"
)

func GetData()([]dto.MasterCompany , error) {
	var company []dto.MasterCompany
	result := db.DBMySQL.Find(&company)
	return company , result.Error
}