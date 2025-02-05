package db

import (
	"context"
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type Resume struct {
	Section1Data  []Section1  `json:"section_1_data"`
	Section2Data  []Section2  `json:"section_2_data"`
	Section3Data  []Section3  `json:"section_3_data"`
	Section4Data  []Section4  `json:"section_4_data"`
	Section5Data  []Section5  `json:"section_5_data"`
	Section6Data  []Section6  `json:"section_6_data"`
	Section7Data  []Section7  `json:"section_7_data"`
	Section8Data  []Section8  `json:"section_8_data"`
	Section9Data  []Section9  `json:"section_9_data"`
	Section10Data []Section10 `json:"section_10_data"`
	Section11Data []Section11 `json:"section_11_data"`
	Section12Data []Section12 `json:"section_12_data"`
	Section13Data []Section13 `json:"section_13_data"`
	Section14Data []Section14 `json:"section_14_data"`
}

type Section1 struct {
	ID string `json:"id"`
	Section int `json:"section"`
	Year string `json:"year"`
	Grade int `json:"grade"`
	ClubName string `json:"club_name"`
	NumInClub int `json:"num_in_club"`
	ClubLeader string `json:"club_leader"`
	MeetingsHeld int `json:"meetings_held"`
	MeetingsAttended int `json:"meetings_attended"`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

type Section2 struct {
	ID string `json:"id"`
	Section int `json:"section"`
	Year string `json:"year"`
	ProjectName string `json:"project_name"`
	ProjectScope string `json:"project_scope"`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

type Section3 struct {
	ID string `json:"id"`
	Section int `json:"section"`
	Year string `json:"year"`
	ActivityKind string `json:"activity_kind"`
	ThingsLearned string `json:"things_learned"`
	Level string `json:"level"`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

type Section4 struct {
	ID string `json:"id"`
	Section int `json:"section"`
	Year string `json:"year"`
	ActivityKind string `json:"activity_kind"`
	Scope string `json:"scope"`
	Level string `json:"level"`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

type Section5 struct {
	ID string `json:"id"`
	Section int `json:"section"`
	Year string `json:"year"`
	LeadershipRole string `json:"leadership_role"`
	HoursSpent int `json:"hours_spent"`
	NumPeopleReached int `json:"num_people_reached"`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

type Section6 struct {
	ID string `json:"id"`
	Section int `json:"section"`
	Year string `json:"year"`
	OrganizationName string `json:"organization_name"`
	LeadershipRole string `json:"leadership_role"`
	HoursSpent int `json:"hours_spent"`
	NumPeopleReached int `json:"num_people_reached"`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

type Section7 struct {
	ID string `json:"id"`
	Section int `json:"section"`
	Year string `json:"year"`
	ClubMemberActivities string `json:"club_member_activities"`
	HoursSpent int `json:"hours_spent"`
	NumPeopleReached int `json:"num_people_reached"`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

type Section8 struct {
	ID string `json:"id"`
	Section int `json:"section"`
	Year string `json:"year"`
	IndividualGroupActivities string `json:"individual_group_activities"`
	HoursSpent int `json:"hours_spent"`
	NumPeopleReached int `json:"num_people_reached"`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

type Section9 struct {
	ID string `json:"id"`
	Section int `json:"section"`
	Year string `json:"year"`
	CommunicationType string `json:"communication_type"`
	Topic string `json:"topic"`
	TimesGiven int `json:"times_given"`
	Location string `json:"location"`
	AudienceSize int `json:"audience_size"`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

type Section10 struct {
	ID string `json:"id"`
	Section int `json:"section"`
	Year string `json:"year"`
	CommunicationType string `json:"communication_type"`
	Topic string `json:"topic"`
	TimesGiven int `json:"times_given"`
	Location string `json:"location"`
	AudienceSize string `json:"audience_size"`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

type Section11 struct {
	ID string `json:"id"`
	Section int `json:"section"`
	Year string `json:"year"`
	EventAndLevel string `json:"event_and_level"`
	ExhibitsOrDivision string `json:"exhibits_or_division"`
	RibbonOrPlacings string `json:"ribbon_or_placings"`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

type Section12 struct {
	ID string `json:"id"`
	Section int `json:"section"`
	Year string `json:"year"`
	ContestOrEvent string `json:"contest_or_event"`
	RecognitionReceived string `json:"recognition_received"`
	Level string `json:"level"`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

type Section13 struct {
	ID string `json:"id"`
	Section int `json:"section"`
	Year string `json:"year"`
	RecognitionType string `json:"recognition_type"`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

type Section14 struct {
	ID string `json:"id"`
	Section int `json:"section"`
	Year string `json:"year"`
	RecognitionType string `json:"recognition_type"`
	UserID string `json:"userid"`
	GenericDatabaseInfo
}

/*******************************
* ALL SECTIONS
********************************/
func (env *env) GetResume(ctx context.Context, userid string) (Resume, error) {

	env.logger.Info("Getting resume")
	resume := Resume{}
	
	section1Data, err := env.GetSection1(ctx, userid)
	if err != nil {
		return resume, err
	}

	section2Data, err := env.GetSection2(ctx, userid)
	if err != nil {
		return resume, err
	}

	section3Data, err := env.GetSection3(ctx, userid)
	if err != nil {
		return resume, err
	}

	section4Data, err := env.GetSection4(ctx, userid)
	if err != nil {
		return resume, err
	}

	section5Data, err := env.GetSection5(ctx, userid)
	if err != nil {
		return resume, err
	}

	section6Data, err := env.GetSection6(ctx, userid)
	if err != nil {
		return resume, err
	}

	section7Data, err := env.GetSection7(ctx, userid)
	if err != nil {
		return resume, err
	}

	section8Data, err := env.GetSection8(ctx, userid)
	if err != nil {
		return resume, err
	}

	section9Data, err := env.GetSection9(ctx, userid)
	if err != nil {
		return resume, err
	}

	section10Data, err := env.GetSection10(ctx, userid)
	if err != nil {
		return resume, err
	}

	section11Data, err := env.GetSection11(ctx, userid)
	if err != nil {
		return resume, err
	}

	section12Data, err := env.GetSection12(ctx, userid)
	if err != nil {
		return resume, err
	}

	section13Data, err := env.GetSection13(ctx, userid)
	if err != nil {
		return resume, err
	}

	section14Data, err := env.GetSection14(ctx, userid)
	if err != nil {
		return resume, err
	}

	resume.Section1Data = section1Data
	resume.Section2Data = section2Data
	resume.Section3Data = section3Data
	resume.Section4Data = section4Data
	resume.Section5Data = section5Data
	resume.Section6Data = section6Data
	resume.Section7Data = section7Data
	resume.Section8Data = section8Data
	resume.Section9Data = section9Data
	resume.Section10Data = section10Data
	resume.Section11Data = section11Data
	resume.Section12Data = section12Data
	resume.Section13Data = section13Data
	resume.Section14Data = section14Data

	return resume, nil

}

/*******************************
* SECTION 1
********************************/

func (env *env) GetSection1ByID(ctx context.Context, userid string, sectionid string) (Section1, error) {

	env.logger.Info("Getting Section 1 by ID")
	section := Section1{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	response, err := container.ReadItem(ctx, partitionKey, sectionid, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	return section, nil

}

func (env *env) GetSection1(ctx context.Context, userid string) ([]Section1, error) {

	env.logger.Info("Getting all Section 1 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section1{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	query := "SELECT * FROM sections s WHERE s.userid = @id AND s.section = @section"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@id", Value: userid},
			{Name: "@section", Value: 1},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section1{}

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return []Section1{}, err
		}

		for _, bytes := range response.Items {
			section := Section1{}
			err := json.Unmarshal(bytes, &section)
			if err != nil {
				return []Section1{}, err
			}
			sections = append(sections, section)
		}
	}

	return sections, nil

}

func (env *env) UpsertSection1(ctx context.Context, section Section1) (interface{}, error) {

	env.logger.Info("Upserting section 1")

	container, err := env.client.NewContainer("sections")

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return nil, err
	}

	response, err := container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}

/*******************************
* SECTION 2
********************************/

func (env *env) GetSection2ByID(ctx context.Context, userid string, sectionid string) (Section2, error) {

	env.logger.Info("Getting Section 2 by ID")
	section := Section2{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	response, err := container.ReadItem(ctx, partitionKey, sectionid, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	return section, nil

}

func (env *env) GetSection2(ctx context.Context, userid string) ([]Section2, error) {

	env.logger.Info("Getting all Section 2 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section2{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	query := "SELECT * FROM sections s WHERE s.userid = @id AND s.section = @section"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@id", Value: userid},
			{Name: "@section", Value: 2},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section2{}

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return []Section2{}, err
		}

		for _, bytes := range response.Items {
			section := Section2{}
			err := json.Unmarshal(bytes, &section)
			if err != nil {
				return []Section2{}, err
			}
			sections = append(sections, section)
		}
	}

	return sections, nil

}

func (env *env) UpsertSection2(ctx context.Context, section Section2) (interface{}, error) {

	env.logger.Info("Upserting section 2")

	container, err := env.client.NewContainer("sections")

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return nil, err
	}

	response, err := container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}

/*******************************
* SECTION 3
********************************/

func (env *env) GetSection3ByID(ctx context.Context, userid string, sectionid string) (Section3, error) {

	env.logger.Info("Getting Section 3 by ID")
	section := Section3{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	response, err := container.ReadItem(ctx, partitionKey, sectionid, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	return section, nil

}

func (env *env) GetSection3(ctx context.Context, userid string) ([]Section3, error) {

	env.logger.Info("Getting all Section 3 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section3{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	query := "SELECT * FROM sections s WHERE s.userid = @id AND s.section = @section"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@id", Value: userid},
			{Name: "@section", Value: 3},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section3{}

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return []Section3{}, err
		}

		for _, bytes := range response.Items {
			section := Section3{}
			err := json.Unmarshal(bytes, &section)
			if err != nil {
				return []Section3{}, err
			}
			sections = append(sections, section)
		}
	}

	return sections, nil

}

func (env *env) UpsertSection3(ctx context.Context, section Section3) (interface{}, error) {

	env.logger.Info("Upserting section 3")

	container, err := env.client.NewContainer("sections")

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return nil, err
	}

	response, err := container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}

/*******************************
* SECTION 4
********************************/

func (env *env) GetSection4ByID(ctx context.Context, userid string, sectionid string) (Section4, error) {

	env.logger.Info("Getting Section 4 by ID")
	section := Section4{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	response, err := container.ReadItem(ctx, partitionKey, sectionid, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	return section, nil

}

func (env *env) GetSection4(ctx context.Context, userid string) ([]Section4, error) {

	env.logger.Info("Getting all Section 4 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section4{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	query := "SELECT * FROM sections s WHERE s.userid = @id AND s.section = @section"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@id", Value: userid},
			{Name: "@section", Value: 4},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section4{}

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return []Section4{}, err
		}

		for _, bytes := range response.Items {
			section := Section4{}
			err := json.Unmarshal(bytes, &section)
			if err != nil {
				return []Section4{}, err
			}
			sections = append(sections, section)
		}
	}

	return sections, nil

}

func (env *env) UpsertSection4(ctx context.Context, section Section4) (interface{}, error) {

	env.logger.Info("Upserting section 4")

	container, err := env.client.NewContainer("sections")

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return nil, err
	}

	response, err := container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}

/*******************************
* SECTION 5
********************************/

func (env *env) GetSection5ByID(ctx context.Context, userid string, sectionid string) (Section5, error) {

	env.logger.Info("Getting Section 5 by ID")
	section := Section5{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	response, err := container.ReadItem(ctx, partitionKey, sectionid, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	return section, nil

}

func (env *env) GetSection5(ctx context.Context, userid string) ([]Section5, error) {

	env.logger.Info("Getting all Section 5 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section5{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	query := "SELECT * FROM sections s WHERE s.userid = @id AND s.section = @section"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@id", Value: userid},
			{Name: "@section", Value: 5},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section5{}

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return []Section5{}, err
		}

		for _, bytes := range response.Items {
			section := Section5{}
			err := json.Unmarshal(bytes, &section)
			if err != nil {
				return []Section5{}, err
			}
			sections = append(sections, section)
		}
	}

	return sections, nil

}

func (env *env) UpsertSection5(ctx context.Context, section Section5) (interface{}, error) {

	env.logger.Info("Upserting section 5")

	container, err := env.client.NewContainer("sections")

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return nil, err
	}

	response, err := container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}

/*******************************
* SECTION 6
********************************/

func (env *env) GetSection6ByID(ctx context.Context, userid string, sectionid string) (Section6, error) {

	env.logger.Info("Getting Section 6 by ID")
	section := Section6{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	response, err := container.ReadItem(ctx, partitionKey, sectionid, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	return section, nil

}

func (env *env) GetSection6(ctx context.Context, userid string) ([]Section6, error) {

	env.logger.Info("Getting all Section 6 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section6{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	query := "SELECT * FROM sections s WHERE s.userid = @id AND s.section = @section"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@id", Value: userid},
			{Name: "@section", Value: 6},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section6{}

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return []Section6{}, err
		}

		for _, bytes := range response.Items {
			section := Section6{}
			err := json.Unmarshal(bytes, &section)
			if err != nil {
				return []Section6{}, err
			}
			sections = append(sections, section)
		}
	}

	return sections, nil

}

func (env *env) UpsertSection6(ctx context.Context, section Section6) (interface{}, error) {

	env.logger.Info("Upserting section 6")

	container, err := env.client.NewContainer("sections")

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return nil, err
	}

	response, err := container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}

/*******************************
* SECTION 7
********************************/

func (env *env) GetSection7ByID(ctx context.Context, userid string, sectionid string) (Section7, error) {

	env.logger.Info("Getting Section 7 by ID")
	section := Section7{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	response, err := container.ReadItem(ctx, partitionKey, sectionid, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	return section, nil

}

func (env *env) GetSection7(ctx context.Context, userid string) ([]Section7, error) {

	env.logger.Info("Getting all Section 7 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section7{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	query := "SELECT * FROM sections s WHERE s.userid = @id AND s.section = @section"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@id", Value: userid},
			{Name: "@section", Value: 7},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section7{}

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return []Section7{}, err
		}

		for _, bytes := range response.Items {
			section := Section7{}
			err := json.Unmarshal(bytes, &section)
			if err != nil {
				return []Section7{}, err
			}
			sections = append(sections, section)
		}
	}

	return sections, nil

}

func (env *env) UpsertSection7(ctx context.Context, section Section7) (interface{}, error) {

	env.logger.Info("Upserting section 7")

	container, err := env.client.NewContainer("sections")

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return nil, err
	}

	response, err := container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}

/*******************************
* SECTION 8
********************************/

func (env *env) GetSection8ByID(ctx context.Context, userid string, sectionid string) (Section8, error) {

	env.logger.Info("Getting Section 8 by ID")
	section := Section8{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	response, err := container.ReadItem(ctx, partitionKey, sectionid, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	return section, nil

}

func (env *env) GetSection8(ctx context.Context, userid string) ([]Section8, error) {

	env.logger.Info("Getting all Section 8 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section8{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	query := "SELECT * FROM sections s WHERE s.userid = @id AND s.section = @section"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@id", Value: userid},
			{Name: "@section", Value: 8},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section8{}

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return []Section8{}, err
		}

		for _, bytes := range response.Items {
			section := Section8{}
			err := json.Unmarshal(bytes, &section)
			if err != nil {
				return []Section8{}, err
			}
			sections = append(sections, section)
		}
	}

	return sections, nil

}

func (env *env) UpsertSection8(ctx context.Context, section Section8) (interface{}, error) {

	env.logger.Info("Upserting section 8")

	container, err := env.client.NewContainer("sections")

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return nil, err
	}

	response, err := container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}

/*******************************
* SECTION 9
********************************/

func (env *env) GetSection9ByID(ctx context.Context, userid string, sectionid string) (Section9, error) {

	env.logger.Info("Getting Section 9 by ID")
	section := Section9{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	response, err := container.ReadItem(ctx, partitionKey, sectionid, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	return section, nil

}

func (env *env) GetSection9(ctx context.Context, userid string) ([]Section9, error) {

	env.logger.Info("Getting all Section 9 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section9{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	query := "SELECT * FROM sections s WHERE s.userid = @id AND s.section = @section"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@id", Value: userid},
			{Name: "@section", Value: 9},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section9{}

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return []Section9{}, err
		}

		for _, bytes := range response.Items {
			section := Section9{}
			err := json.Unmarshal(bytes, &section)
			if err != nil {
				return []Section9{}, err
			}
			sections = append(sections, section)
		}
	}

	return sections, nil

}

func (env *env) UpsertSection9(ctx context.Context, section Section9) (interface{}, error) {

	env.logger.Info("Upserting section 9")

	container, err := env.client.NewContainer("sections")

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return nil, err
	}

	response, err := container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}

/*******************************
* SECTION 10
********************************/

func (env *env) GetSection10ByID(ctx context.Context, userid string, sectionid string) (Section10, error) {

	env.logger.Info("Getting Section 10 by ID")
	section := Section10{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	response, err := container.ReadItem(ctx, partitionKey, sectionid, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	return section, nil

}

func (env *env) GetSection10(ctx context.Context, userid string) ([]Section10, error) {

	env.logger.Info("Getting all Section 10 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section10{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	query := "SELECT * FROM sections s WHERE s.userid = @id AND s.section = @section"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@id", Value: userid},
			{Name: "@section", Value: 10},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section10{}

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return []Section10{}, err
		}

		for _, bytes := range response.Items {
			section := Section10{}
			err := json.Unmarshal(bytes, &section)
			if err != nil {
				return []Section10{}, err
			}
			sections = append(sections, section)
		}
	}

	return sections, nil

}

func (env *env) UpsertSection10(ctx context.Context, section Section10) (interface{}, error) {

	env.logger.Info("Upserting section 10")

	container, err := env.client.NewContainer("sections")

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return nil, err
	}

	response, err := container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}

/*******************************
* SECTION 11
********************************/

func (env *env) GetSection11ByID(ctx context.Context, userid string, sectionid string) (Section11, error) {

	env.logger.Info("Getting Section 11 by ID")
	section := Section11{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	response, err := container.ReadItem(ctx, partitionKey, sectionid, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	return section, nil

}

func (env *env) GetSection11(ctx context.Context, userid string) ([]Section11, error) {

	env.logger.Info("Getting all Section 11 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section11{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	query := "SELECT * FROM sections s WHERE s.userid = @id AND s.section = @section"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@id", Value: userid},
			{Name: "@section", Value: 11},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section11{}

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return []Section11{}, err
		}

		for _, bytes := range response.Items {
			section := Section11{}
			err := json.Unmarshal(bytes, &section)
			if err != nil {
				return []Section11{}, err
			}
			sections = append(sections, section)
		}
	}

	return sections, nil

}

func (env *env) UpsertSection11(ctx context.Context, section Section11) (interface{}, error) {

	env.logger.Info("Upserting section 11")

	container, err := env.client.NewContainer("sections")

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return nil, err
	}

	response, err := container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}

/*******************************
* SECTION 12
********************************/

func (env *env) GetSection12ByID(ctx context.Context, userid string, sectionid string) (Section12, error) {

	env.logger.Info("Getting Section 12 by ID")
	section := Section12{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	response, err := container.ReadItem(ctx, partitionKey, sectionid, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	return section, nil

}

func (env *env) GetSection12(ctx context.Context, userid string) ([]Section12, error) {

	env.logger.Info("Getting all Section 10 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section12{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	query := "SELECT * FROM sections s WHERE s.userid = @id AND s.section = @section"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@id", Value: userid},
			{Name: "@section", Value: 12},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section12{}

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return []Section12{}, err
		}

		for _, bytes := range response.Items {
			section := Section12{}
			err := json.Unmarshal(bytes, &section)
			if err != nil {
				return []Section12{}, err
			}
			sections = append(sections, section)
		}
	}

	return sections, nil

}

func (env *env) UpsertSection12(ctx context.Context, section Section12) (interface{}, error) {

	env.logger.Info("Upserting section 12")

	container, err := env.client.NewContainer("sections")

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return nil, err
	}

	response, err := container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}

/*******************************
* SECTION 13
********************************/

func (env *env) GetSection13ByID(ctx context.Context, userid string, sectionid string) (Section13, error) {

	env.logger.Info("Getting Section 13 by ID")
	section := Section13{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	response, err := container.ReadItem(ctx, partitionKey, sectionid, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	return section, nil

}

func (env *env) GetSection13(ctx context.Context, userid string) ([]Section13, error) {

	env.logger.Info("Getting all Section 13 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section13{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	query := "SELECT * FROM sections s WHERE s.userid = @id AND s.section = @section"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@id", Value: userid},
			{Name: "@section", Value: 13},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section13{}

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return []Section13{}, err
		}

		for _, bytes := range response.Items {
			section := Section13{}
			err := json.Unmarshal(bytes, &section)
			if err != nil {
				return []Section13{}, err
			}
			sections = append(sections, section)
		}
	}

	return sections, nil

}

func (env *env) UpsertSection13(ctx context.Context, section Section13) (interface{}, error) {

	env.logger.Info("Upserting section 13")

	container, err := env.client.NewContainer("sections")

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return nil, err
	}

	response, err := container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}

/*******************************
* SECTION 13
********************************/

func (env *env) GetSection14ByID(ctx context.Context, userid string, sectionid string) (Section14, error) {

	env.logger.Info("Getting Section 14 by ID")
	section := Section14{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	response, err := container.ReadItem(ctx, partitionKey, sectionid, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	return section, nil

}

func (env *env) GetSection14(ctx context.Context, userid string) ([]Section14, error) {

	env.logger.Info("Getting all Section 14 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section14{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	query := "SELECT * FROM sections s WHERE s.userid = @id AND s.section = @section"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@id", Value: userid},
			{Name: "@section", Value: 14},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section14{}

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return []Section14{}, err
		}

		for _, bytes := range response.Items {
			section := Section14{}
			err := json.Unmarshal(bytes, &section)
			if err != nil {
				return []Section14{}, err
			}
			sections = append(sections, section)
		}
	}

	return sections, nil

}

func (env *env) UpsertSection14(ctx context.Context, section Section14) (interface{}, error) {

	env.logger.Info("Upserting section 14")

	container, err := env.client.NewContainer("sections")

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return nil, err
	}

	response, err := container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}

/*******************************
* DELETING
********************************/

func (env *env) RemoveSection(ctx context.Context, userid string, sectionid string) (interface{}, error) {
	
	env.logger.Info("Removing section")

	container, err := env.client.NewContainer("sections")

	partitionKey := azcosmos.NewPartitionKeyString(userid)

	response, err := container.DeleteItem(ctx, partitionKey, sectionid, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}