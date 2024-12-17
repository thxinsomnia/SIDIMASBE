package models

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
)

var DB *gorm.DB

func ConnectDatabase() {
    // Replace with your actual Supabase connection string
    dbURL := "postgresql://postgres:Loveyouelysia3399@db.elsqjtnetkaltswezoxh.supabase.co:5432/postgres"

    db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to the database: %v", err)
    }

    db.AutoMigrate(&User{})
	DB = db
}
