package services

import (
	"gorm.io/gorm"
)

type ServiceUserService struct {
	DB *gorm.DB
}

func NewServiceUserService(db *gorm.DB) *ServiceUserService {
	return &ServiceUserService{DB: db}
}

type ServiceUsers struct {
	ID                int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	FullName          string    `gorm:"column:full_name;not null" json:"full_name"`
}

func (s *ServiceUserService) GetByIwoTemplateID(iwoTemplateID int32) ([]ServiceUsers, error) {
	var users []ServiceUsers

	err := s.DB.
		Select("id, full_name").
		Where("iwo_templates_id = ?", iwoTemplateID).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}
