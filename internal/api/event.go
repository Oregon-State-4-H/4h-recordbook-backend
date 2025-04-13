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

type UpsertEventInput struct {
	Name        string `json:"name"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

type UpsertEventOutput struct {
	Event db.Event `json:"event"`
}

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

type GetEventWithSectionsOutput struct {
	Event    db.Event `json:"event"`
	Sections []any    `json:"sections"`
}

type UpsertEventSectionInput struct {
	SectionNumber *int   `json:"section_number"`
	SectionID     string `json:"section_id"`
}

type UpsertEventSectionOutput struct {
	EventSection db.EventSection `json:"event_section"`
}

// GetEventWithSections godoc
// @Summary Get an event with sections
// @Description Get a user's event by ID and includes relevant section data
// @Tags Event
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param eventID path string true "Event ID"
// @Success 200 {object} api.GetEventWithSectionsOutput
// @Failure 401
// @Failure 404
// @Router /event/{eventID} [get]
func (e *env) getEventWithSections(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	eventID := c.Param("eventID")

	var output GetEventWithSectionsOutput

	output.Event, err = e.db.GetEventByID(context.TODO(), claims.ID, eventID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	sections, err := e.db.GetEventSectionsByEvent(context.TODO(), claims.ID, eventID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	for _, section := range sections {
		switch section.SectionNumber {
		case 1:
			sectionInterface, err := e.db.GetSection1ByID(context.TODO(), claims.ID, section.SectionID)
			if err == nil {
				output.Sections = append(output.Sections, sectionInterface)
			}
		case 2:
			sectionInterface, err := e.db.GetSection2ByID(context.TODO(), claims.ID, section.SectionID)
			if err == nil {
				output.Sections = append(output.Sections, sectionInterface)
			}
		case 3:
			sectionInterface, err := e.db.GetSection3ByID(context.TODO(), claims.ID, section.SectionID)
			if err == nil {
				output.Sections = append(output.Sections, sectionInterface)
			}
		case 4:
			sectionInterface, err := e.db.GetSection4ByID(context.TODO(), claims.ID, section.SectionID)
			if err == nil {
				output.Sections = append(output.Sections, sectionInterface)
			}
		case 5:
			sectionInterface, err := e.db.GetSection5ByID(context.TODO(), claims.ID, section.SectionID)
			if err == nil {
				output.Sections = append(output.Sections, sectionInterface)
			}
		case 6:
			sectionInterface, err := e.db.GetSection6ByID(context.TODO(), claims.ID, section.SectionID)
			if err == nil {
				output.Sections = append(output.Sections, sectionInterface)
			}
		case 7:
			sectionInterface, err := e.db.GetSection7ByID(context.TODO(), claims.ID, section.SectionID)
			if err == nil {
				output.Sections = append(output.Sections, sectionInterface)
			}
		case 8:
			sectionInterface, err := e.db.GetSection8ByID(context.TODO(), claims.ID, section.SectionID)
			if err == nil {
				output.Sections = append(output.Sections, sectionInterface)
			}
		case 9:
			sectionInterface, err := e.db.GetSection9ByID(context.TODO(), claims.ID, section.SectionID)
			if err == nil {
				output.Sections = append(output.Sections, sectionInterface)
			}
		case 10:
			sectionInterface, err := e.db.GetSection10ByID(context.TODO(), claims.ID, section.SectionID)
			if err == nil {
				output.Sections = append(output.Sections, sectionInterface)
			}
		case 11:
			sectionInterface, err := e.db.GetSection11ByID(context.TODO(), claims.ID, section.SectionID)
			if err == nil {
				output.Sections = append(output.Sections, sectionInterface)
			}
		case 12:
			sectionInterface, err := e.db.GetSection12ByID(context.TODO(), claims.ID, section.SectionID)
			if err == nil {
				output.Sections = append(output.Sections, sectionInterface)
			}
		case 13:
			sectionInterface, err := e.db.GetSection13ByID(context.TODO(), claims.ID, section.SectionID)
			if err == nil {
				output.Sections = append(output.Sections, sectionInterface)
			}
		case 14:
			sectionInterface, err := e.db.GetSection14ByID(context.TODO(), claims.ID, section.SectionID)
			if err == nil {
				output.Sections = append(output.Sections, sectionInterface)
			}
		}
	}

	c.JSON(200, output)

}

// AddEventSection godoc
// @Summary Adds an event section
// @Description Adds an event section to a user's personal records
// @Tags Event
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param eventID path string true "Event ID"
// @Param UpsertEventSectionInput body api.UpsertEventSectionInput true "Identifying section information"
// @Success 201 {object} api.UpsertEventSectionOutput
// @Failure 400
// @Failure 401
// @Failure 409
// @Router /event/{eventID} [post]
func (e *env) addEventSection(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertEventSectionInput
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

	eventID := c.Param("eventID")

	//verify event exists
	event, err := e.db.GetEventByID(context.TODO(), claims.ID, eventID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	//verify section exists
	if *input.SectionNumber < 1 || *input.SectionNumber > 14 {
		c.JSON(400, gin.H{
			"message": ErrInvalidSectionNumber,
		})
		return
	}

	switch *input.SectionNumber {
	case 1:
		section, err := e.db.GetSection1ByID(context.TODO(), claims.ID, input.SectionID)
		if err != nil {
			response := InterpretCosmosError(err)
			c.JSON(response.Code, gin.H{
				"message": response.Message,
			})
			return
		}
		if section.Section != 1 {
			c.JSON(404, gin.H{
				"message": ErrNotFound,
			})
			return
		}
	case 2:
		section, err := e.db.GetSection2ByID(context.TODO(), claims.ID, input.SectionID)
		if err != nil {
			response := InterpretCosmosError(err)
			c.JSON(response.Code, gin.H{
				"message": response.Message,
			})
			return
		}
		if section.Section != 2 {
			c.JSON(404, gin.H{
				"message": ErrNotFound,
			})
			return
		}
	case 3:
		section, err := e.db.GetSection3ByID(context.TODO(), claims.ID, input.SectionID)
		if err != nil {
			response := InterpretCosmosError(err)
			c.JSON(response.Code, gin.H{
				"message": response.Message,
			})
			return
		}
		if section.Section != 3 {
			c.JSON(404, gin.H{
				"message": ErrNotFound,
			})
			return
		}
	case 4:
		section, err := e.db.GetSection4ByID(context.TODO(), claims.ID, input.SectionID)
		if err != nil {
			response := InterpretCosmosError(err)
			c.JSON(response.Code, gin.H{
				"message": response.Message,
			})
			return
		}
		if section.Section != 4 {
			c.JSON(404, gin.H{
				"message": ErrNotFound,
			})
			return
		}
	case 5:
		section, err := e.db.GetSection5ByID(context.TODO(), claims.ID, input.SectionID)
		if err != nil {
			response := InterpretCosmosError(err)
			c.JSON(response.Code, gin.H{
				"message": response.Message,
			})
			return
		}
		if section.Section != 5 {
			c.JSON(404, gin.H{
				"message": ErrNotFound,
			})
			return
		}
	case 6:
		section, err := e.db.GetSection6ByID(context.TODO(), claims.ID, input.SectionID)
		if err != nil {
			response := InterpretCosmosError(err)
			c.JSON(response.Code, gin.H{
				"message": response.Message,
			})
			return
		}
		if section.Section != 6 {
			c.JSON(404, gin.H{
				"message": ErrNotFound,
			})
			return
		}
	case 7:
		section, err := e.db.GetSection7ByID(context.TODO(), claims.ID, input.SectionID)
		if err != nil {
			response := InterpretCosmosError(err)
			c.JSON(response.Code, gin.H{
				"message": response.Message,
			})
			return
		}
		if section.Section != 7 {
			c.JSON(404, gin.H{
				"message": ErrNotFound,
			})
			return
		}
	case 8:
		section, err := e.db.GetSection8ByID(context.TODO(), claims.ID, input.SectionID)
		if err != nil {
			response := InterpretCosmosError(err)
			c.JSON(response.Code, gin.H{
				"message": response.Message,
			})
			return
		}
		if section.Section != 8 {
			c.JSON(404, gin.H{
				"message": ErrNotFound,
			})
			return
		}
	case 9:
		section, err := e.db.GetSection9ByID(context.TODO(), claims.ID, input.SectionID)
		if err != nil {
			response := InterpretCosmosError(err)
			c.JSON(response.Code, gin.H{
				"message": response.Message,
			})
			return
		}
		if section.Section != 9 {
			c.JSON(404, gin.H{
				"message": ErrNotFound,
			})
			return
		}
	case 10:
		section, err := e.db.GetSection10ByID(context.TODO(), claims.ID, input.SectionID)
		if err != nil {
			response := InterpretCosmosError(err)
			c.JSON(response.Code, gin.H{
				"message": response.Message,
			})
			return
		}
		if section.Section != 10 {
			c.JSON(404, gin.H{
				"message": ErrNotFound,
			})
			return
		}
	case 11:
		section, err := e.db.GetSection11ByID(context.TODO(), claims.ID, input.SectionID)
		if err != nil {
			response := InterpretCosmosError(err)
			c.JSON(response.Code, gin.H{
				"message": response.Message,
			})
			return
		}
		if section.Section != 11 {
			c.JSON(404, gin.H{
				"message": ErrNotFound,
			})
			return
		}
	case 12:
		section, err := e.db.GetSection12ByID(context.TODO(), claims.ID, input.SectionID)
		if err != nil {
			response := InterpretCosmosError(err)
			c.JSON(response.Code, gin.H{
				"message": response.Message,
			})
			return
		}
		if section.Section != 12 {
			c.JSON(404, gin.H{
				"message": ErrNotFound,
			})
			return
		}
	case 13:
		section, err := e.db.GetSection13ByID(context.TODO(), claims.ID, input.SectionID)
		if err != nil {
			response := InterpretCosmosError(err)
			c.JSON(response.Code, gin.H{
				"message": response.Message,
			})
			return
		}
		if section.Section != 13 {
			c.JSON(404, gin.H{
				"message": ErrNotFound,
			})
			return
		}
	case 14:
		section, err := e.db.GetSection14ByID(context.TODO(), claims.ID, input.SectionID)
		if err != nil {
			response := InterpretCosmosError(err)
			c.JSON(response.Code, gin.H{
				"message": response.Message,
			})
			return
		}
		if section.Section != 14 {
			c.JSON(404, gin.H{
				"message": ErrNotFound,
			})
			return
		}
	default:
		c.JSON(400, gin.H{
			"message": ErrBadRequest,
		})
		return
	}

	//verify eventSection doesn't already exist
	existingEventSection, err := e.db.GetEventSectionByIDs(context.TODO(), claims.ID, eventID, input.SectionID)
	if existingEventSection != (db.EventSection{}) {
		c.JSON(409, gin.H{
			"message": ErrEventSectionConflict,
		})
		return
	}
	if err != nil {
		response := InterpretCosmosError(err)
		if response.Code != 404 {
			c.JSON(response.Code, gin.H{
				"message": response.Message,
			})
			return
		}
	}

	g := guid.New()
	timestamp := utils.TimeNow()

	eventSection := db.EventSection{
		ID:            g.String(),
		UserID:        claims.ID,
		EventID:       event.ID,
		SectionNumber: *input.SectionNumber,
		SectionID:     input.SectionID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: timestamp.String(),
			Updated: timestamp.String(),
		},
	}

	var output UpsertEventSectionOutput

	output.EventSection, err = e.db.UpsertEventSection(context.TODO(), eventSection)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(201, output)

}

// DeleteEventSection godoc
// @Summary Removes an event section
// @Description Deletes a user's event section given the event ID
// @Tags Event
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param eventID path string true "Event ID"
// @Param sectionID path string true "Section ID"
// @Success 204
// @Failure 401
// @Failure 404
// @Router /event/{eventID}/{sectionID} [delete]
func (e *env) deleteEventSection(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	eventID := c.Param("eventID")
	sectionID := c.Param("sectionID")

	eventSection, err := e.db.GetEventSectionByIDs(context.TODO(), claims.ID, eventID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}
	if eventSection == (db.EventSection{}) {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	response, err := e.db.RemoveEventSection(context.TODO(), claims.ID, eventSection.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}
