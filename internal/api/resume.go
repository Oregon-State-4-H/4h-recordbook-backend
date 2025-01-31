package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"4h-recordbook-backend/pkg/db"
	"4h-recordbook-backend/internal/utils"
	"github.com/beevik/guid"
)

type UpsertSection1Input struct {
	Year string `json:"year" validate:"required"`
	Grade int `json:"grade" validate:"required"`
	ClubName string `json:"club_name" validate:"required"`
	NumInClub int `json:"num_in_club" validate:"required"`
	ClubLeader string `json:"club_leader" validate:"required"`
	MeetingsHeld int `json:"meetings_held" validate:"required"`
	MeetingsAttended int `json:"meetings_attended" validate:"required"`
}

type UpsertSection2Input struct {
	Year string `json:"year" validate:"required"`
	ProjectName string `json:"project_name" validate:"required"`
	ProjectScope string `json:"project_scope" validate:"required"`
}

type UpsertSection3Input struct {
	Year string `json:"year" validate:"required"`
	ActivityKind string `json:"activity_kind" validate:"required"`
	ThingsLearned string `json:"things_learned" validate:"required"`
	Level string `json:"level" validate:"required"`
}

type UpsertSection4Input struct {
	Year string `json:"year" validate:"required"`
	ActivityKind string `json:"activity_kind" validate:"required"`
	Scope string `json:"scope" validate:"required"`
	Level string `json:"level" validate:"required"`
}

type UpsertSection5Input struct {
	Year string `json:"year" validate:"required"`
	LeadershipRole string `json:"leadership_role" validate:"required"`
	HoursSpent int `json:"hours_spent" validate:"required"`
	NumPeopleReached int `json:"num_people_reached" validate:"required"`
}

type UpsertSection6Input struct {
	Year string `json:"year" validate:"required"`
	OrganizationName string `json:"organization_name" validate:"required"`
	LeadershipRole string `json:"leadership_role" validate:"required"`
	HoursSpent int `json:"hours_spent" validate:"required"`
	NumPeopleReached int `json:"num_people_reached" validate:"required"`
}

type UpsertSection7Input struct {
	Year string `json:"year" validate:"required"`
	ClubMemberActivities string `json:"club_member_activities" validate:"required"`
	HoursSpent int `json:"hours_spent" validate:"required"`
	NumPeopleReached int `json:"num_people_reached" validate:"required"`
}

type UpsertSection8Input struct {
	Year string `json:"year" validate:"required"`
	IndividualGroupActivities string `json:"individual_group_activities" validate:"required"`
	HoursSpent int `json:"hours_spent" validate:"required"`
	NumPeopleReached int `json:"num_people_reached" validate:"required"`
}

type UpsertSection9Input struct {
	Year string `json:"year" validate:"required"`
	CommunicationType string `json:"communication_type" validate:"required"`
	Topic string `json:"topic" validate:"required"`
	TimesGiven int `json:"times_given" validate:"required"`
	Location string `json:"location" validate:"required"`
	AudienceSize int `json:"audience_size" validate:"required"`
}

type UpsertSection10Input struct {
	Year string `json:"year" validate:"required"`
	CommunicationType string `json:"communication_type" validate:"required"`
	Topic string `json:"topic" validate:"required"`
	TimesGiven int `json:"times_given" validate:"required"`
	Location string `json:"location" validate:"required"`
	AudienceSize string `json:"audience_size" validate:"required"`
}

type UpsertSection11Input struct {
	Year string `json:"year" validate:"required"`
	EventAndLevel string `json:"event_and_level" validate:"required"`
	ExhibitsOrDivision string `json:"exhibits_or_division" validate:"required"`
	RibbonOrPlacings string `json:"ribbon_or_placings" validate:"required"`
}

type UpsertSection12Input struct {
	Year string `json:"year" validate:"required"`
	ContestOrEvent string `json:"contest_or_event" validate:"required"`
	RecognitionReceived string `json:"recognition_received" validate:"required"`
	Level string `json:"level" validate:"required"`
}

type UpsertSection13Input struct {
	Year string `json:"year" validate:"required"`
	RecognitionType string `json:"recognition_type" validate:"required"`
}

type UpsertSection14Input struct {
	Year string `json:"year" validate:"required"`
	RecognitionType string `json:"recognition_type" validate:"required"`
}

/*******************************
* GET FULL RESUME
********************************/

func (e *env) getResume(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	resume, err := e.db.GetResume(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, resume)

}

/*******************************
* SECTION 1
********************************/

func (e *env) getSection1(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	sections, err := e.db.GetSection1(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, sections)

}

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

func (e *env) getSection2(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	sections, err := e.db.GetSection2(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, sections)

}

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

func (e *env) getSection3(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	sections, err := e.db.GetSection3(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, sections)

}

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

func (e *env) getSection4(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	sections, err := e.db.GetSection4(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, sections)

}

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

func (e *env) getSection5(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	sections, err := e.db.GetSection5(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, sections)

}

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

func (e *env) getSection6(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	sections, err := e.db.GetSection6(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, sections)

}

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

func (e *env) getSection7(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	sections, err := e.db.GetSection7(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, sections)

}

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

func (e *env) getSection8(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	sections, err := e.db.GetSection8(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, sections)

}

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

func (e *env) getSection9(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	sections, err := e.db.GetSection9(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, sections)

}

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

func (e *env) getSection10(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	sections, err := e.db.GetSection10(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, sections)

}

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

func (e *env) getSection11(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	sections, err := e.db.GetSection11(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, sections)

}

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

func (e *env) getSection12(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	sections, err := e.db.GetSection12(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, sections)

}

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

func (e *env) getSection13(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	sections, err := e.db.GetSection13(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, sections)

}

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

func (e *env) getSection14(c *gin.Context) {

	cookie, err := c.Cookie("login_cookie")
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	sections, err := e.db.GetSection14(context.TODO(), cookie)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, sections)

}

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