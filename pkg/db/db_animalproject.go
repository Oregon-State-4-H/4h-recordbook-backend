package db

import (
	_"context"
	_ "encoding/json"
	_ "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type Animal struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Species string `json:"species"`
	BirthDate string `json:"birth_date"`
	PurchaseDate string `json:"purchase_date"`
	SireBreed string `json:"sire_breed"`
	DamBreed string `json:"dam_breed"`
	BeginningWeight float64 `json:"beginning_weight"`
	EndWeight float64 `json:"end_weight"`
	EndDate string `json:"end_date"`
	AnimalCost string `json:"animal_cost"`
	SalePrice string `json:"sale_price"`
	YieldGrade string `json:"yield_grade"`
	QualityGrade string `json:"quality_grade"`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

type AnimalProjectRelation struct {
	ID string `json:"id"`
	AnimalID string `json:"animalid"`
	ProjectID string `json:"projectid"`
	GenericDatabaseInfo
}
