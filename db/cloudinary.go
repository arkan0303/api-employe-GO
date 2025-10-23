package db

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
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
