package main_test

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/mixpeal/go-dataset/storage"
	"github.com/mixpeal/go-dataset"
)

var a main.Repository

func TestMain(m *testing.M) {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	storage.NewConnection(config)

	code := m.Run()
	os.Exit(code)
}