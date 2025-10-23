package services

import (
	"api-rect-go/db"
	models "api-rect-go/modals"
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func GetRequestDrivers() ([]models.TbRequestDriver, error) {
	var request []models.TbRequestDriver
	result := db.DB.Find(&request)
	return request, result.Error
}

// parameter foto di sini akan datang dari controller (form file)
func CreateRequestDriver(request *models.TbRequestDriver, fotoFile multipart.File, fotoHeader *multipart.FileHeader) error {
	ctx := context.Background()

	// upload ke Cloudinary
	uploadResult, err := db.Cloud.Upload.Upload(ctx, fotoFile, uploader.UploadParams{
		Folder: "foto_mobil", // folder Cloudinary kamu
	})
	if err != nil {
		return err
	}

	// simpan URL hasil upload ke kolom foto_mobil
	request.FotoMobil = uploadResult.SecureURL

	// simpan ke DB
	result := db.DB.Create(request)
	return result.Error
}
