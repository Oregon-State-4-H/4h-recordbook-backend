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

type Auth0 struct {
	Secret         string
	ClientID       string
	ClientSecret   string
	BaseURL        string
	IssuerBaseURL  string
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
	auth0_secret = os.Getenv("AUTH0_SECRET")
	auth0_clientID = os.Getenv("AUTH0_BASE_URL")
	auth0_clientSecret = os.Getenv("AUTH0_ISSUER_BASE_URL")
	auth0_baseURL = os.Getenv("AUTH0_CLIENT_ID")
	auth0_issuerBaseURL = os.Getenv("AUTH0_CLIENT_SECRET")

	c.Server.Port = ":8000"
	c.Database.Uri = uri
	c.Database.DatabaseName = "demo"
	
	c.Auth0 = Auth0{
		auth0_secret,
		auth0_clientID,
                auth0_clientSecret,
                auth0_baseURL,
                auth0_issuerBaseURL
	}
}
