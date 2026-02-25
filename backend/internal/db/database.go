// Package db manages the lifecycle of the database connection and schema migrations
package db

import (
	"log"
	"os"

	"github.com/IZO-Ong/gridgo/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is a global singleton instance of the GORM database connection
var DB *gorm.DB

// InitDB establishes a connection to the PostgreSQL database using the 
// DATABASE_URL environment variable and performs a schema auto-migration
//
// If the connection fails or the migration cannot be completed, the 
// application will log a fatal error and exit
func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	var err error
	
	// Initialise connection using Postgres driver
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Running database migrations...")
	err = DB.AutoMigrate(
		&models.User{},
		&models.Maze{},
		&models.PendingUser{},
		&models.Post{},
		&models.Comment{},
		&models.PostVote{},
		&models.CommentVote{},
	)
	
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	
	log.Println("Database connection and migration successful.")
}