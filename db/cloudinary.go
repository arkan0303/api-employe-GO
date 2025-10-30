package db

import (
	"context"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var Cloud *cloudinary.Cloudinary

func InitCloudinary() {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		log.Fatalf("❌ Gagal konek ke Cloudinary: %v", err)
	}
	Cloud = cld
	log.Println("✅ Cloudinary connected")
}

// UploadFile mengupload file ke Cloudinary dan mengembalikan URL-nya
func UploadFile(fileHeader *multipart.FileHeader, folder string) (string, error) {
	// Buka file
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Buat file sementara
	tempFile, err := os.CreateTemp("", "upload-*"+filepath.Ext(fileHeader.Filename))
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFile.Name())

	// Salin isi file ke file sementara
	if _, err = io.Copy(tempFile, src); err != nil {
		return "", err
	}

	// Upload ke Cloudinary
	uploadResult, err := Cloud.Upload.Upload(
		context.Background(),
		tempFile.Name(),
		uploader.UploadParams{
			Folder:   folder,
			PublicID: strings.TrimSuffix(fileHeader.Filename, filepath.Ext(fileHeader.Filename)),
		})

	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}

// UploadFromBase64 mengupload file dari base64 string ke Cloudinary
func UploadFromBase64(base64Data, folder, filename string) (string, error) {
	ctx := context.Background()
	
	// Upload ke Cloudinary
	uploadResult, err := Cloud.Upload.Upload(
		ctx,
		base64Data,
		uploader.UploadParams{
			Folder:   folder,
			PublicID: filename,
		})

	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
