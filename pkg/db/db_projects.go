package db 

import (
	_ "context"
	_ "encoding/json"
	_ "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type Project struct {
	ID			string `json:"id"`
	Year 		string `json:"year"`
	Name 		string `json:"name"`
	Description string `json:"description"`
	Type 		string `json:"type"`
	StartDate   string `json:"start_date"`
	EndDate		string `json:"end_date"`
}

