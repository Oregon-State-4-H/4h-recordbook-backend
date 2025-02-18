package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"4h-recordbook-backend/internal/utils"
	"4h-recordbook-backend/pkg/db"
	"github.com/beevik/guid"
)

type UpsertAnimalInput struct {
	Name string `json:"name" validate:"required"`
	Species string `json:"species" validate:"required"`
	BirthDate string `json:"birth_date" validate:"required"`
	PurchaseDate string `json:"purchase_date" validate:"required"`
	SireBreed string `json:"sire_breed" validate:"required"`
	DamBreed string `json:"dam_breed" validate:"required"`
	AnimalCost string `json:"animal_cost" validate:"required"`
	SalePrice string `json:"sale_price" validate:"required"`
	YieldGrade string `json:"yield_grade" validate:"required"`
	QualityGrade string `json:"quality_grade" validate:"required"`
	ProjectID string `json:"projectid" validate:"required"` 
}

type UpdateRateOfGainInput struct {
	BeginningWeight *float64 `json:"beginning_weight" validate:"required"`
	BeginningDate string `json:"beginning_date" validate:"required"`
	EndWeight *float64 `json:"end_weight" validate:"required"`
	EndDate string `json:"end_date" validate:"required"`
}

type GetAnimalsOutput struct {
	Animals []db.Animal `json:"animals"`
}

type GetAnimalOutput struct {
	Animal db.Animal `json:"animal"`
}

// GetAnimals godoc
// @Summary Get animals by project
// @Description Gets all of a user's animals for a given project
// @Tags Animal
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param projectId query string true "Project ID"
// @Success 200 {object} api.GetAnimalsOutput
// @Failure 400
// @Failure 401
// @Router /animal [get]
func (e *env) getAnimals(c *gin.Context) {
	
	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	projectID := c.DefaultQuery("projectId", "")
	if projectID == "" {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	var output GetAnimalsOutput

	output.Animals, err = e.db.GetAnimalsByProject(context.TODO(), claims.ID, projectID)
	if err != nil {
		e.logger.Info(err)
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// GetAnimal godoc
// @Summary Get an animal
// @Description Get a user's animal by ID
// @Tags Animal
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param animalId path string true "Animal ID"
// @Success 200 {object} api.GetAnimalOutput
// @Failure 401
// @Failure 404
// @Router /animal/{animalId} [get]
func (e *env) getAnimal(c *gin.Context) {
	
	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	id := c.Param("animalId")

	var output GetAnimalOutput

	output.Animal, err = e.db.GetAnimalByID(context.TODO(), claims.ID, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// AddAnimal godoc
// @Summary Add an animal
// @Description Adds an animal to a user's personal records
// @Tags Animal
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UpsertAnimalInput body api.UpsertAnimalInput true "Animal information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /animal [post]
func (e *env) addAnimal(c *gin.Context) {
	
	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertAnimalInput
	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	err = e.validator.Struct(input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	birthDate, err := utils.StringToTimestamp(input.BirthDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	purchaseDate, err := utils.StringToTimestamp(input.PurchaseDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	g := guid.New()
	timestamp := utils.TimeNow()

	animal := db.Animal{
		ID: g.String(),
		Name: input.Name,
		Species: input.Species,
		BirthDate: birthDate.ToString(),
		PurchaseDate: purchaseDate.ToString(),
		SireBreed: input.SireBreed,
		DamBreed: input.DamBreed,
		AnimalCost: input.AnimalCost,
		SalePrice: input.SalePrice,
		YieldGrade: input.YieldGrade,
		QualityGrade: input.QualityGrade,
		BeginningWeight: 0,
		BeginningDate: "",
		EndWeight: 0,
		EndDate: "",
		ProjectID: input.ProjectID,
		UserID: claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo {
			Created: timestamp.ToString(),
			Updated: timestamp.ToString(),
		},
	}

	response, err := e.db.UpsertAnimal(context.TODO(), animal)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

// UpdateAnimal godoc
// @Summary Update an animal
// @Description Updates a user's animal information
// @Tags Animal
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param animalId path string true "Animal ID"
// @Param UpsertAnimalInput body api.UpsertAnimalInput true "Animal information"
// @Success 204 
// @Failure 400
// @Failure 401
// @Failure 404
// @Router /animal/{animalId} [put]
func (e *env) updateAnimal(c *gin.Context) {
	
	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertAnimalInput
	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}
	
	err = e.validator.Struct(input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	birthDate, err := utils.StringToTimestamp(input.BirthDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	purchaseDate, err := utils.StringToTimestamp(input.PurchaseDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	id := c.Param("animalId")

	animal, err := e.db.GetAnimalByID(context.TODO(), claims.ID, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedAnimal := db.Animal{
		ID: animal.ID,
		Name: input.Name,
		Species: input.Species,
		BirthDate: birthDate.ToString(),
		PurchaseDate: purchaseDate.ToString(),
		SireBreed: input.SireBreed,
		DamBreed: input.DamBreed,
		AnimalCost: input.AnimalCost,
		SalePrice: input.SalePrice,
		YieldGrade: input.YieldGrade,
		QualityGrade: input.QualityGrade,
		BeginningWeight: animal.BeginningWeight,
		BeginningDate: animal.BeginningDate,
		EndWeight: animal.EndWeight,
		EndDate: animal.EndDate,
		ProjectID: animal.ProjectID,
		UserID: claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo {
			Created: animal.Created,
			Updated: timestamp.ToString(),
		},
	}

	response, err := e.db.UpsertAnimal(context.TODO(), updatedAnimal)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

// UpdateRateOfGain godoc
// @Summary Update an animal's rate of gain
// @Description Updates a user's animal rate of gain information
// @Tags Animal
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param animalId path string true "Animal ID"
// @Param UpdateRateOfGainInput body api.UpdateRateOfGainInput true "Animal rate of gain information"
// @Success 204 
// @Failure 400
// @Failure 401
// @Failure 404
// @Router /rate-of-gain/{animalId} [put]
func (e *env) updateRateOfGain(c *gin.Context) {
	
	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpdateRateOfGainInput
	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	err = e.validator.Struct(input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	beginningDate, err := utils.StringToTimestamp(input.BeginningDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	endDate, err := utils.StringToTimestamp(input.EndDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	id := c.Param("animalId")

	animal, err := e.db.GetAnimalByID(context.TODO(), claims.ID, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedAnimal := db.Animal{
		ID: animal.ID,
		Name: animal.Name,
		Species: animal.Species,
		BirthDate: animal.BirthDate,
		PurchaseDate: animal.PurchaseDate,
		SireBreed: animal.SireBreed,
		DamBreed: animal.DamBreed,
		AnimalCost: animal.AnimalCost,
		SalePrice: animal.SalePrice,
		YieldGrade: animal.YieldGrade,
		QualityGrade: animal.QualityGrade,
		BeginningWeight: *input.BeginningWeight,
		BeginningDate: beginningDate.ToString(),
		EndWeight: *input.EndWeight,
		EndDate: endDate.ToString(),
		ProjectID: animal.ProjectID,
		UserID: claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo {
			Created: animal.Created,
			Updated: timestamp.ToString(),
		},
	}

	response, err := e.db.UpsertAnimal(context.TODO(), updatedAnimal)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}