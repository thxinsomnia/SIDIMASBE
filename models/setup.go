package models

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
    // Replace with your actual Supabase connection string
    dbURL := "postgresql://postgres.xlrzwrrmqiehdbemlkfq:loveyouelysia3399@aws-0-ap-southeast-1.pooler.supabase.com:5432/postgres"
    db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to the database: %v", err)
    }

    db.AutoMigrate(&User{})
	DB = db
}
