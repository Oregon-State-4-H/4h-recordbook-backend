package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"4h-recordbook-backend/pkg/db"
	"4h-recordbook-backend/internal/utils"
	"github.com/beevik/guid"
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
// @Success 200 {object} api.GetResumeOutput
// @Failure 401
// @Router /resume [get]
func (e *env) getResume(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var output GetResumeOutput

	output.Resume, err = e.db.GetResume(context.TODO(), cookie)
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

type UpsertSection1Input struct {
	Year string `json:"year" validate:"required"`
	Grade int `json:"grade" validate:"required"`
	ClubName string `json:"club_name" validate:"required"`
	NumInClub int `json:"num_in_club" validate:"required"`
	ClubLeader string `json:"club_leader" validate:"required"`
	MeetingsHeld int `json:"meetings_held" validate:"required"`
	MeetingsAttended int `json:"meetings_attended" validate:"required"`
}

type GetSection1Output struct {
	Sections []db.Section1 `json:"section_1_data"`
}

// GetSection1 godoc
// @Summary Gets all Section 1 entries
// @Description Gets all of a user's Section 1 entries
// @Tags Resume Section 1
// @Accept json 
// @Produce json
// @Success 200 {object} api.GetSection1Output
// @Failure 401
// @Router /section1 [get]
func (e *env) getSection1(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var output GetSection1Output

	output.Sections, err = e.db.GetSection1(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// AddSection1 godoc
// @Summary Add a Section 1 entry
// @Description Adds a Section 1 entry to a user's personal records
// @Tags Resume Section 1
// @Accept json 
// @Produce json
// @Param UpsertSection1Input body api.UpsertSection1Input true "Section 1 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section1 [post]
func (e *env) addSection1(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertSection1Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section1 {
		ID: g.String(),
		Section: 1,
		Year: input.Year,
		Grade: input.Grade,
		ClubName: input.ClubName,
		NumInClub: input.NumInClub,
		ClubLeader: input.ClubLeader,
		MeetingsHeld: input.MeetingsHeld,
		MeetingsAttended: input.MeetingsAttended,
		UserID: cookie,
		Created: timestamp.ToString(),
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection1(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

// UpdateSection1 godoc
// @Summary Updates a Section 1 entry
// @Description Updates a user's Section 1 entry information
// @Tags Resume Section 1
// @Accept json 
// @Produce json
// @Param UpsertSection1Input body api.UpsertSection1Input true "Section 1 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section1 [put]
func (e *env) updateSection1(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	id := c.Param("sectionId")

	var input UpsertSection1Input
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

	existingSection, err := e.db.GetSection1ByID(context.TODO(), cookie, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section1 {
		ID: existingSection.ID,
		Section: 1,
		Year: input.Year,
		Grade: input.Grade,
		ClubName: input.ClubName,
		NumInClub: input.NumInClub,
		ClubLeader: input.ClubLeader,
		MeetingsHeld: input.MeetingsHeld,
		MeetingsAttended: input.MeetingsAttended,
		UserID: cookie,
		Created: existingSection.Created,
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection1(context.TODO(), updatedSection)
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
* SECTION 2
********************************/

type UpsertSection2Input struct {
	Year string `json:"year" validate:"required"`
	ProjectName string `json:"project_name" validate:"required"`
	ProjectScope string `json:"project_scope" validate:"required"`
}

type GetSection2Output struct {
	Sections []db.Section2 `json:"section_2_data"`
}

// GetSection2 godoc
// @Summary Gets all Section 2 entries
// @Description Gets all of a user's Section 2 entries
// @Tags Resume Section 2
// @Accept json 
// @Produce json
// @Success 200 {object} api.GetSection2Output
// @Failure 401
// @Router /section2 [get]
func (e *env) getSection2(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var output GetSection2Output

	output.Sections, err = e.db.GetSection2(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// AddSection2 godoc
// @Summary Add a Section 2 entry
// @Description Adds a Section 2 entry to a user's personal records
// @Tags Resume Section 2
// @Accept json 
// @Produce json
// @Param UpsertSection2Input body api.UpsertSection2Input true "Section 2 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section2 [post]
func (e *env) addSection2(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertSection2Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section2 {
		ID: g.String(),
		Section: 2,
		Year: input.Year,
		ProjectName: input.ProjectName,
		ProjectScope: input.ProjectScope,
		UserID: cookie,
		Created: timestamp.ToString(),
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection2(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

// UpdateSection2 godoc
// @Summary Updates a Section 2 entry
// @Description Updates a user's Section 2 entry information
// @Tags Resume Section 2
// @Accept json 
// @Produce json
// @Param UpsertSection2Input body api.UpsertSection2Input true "Section 2 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section2 [put]
func (e *env) updateSection2(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	id := c.Param("sectionId")

	var input UpsertSection2Input
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

	existingSection, err := e.db.GetSection2ByID(context.TODO(), cookie, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section2 {
		ID: existingSection.ID,
		Section: 2,
		Year: input.Year,
		ProjectName: input.ProjectName,
		ProjectScope: input.ProjectScope,
		UserID: cookie,
		Created: existingSection.Created,
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection2(context.TODO(), updatedSection)
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
* SECTION 3
********************************/

type UpsertSection3Input struct {
	Year string `json:"year" validate:"required"`
	ActivityKind string `json:"activity_kind" validate:"required"`
	ThingsLearned string `json:"things_learned" validate:"required"`
	Level string `json:"level" validate:"required"`
}

type GetSection3Output struct {
	Sections []db.Section3 `json:"section_3_data"`
}

// GetSection3 godoc
// @Summary Gets all Section 3 entries
// @Description Gets all of a user's Section 3 entries
// @Tags Resume Section 3
// @Accept json 
// @Produce json
// @Success 200 {object} api.GetSection3Output
// @Failure 401
// @Router /section3 [get]
func (e *env) getSection3(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var output GetSection3Output

	output.Sections, err = e.db.GetSection3(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// AddSection3 godoc
// @Summary Add a Section 3 entry
// @Description Adds a Section 3 entry to a user's personal records
// @Tags Resume Section 3
// @Accept json 
// @Produce json
// @Param UpsertSection3Input body api.UpsertSection3Input true "Section 3 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section3 [post]
func (e *env) addSection3(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertSection3Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section3 {
		ID: g.String(),
		Section: 3,
		Year: input.Year,
		ActivityKind: input.ActivityKind,
		ThingsLearned: input.ThingsLearned,
		Level: input.Level,
		UserID: cookie,
		Created: timestamp.ToString(),
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection3(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

// UpdateSection3 godoc
// @Summary Updates a Section 3 entry
// @Description Updates a user's Section 3 entry information
// @Tags Resume Section 3
// @Accept json 
// @Produce json
// @Param UpsertSection3Input body api.UpsertSection3Input true "Section 3 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section3 [put]
func (e *env) updateSection3(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	id := c.Param("sectionId")

	var input UpsertSection3Input
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

	existingSection, err := e.db.GetSection3ByID(context.TODO(), cookie, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section3 {
		ID: existingSection.ID,
		Section: 3,
		Year: input.Year,
		ActivityKind: input.ActivityKind,
		ThingsLearned: input.ThingsLearned,
		Level: input.Level,
		UserID: cookie,
		Created: existingSection.Created,
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection3(context.TODO(), updatedSection)
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
* SECTION 4
********************************/

type UpsertSection4Input struct {
	Year string `json:"year" validate:"required"`
	ActivityKind string `json:"activity_kind" validate:"required"`
	Scope string `json:"scope" validate:"required"`
	Level string `json:"level" validate:"required"`
}

type GetSection4Output struct {
	Sections []db.Section4 `json:"section_4_data"`
}

// GetSection4 godoc
// @Summary Gets all Section 4 entries
// @Description Gets all of a user's Section 4 entries
// @Tags Resume Section 4
// @Accept json 
// @Produce json
// @Success 200 {object} api.GetSection4Output
// @Failure 401
// @Router /section4 [get]
func (e *env) getSection4(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var output GetSection4Output

	output.Sections, err = e.db.GetSection4(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// AddSection4 godoc
// @Summary Add a Section 4 entry
// @Description Adds a Section 4 entry to a user's personal records
// @Tags Resume Section 4
// @Accept json 
// @Produce json
// @Param UpsertSection4Input body api.UpsertSection4Input true "Section 4 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section4 [post]
func (e *env) addSection4(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertSection4Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section4 {
		ID: g.String(),
		Section: 4,
		Year: input.Year,
		ActivityKind: input.ActivityKind,
		Scope: input.Scope,
		Level: input.Level,
		UserID: cookie,
		Created: timestamp.ToString(),
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection4(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

// UpdateSection4 godoc
// @Summary Updates a Section 4 entry
// @Description Updates a user's Section 4 entry information
// @Tags Resume Section 4
// @Accept json 
// @Produce json
// @Param UpsertSection4Input body api.UpsertSection4Input true "Section 4 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section4 [put]
func (e *env) updateSection4(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	id := c.Param("sectionId")

	var input UpsertSection4Input
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

	existingSection, err := e.db.GetSection4ByID(context.TODO(), cookie, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section4 {
		ID: existingSection.ID,
		Section: 4,
		Year: input.Year,
		ActivityKind: input.ActivityKind,
		Scope: input.Scope,
		Level: input.Level,
		UserID: cookie,
		Created: existingSection.Created,
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection4(context.TODO(), updatedSection)
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
* SECTION 5
********************************/

type UpsertSection5Input struct {
	Year string `json:"year" validate:"required"`
	LeadershipRole string `json:"leadership_role" validate:"required"`
	HoursSpent int `json:"hours_spent" validate:"required"`
	NumPeopleReached int `json:"num_people_reached" validate:"required"`
}

type GetSection5Output struct {
	Sections []db.Section5 `json:"section_5_data"`
}

// GetSection5 godoc
// @Summary Gets all Section 5 entries
// @Description Gets all of a user's Section 5 entries
// @Tags Resume Section 5
// @Accept json 
// @Produce json
// @Success 200 {object} api.GetSection5Output
// @Failure 401
// @Router /section5 [get]
func (e *env) getSection5(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var output GetSection5Output

	output.Sections, err = e.db.GetSection5(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// AddSection5 godoc
// @Summary Add a Section 5 entry
// @Description Adds a Section 5 entry to a user's personal records
// @Tags Resume Section 5
// @Accept json 
// @Produce json
// @Param UpsertSection5Input body api.UpsertSection5Input true "Section 5 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section5 [post]
func (e *env) addSection5(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertSection5Input
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

	g := guid.New()	
	timestamp := utils.TimeNow()

	section := db.Section5 {
		ID: g.String(),
		Section: 5,
		Year: input.Year,
		LeadershipRole: input.LeadershipRole,
		HoursSpent: input.HoursSpent,
		NumPeopleReached: input.NumPeopleReached,
		UserID: cookie,
		Created: timestamp.ToString(),
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection5(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

// UpdateSection5 godoc
// @Summary Updates a Section 5 entry
// @Description Updates a user's Section 5 entry information
// @Tags Resume Section 5
// @Accept json 
// @Produce json
// @Param UpsertSection5Input body api.UpsertSection5Input true "Section 5 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section5 [put]
func (e *env) updateSection5(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	id := c.Param("sectionId")

	var input UpsertSection5Input
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

	existingSection, err := e.db.GetSection5ByID(context.TODO(), cookie, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section5 {
		ID: existingSection.ID,
		Section: 5,
		Year: input.Year,
		LeadershipRole: input.LeadershipRole,
		HoursSpent: input.HoursSpent,
		NumPeopleReached: input.NumPeopleReached,
		UserID: cookie,
		Created: existingSection.Created,
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection5(context.TODO(), updatedSection)
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
* SECTION 6
********************************/

type UpsertSection6Input struct {
	Year string `json:"year" validate:"required"`
	OrganizationName string `json:"organization_name" validate:"required"`
	LeadershipRole string `json:"leadership_role" validate:"required"`
	HoursSpent int `json:"hours_spent" validate:"required"`
	NumPeopleReached int `json:"num_people_reached" validate:"required"`
}

type GetSection6Output struct {
	Sections []db.Section6 `json:"section_6_data"`
}

// GetSection6 godoc
// @Summary Gets all Section 6 entries
// @Description Gets all of a user's Section 6 entries
// @Tags Resume Section 6
// @Accept json 
// @Produce json
// @Success 200 {object} api.GetSection6Output
// @Failure 401
// @Router /section6 [get]
func (e *env) getSection6(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var output GetSection6Output

	output.Sections, err = e.db.GetSection6(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// AddSection6 godoc
// @Summary Add a Section 6 entry
// @Description Adds a Section 6 entry to a user's personal records
// @Tags Resume Section 6
// @Accept json 
// @Produce json
// @Param UpsertSection6Input body api.UpsertSection6Input true "Section 6 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section6 [post]
func (e *env) addSection6(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertSection6Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section6 {
		ID: g.String(),
		Section: 6,
		Year: input.Year,
		OrganizationName: input.OrganizationName,
		LeadershipRole: input.LeadershipRole,
		HoursSpent: input.HoursSpent,
		NumPeopleReached: input.NumPeopleReached,
		UserID: cookie,
		Created: timestamp.ToString(),
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection6(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

// UpdateSection6 godoc
// @Summary Updates a Section 6 entry
// @Description Updates a user's Section 6 entry information
// @Tags Resume Section 6
// @Accept json 
// @Produce json
// @Param UpsertSection6Input body api.UpsertSection6Input true "Section 6 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section6 [put]
func (e *env) updateSection6(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	id := c.Param("sectionId")

	var input UpsertSection6Input
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

	existingSection, err := e.db.GetSection6ByID(context.TODO(), cookie, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section6 {
		ID: existingSection.ID,
		Section: 6,
		Year: input.Year,
		OrganizationName: input.OrganizationName,
		LeadershipRole: input.LeadershipRole,
		HoursSpent: input.HoursSpent,
		NumPeopleReached: input.NumPeopleReached,
		UserID: cookie,
		Created: existingSection.Created,
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection6(context.TODO(), updatedSection)
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
* SECTION 7
********************************/

type UpsertSection7Input struct {
	Year string `json:"year" validate:"required"`
	ClubMemberActivities string `json:"club_member_activities" validate:"required"`
	HoursSpent int `json:"hours_spent" validate:"required"`
	NumPeopleReached int `json:"num_people_reached" validate:"required"`
}

type GetSection7Output struct {
	Sections []db.Section7 `json:"section_7_data"`
}

// GetSection7 godoc
// @Summary Gets all Section 7 entries
// @Description Gets all of a user's Section 7 entries
// @Tags Resume Section 7
// @Accept json 
// @Produce json
// @Success 200 {object} api.GetSection7Output
// @Failure 401
// @Router /section7 [get]
func (e *env) getSection7(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var output GetSection7Output

	output.Sections, err = e.db.GetSection7(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// AddSection7 godoc
// @Summary Add a Section 7 entry
// @Description Adds a Section 7 entry to a user's personal records
// @Tags Resume Section 7
// @Accept json 
// @Produce json
// @Param UpsertSection7Input body api.UpsertSection7Input true "Section 7 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section7 [post]
func (e *env) addSection7(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertSection7Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section7 {
		ID: g.String(),
		Section: 7,
		Year: input.Year,
		ClubMemberActivities: input.ClubMemberActivities,
		HoursSpent: input.HoursSpent,
		NumPeopleReached: input.NumPeopleReached,
		UserID: cookie,
		Created: timestamp.ToString(),
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection7(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

// UpdateSection7 godoc
// @Summary Updates a Section 7 entry
// @Description Updates a user's Section 7 entry information
// @Tags Resume Section 7
// @Accept json 
// @Produce json
// @Param UpsertSection7Input body api.UpsertSection7Input true "Section 7 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section7 [put]
func (e *env) updateSection7(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	id := c.Param("sectionId")

	var input UpsertSection7Input
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

	existingSection, err := e.db.GetSection7ByID(context.TODO(), cookie, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section7 {
		ID: existingSection.ID,
		Section: 7,
		Year: input.Year,
		ClubMemberActivities: input.ClubMemberActivities,
		HoursSpent: input.HoursSpent,
		NumPeopleReached: input.NumPeopleReached,
		UserID: cookie,
		Created: existingSection.Created,
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection7(context.TODO(), updatedSection)
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
* SECTION 8
********************************/

type UpsertSection8Input struct {
	Year string `json:"year" validate:"required"`
	IndividualGroupActivities string `json:"individual_group_activities" validate:"required"`
	HoursSpent int `json:"hours_spent" validate:"required"`
	NumPeopleReached int `json:"num_people_reached" validate:"required"`
}

type GetSection8Output struct {
	Sections []db.Section8 `json:"section_8_data"`
}

// GetSection8 godoc
// @Summary Gets all Section 8 entries
// @Description Gets all of a user's Section 8 entries
// @Tags Resume Section 8
// @Accept json 
// @Produce json
// @Success 200 {object} api.GetSection8Output
// @Failure 401
// @Router /section8 [get]
func (e *env) getSection8(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var output GetSection8Output

	output.Sections, err = e.db.GetSection8(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// AddSection8 godoc
// @Summary Add a Section 8 entry
// @Description Adds a Section 8 entry to a user's personal records
// @Tags Resume Section 8
// @Accept json 
// @Produce json
// @Param UpsertSection8Input body api.UpsertSection8Input true "Section 8 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section8 [post]
func (e *env) addSection8(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertSection8Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section8 {
		ID: g.String(),
		Section: 8,
		Year: input.Year,
		IndividualGroupActivities: input.IndividualGroupActivities,
		HoursSpent: input.HoursSpent,
		NumPeopleReached: input.NumPeopleReached,
		UserID: cookie,
		Created: timestamp.ToString(),
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection8(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

// UpdateSection8 godoc
// @Summary Updates a Section 8 entry
// @Description Updates a user's Section 8 entry information
// @Tags Resume Section 8
// @Accept json 
// @Produce json
// @Param UpsertSection8Input body api.UpsertSection8Input true "Section 8 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section8 [put]
func (e *env) updateSection8(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	id := c.Param("sectionId")

	var input UpsertSection8Input
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

	existingSection, err := e.db.GetSection8ByID(context.TODO(), cookie, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section8 {
		ID: existingSection.ID,
		Section: 8,
		Year: input.Year,
		IndividualGroupActivities: input.IndividualGroupActivities,
		HoursSpent: input.HoursSpent,
		NumPeopleReached: input.NumPeopleReached,
		UserID: cookie,
		Created: existingSection.Created,
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection8(context.TODO(), updatedSection)
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
* SECTION 9
********************************/

type UpsertSection9Input struct {
	Year string `json:"year" validate:"required"`
	CommunicationType string `json:"communication_type" validate:"required"`
	Topic string `json:"topic" validate:"required"`
	TimesGiven int `json:"times_given" validate:"required"`
	Location string `json:"location" validate:"required"`
	AudienceSize int `json:"audience_size" validate:"required"`
}

type GetSection9Output struct {
	Sections []db.Section9 `json:"section_9_data"`
}

// GetSection9 godoc
// @Summary Gets all Section 9 entries
// @Description Gets all of a user's Section 9 entries
// @Tags Resume Section 9
// @Accept json 
// @Produce json
// @Success 200 {object} api.GetSection9Output
// @Failure 401
// @Router /section9 [get]
func (e *env) getSection9(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var output GetSection9Output

	output.Sections, err = e.db.GetSection9(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// AddSection9 godoc
// @Summary Add a Section 9 entry
// @Description Adds a Section 9 entry to a user's personal records
// @Tags Resume Section 9
// @Accept json 
// @Produce json
// @Param UpsertSection9Input body api.UpsertSection9Input true "Section 9 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section9 [post]
func (e *env) addSection9(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertSection9Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section9 {
		ID: g.String(),
		Section: 9,
		Year: input.Year,
		CommunicationType: input.CommunicationType,
		Topic: input.Topic,
		TimesGiven: input.TimesGiven,
		Location: input.Location,
		AudienceSize: input.AudienceSize,
		UserID: cookie,
		Created: timestamp.ToString(),
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection9(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

// UpdateSection9 godoc
// @Summary Updates a Section 9 entry
// @Description Updates a user's Section 9 entry information
// @Tags Resume Section 9
// @Accept json 
// @Produce json
// @Param UpsertSection9Input body api.UpsertSection9Input true "Section 9 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section9 [put]
func (e *env) updateSection9(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	id := c.Param("sectionId")

	var input UpsertSection9Input
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

	existingSection, err := e.db.GetSection9ByID(context.TODO(), cookie, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section9 {
		ID: existingSection.ID,
		Section: 9,
		Year: input.Year,
		CommunicationType: input.CommunicationType,
		Topic: input.Topic,
		TimesGiven: input.TimesGiven,
		Location: input.Location,
		AudienceSize: input.AudienceSize,
		UserID: cookie,
		Created: existingSection.Created,
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection9(context.TODO(), updatedSection)
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
* SECTION 10
********************************/

type UpsertSection10Input struct {
	Year string `json:"year" validate:"required"`
	CommunicationType string `json:"communication_type" validate:"required"`
	Topic string `json:"topic" validate:"required"`
	TimesGiven int `json:"times_given" validate:"required"`
	Location string `json:"location" validate:"required"`
	AudienceSize string `json:"audience_size" validate:"required"`
}

type GetSection10Output struct {
	Sections []db.Section10 `json:"section_10_data"`
}

// GetSection10 godoc
// @Summary Gets all Section 10 entries
// @Description Gets all of a user's Section 10 entries
// @Tags Resume Section 10
// @Accept json 
// @Produce json
// @Success 200 {object} api.GetSection10Output
// @Failure 401
// @Router /section10 [get]
func (e *env) getSection10(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var output GetSection10Output

	output.Sections, err = e.db.GetSection10(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
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
// @Param UpsertSection10Input body api.UpsertSection10Input true "Section 10 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section10 [post]
func (e *env) addSection10(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertSection10Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section10 {
		ID: g.String(),
		Section: 10,
		Year: input.Year,
		CommunicationType: input.CommunicationType,
		Topic: input.Topic,
		TimesGiven: input.TimesGiven,
		Location: input.Location,
		AudienceSize: input.AudienceSize,
		UserID: cookie,
		Created: timestamp.ToString(),
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection10(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

// UpdateSection10 godoc
// @Summary Updates a Section 10 entry
// @Description Updates a user's Section 10 entry information
// @Tags Resume Section 10
// @Accept json 
// @Produce json
// @Param UpsertSection10Input body api.UpsertSection10Input true "Section 10 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section10 [put]
func (e *env) updateSection10(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	id := c.Param("sectionId")

	var input UpsertSection10Input
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

	existingSection, err := e.db.GetSection10ByID(context.TODO(), cookie, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section10 {
		ID: existingSection.ID,
		Section: 10,
		Year: input.Year,
		CommunicationType: input.CommunicationType,
		Topic: input.Topic,
		TimesGiven: input.TimesGiven,
		Location: input.Location,
		AudienceSize: input.AudienceSize,
		UserID: cookie,
		Created: existingSection.Created,
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection10(context.TODO(), updatedSection)
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
* SECTION 11
********************************/

type UpsertSection11Input struct {
	Year string `json:"year" validate:"required"`
	EventAndLevel string `json:"event_and_level" validate:"required"`
	ExhibitsOrDivision string `json:"exhibits_or_division" validate:"required"`
	RibbonOrPlacings string `json:"ribbon_or_placings" validate:"required"`
}

type GetSection11Output struct {
	Sections []db.Section11 `json:"section_11_data"`
}

// GetSection11 godoc
// @Summary Gets all Section 11 entries
// @Description Gets all of a user's Section 11 entries
// @Tags Resume Section 11
// @Accept json 
// @Produce json
// @Success 200 {object} api.GetSection11Output
// @Failure 401
// @Router /section11 [get]
func (e *env) getSection11(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var output GetSection11Output

	output.Sections, err = e.db.GetSection11(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
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
// @Param UpsertSection11Input body api.UpsertSection11Input true "Section 11 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section11 [post]
func (e *env) addSection11(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertSection11Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section11 {
		ID: g.String(),
		Section: 11,
		Year: input.Year,
		EventAndLevel: input.EventAndLevel,
		ExhibitsOrDivision: input.ExhibitsOrDivision,
		RibbonOrPlacings: input.RibbonOrPlacings,
		UserID: cookie,
		Created: timestamp.ToString(),
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection11(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

// UpdateSection11 godoc
// @Summary Updates a Section 11 entry
// @Description Updates a user's Section 11 entry information
// @Tags Resume Section 11
// @Accept json 
// @Produce json
// @Param UpsertSection11Input body api.UpsertSection11Input true "Section 11 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section11 [put]
func (e *env) updateSection11(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	id := c.Param("sectionId")

	var input UpsertSection11Input
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

	existingSection, err := e.db.GetSection11ByID(context.TODO(), cookie, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section11 {
		ID: existingSection.ID,
		Section: 11,
		Year: input.Year,
		EventAndLevel: input.EventAndLevel,
		ExhibitsOrDivision: input.ExhibitsOrDivision,
		RibbonOrPlacings: input.RibbonOrPlacings,
		UserID: cookie,
		Created: existingSection.Created,
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection11(context.TODO(), updatedSection)
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
* SECTION 12
********************************/

type UpsertSection12Input struct {
	Year string `json:"year" validate:"required"`
	ContestOrEvent string `json:"contest_or_event" validate:"required"`
	RecognitionReceived string `json:"recognition_received" validate:"required"`
	Level string `json:"level" validate:"required"`
}

type GetSection12Output struct {
	Sections []db.Section12 `json:"section_1_data"`
}

// GetSection12 godoc
// @Summary Gets all Section 12 entries
// @Description Gets all of a user's Section 12 entries
// @Tags Resume Section 12
// @Accept json 
// @Produce json
// @Success 200 {object} api.GetSection12Output
// @Failure 401
// @Router /section12 [get]
func (e *env) getSection12(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var output GetSection12Output

	output.Sections, err = e.db.GetSection12(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
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
// @Param UpsertSection12Input body api.UpsertSection12Input true "Section 12 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section12 [post]
func (e *env) addSection12(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertSection12Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section12 {
		ID: g.String(),
		Section: 12,
		Year: input.Year,
		ContestOrEvent: input.ContestOrEvent,
		RecognitionReceived: input.RecognitionReceived,
		Level: input.Level,
		UserID: cookie,
		Created: timestamp.ToString(),
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection12(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

// UpdateSection12 godoc
// @Summary Updates a Section 12 entry
// @Description Updates a user's Section 12 entry information
// @Tags Resume Section 12
// @Accept json 
// @Produce json
// @Param UpsertSection12Input body api.UpsertSection12Input true "Section 12 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section12 [put]
func (e *env) updateSection12(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	id := c.Param("sectionId")

	var input UpsertSection12Input
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

	existingSection, err := e.db.GetSection12ByID(context.TODO(), cookie, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section12 {
		ID: existingSection.ID,
		Section: 12,
		Year: input.Year,
		ContestOrEvent: input.ContestOrEvent,
		RecognitionReceived: input.RecognitionReceived,
		Level: input.Level,
		UserID: cookie,
		Created: existingSection.Created,
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection12(context.TODO(), updatedSection)
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
* SECTION 13
********************************/

type UpsertSection13Input struct {
	Year string `json:"year" validate:"required"`
	RecognitionType string `json:"recognition_type" validate:"required"`
}

type GetSection13Output struct {
	Sections []db.Section13 `json:"section_13_data"`
}

// GetSection13 godoc
// @Summary Gets all Section 13 entries
// @Description Gets all of a user's Section 13 entries
// @Tags Resume Section 13
// @Accept json 
// @Produce json
// @Success 200 {object} api.GetSection13Output
// @Failure 401
// @Router /section13 [get]
func (e *env) getSection13(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var output GetSection13Output

	output.Sections, err = e.db.GetSection13(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
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
// @Param UpsertSection13Input body api.UpsertSection13Input true "Section 13 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section13 [post]
func (e *env) addSection13(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertSection13Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section13 {
		ID: g.String(),
		Section: 13,
		Year: input.Year,
		RecognitionType: input.RecognitionType,
		UserID: cookie,
		Created: timestamp.ToString(),
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection13(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

// UpdateSection13 godoc
// @Summary Updates a Section 13 entry
// @Description Updates a user's Section 13 entry information
// @Tags Resume Section 13
// @Accept json 
// @Produce json
// @Param UpsertSection13Input body api.UpsertSection13Input true "Section 13 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section13 [put]
func (e *env) updateSection13(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	id := c.Param("sectionId")

	var input UpsertSection13Input
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

	existingSection, err := e.db.GetSection13ByID(context.TODO(), cookie, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section13 {
		ID: existingSection.ID,
		Section: 13,
		Year: input.Year,
		RecognitionType: input.RecognitionType,
		UserID: cookie,
		Created: existingSection.Created,
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection13(context.TODO(), updatedSection)
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
* SECTION 14
********************************/

type UpsertSection14Input struct {
	Year string `json:"year" validate:"required"`
	RecognitionType string `json:"recognition_type" validate:"required"`
}

type GetSection14Output struct {
	Sections []db.Section14 `json:"section_14_data"`
}

// GetSection14 godoc
// @Summary Gets all Section 14 entries
// @Description Gets all of a user's Section 14 entries
// @Tags Resume Section 14
// @Accept json 
// @Produce json
// @Success 200 {object} api.GetSection14Output
// @Failure 401
// @Router /section14 [get]
func (e *env) getSection14(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var output GetSection14Output

	output.Sections, err = e.db.GetSection14(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
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
// @Param UpsertSection14Input body api.UpsertSection14Input true "Section 14 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section14 [post]
func (e *env) addSection14(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertSection14Input
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

	g := guid.New()
	timestamp := utils.TimeNow()

	section := db.Section14 {
		ID: g.String(),
		Section: 14,
		Year: input.Year,
		RecognitionType: input.RecognitionType,
		UserID: cookie,
		Created: timestamp.ToString(),
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection14(context.TODO(), section)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

// UpdateSection14 godoc
// @Summary Updates a Section 14 entry
// @Description Updates a user's Section 14 entry information
// @Tags Resume Section 14
// @Accept json 
// @Produce json
// @Param UpsertSection14Input body api.UpsertSection14Input true "Section 14 information"
// @Success 204
// @Failure 400
// @Failure 401
// @Router /section14 [put]
func (e *env) updateSection14(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	id := c.Param("sectionId")

	var input UpsertSection14Input
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

	existingSection, err := e.db.GetSection14ByID(context.TODO(), cookie, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section14 {
		ID: existingSection.ID,
		Section: 14,
		Year: input.Year,
		RecognitionType: input.RecognitionType,
		UserID: cookie,
		Created: existingSection.Created,
		Updated: timestamp.ToString(),
	}

	response, err := e.db.UpsertSection14(context.TODO(), updatedSection)
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
* DELETING
********************************/

// DeleteSection godoc
// @Summary Removes a resume section
// @Description Deletes a user's resume section given the section ID. Can be any resume section
// @Tags Resume
// @Accept json
// @Produce json
// @Param sectionId path string true "Section ID"
// @Success 204
// @Failure 401
// @Failure 404 
// @Router /section/{sectionId} [delete]
func (e *env) deleteSection(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	id := c.Param("sectionId")

	response, err := e.db.RemoveSection(context.TODO(), cookie, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}