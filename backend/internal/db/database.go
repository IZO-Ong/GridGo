package db

import (
	"log"
	"os"

	"github.com/IZO-Ong/gridgo/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
    dsn := os.Getenv("DATABASE_URL")
    var err error
    
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    err = DB.AutoMigrate(&models.User{}, &models.Maze{})
    if err != nil {
        log.Fatal("Migration failed:", err)
    }
}