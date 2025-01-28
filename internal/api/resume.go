package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"4h-recordbook-backend/pkg/db"
	"4h-recordbook-backend/internal/utils"
)

type Section1Req struct {
	ID string `json:"id"`
	Year string `json:"year"`
	Grade int `json:"grade,omitempty" validate:"required"`
	ClubName string `json:"club_name"`
	NumInClub int `json:"num_in_club"`
	ClubLeader string `json:"club_leader"`
	MeetingsHeld int `json:"meetings_held"`
	MeetingsAttended int `json:"meetings_attended"`
}

type Section2Req struct {
	ID string `json:"id"`
	Year string `json:"year"`
	ProjectName string `json:"project_name"`
	ProjectScope string `json:"project_scope"`
}

type Section3Req struct {
	ID string `json:"id"`
	Year string `json:"year"`
	ActivityKind string `json:"activity_kind"`
	ThingsLearned string `json:"things_learned"`
	Level string `json:"level"`
}

type Section4Req struct {
	ID string `json:"id"`
	Year string `json:"year"`
	ActivityKind string `json:"activity_kind"`
	Scope string `json:"scope"`
	Level string `json:"level"`
}

type Section5Req struct {
	ID string `json:"id"`
	Year string `json:"year"`
	LeadershipRole string `json:"leadership_role"`
	HoursSpent int `json:"hours_spent"`
	NumPeopleReached int `json:"num_people_reached"`
}

type Section6Req struct {
	ID string `json:"id"`
	Year string `json:"year"`
	OrganizationName string `json:"organization_name"`
	LeadershipRole string `json:"leadership_role"`
	HoursSpent int `json:"hours_spent"`
	NumPeopleReached int `json:"num_people_reached"`
}

type Section7Req struct {
	ID string `json:"id"`
	Year string `json:"year"`
	ClubMemberActivities string `json:"club_member_activities"`
	HoursSpent int `json:"hours_spent"`
	NumPeopleReached int `json:"num_people_reached"`
}

type Section8Req struct {
	ID string `json:"id"`
	Year string `json:"year"`
	IndividualGroupActivities string `json:"individual_group_activities"`
	HoursSpent int `json:"hours_spent"`
	NumPeopleReached int `json:"num_people_reached"`
}

type Section9Req struct {
	ID string `json:"id"`
	Year string `json:"year"`
	CommunicationType string `json:"communication_type"`
	Topic string `json:"topic"`
	TimesGiven int `json:"times_given"`
	Location string `json:"location"`
	AudienceSize int `json:"audience_size"`
}

type Section10Req struct {
	ID string `json:"id"`
	Year string `json:"year"`
	CommunicationType string `json:"communication_type"`
	Topic string `json:"topic"`
	TimesGiven int `json:"times_given"`
	Location string `json:"location"`
	AudienceSize string `json:"audience_size"`
}

type Section11Req struct {
	ID string `json:"id"`
	Year string `json:"year"`
	EventAndLevel string `json:"event_and_level"`
	ExhibitsOrDivision string `json:"exhibits_or_division"`
	RibbonOrPlacings string `json:"ribbon_or_placings"`
}

type Section12Req struct {
	ID string `json:"id"`
	Year string `json:"year"`
	ContestOrEvent string `json:"contest_or_event"`
	RecognitionReceived string `json:"recognition_received"`
	Level string `json:"level"`
}

type Section13Req struct {
	ID string `json:"id"`
	Year string `json:"year"`
	RecognitionType string `json:"recognition_type"`
}

