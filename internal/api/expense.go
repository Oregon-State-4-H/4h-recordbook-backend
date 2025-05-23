package api

import (
	"4h-recordbook-backend/internal/utils"
	"4h-recordbook-backend/pkg/db"
	"strconv"

	"github.com/beevik/guid"
	"github.com/gin-gonic/gin"
)

type GetExpensesOutput struct {
	Expenses []db.Expense `json:"expenses"`
	Next     string       `json:"next"`
}

type GetExpenseOutput struct {
	Expense db.Expense `json:"expense"`
}

type UpsertExpenseInput struct {
	Date      string   `json:"date" validate:"required"`
	Items     string   `json:"items" validate:"required"`
	Quantity  *float64 `json:"quantity" validate:"required"`
	Cost      *float64 `json:"cost" validate:"required"`
	ProjectID string   `json:"project_id" validate:"required"`
}

type UpsertExpenseOutput GetExpenseOutput

// GetExpenses godoc
// @Summary Get expenses by project
// @Description Gets all of a user's expenses given a project ID
// @Tags Expense
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param projectID path string true "Project ID"
// @Param page query int false "Page number, default 0"
// @Param per_page query int false "Max number of items to return. Can be [1-200], default 100"
// @Param sort_by_newest query bool false "Sort results by most recently added, default false"
// @Success 200 {object} api.GetExpensesOutput
// @Failure 400
// @Failure 401
// @Router /project/{projectID}/expense [get]
func (e *env) getExpenses(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	projectID := c.Param("projectID")

	var output GetExpensesOutput

	paginationOptions := db.PaginationOptions{
		Page:         c.GetInt(CONTEXT_KEY_PAGE),
		PerPage:      c.GetInt(CONTEXT_KEY_PER_PAGE),
		SortByNewest: c.GetBool(CONTEXT_KEY_SORT_BY_NEWEST),
	}

	output.Expenses, err = e.db.GetExpensesByProject(c.Request.Context(), claims.ID, projectID, paginationOptions)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if len(output.Expenses) == paginationOptions.PerPage {

		queryParamsMap := make(map[string]string)
		queryParamsMap[CONTEXT_KEY_PAGE] = strconv.Itoa(paginationOptions.Page + 1)
		queryParamsMap[CONTEXT_KEY_PER_PAGE] = strconv.Itoa(paginationOptions.PerPage)
		queryParamsMap[CONTEXT_KEY_SORT_BY_NEWEST] = strconv.FormatBool(paginationOptions.SortByNewest)

		nextUrlInput := utils.NextUrlInput{
			Context:     c,
			QueryParams: queryParamsMap,
		}

		output.Next = utils.BuildNextUrl(nextUrlInput)
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
// @Param expenseID path string true "Expense ID"
// @Success 200 {object} api.GetExpenseOutput
// @Failure 401
// @Failure 404
// @Router /expense/{expenseID} [get]
func (e *env) getExpense(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	expenseID := c.Param("expenseID")

	var output GetExpenseOutput

	output.Expense, err = e.db.GetExpenseByID(c.Request.Context(), claims.ID, expenseID)
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
// @Success 201 {object} api.UpsertExpenseOutput
// @Failure 400
// @Failure 401
// @Router /expense [post]
func (e *env) addExpense(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertExpenseInput
	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrBadRequest,
		})
		return
	}

	err = e.validator.Struct(input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrMissingFields,
		})
		return
	}

	date, err := utils.StringToTimestamp(input.Date)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrBadDate,
		})
		return
	}

	g := guid.New()
	timestamp := utils.TimeNow()

	expense := db.Expense{
		ID:        g.String(),
		Date:      date.String(),
		Items:     input.Items,
		Quantity:  *input.Quantity,
		Cost:      *input.Cost,
		ProjectID: input.ProjectID,
		UserID:    claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: timestamp.String(),
			Updated: timestamp.String(),
		},
	}

	var output UpsertExpenseOutput

	output.Expense, err = e.db.UpsertExpense(c.Request.Context(), expense)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(201, output)

}

// UpdateExpense godoc
// @Summary Update an expense
// @Description Updates a user's expense information
// @Tags Expense
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param expenseID path string true "Expense ID"
// @Param UpsertExpenseInput body api.UpsertExpenseInput true "Expense information"
// @Success 200 {object} api.UpsertExpenseOutput
// @Failure 400
// @Failure 401
// @Failure 404
// @Router /expense/{expenseID} [put]
func (e *env) updateExpense(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertExpenseInput
	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrBadRequest,
		})
		return
	}

	err = e.validator.Struct(input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrMissingFields,
		})
		return
	}

	date, err := utils.StringToTimestamp(input.Date)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrBadDate,
		})
		return
	}

	expenseID := c.Param("expenseID")

	expense, err := e.db.GetExpenseByID(c.Request.Context(), claims.ID, expenseID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedExpense := db.Expense{
		ID:        expense.ID,
		Date:      date.String(),
		Items:     input.Items,
		Quantity:  *input.Quantity,
		Cost:      *input.Cost,
		ProjectID: expense.ProjectID,
		UserID:    claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: expense.Created,
			Updated: timestamp.String(),
		},
	}

	var output UpsertExpenseOutput

	output.Expense, err = e.db.UpsertExpense(c.Request.Context(), updatedExpense)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// DeleteExpense godoc
// @Summary Removes an expense
// @Description Deletes a user's expense given the expense ID
// @Tags Expense
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param expenseID path string true "Expense ID"
// @Success 204
// @Failure 401
// @Failure 404
// @Router /expense/{expenseID} [delete]
func (e *env) deleteExpense(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	expenseID := c.Param("expenseID")

	response, err := e.db.RemoveExpense(c.Request.Context(), claims.ID, expenseID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}
