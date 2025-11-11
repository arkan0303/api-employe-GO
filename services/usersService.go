package services

import (
	models "api-rect-go/modals/mysql"
	"io"

	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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


func (s *ServiceUserService) CreateUser(user *models.ServiceUser) error {
    // 1. Simpan ke database terlebih dahulu
    if err := s.DB.Create(user).Error; err != nil {
        return fmt.Errorf("gagal menyimpan user ke database: %v", err)
    }

    // 2. Kirim ke API eksternal secara asynchronous
    go func() {
        payload := map[string]interface{}{
            "id":                  user.ID,
            "master_companies_id": user.MasterCompaniesID,
            "iwo_templates_id":    user.IwoTemplatesID,
            "full_name":           user.FullName,
            "home_address":        user.HomeAddress,
            "work_address":        user.WorkAddress,
            "email":               user.Email,
            "phone_number":        user.PhoneNumber,
            "status":             user.Status,
            "created_at":         user.CreatedAt.Format(time.RFC3339),
            "updated_at":         user.UpdatedAt.Format(time.RFC3339),
        }

        payloadBytes, err := json.Marshal(payload)
        if err != nil {
            log.Printf("Gagal mengencode payload: %v", err)
            return
        }

        req, err := http.NewRequest("POST", "https://backend.sigapdriver.com/api/create_users_recruitment", 
            bytes.NewBuffer(payloadBytes))
        if err != nil {
            log.Printf("Gagal membuat request ke API eksternal: %v", err)
            return
        }
        req.Header.Set("Content-Type", "application/json")

        client := &http.Client{Timeout: 10 * time.Second}
        resp, err := client.Do(req)
        if err != nil {
            log.Printf("Gagal mengirim data ke API eksternal: %v", err)
            return
        }
        defer resp.Body.Close()

        // Baca response untuk logging
        body, _ := io.ReadAll(resp.Body)
        log.Printf("Response dari API eksternal - Status: %d, Body: %s", resp.StatusCode, string(body))
    }()

    return nil
}
