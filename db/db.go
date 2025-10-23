package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	DBMySQL    *gorm.DB
)

// ConnectAll membuka koneksi ke PostgreSQL dan MySQL
func ConnectAll() {
	connectPostgres()
	connectMySQL()
}

func connectPostgres() {
	dsnPostgres := os.Getenv("DB_URL")

	db, err := gorm.Open(postgres.Open(dsnPostgres), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Gagal konek ke PostgreSQL: %v", err)
	}
	DB = db
	fmt.Println("✅ PostgreSQL connection successful")
}

func connectMySQL() {
	dsnMySQL := os.Getenv("DB_URL_MYSQL")

	db, err := gorm.Open(mysql.Open(dsnMySQL), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Gagal konek ke MySQL: %v", err)
	}
	DBMySQL = db
	fmt.Println("✅ MySQL connection successful")
}
