package models

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	ID       uint    `gorm:"primary key;autoIncrement" json: "id"`
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
	Date     time.Time `json:"date"`
	Company  *string `json:"company"`
}

func MigrateUsers(db *gorm.DB) error {
	err := db.AutoMigrate(&Users{})
	return err
}