type Section14Req struct {
	ID string `json:"id"`
	Year string `json:"year"`
	RecognitionType string `json:"recognition_type"`
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

	var req Section1Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	timestamp := utils.TimeNow()

	section := db.Section1 {
		ID: req.ID, //temporary
		Section: 1,
		Year: req.Year,
		Grade: req.Grade,
		ClubName: req.ClubName,
		NumInClub: req.NumInClub,
		ClubLeader: req.ClubLeader,
		MeetingsHeld: req.MeetingsHeld,
		MeetingsAttended: req.MeetingsAttended,
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

	var req Section1Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	existingSection, err := e.db.GetSection1ByID(context.TODO(), cookie, req.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section1 {
		ID: existingSection.ID, //temporary
		Section: 1,
		Year: req.Year,
		Grade: req.Grade,
		ClubName: req.ClubName,
		NumInClub: req.NumInClub,
		ClubLeader: req.ClubLeader,
		MeetingsHeld: req.MeetingsHeld,
		MeetingsAttended: req.MeetingsAttended,
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

	var req Section2Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	timestamp := utils.TimeNow()

	section := db.Section2 {
		ID: req.ID, //temporary
		Section: 2,
		Year: req.Year,
		ProjectName: req.ProjectName,
		ProjectScope: req.ProjectScope,
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

	var req Section2Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	existingSection, err := e.db.GetSection2ByID(context.TODO(), cookie, req.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section2 {
		ID: existingSection.ID, //temporary
		Section: 2,
		Year: req.Year,
		ProjectName: req.ProjectName,
		ProjectScope: req.ProjectScope,
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

	var req Section3Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	timestamp := utils.TimeNow()

	section := db.Section3 {
		ID: req.ID, //temporary
		Section: 3,
		Year: req.Year,
		ActivityKind: req.ActivityKind,
		ThingsLearned: req.ThingsLearned,
		Level: req.Level,
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

	var req Section3Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	existingSection, err := e.db.GetSection3ByID(context.TODO(), cookie, req.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section3 {
		ID: existingSection.ID, //temporary
		Section: 3,
		Year: req.Year,
		ActivityKind: req.ActivityKind,
		ThingsLearned: req.ThingsLearned,
		Level: req.Level,
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

	var req Section4Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	timestamp := utils.TimeNow()

	section := db.Section4 {
		ID: req.ID, //temporary
		Section: 4,
		Year: req.Year,
		ActivityKind: req.ActivityKind,
		Scope: req.Scope,
		Level: req.Level,
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

	var req Section4Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	existingSection, err := e.db.GetSection4ByID(context.TODO(), cookie, req.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section4 {
		ID: existingSection.ID, //temporary
		Section: 4,
		Year: req.Year,
		ActivityKind: req.ActivityKind,
		Scope: req.Scope,
		Level: req.Level,
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

	var req Section5Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	timestamp := utils.TimeNow()

	section := db.Section5 {
		ID: req.ID, //temporary
		Section: 5,
		Year: req.Year,
		LeadershipRole: req.LeadershipRole,
		HoursSpent: req.HoursSpent,
		NumPeopleReached: req.NumPeopleReached,
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

	var req Section5Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	existingSection, err := e.db.GetSection5ByID(context.TODO(), cookie, req.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section5 {
		ID: existingSection.ID, //temporary
		Section: 5,
		Year: req.Year,
		LeadershipRole: req.LeadershipRole,
		HoursSpent: req.HoursSpent,
		NumPeopleReached: req.NumPeopleReached,
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

	var req Section6Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	timestamp := utils.TimeNow()

	section := db.Section6 {
		ID: req.ID, //temporary
		Section: 6,
		Year: req.Year,
		OrganizationName: req.OrganizationName,
		LeadershipRole: req.LeadershipRole,
		HoursSpent: req.HoursSpent,
		NumPeopleReached: req.NumPeopleReached,
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

	var req Section6Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	existingSection, err := e.db.GetSection6ByID(context.TODO(), cookie, req.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section6 {
		ID: existingSection.ID, //temporary
		Section: 6,
		Year: req.Year,
		OrganizationName: req.OrganizationName,
		LeadershipRole: req.LeadershipRole,
		HoursSpent: req.HoursSpent,
		NumPeopleReached: req.NumPeopleReached,
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

	var req Section7Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	timestamp := utils.TimeNow()

	section := db.Section7 {
		ID: req.ID, //temporary
		Section: 7,
		Year: req.Year,
		ClubMemberActivities: req.ClubMemberActivities,
		HoursSpent: req.HoursSpent,
		NumPeopleReached: req.NumPeopleReached,
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

	var req Section7Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	existingSection, err := e.db.GetSection7ByID(context.TODO(), cookie, req.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section7 {
		ID: existingSection.ID, //temporary
		Section: 7,
		Year: req.Year,
		ClubMemberActivities: req.ClubMemberActivities,
		HoursSpent: req.HoursSpent,
		NumPeopleReached: req.NumPeopleReached,
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

	var req Section8Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	timestamp := utils.TimeNow()

	section := db.Section8 {
		ID: req.ID, //temporary
		Section: 8,
		Year: req.Year,
		IndividualGroupActivities: req.IndividualGroupActivities,
		HoursSpent: req.HoursSpent,
		NumPeopleReached: req.NumPeopleReached,
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

	var req Section8Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	existingSection, err := e.db.GetSection8ByID(context.TODO(), cookie, req.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section8 {
		ID: existingSection.ID, //temporary
		Section: 8,
		Year: req.Year,
		IndividualGroupActivities: req.IndividualGroupActivities,
		HoursSpent: req.HoursSpent,
		NumPeopleReached: req.NumPeopleReached,
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

	var req Section9Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	timestamp := utils.TimeNow()

	section := db.Section9 {
		ID: req.ID, //temporary
		Section: 9,
		Year: req.Year,
		CommunicationType: req.CommunicationType,
		Topic: req.Topic,
		TimesGiven: req.TimesGiven,
		Location: req.Location,
		AudienceSize: req.AudienceSize,
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

	var req Section9Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	existingSection, err := e.db.GetSection9ByID(context.TODO(), cookie, req.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section9 {
		ID: existingSection.ID, //temporary
		Section: 9,
		Year: req.Year,
		CommunicationType: req.CommunicationType,
		Topic: req.Topic,
		TimesGiven: req.TimesGiven,
		Location: req.Location,
		AudienceSize: req.AudienceSize,
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

	var req Section10Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	timestamp := utils.TimeNow()

	section := db.Section10 {
		ID: req.ID, //temporary
		Section: 10,
		Year: req.Year,
		CommunicationType: req.CommunicationType,
		Topic: req.Topic,
		TimesGiven: req.TimesGiven,
		Location: req.Location,
		AudienceSize: req.AudienceSize,
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

	var req Section10Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	existingSection, err := e.db.GetSection10ByID(context.TODO(), cookie, req.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section10 {
		ID: existingSection.ID, //temporary
		Section: 10,
		Year: req.Year,
		CommunicationType: req.CommunicationType,
		Topic: req.Topic,
		TimesGiven: req.TimesGiven,
		Location: req.Location,
		AudienceSize: req.AudienceSize,
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

	var req Section11Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	timestamp := utils.TimeNow()

	section := db.Section11 {
		ID: req.ID, //temporary
		Section: 11,
		Year: req.Year,
		EventAndLevel: req.EventAndLevel,
		ExhibitsOrDivision: req.ExhibitsOrDivision,
		RibbonOrPlacings: req.RibbonOrPlacings,
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

	var req Section11Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	existingSection, err := e.db.GetSection11ByID(context.TODO(), cookie, req.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section11 {
		ID: existingSection.ID, //temporary
		Section: 11,
		Year: req.Year,
		EventAndLevel: req.EventAndLevel,
		ExhibitsOrDivision: req.ExhibitsOrDivision,
		RibbonOrPlacings: req.RibbonOrPlacings,
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

	var req Section12Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	timestamp := utils.TimeNow()

	section := db.Section12 {
		ID: req.ID, //temporary
		Section: 12,
		Year: req.Year,
		ContestOrEvent: req.ContestOrEvent,
		RecognitionReceived: req.RecognitionReceived,
		Level: req.Level,
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

	var req Section12Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	existingSection, err := e.db.GetSection12ByID(context.TODO(), cookie, req.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section12 {
		ID: existingSection.ID, //temporary
		Section: 12,
		Year: req.Year,
		ContestOrEvent: req.ContestOrEvent,
		RecognitionReceived: req.RecognitionReceived,
		Level: req.Level,
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

	var req Section13Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	timestamp := utils.TimeNow()

	section := db.Section13 {
		ID: req.ID, //temporary
		Section: 13,
		Year: req.Year,
		RecognitionType: req.RecognitionType,
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

	var req Section13Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	existingSection, err := e.db.GetSection13ByID(context.TODO(), cookie, req.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section13 {
		ID: existingSection.ID, //temporary
		Section: 13,
		Year: req.Year,
		RecognitionType: req.RecognitionType,
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

	var req Section14Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	timestamp := utils.TimeNow()

	section := db.Section14 {
		ID: req.ID, //temporary
		Section: 14,
		Year: req.Year,
		RecognitionType: req.RecognitionType,
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

	var req Section14Req
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	existingSection, err := e.db.GetSection14ByID(context.TODO(), cookie, req.ID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSection := db.Section14 {
		ID: existingSection.ID, //temporary
		Section: 14,
		Year: req.Year,
		RecognitionType: req.RecognitionType,
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
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

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