package api

import (
	"4h-recordbook-backend/internal/utils"
	"4h-recordbook-backend/pkg/db"
	"context"

	"github.com/beevik/guid"
	"github.com/gin-gonic/gin"
)

/*******************************
* GET FULL RESUME
********************************/

type GetResumeOutput struct {
	Resume db.Resume `json:"resume"`
}

// GetResume godoc
// @Summary Gets full resume
// @Description Gets all of a user's entries for every resume section
// @Tags Resume
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} api.GetResumeOutput
// @Failure 401
// @Router /resume [get]
func (e *env) getResume(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var output GetResumeOutput

	output.Resume, err = e.db.GetResume(context.TODO(), claims.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

/*******************************
* SECTION 1
********************************/

type GetSection1sOutput struct {
	Sections []db.Section1 `json:"section_1_data"`
}

type GetSection1Output struct {
	Section db.Section1 `json:"section_1"`
}

type UpsertSection1Input struct {
	Nickname         string `json:"nickname" validate:"required"`
	Year             string `json:"year" validate:"required"`
	Grade            *int   `json:"grade" validate:"required"`
	ClubName         string `json:"club_name" validate:"required"`
	NumInClub        *int   `json:"num_in_club" validate:"required"`
	ClubLeader       string `json:"club_leader" validate:"required"`
	MeetingsHeld     *int   `json:"meetings_held" validate:"required"`
	MeetingsAttended *int   `json:"meetings_attended" validate:"required"`
}

type UpsertSection1Output GetSection1Output

// GetSection1s godoc
// @Summary Gets all Section 1 entries
// @Description Gets all of a user's Section 1 entries
// @Tags Resume Section 01
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} api.GetSection1sOutput
// @Failure 401
// @Failure 404
// @Router /section1 [get]
func (e *env) getSection1s(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var output GetSection1sOutput

	output.Sections, err = e.db.GetSection1sByUser(context.TODO(), claims.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// GetSection1 godoc
// @Summary Get a Section 1
// @Description Gets a user's Section 1 by ID
// @Tags Resume Section 01
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Success 200 {object} api.GetSection1Output
// @Failure 401
// @Router /section1/{sectionID} [get]
func (e *env) getSection1(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var output GetSection1Output

	output.Section, err = e.db.GetSection1ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	//todo: filter this in the db package
	if output.Section.Section != 1 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	c.JSON(200, output)

}

// AddSection1 godoc
// @Summary Add a Section 1 entry, return added Section 1 entry
// @Description Adds a Section 1 entry to a user's personal records
// @Tags Resume Section 01
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UpsertSection1Input body api.UpsertSection1Input true "Section 1 information"
// @Success 201 {object} api.UpsertSection1Output
// @Failure 400
// @Failure 401
// @Router /section1 [post]
func (e *env) addSection1(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertSection1Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section1{
		ID:               g.String(),
		Section:          1,
		Nickname:         input.Nickname,
		Year:             input.Year,
		Grade:            *input.Grade,
		ClubName:         input.ClubName,
		NumInClub:        *input.NumInClub,
		ClubLeader:       input.ClubLeader,
		MeetingsHeld:     *input.MeetingsHeld,
		MeetingsAttended: *input.MeetingsAttended,
		UserID:           claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: timestamp.String(),
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection1Output

	output.Section, err = e.db.UpsertSection1(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(201, output)
}

// UpdateSection1 godoc
// @Summary Updates a Section 1 entry
// @Description Updates a user's Section 1 entry information
// @Tags Resume Section 01
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Param UpsertSection1Input body api.UpsertSection1Input true "Section 1 information"
// @Success 200 {object} api.UpsertSection1Output
// @Failure 400
// @Failure 401
// @Router /section1/{sectionID} [put]
func (e *env) updateSection1(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var input UpsertSection1Input
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

	existingSection, err := e.db.GetSection1ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if existingSection.Section != 1 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section1{
		ID:               existingSection.ID,
		Section:          1,
		Nickname:         input.Nickname,
		Year:             input.Year,
		Grade:            *input.Grade,
		ClubName:         input.ClubName,
		NumInClub:        *input.NumInClub,
		ClubLeader:       input.ClubLeader,
		MeetingsHeld:     *input.MeetingsHeld,
		MeetingsAttended: *input.MeetingsAttended,
		UserID:           claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: existingSection.Created,
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection1Output

	output.Section, err = e.db.UpsertSection1(context.TODO(), updatedSection)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)
}

/*******************************
* SECTION 2
********************************/

type GetSection2sOutput struct {
	Sections []db.Section2 `json:"section_2_data"`
}

type GetSection2Output struct {
	Section db.Section2 `json:"section_2"`
}

type UpsertSection2Input struct {
	Year         string `json:"year" validate:"required"`
	ProjectName  string `json:"project_name" validate:"required"`
	ProjectScope string `json:"project_scope" validate:"required"`
}

type UpsertSection2Output GetSection2Output

// GetSection2s godoc
// @Summary Gets all Section 2 entries
// @Description Gets all of a user's Section 2 entries
// @Tags Resume Section 02
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} api.GetSection2sOutput
// @Failure 401
// @Router /section2 [get]
func (e *env) getSection2s(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var output GetSection2sOutput

	output.Sections, err = e.db.GetSection2sByUser(context.TODO(), claims.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// GetSection2 godoc
// @Summary Get a Section 2
// @Description Gets a user's Section 2 by ID
// @Tags Resume Section 02
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Success 200 {object} api.GetSection2Output
// @Failure 401
// @Failure 404
// @Router /section2/{sectionID} [get]
func (e *env) getSection2(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var output GetSection2Output

	output.Section, err = e.db.GetSection2ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if output.Section.Section != 2 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	c.JSON(200, output)

}

// AddSection2 godoc
// @Summary Add a Section 2 entry
// @Description Adds a Section 2 entry to a user's personal records
// @Tags Resume Section 02
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UpsertSection2Input body api.UpsertSection2Input true "Section 2 information"
// @Success 201 {object} api.UpsertSection2Output
// @Failure 400
// @Failure 401
// @Router /section2 [post]
func (e *env) addSection2(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertSection2Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section2{
		ID:           g.String(),
		Section:      2,
		Year:         input.Year,
		ProjectName:  input.ProjectName,
		ProjectScope: input.ProjectScope,
		UserID:       claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: timestamp.String(),
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection2Output

	output.Section, err = e.db.UpsertSection2(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(201, output)

}

// UpdateSection2 godoc
// @Summary Updates a Section 2 entry
// @Description Updates a user's Section 2 entry information
// @Tags Resume Section 02
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Param UpsertSection2Input body api.UpsertSection2Input true "Section 2 information"
// @Success 200 {object} api.UpsertSection2Output
// @Failure 400
// @Failure 401
// @Router /section2/{sectionID} [put]
func (e *env) updateSection2(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var input UpsertSection2Input
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

	existingSection, err := e.db.GetSection2ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if existingSection.Section != 2 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section2{
		ID:           existingSection.ID,
		Section:      2,
		Year:         input.Year,
		ProjectName:  input.ProjectName,
		ProjectScope: input.ProjectScope,
		UserID:       claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: existingSection.Created,
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection2Output

	output.Section, err = e.db.UpsertSection2(context.TODO(), updatedSection)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

/*******************************
* SECTION 3
********************************/

type GetSection3sOutput struct {
	Sections []db.Section3 `json:"section_3_data"`
}

type GetSection3Output struct {
	Section db.Section3 `json:"section_3"`
}

type UpsertSection3Input struct {
	Nickname      string `json:"nickname" validate:"required"`
	Year          string `json:"year" validate:"required"`
	ActivityKind  string `json:"activity_kind" validate:"required"`
	ThingsLearned string `json:"things_learned" validate:"required"`
	Level         string `json:"level" validate:"required"`
}

type UpsertSection3Output GetSection3Output

// GetSection3s godoc
// @Summary Gets all Section 3 entries
// @Description Gets all of a user's Section 3 entries
// @Tags Resume Section 03
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} api.GetSection3sOutput
// @Failure 401
// @Router /section3 [get]
func (e *env) getSection3s(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var output GetSection3sOutput

	output.Sections, err = e.db.GetSection3sByUser(context.TODO(), claims.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// GetSection3 godoc
// @Summary Get a Section 3
// @Description Gets a user's Section 3 by ID
// @Tags Resume Section 03
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Success 200 {object} api.GetSection3Output
// @Failure 401
// @Failure 404
// @Router /section3/{sectionID} [get]
func (e *env) getSection3(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var output GetSection3Output

	output.Section, err = e.db.GetSection3ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if output.Section.Section != 3 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	c.JSON(200, output)

}

// AddSection3 godoc
// @Summary Add a Section 3 entry
// @Description Adds a Section 3 entry to a user's personal records
// @Tags Resume Section 03
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UpsertSection3Input body api.UpsertSection3Input true "Section 3 information"
// @Success 201 {object} api.UpsertSection3Output
// @Failure 400
// @Failure 401
// @Router /section3 [post]
func (e *env) addSection3(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertSection3Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section3{
		ID:            g.String(),
		Section:       3,
		Nickname:      input.Nickname,
		Year:          input.Year,
		ActivityKind:  input.ActivityKind,
		ThingsLearned: input.ThingsLearned,
		Level:         input.Level,
		UserID:        claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: timestamp.String(),
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection3Output

	output.Section, err = e.db.UpsertSection3(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(201, output)

}

// UpdateSection3 godoc
// @Summary Updates a Section 3 entry
// @Description Updates a user's Section 3 entry information
// @Tags Resume Section 03
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Param UpsertSection3Input body api.UpsertSection3Input true "Section 3 information"
// @Success 200 {object} api.UpsertSection3Output
// @Failure 400
// @Failure 401
// @Router /section3/{sectionID} [put]
func (e *env) updateSection3(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var input UpsertSection3Input
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

	existingSection, err := e.db.GetSection3ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if existingSection.Section != 3 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section3{
		ID:            existingSection.ID,
		Section:       3,
		Nickname:      input.Nickname,
		Year:          input.Year,
		ActivityKind:  input.ActivityKind,
		ThingsLearned: input.ThingsLearned,
		Level:         input.Level,
		UserID:        claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: existingSection.Created,
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection3Output

	output.Section, err = e.db.UpsertSection3(context.TODO(), updatedSection)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

/*******************************
* SECTION 4
********************************/

type GetSection4sOutput struct {
	Sections []db.Section4 `json:"section_4_data"`
}

type GetSection4Output struct {
	Section db.Section4 `json:"section_4"`
}

type UpsertSection4Input struct {
	Nickname     string `json:"nickname" validate:"required"`
	Year         string `json:"year" validate:"required"`
	ActivityKind string `json:"activity_kind" validate:"required"`
	Scope        string `json:"scope" validate:"required"`
	Level        string `json:"level" validate:"required"`
}

type UpsertSection4Output GetSection4Output

// GetSection4s godoc
// @Summary Gets all Section 4 entries
// @Description Gets all of a user's Section 4 entries
// @Tags Resume Section 04
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} api.GetSection4sOutput
// @Failure 401
// @Router /section4 [get]
func (e *env) getSection4s(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var output GetSection4sOutput

	output.Sections, err = e.db.GetSection4sByUser(context.TODO(), claims.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// GetSection4 godoc
// @Summary Get a Section 4
// @Description Gets a user's Section 4 by ID
// @Tags Resume Section 04
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Success 200 {object} api.GetSection4Output
// @Failure 401
// @Failure 404
// @Router /section4/{sectionID} [get]
func (e *env) getSection4(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var output GetSection4Output

	output.Section, err = e.db.GetSection4ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if output.Section.Section != 4 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	c.JSON(200, output)

}

// AddSection4 godoc
// @Summary Add a Section 4 entry
// @Description Adds a Section 4 entry to a user's personal records
// @Tags Resume Section 04
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UpsertSection4Input body api.UpsertSection4Input true "Section 4 information"
// @Success 201 {object} api.UpsertSection4Output
// @Failure 400
// @Failure 401
// @Router /section4 [post]
func (e *env) addSection4(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertSection4Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section4{
		ID:           g.String(),
		Section:      4,
		Nickname:     input.Nickname,
		Year:         input.Year,
		ActivityKind: input.ActivityKind,
		Scope:        input.Scope,
		Level:        input.Level,
		UserID:       claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: timestamp.String(),
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection4Output

	output.Section, err = e.db.UpsertSection4(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(201, output)

}

// UpdateSection4 godoc
// @Summary Updates a Section 4 entry
// @Description Updates a user's Section 4 entry information
// @Tags Resume Section 04
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Param UpsertSection4Input body api.UpsertSection4Input true "Section 4 information"
// @Success 200 {object} api.UpsertSection4Output
// @Failure 400
// @Failure 401
// @Router /section4/{sectionID} [put]
func (e *env) updateSection4(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var input UpsertSection4Input
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

	existingSection, err := e.db.GetSection4ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if existingSection.Section != 4 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section4{
		ID:           existingSection.ID,
		Section:      4,
		Nickname:     input.Nickname,
		Year:         input.Year,
		ActivityKind: input.ActivityKind,
		Scope:        input.Scope,
		Level:        input.Level,
		UserID:       claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: existingSection.Created,
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection4Output

	output.Section, err = e.db.UpsertSection4(context.TODO(), updatedSection)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

/*******************************
* SECTION 5
********************************/

type GetSection5sOutput struct {
	Sections []db.Section5 `json:"section_5_data"`
}

type GetSection5Output struct {
	Section db.Section5 `json:"section_5"`
}

type UpsertSection5Input struct {
	Nickname         string `json:"nickname" validate:"required"`
	Year             string `json:"year" validate:"required"`
	LeadershipRole   string `json:"leadership_role" validate:"required"`
	HoursSpent       *int   `json:"hours_spent" validate:"required"`
	NumPeopleReached *int   `json:"num_people_reached" validate:"required"`
}

type UpsertSection5Output GetSection5Output

// GetSection5s godoc
// @Summary Gets all Section 5 entries
// @Description Gets all of a user's Section 5 entries
// @Tags Resume Section 05
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} api.GetSection5sOutput
// @Failure 401
// @Router /section5 [get]
func (e *env) getSection5s(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var output GetSection5sOutput

	output.Sections, err = e.db.GetSection5sByUser(context.TODO(), claims.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// GetSection5 godoc
// @Summary Get a Section 5
// @Description Gets a user's Section 5 by ID
// @Tags Resume Section 05
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Success 200 {object} api.GetSection5Output
// @Failure 401
// @Failure 404
// @Router /section5/{sectionID} [get]
func (e *env) getSection5(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var output GetSection5Output

	output.Section, err = e.db.GetSection5ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if output.Section.Section != 5 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	c.JSON(200, output)

}

// AddSection5 godoc
// @Summary Add a Section 5 entry
// @Description Adds a Section 5 entry to a user's personal records
// @Tags Resume Section 05
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UpsertSection5Input body api.UpsertSection5Input true "Section 5 information"
// @Success 201 {object} api.UpsertSection5Output
// @Failure 400
// @Failure 401
// @Router /section5 [post]
func (e *env) addSection5(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertSection5Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section5{
		ID:               g.String(),
		Section:          5,
		Nickname:         input.Nickname,
		Year:             input.Year,
		LeadershipRole:   input.LeadershipRole,
		HoursSpent:       *input.HoursSpent,
		NumPeopleReached: *input.NumPeopleReached,
		UserID:           claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: timestamp.String(),
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection5Output

	output.Section, err = e.db.UpsertSection5(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(201, output)

}

// UpdateSection5 godoc
// @Summary Updates a Section 5 entry
// @Description Updates a user's Section 5 entry information
// @Tags Resume Section 05
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Param UpsertSection5Input body api.UpsertSection5Input true "Section 5 information"
// @Success 200 {object} api.UpsertSection5Output
// @Failure 400
// @Failure 401
// @Router /section5/{sectionID} [put]
func (e *env) updateSection5(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var input UpsertSection5Input
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

	existingSection, err := e.db.GetSection5ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if existingSection.Section != 5 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section5{
		ID:               existingSection.ID,
		Section:          5,
		Nickname:         input.Nickname,
		Year:             input.Year,
		LeadershipRole:   input.LeadershipRole,
		HoursSpent:       *input.HoursSpent,
		NumPeopleReached: *input.NumPeopleReached,
		UserID:           claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: existingSection.Created,
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection5Output

	output.Section, err = e.db.UpsertSection5(context.TODO(), updatedSection)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

/*******************************
* SECTION 6
********************************/

type GetSection6sOutput struct {
	Sections []db.Section6 `json:"section_6_data"`
}

type GetSection6Output struct {
	Section db.Section6 `json:"section_6"`
}

type UpsertSection6Input struct {
	Nickname         string `json:"nickname" validate:"required"`
	Year             string `json:"year" validate:"required"`
	OrganizationName string `json:"organization_name" validate:"required"`
	LeadershipRole   string `json:"leadership_role" validate:"required"`
	HoursSpent       *int   `json:"hours_spent" validate:"required"`
	NumPeopleReached *int   `json:"num_people_reached" validate:"required"`
}

type UpsertSection6Output GetSection6Output

// GetSection6s godoc
// @Summary Gets all Section 6 entries
// @Description Gets all of a user's Section 6 entries
// @Tags Resume Section 06
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} api.GetSection6sOutput
// @Failure 401
// @Router /section6 [get]
func (e *env) getSection6s(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var output GetSection6sOutput

	output.Sections, err = e.db.GetSection6sByUser(context.TODO(), claims.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// GetSection6 godoc
// @Summary Get a Section 6
// @Description Gets a user's Section 6 by ID
// @Tags Resume Section 06
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Success 200 {object} api.GetSection6Output
// @Failure 401
// @Failure 404
// @Router /section6/{sectionID} [get]
func (e *env) getSection6(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var output GetSection6Output

	output.Section, err = e.db.GetSection6ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if output.Section.Section != 6 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	c.JSON(200, output)

}

// AddSection6 godoc
// @Summary Add a Section 6 entry
// @Description Adds a Section 6 entry to a user's personal records
// @Tags Resume Section 06
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UpsertSection6Input body api.UpsertSection6Input true "Section 6 information"
// @Success 201 {object} api.UpsertSection6Output
// @Failure 400
// @Failure 401
// @Router /section6 [post]
func (e *env) addSection6(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertSection6Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section6{
		ID:               g.String(),
		Section:          6,
		Nickname:         input.Nickname,
		Year:             input.Year,
		OrganizationName: input.OrganizationName,
		LeadershipRole:   input.LeadershipRole,
		HoursSpent:       *input.HoursSpent,
		NumPeopleReached: *input.NumPeopleReached,
		UserID:           claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: timestamp.String(),
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection6Output

	output.Section, err = e.db.UpsertSection6(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(201, output)

}

// UpdateSection6 godoc
// @Summary Updates a Section 6 entry
// @Description Updates a user's Section 6 entry information
// @Tags Resume Section 06
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Param UpsertSection6Input body api.UpsertSection6Input true "Section 6 information"
// @Success 200 {object} api.UpsertSection6Output
// @Failure 400
// @Failure 401
// @Router /section6/{sectionID} [put]
func (e *env) updateSection6(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var input UpsertSection6Input
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

	existingSection, err := e.db.GetSection6ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if existingSection.Section != 6 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section6{
		ID:               existingSection.ID,
		Section:          6,
		Nickname:         input.Nickname,
		Year:             input.Year,
		OrganizationName: input.OrganizationName,
		LeadershipRole:   input.LeadershipRole,
		HoursSpent:       *input.HoursSpent,
		NumPeopleReached: *input.NumPeopleReached,
		UserID:           claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: existingSection.Created,
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection6Output

	output.Section, err = e.db.UpsertSection6(context.TODO(), updatedSection)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

/*******************************
* SECTION 7
********************************/

type GetSection7sOutput struct {
	Sections []db.Section7 `json:"section_7_data"`
}

type GetSection7Output struct {
	Section db.Section7 `json:"section_7"`
}

type UpsertSection7Input struct {
	Nickname             string `json:"nickname" validate:"required"`
	Year                 string `json:"year" validate:"required"`
	ClubMemberActivities string `json:"club_member_activities" validate:"required"`
	HoursSpent           *int   `json:"hours_spent" validate:"required"`
	NumPeopleReached     *int   `json:"num_people_reached" validate:"required"`
}

type UpsertSection7Output GetSection7Output

// GetSection7s godoc
// @Summary Gets all Section 7 entries
// @Description Gets all of a user's Section 7 entries
// @Tags Resume Section 07
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} api.GetSection7sOutput
// @Failure 401
// @Router /section7 [get]
func (e *env) getSection7s(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var output GetSection7sOutput

	output.Sections, err = e.db.GetSection7sByUser(context.TODO(), claims.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// GetSection7 godoc
// @Summary Get a Section 7
// @Description Gets a user's Section 7 by ID
// @Tags Resume Section 07
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Success 200 {object} api.GetSection7Output
// @Failure 401
// @Failure 404
// @Router /section7/{sectionID} [get]
func (e *env) getSection7(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var output GetSection7Output

	output.Section, err = e.db.GetSection7ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if output.Section.Section != 7 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	c.JSON(200, output)

}

// AddSection7 godoc
// @Summary Add a Section 7 entry
// @Description Adds a Section 7 entry to a user's personal records
// @Tags Resume Section 07
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UpsertSection7Input body api.UpsertSection7Input true "Section 7 information"
// @Success 201 {object} api.UpsertSection7Output
// @Failure 400
// @Failure 401
// @Router /section7 [post]
func (e *env) addSection7(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertSection7Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section7{
		ID:                   g.String(),
		Section:              7,
		Nickname:             input.Nickname,
		Year:                 input.Year,
		ClubMemberActivities: input.ClubMemberActivities,
		HoursSpent:           *input.HoursSpent,
		NumPeopleReached:     *input.NumPeopleReached,
		UserID:               claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: timestamp.String(),
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection7Output

	output.Section, err = e.db.UpsertSection7(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(201, output)

}

// UpdateSection7 godoc
// @Summary Updates a Section 7 entry
// @Description Updates a user's Section 7 entry information
// @Tags Resume Section 07
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Param UpsertSection7Input body api.UpsertSection7Input true "Section 7 information"
// @Success 200 {object} api.UpsertSection7Output
// @Failure 400
// @Failure 401
// @Router /section7/{sectionID} [put]
func (e *env) updateSection7(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var input UpsertSection7Input
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

	existingSection, err := e.db.GetSection7ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if existingSection.Section != 7 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section7{
		ID:                   existingSection.ID,
		Section:              7,
		Nickname:             input.Nickname,
		Year:                 input.Year,
		ClubMemberActivities: input.ClubMemberActivities,
		HoursSpent:           *input.HoursSpent,
		NumPeopleReached:     *input.NumPeopleReached,
		UserID:               claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: existingSection.Created,
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection7Output

	output.Section, err = e.db.UpsertSection7(context.TODO(), updatedSection)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

/*******************************
* SECTION 8
********************************/

type GetSection8sOutput struct {
	Sections []db.Section8 `json:"section_8_data"`
}

type GetSection8Output struct {
	Section db.Section8 `json:"section_8"`
}

type UpsertSection8Input struct {
	Nickname                  string `json:"nickname" validate:"required"`
	Year                      string `json:"year" validate:"required"`
	IndividualGroupActivities string `json:"individual_group_activities" validate:"required"`
	HoursSpent                *int   `json:"hours_spent" validate:"required"`
	NumPeopleReached          *int   `json:"num_people_reached" validate:"required"`
}

type UpsertSection8Output GetSection8Output

// GetSection8s godoc
// @Summary Gets all Section 8 entries
// @Description Gets all of a user's Section 8 entries
// @Tags Resume Section 08
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} api.GetSection8sOutput
// @Failure 401
// @Router /section8 [get]
func (e *env) getSection8s(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var output GetSection8sOutput

	output.Sections, err = e.db.GetSection8sByUser(context.TODO(), claims.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// GetSection8 godoc
// @Summary Get a Section 8
// @Description Gets a user's Section 8 by ID
// @Tags Resume Section 08
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Success 200 {object} api.GetSection8Output
// @Failure 401
// @Failure 404
// @Router /section8/{sectionID} [get]
func (e *env) getSection8(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var output GetSection8Output

	output.Section, err = e.db.GetSection8ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if output.Section.Section != 8 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	c.JSON(200, output)

}

// AddSection8 godoc
// @Summary Add a Section 8 entry
// @Description Adds a Section 8 entry to a user's personal records
// @Tags Resume Section 08
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UpsertSection8Input body api.UpsertSection8Input true "Section 8 information"
// @Success 201 {object} api.UpsertSection8Output
// @Failure 400
// @Failure 401
// @Router /section8 [post]
func (e *env) addSection8(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertSection8Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section8{
		ID:                        g.String(),
		Section:                   8,
		Nickname:                  input.Nickname,
		Year:                      input.Year,
		IndividualGroupActivities: input.IndividualGroupActivities,
		HoursSpent:                *input.HoursSpent,
		NumPeopleReached:          *input.NumPeopleReached,
		UserID:                    claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: timestamp.String(),
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection8Output

	output.Section, err = e.db.UpsertSection8(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(201, output)

}

// UpdateSection8 godoc
// @Summary Updates a Section 8 entry
// @Description Updates a user's Section 8 entry information
// @Tags Resume Section 08
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Param UpsertSection8Input body api.UpsertSection8Input true "Section 8 information"
// @Success 200 {object} api.UpsertSection8Output
// @Failure 400
// @Failure 401
// @Router /section8/{sectionID} [put]
func (e *env) updateSection8(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var input UpsertSection8Input
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

	existingSection, err := e.db.GetSection8ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if existingSection.Section != 8 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section8{
		ID:                        existingSection.ID,
		Section:                   8,
		Nickname:                  input.Nickname,
		Year:                      input.Year,
		IndividualGroupActivities: input.IndividualGroupActivities,
		HoursSpent:                *input.HoursSpent,
		NumPeopleReached:          *input.NumPeopleReached,
		UserID:                    claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: existingSection.Created,
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection8Output

	output.Section, err = e.db.UpsertSection8(context.TODO(), updatedSection)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

/*******************************
* SECTION 9
********************************/

type GetSection9sOutput struct {
	Sections []db.Section9 `json:"section_9_data"`
}

type GetSection9Output struct {
	Section db.Section9 `json:"section_9"`
}

type UpsertSection9Input struct {
	Nickname          string `json:"nickname" validate:"required"`
	Year              string `json:"year" validate:"required"`
	CommunicationType string `json:"communication_type" validate:"required"`
	Topic             string `json:"topic" validate:"required"`
	TimesGiven        *int   `json:"times_given" validate:"required"`
	Location          string `json:"location" validate:"required"`
	AudienceSize      *int   `json:"audience_size" validate:"required"`
}

type UpsertSection9Output GetSection9Output

// GetSection9s godoc
// @Summary Gets all Section 9 entries
// @Description Gets all of a user's Section 9 entries
// @Tags Resume Section 09
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} api.GetSection9sOutput
// @Failure 401
// @Router /section9 [get]
func (e *env) getSection9s(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var output GetSection9sOutput

	output.Sections, err = e.db.GetSection9sByUser(context.TODO(), claims.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// GetSection9 godoc
// @Summary Get a Section 9
// @Description Gets a user's Section 9 by ID
// @Tags Resume Section 09
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Success 200 {object} api.GetSection9Output
// @Failure 401
// @Failure 404
// @Router /section9/{sectionID} [get]
func (e *env) getSection9(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var output GetSection9Output

	output.Section, err = e.db.GetSection9ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if output.Section.Section != 9 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	c.JSON(200, output)

}

// AddSection9 godoc
// @Summary Add a Section 9 entry
// @Description Adds a Section 9 entry to a user's personal records
// @Tags Resume Section 09
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UpsertSection9Input body api.UpsertSection9Input true "Section 9 information"
// @Success 201 {object} api.UpsertSection9Output
// @Failure 400
// @Failure 401
// @Router /section9 [post]
func (e *env) addSection9(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertSection9Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section9{
		ID:                g.String(),
		Section:           9,
		Nickname:          input.Nickname,
		Year:              input.Year,
		CommunicationType: input.CommunicationType,
		Topic:             input.Topic,
		TimesGiven:        *input.TimesGiven,
		Location:          input.Location,
		AudienceSize:      *input.AudienceSize,
		UserID:            claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: timestamp.String(),
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection9Output

	output.Section, err = e.db.UpsertSection9(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(201, output)

}

// UpdateSection9 godoc
// @Summary Updates a Section 9 entry
// @Description Updates a user's Section 9 entry information
// @Tags Resume Section 09
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Param UpsertSection9Input body api.UpsertSection9Input true "Section 9 information"
// @Success 200 {object} api.UpsertSection9Output
// @Failure 400
// @Failure 401
// @Router /section9/{sectionID} [put]
func (e *env) updateSection9(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var input UpsertSection9Input
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

	existingSection, err := e.db.GetSection9ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if existingSection.Section != 9 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section9{
		ID:                existingSection.ID,
		Section:           9,
		Nickname:          input.Nickname,
		Year:              input.Year,
		CommunicationType: input.CommunicationType,
		Topic:             input.Topic,
		TimesGiven:        *input.TimesGiven,
		Location:          input.Location,
		AudienceSize:      *input.AudienceSize,
		UserID:            claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: existingSection.Created,
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection9Output

	output.Section, err = e.db.UpsertSection9(context.TODO(), updatedSection)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

/*******************************
* SECTION 10
********************************/

type GetSection10sOutput struct {
	Sections []db.Section10 `json:"section_10_data"`
}

type GetSection10Output struct {
	Section db.Section10 `json:"section_10"`
}

type UpsertSection10Input struct {
	Nickname          string `json:"nickname" validate:"required"`
	Year              string `json:"year" validate:"required"`
	CommunicationType string `json:"communication_type" validate:"required"`
	Topic             string `json:"topic" validate:"required"`
	TimesGiven        *int   `json:"times_given" validate:"required"`
	Location          string `json:"location" validate:"required"`
	AudienceSize      *int   `json:"audience_size" validate:"required"`
}

type UpsertSection10Output GetSection10Output

// GetSection10s godoc
// @Summary Gets all Section 10 entries
// @Description Gets all of a user's Section 10 entries
// @Tags Resume Section 10
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} api.GetSection10sOutput
// @Failure 401
// @Router /section10 [get]
func (e *env) getSection10s(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var output GetSection10sOutput

	output.Sections, err = e.db.GetSection10sByUser(context.TODO(), claims.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// GetSection10 godoc
// @Summary Get a Section 10
// @Description Gets a user's Section 10 by ID
// @Tags Resume Section 10
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Success 200 {object} api.GetSection10Output
// @Failure 401
// @Failure 404
// @Router /section10/{sectionID} [get]
func (e *env) getSection10(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var output GetSection10Output

	output.Section, err = e.db.GetSection10ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if output.Section.Section != 10 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	c.JSON(200, output)

}

// AddSection10 godoc
// @Summary Add a Section 10 entry
// @Description Adds a Section 10 entry to a user's personal records
// @Tags Resume Section 10
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UpsertSection10Input body api.UpsertSection10Input true "Section 10 information"
// @Success 201 {object} api.UpsertSection10Output
// @Failure 400
// @Failure 401
// @Router /section10 [post]
func (e *env) addSection10(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertSection10Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section10{
		ID:                g.String(),
		Section:           10,
		Nickname:          input.Nickname,
		Year:              input.Year,
		CommunicationType: input.CommunicationType,
		Topic:             input.Topic,
		TimesGiven:        *input.TimesGiven,
		Location:          input.Location,
		AudienceSize:      *input.AudienceSize,
		UserID:            claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: timestamp.String(),
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection10Output

	output.Section, err = e.db.UpsertSection10(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(201, output)

}

// UpdateSection10 godoc
// @Summary Updates a Section 10 entry
// @Description Updates a user's Section 10 entry information
// @Tags Resume Section 10
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Param UpsertSection10Input body api.UpsertSection10Input true "Section 10 information"
// @Success 200 {object} api.UpsertSection10Output
// @Failure 400
// @Failure 401
// @Router /section10/{sectionID} [put]
func (e *env) updateSection10(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var input UpsertSection10Input
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

	existingSection, err := e.db.GetSection10ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if existingSection.Section != 10 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section10{
		ID:                existingSection.ID,
		Section:           10,
		Nickname:          input.Nickname,
		Year:              input.Year,
		CommunicationType: input.CommunicationType,
		Topic:             input.Topic,
		TimesGiven:        *input.TimesGiven,
		Location:          input.Location,
		AudienceSize:      *input.AudienceSize,
		UserID:            claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: existingSection.Created,
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection10Output

	output.Section, err = e.db.UpsertSection10(context.TODO(), updatedSection)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

/*******************************
* SECTION 11
********************************/

type GetSection11sOutput struct {
	Sections []db.Section11 `json:"section_11_data"`
}

type GetSection11Output struct {
	Section db.Section11 `json:"section_11"`
}

type UpsertSection11Input struct {
	Nickname           string `json:"nickname" validate:"required"`
	Year               string `json:"year" validate:"required"`
	EventAndLevel      string `json:"event_and_level" validate:"required"`
	ExhibitsOrDivision string `json:"exhibits_or_division" validate:"required"`
	RibbonOrPlacings   string `json:"ribbon_or_placings" validate:"required"`
}

type UpsertSection11Output GetSection11Output

// GetSection11s godoc
// @Summary Gets all Section 11 entries
// @Description Gets all of a user's Section 11 entries
// @Tags Resume Section 11
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} api.GetSection11sOutput
// @Failure 401
// @Router /section11 [get]
func (e *env) getSection11s(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var output GetSection11sOutput

	output.Sections, err = e.db.GetSection11sByUser(context.TODO(), claims.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// GetSection11 godoc
// @Summary Get a Section 11
// @Description Gets a user's Section 11 by ID
// @Tags Resume Section 11
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Success 200 {object} api.GetSection11Output
// @Failure 401
// @Failure 404
// @Router /section11/{sectionID} [get]
func (e *env) getSection11(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var output GetSection11Output

	output.Section, err = e.db.GetSection11ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if output.Section.Section != 11 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	c.JSON(200, output)

}

// AddSection11 godoc
// @Summary Add a Section 11 entry
// @Description Adds a Section 11 entry to a user's personal records
// @Tags Resume Section 11
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UpsertSection11Input body api.UpsertSection11Input true "Section 11 information"
// @Success 201 {object} api.UpsertSection11Output
// @Failure 400
// @Failure 401
// @Router /section11 [post]
func (e *env) addSection11(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertSection11Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section11{
		ID:                 g.String(),
		Section:            11,
		Nickname:           input.Nickname,
		Year:               input.Year,
		EventAndLevel:      input.EventAndLevel,
		ExhibitsOrDivision: input.ExhibitsOrDivision,
		RibbonOrPlacings:   input.RibbonOrPlacings,
		UserID:             claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: timestamp.String(),
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection11Output

	output.Section, err = e.db.UpsertSection11(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(201, output)

}

// UpdateSection11 godoc
// @Summary Updates a Section 11 entry
// @Description Updates a user's Section 11 entry information
// @Tags Resume Section 11
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Param UpsertSection11Input body api.UpsertSection11Input true "Section 11 information"
// @Success 200 {object} api.UpsertSection11Output
// @Failure 400
// @Failure 401
// @Router /section11/{sectionID} [put]
func (e *env) updateSection11(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var input UpsertSection11Input
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

	existingSection, err := e.db.GetSection11ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if existingSection.Section != 11 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section11{
		ID:                 existingSection.ID,
		Section:            11,
		Nickname:           input.Nickname,
		Year:               input.Year,
		EventAndLevel:      input.EventAndLevel,
		ExhibitsOrDivision: input.ExhibitsOrDivision,
		RibbonOrPlacings:   input.RibbonOrPlacings,
		UserID:             claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: existingSection.Created,
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection11Output

	output.Section, err = e.db.UpsertSection11(context.TODO(), updatedSection)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

/*******************************
* SECTION 12
********************************/

type GetSection12sOutput struct {
	Sections []db.Section12 `json:"section_12_data"`
}

type GetSection12Output struct {
	Section db.Section12 `json:"section_12"`
}

type UpsertSection12Input struct {
	Nickname            string `json:"nickname" validate:"required"`
	Year                string `json:"year" validate:"required"`
	ContestOrEvent      string `json:"contest_or_event" validate:"required"`
	RecognitionReceived string `json:"recognition_received" validate:"required"`
	Level               string `json:"level" validate:"required"`
}

type UpsertSection12Output GetSection12Output

// GetSection12s godoc
// @Summary Gets all Section 12 entries
// @Description Gets all of a user's Section 12 entries
// @Tags Resume Section 12
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} api.GetSection12sOutput
// @Failure 401
// @Router /section12 [get]
func (e *env) getSection12s(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var output GetSection12sOutput

	output.Sections, err = e.db.GetSection12sByUser(context.TODO(), claims.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// GetSection12 godoc
// @Summary Get a Section 12
// @Description Gets a user's Section 12 by ID
// @Tags Resume Section 12
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Success 200 {object} api.GetSection12Output
// @Failure 401
// @Failure 404
// @Router /section12/{sectionID} [get]
func (e *env) getSection12(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var output GetSection12Output

	output.Section, err = e.db.GetSection12ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if output.Section.Section != 12 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	c.JSON(200, output)

}

// AddSection12 godoc
// @Summary Add a Section 12 entry
// @Description Adds a Section 12 entry to a user's personal records
// @Tags Resume Section 12
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UpsertSection12Input body api.UpsertSection12Input true "Section 12 information"
// @Success 201 {object} api.UpsertSection12Output
// @Failure 400
// @Failure 401
// @Router /section12 [post]
func (e *env) addSection12(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertSection12Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section12{
		ID:                  g.String(),
		Section:             12,
		Nickname:            input.Nickname,
		Year:                input.Year,
		ContestOrEvent:      input.ContestOrEvent,
		RecognitionReceived: input.RecognitionReceived,
		Level:               input.Level,
		UserID:              claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: timestamp.String(),
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection12Output

	output.Section, err = e.db.UpsertSection12(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(201, output)

}

// UpdateSection12 godoc
// @Summary Updates a Section 12 entry
// @Description Updates a user's Section 12 entry information
// @Tags Resume Section 12
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Param UpsertSection12Input body api.UpsertSection12Input true "Section 12 information"
// @Success 200 {object} api.UpsertSection12Output
// @Failure 400
// @Failure 401
// @Router /section12/{sectionID} [put]
func (e *env) updateSection12(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var input UpsertSection12Input
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

	existingSection, err := e.db.GetSection12ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if existingSection.Section != 12 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section12{
		ID:                  existingSection.ID,
		Section:             12,
		Nickname:            input.Nickname,
		Year:                input.Year,
		ContestOrEvent:      input.ContestOrEvent,
		RecognitionReceived: input.RecognitionReceived,
		Level:               input.Level,
		UserID:              claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: existingSection.Created,
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection12Output

	output.Section, err = e.db.UpsertSection12(context.TODO(), updatedSection)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

/*******************************
* SECTION 13
********************************/

type GetSection13sOutput struct {
	Sections []db.Section13 `json:"section_13_data"`
}

type GetSection13Output struct {
	Section db.Section13 `json:"section_13"`
}

type UpsertSection13Input struct {
	Nickname        string `json:"nickname" validate:"required"`
	Year            string `json:"year" validate:"required"`
	RecognitionType string `json:"recognition_type" validate:"required"`
}

type UpsertSection13Output GetSection13Output

// GetSection13s godoc
// @Summary Gets all Section 13 entries
// @Description Gets all of a user's Section 13 entries
// @Tags Resume Section 13
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} api.GetSection13sOutput
// @Failure 401
// @Router /section13 [get]
func (e *env) getSection13s(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var output GetSection13sOutput

	output.Sections, err = e.db.GetSection13sByUser(context.TODO(), claims.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// GetSection13 godoc
// @Summary Get a Section 13
// @Description Gets a user's Section 13 by ID
// @Tags Resume Section 13
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Success 200 {object} api.GetSection13Output
// @Failure 401
// @Failure 404
// @Router /section13/{sectionID} [get]
func (e *env) getSection13(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var output GetSection13Output

	output.Section, err = e.db.GetSection13ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if output.Section.Section != 13 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	c.JSON(200, output)

}

// AddSection13 godoc
// @Summary Add a Section 13 entry
// @Description Adds a Section 13 entry to a user's personal records
// @Tags Resume Section 13
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UpsertSection13Input body api.UpsertSection13Input true "Section 13 information"
// @Success 201 {object} api.UpsertSection13Output
// @Failure 400
// @Failure 401
// @Router /section13 [post]
func (e *env) addSection13(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertSection13Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section13{
		ID:              g.String(),
		Section:         13,
		Nickname:        input.Nickname,
		Year:            input.Year,
		RecognitionType: input.RecognitionType,
		UserID:          claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: timestamp.String(),
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection13Output

	output.Section, err = e.db.UpsertSection13(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(201, output)

}

// UpdateSection13 godoc
// @Summary Updates a Section 13 entry
// @Description Updates a user's Section 13 entry information
// @Tags Resume Section 13
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Param UpsertSection13Input body api.UpsertSection13Input true "Section 13 information"
// @Success 200 {object} api.UpsertSection13Output
// @Failure 400
// @Failure 401
// @Router /section13/{sectionID} [put]
func (e *env) updateSection13(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var input UpsertSection13Input
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

	existingSection, err := e.db.GetSection13ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if existingSection.Section != 13 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section13{
		ID:              existingSection.ID,
		Section:         13,
		Nickname:        input.Nickname,
		Year:            input.Year,
		RecognitionType: input.RecognitionType,
		UserID:          claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: existingSection.Created,
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection13Output

	output.Section, err = e.db.UpsertSection13(context.TODO(), updatedSection)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)
}

/*******************************
* SECTION 14
********************************/

type GetSection14sOutput struct {
	Sections []db.Section14 `json:"section_14_data"`
}

type GetSection14Output struct {
	Section db.Section14 `json:"section_14"`
}

type UpsertSection14Input struct {
	Nickname        string `json:"nickname" validate:"required"`
	Year            string `json:"year" validate:"required"`
	RecognitionType string `json:"recognition_type" validate:"required"`
}

type UpsertSection14Output GetSection14Output

// GetSection14s godoc
// @Summary Gets all Section 14 entries
// @Description Gets all of a user's Section 14 entries
// @Tags Resume Section 14
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} api.GetSection14sOutput
// @Failure 401
// @Router /section14 [get]
func (e *env) getSection14s(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var output GetSection14sOutput

	output.Sections, err = e.db.GetSection14sByUser(context.TODO(), claims.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// GetSection14 godoc
// @Summary Get a Section 14
// @Description Gets a user's Section 14 by ID
// @Tags Resume Section 14
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Success 200 {object} api.GetSection14Output
// @Failure 401
// @Failure 404
// @Router /section14/{sectionID} [get]
func (e *env) getSection14(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var output GetSection14Output

	output.Section, err = e.db.GetSection14ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if output.Section.Section != 14 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	c.JSON(200, output)

}

// AddSection14 godoc
// @Summary Add a Section 14 entry
// @Description Adds a Section 14 entry to a user's personal records
// @Tags Resume Section 14
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UpsertSection14Input body api.UpsertSection14Input true "Section 14 information"
// @Success 201 {object} api.UpsertSection14Output
// @Failure 400
// @Failure 401
// @Router /section14 [post]
func (e *env) addSection14(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertSection14Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section14{
		ID:              g.String(),
		Section:         14,
		Nickname:        input.Nickname,
		Year:            input.Year,
		RecognitionType: input.RecognitionType,
		UserID:          claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: timestamp.String(),
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection14Output

	output.Section, err = e.db.UpsertSection14(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(201, output)

}

// UpdateSection14 godoc
// @Summary Updates a Section 14 entry
// @Description Updates a user's Section 14 entry information
// @Tags Resume Section 14
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Param UpsertSection14Input body api.UpsertSection14Input true "Section 14 information"
// @Success 200 {object} api.UpsertSection14Output
// @Failure 400
// @Failure 401
// @Router /section14/{sectionID} [put]
func (e *env) updateSection14(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	var input UpsertSection14Input
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

	existingSection, err := e.db.GetSection14ByID(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if existingSection.Section != 14 {
		c.JSON(404, gin.H{
			"message": ErrNotFound,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section14{
		ID:              existingSection.ID,
		Section:         14,
		Nickname:        input.Nickname,
		Year:            input.Year,
		RecognitionType: input.RecognitionType,
		UserID:          claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: existingSection.Created,
			Updated: timestamp.String(),
		},
	}

	var output UpsertSection14Output

	output.Section, err = e.db.UpsertSection14(context.TODO(), updatedSection)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

/*******************************
* DELETING
********************************/

// DeleteSection godoc
// @Summary Removes a resume section
// @Description Deletes a user's resume section given the section ID. Can be any resume section
// @Tags Resume
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sectionID path string true "Section ID"
// @Success 204
// @Failure 401
// @Failure 404
// @Router /section/{sectionID} [delete]
func (e *env) deleteSection(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	sectionID := c.Param("sectionID")

	response, err := e.db.RemoveSection(context.TODO(), claims.ID, sectionID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}
