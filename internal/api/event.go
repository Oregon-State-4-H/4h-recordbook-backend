package api

import (
	"4h-recordbook-backend/internal/utils"
	"4h-recordbook-backend/pkg/db"
	"context"

	"github.com/beevik/guid"
	"github.com/gin-gonic/gin"
)

type GetEventsOutput struct {
	Events []db.Event `json:"events"`
}

type GetEventOutput struct {
	Event db.Event `json:"event"`
}

type UpsertEventInput struct {
	Name        string `json:"name"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

type UpsertEventOutput GetEventOutput

/*******************************
* FULL EVENTS
********************************/

// GetEvents godoc
// @Summary Get events by user
// @Description Gets all of a user's events
// @Tags Event
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} api.GetEventsOutput
// @Failure 400
// @Failure 401
// @Router /event [get]
func (e *env) getEvents(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var output GetEventsOutput

	output.Events, err = e.db.GetEventsByUser(context.TODO(), claims.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// GetEvent godoc
// @Summary Get an event
// @Description Get a user's event by ID
// @Tags Event
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param eventID path string true "Event ID"
// @Success 200 {object} api.GetEventOutput
// @Failure 401
// @Failure 404
// @Router /event/{eventID} [get]
func (e *env) getEvent(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	eventID := c.Param("eventID")

	var output GetEventOutput

	output.Event, err = e.db.GetEventByID(context.TODO(), claims.ID, eventID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// AddEvent godoc
// @Summary Adds an event
// @Description Adds an event to a user's personal records
// @Tags Event
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UpsertEventInput body api.UpsertEventInput true "General event information"
// @Success 201 {object} api.UpsertEventOutput
// @Failure 400
// @Failure 401
// @Router /event [post]
func (e *env) addEvent(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertEventInput
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

	startDate, err := utils.StringToTimestamp(input.StartDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrBadDate,
		})
		return
	}

	endDate, err := utils.StringToTimestamp(input.EndDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrBadDate,
		})
		return
	}

	g := guid.New()
	timestamp := utils.TimeNow()

	event := db.Event{
		ID:          g.String(),
		Name:        input.Name,
		StartDate:   startDate.String(),
		EndDate:     endDate.String(),
		Location:    input.Location,
		Description: input.Description,
		Sections:    make(map[string]int),
		UserID:      claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: timestamp.String(),
			Updated: timestamp.String(),
		},
	}

	var output UpsertEventOutput

	output.Event, err = e.db.UpsertEvent(context.TODO(), event)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(201, output)

}

// UpdateEvent godoc
// @Summary Update an event
// @Description Updates a user's event information
// @Tags Event
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param eventID path string true "Event ID"
// @Param UpsertEventInput body api.UpsertEventInput true "Event information"
// @Success 200 {object} api.UpsertEventOutput
// @Failure 400
// @Failure 401
// @Failure 404
// @Router /event/{eventID} [put]
func (e *env) updateEvent(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertEventInput
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

	startDate, err := utils.StringToTimestamp(input.StartDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrBadDate,
		})
		return
	}

	endDate, err := utils.StringToTimestamp(input.EndDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrBadDate,
		})
		return
	}

	eventID := c.Param("eventID")

	event, err := e.db.GetEventByID(context.TODO(), claims.ID, eventID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedEvent := db.Event{
		ID:          event.ID,
		Name:        input.Name,
		StartDate:   startDate.String(),
		EndDate:     endDate.String(),
		Location:    input.Location,
		Description: input.Description,
		Sections:    event.Sections,
		UserID:      claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: event.Created,
			Updated: timestamp.String(),
		},
	}

	var output UpsertEventOutput

	output.Event, err = e.db.UpsertEvent(context.TODO(), updatedEvent)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// DeleteEvent godoc
// @Summary Removes an event
// @Description Deletes a user's event given the event ID
// @Tags Event
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param eventID path string true "Event ID"
// @Success 204
// @Failure 401
// @Failure 404
// @Router /event/{eventID} [delete]
func (e *env) deleteEvent(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	eventID := c.Param("eventID")

	response, err := e.db.RemoveEvent(context.TODO(), claims.ID, eventID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

/*******************************
* EVENT SECTIONS
********************************/

type AddSectionToEventInput struct {
	SectionNumber *int   `json:"section_number"`
	SectionID     string `json:"section_id"`
}

// AddSectionToEvent godoc
// @Summary Adds a section to an event
// @Description Adds an section to a user's personal event records, updating the event
// @Tags Event
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param eventID path string true "Event ID"
// @Param AddSectionToEventInput body api.AddSectionToEventInput true "Identifying section information"
// @Success 200 {object} api.UpsertEventOutput
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 409
// @Router /event/{eventID} [post]
func (e *env) addSectionToEvent(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input AddSectionToEventInput
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

	if *input.SectionNumber < 1 || *input.SectionNumber > 14 {
		c.JSON(400, gin.H{
			"message": ErrBadRequest,
		})
		return
	}

	eventID := c.Param("eventID")

	event, err := e.db.GetEventByID(context.TODO(), claims.ID, eventID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	_, ok := event.Sections[input.SectionID]
	if ok {
		c.JSON(409, gin.H{
			"message": ErrEventSectionConflict,
		})
		return
	}

	updatedMap := event.Sections
	updatedMap[input.SectionID] = *input.SectionNumber

	timestamp := utils.TimeNow()

	updatedEvent := db.Event{
		ID:          event.ID,
		Name:        event.Name,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		Location:    event.Location,
		Description: event.Location,
		Sections:    updatedMap,
		UserID:      claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: event.Created,
			Updated: timestamp.String(),
		},
	}

	var output UpsertEventOutput

	output.Event, err = e.db.UpsertEvent(context.TODO(), updatedEvent)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// DeleteSectionFromEvent godoc
// @Summary Deletes a section from an event
// @Description Deletes a section from a user's personal event records, updating the event
// @Tags Event
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param eventID path string true "Event ID"
// @Param sectionID path string true "Section ID"
// @Success 200 {object} api.UpsertEventOutput
// @Failure 400
// @Failure 401
// @Failure 404
// @Failure 409
// @Router /event/{eventID}/{sectionID} [delete]
func (e *env) deleteSectionFromEvent(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	eventID := c.Param("eventID")
	sectionID := c.Param("sectionID")

	event, err := e.db.GetEventByID(context.TODO(), claims.ID, eventID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	_, ok := event.Sections[sectionID]
	if !ok {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	updatedMap := event.Sections
	delete(updatedMap, sectionID)

	timestamp := utils.TimeNow()

	updatedEvent := db.Event{
		ID:          event.ID,
		Name:        event.Name,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		Location:    event.Location,
		Description: event.Location,
		Sections:    updatedMap,
		UserID:      claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: event.Created,
			Updated: timestamp.String(),
		},
	}

	var output UpsertEventOutput

	output.Event, err = e.db.UpsertEvent(context.TODO(), updatedEvent)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}
