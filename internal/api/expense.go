package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"4h-recordbook-backend/internal/utils"
	"4h-recordbook-backend/pkg/db"
	"github.com/beevik/guid"
)

type UpsertExpenseInput struct {
	Date string `json:"date" validate:"required"`
	Items string `json:"items" validate:"required"`
	Quantity *float64 `json:"quantity" validate:"required"`
	Cost *float64 `json:"cost" validate:"required"`
	ProjectID string `json:"projectid" validate:"required"`
}

type GetExpensesOutput struct {
	Expenses []db.Expense `json:"expenses"`
}

type GetExpenseOutput struct {
	Expense db.Expense `json:"expense"`
}

// GetExpenses godoc
// @Summary Get expenses by project
// @Description Gets all of a user's expenses given a project ID
// @Tags Expense
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param projectId query string true "Project ID"
// @Success 200 {object} api.GetExpensesOutput
// @Failure 400
// @Failure 401
// @Router /expense [get]
func (e *env) getExpenses(c *gin.Context) {

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

	var output GetExpensesOutput

	output.Expenses, err = e.db.GetExpensesByProject(context.TODO(), claims.ID, projectID)
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

// GetExpense godoc
// @Summary Get an expense
// @Description Get a user's expense by ID
// @Tags Expense
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param expenseId path string true "Expense ID"
// @Success 200 {object} api.GetExpenseOutput
// @Failure 401
// @Failure 404
// @Router /expense/{expenseId} [get]
func (e *env) getExpense(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	id := c.Param("expenseId")

	var output GetExpenseOutput

	output.Expense, err = e.db.GetExpenseByID(context.TODO(), claims.ID, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// AddExpense godoc
// @Summary Adds an expense
// @Description Adds an expense to a user's personal records
// @Tags Expense
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UpsertExpenseInput body api.UpsertExpenseInput true "Expense information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /expense [post]
func (e *env) addExpense(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertExpenseInput
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

	date, err := utils.StringToTimestamp(input.Date)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	g := guid.New()
	timestamp := utils.TimeNow()

	expense := db.Expense{
		ID: g.String(),
		Date: date.ToString(),
		Items: input.Items,
		Quantity: *input.Quantity,
		Cost: *input.Cost,
		ProjectID: input.ProjectID,
		UserID: claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo {
			Created: timestamp.ToString(),
			Updated: timestamp.ToString(),
		},
	}

	response, err := e.db.UpsertExpense(context.TODO(), expense)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}
