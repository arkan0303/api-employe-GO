package main

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func mainn() {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("❌ failed to connect to PostgreSQL: %v", err))
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:        "./query",
		ModelPkgPath:   "./models",
		Mode:           gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(db)
	if err := g.GenerateAllTable(); err != nil {
		panic(fmt.Sprintf("❌ failed to generate models and queries: %v", err))
	}

	// if err := g.Execute(); err != nil {
	// 	panic(fmt.Sprintf("❌ failed to execute generation: %v", err))
	// }

	fmt.Println("✅ Generation completed, models and queries files created successfully!")
}