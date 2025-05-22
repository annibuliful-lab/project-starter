package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-jet/jet/v2/generator/postgres"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	fmt.Println("aaaa", os.Getenv("POSTGRES_HOST"))
	err = postgres.Generate("../../backend/src/.gen/",
		postgres.DBConnection{
			Host:       os.Getenv("POSTGRES_HOST"),
			Port:       5432,
			User:       os.Getenv("POSTGRES_USER"),
			Password:   os.Getenv("POSTGRES_PASSWORD"),
			DBName:     os.Getenv("POSTGRES_DB"),
			SchemaName: "public",
			SslMode:    "disable",
		})

	if err != nil {
		panic(err)
	}
}
