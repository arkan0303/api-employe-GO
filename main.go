package main

import (
	"log"
	"time"

	"api-rect-go/db"
	models "api-rect-go/modals"
	"api-rect-go/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Failed to load environment variables")
	}

	// Koneksi ke database
	db.ConnectAll()
	db.InitCloudinary()

	// Auto migrate tabel
	if err := db.DB.AutoMigrate(&models.Product{}, &models.Mobil{}); err != nil {
		log.Fatalf("Gagal migrate tabel: %v", err)
	}

	r := gin.Default()

	// --- Atur Cors---
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))

	routes.RegisterRoutes(r)

	// Jalankan server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
