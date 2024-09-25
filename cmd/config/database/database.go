package database

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email          string
	UserName       string
	HashedPassword string
	User_id        uuid.UUID `gorm:"primaryKey"`
	Roles          string
	CreatedAt      time.Time
}

var DB *gorm.DB

func DatabaseConnection() *gorm.DB {
	var err error
	//
	DB, err = gorm.Open(postgres.Open(fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=%s sslmode=disable search_path=%s", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("PORTDB"),os.Getenv("DB_path"))), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database")

	return DB

}
