package config

import (
	"os"
	"log"
	"github.com/joho/godotenv"
)

type Config struct {
	Server		Server
	Database	Database
}

type Server struct {
	Port 	string
}

type Database struct {
	Uri				string
	DatabaseName	string
}

func (c *Config) Read() {

	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("Set 'MONGODB_URI' environment variable")
	}

	c.Server.Port = ":8000"
	c.Database.Uri = uri
	c.Database.DatabaseName = "demo"

}