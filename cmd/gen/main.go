package main

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	// === 1️⃣ Koneksi ke PostgreSQL ===
	dsnPostgres := os.Getenv("DB_URL")
	dbPostgres, err := gorm.Open(postgres.Open(dsnPostgres), &gorm.Config{})
	if err != nil {
		panic("❌ gagal konek ke PostgreSQL: " + err.Error())
	}
	fmt.Println("✅ Koneksi PostgreSQL berhasil")

	// === 2️⃣ Koneksi ke MySQL ===
	dsnMySQL := os.Getenv("DB_URL_MYSQL")
	dbMySQL, err := gorm.Open(mysql.Open(dsnMySQL), &gorm.Config{})
	if err != nil {
		panic("❌ gagal konek ke MySQL: " + err.Error())
	}
	fmt.Println("✅ Koneksi MySQL berhasil")

	// === 3️⃣ Generate untuk PostgreSQL ===
	gPostgres := gen.NewGenerator(gen.Config{
		OutPath:      "./query",  // hasil query builder PostgreSQL
		ModelPkgPath: "./modals/postgres", // hasil model struct PostgreSQL
		Mode:         gen.WithDefaultQuery | gen.WithQueryInterface,
	})
	gPostgres.UseDB(dbPostgres)
	gPostgres.GenerateAllTable()
	gPostgres.Execute()
	fmt.Println("✅ Generate model PostgreSQL selesai")

	// === 4️⃣ Generate untuk MySQL ===
	gMySQL := gen.NewGenerator(gen.Config{
		OutPath:      "./query",  // hasil query builder MySQL
		ModelPkgPath: "./modals/mysql", // hasil model struct MySQL
		Mode:         gen.WithDefaultQuery | gen.WithQueryInterface,
	})
	gMySQL.UseDB(dbMySQL)
	gMySQL.GenerateAllTable()
	gMySQL.Execute()
	fmt.Println("✅ Generate model MySQL selesai")

	fmt.Println("🎉 Semua generate selesai!")
}
