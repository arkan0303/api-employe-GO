package services

import (
	"api-rect-go/db"
)

type IwoMinimal struct {
	ID          int32  `json:"id"`
	IwoName     string `json:"iwo_name"`
	BasicSalary string `json:"basic_salary"`
}

func GetIwoDataByMasterCompaniesID(masterCompaniesID int32) ([]IwoMinimal, error) {
	var iwos []IwoMinimal
	if err := db.DBMySQL.
		Table("iwo_templates").
		Select("id, iwo_name, basic_salary").
		Where("master_companies_id = ? AND status = ?", masterCompaniesID, 1).
		Find(&iwos).Error; err != nil {
		return nil, err
	}
	return iwos, nil
}
