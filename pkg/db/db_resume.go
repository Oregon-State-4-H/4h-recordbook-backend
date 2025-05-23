package db

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
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

type GenericSectionInfo struct {
	Section int    `json:"section"`
	Year    string `json:"year"`
	UserID  string `json:"user_id"`
	GenericDatabaseInfo
}

/*******************************
* ALL SECTIONS
********************************/
func (env *env) GetResume(ctx context.Context, userID string) (Resume, error) {

	env.logger.Info("Getting resume")
	resume := Resume{}

	var err error

	paginationOptions := PaginationOptions{
		Page:         0,
		PerPage:      500,
		SortByNewest: false,
	}

	resume.Section1Data, err = env.GetSection1sByUser(ctx, userID, paginationOptions)
	if err != nil {
		return resume, err
	}

	resume.Section2Data, err = env.GetSection2sByUser(ctx, userID, paginationOptions)
	if err != nil {
		return resume, err
	}

	resume.Section3Data, err = env.GetSection3sByUser(ctx, userID, paginationOptions)
	if err != nil {
		return resume, err
	}

	resume.Section4Data, err = env.GetSection4sByUser(ctx, userID, paginationOptions)
	if err != nil {
		return resume, err
	}

	resume.Section5Data, err = env.GetSection5sByUser(ctx, userID, paginationOptions)
	if err != nil {
		return resume, err
	}

	resume.Section6Data, err = env.GetSection6sByUser(ctx, userID, paginationOptions)
	if err != nil {
		return resume, err
	}

	resume.Section7Data, err = env.GetSection7sByUser(ctx, userID, paginationOptions)
	if err != nil {
		return resume, err
	}

	resume.Section8Data, err = env.GetSection8sByUser(ctx, userID, paginationOptions)
	if err != nil {
		return resume, err
	}

	resume.Section9Data, err = env.GetSection9sByUser(ctx, userID, paginationOptions)
	if err != nil {
		return resume, err
	}

	resume.Section10Data, err = env.GetSection10sByUser(ctx, userID, paginationOptions)
	if err != nil {
		return resume, err
	}

	resume.Section11Data, err = env.GetSection11sByUser(ctx, userID, paginationOptions)
	if err != nil {
		return resume, err
	}

	resume.Section12Data, err = env.GetSection12sByUser(ctx, userID, paginationOptions)
	if err != nil {
		return resume, err
	}

	resume.Section13Data, err = env.GetSection13sByUser(ctx, userID, paginationOptions)
	if err != nil {
		return resume, err
	}

	resume.Section14Data, err = env.GetSection14sByUser(ctx, userID, paginationOptions)
	if err != nil {
		return resume, err
	}

	return resume, nil

}

/*******************************
* SECTION 1
********************************/

type Section1 struct {
	ID               string `json:"id"`
	Nickname         string `json:"nickname"`
	Grade            int    `json:"grade"`
	ClubName         string `json:"club_name"`
	NumInClub        int    `json:"num_in_club"`
	ClubLeader       string `json:"club_leader"`
	MeetingsHeld     int    `json:"meetings_held"`
	MeetingsAttended int    `json:"meetings_attended"`
	GenericSectionInfo
}

func (env *env) GetSection1ByID(ctx context.Context, userID string, sectionID string) (Section1, error) {

	env.logger.Info("Getting Section 1 by ID")
	section := Section1{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.ReadItem(ctx, partitionKey, sectionID, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	if section.Section != 1 {
		err = &azcore.ResponseError{
			StatusCode: http.StatusNotFound,
		}
		return section, err
	}

	return section, nil

}

func (env *env) GetSection1sByUser(ctx context.Context, userID string, paginationOptions PaginationOptions) ([]Section1, error) {

	env.logger.Info("Getting all Section 1 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section1{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	sortOrder := "ASC"
	if paginationOptions.SortByNewest {
		sortOrder = "DESC"
	}

	query := fmt.Sprintf("SELECT * FROM sections s WHERE s.user_id = @user_id AND s.section = @section ORDER BY s.created %s", sortOrder)

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@section", Value: 1},
		},
		PageSizeHint: int32(paginationOptions.PerPage),
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section1{}
	currentPage := 0

	for pager.More() {

		if currentPage == paginationOptions.Page {
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

			return sections, nil

		} else {
			_, err := pager.NextPage(ctx)
			if err != nil {
				return []Section1{}, err
			}
			currentPage++
		}

	}

	return sections, nil

}

func (env *env) UpsertSection1(ctx context.Context, section Section1) (Section1, error) {

	env.logger.Info("Upserting section 1")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return section, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return section, err
	}

	return section, nil

}

/*******************************
* SECTION 2
********************************/

type Section2 struct {
	ID           string `json:"id"`
	ProjectName  string `json:"project_name"`
	ProjectScope string `json:"project_scope"`
	GenericSectionInfo
}

func (env *env) GetSection2ByID(ctx context.Context, userID string, sectionID string) (Section2, error) {

	env.logger.Info("Getting Section 2 by ID")
	section := Section2{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.ReadItem(ctx, partitionKey, sectionID, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	if section.Section != 2 {
		err = &azcore.ResponseError{
			StatusCode: http.StatusNotFound,
		}
		return section, err
	}

	return section, nil

}

func (env *env) GetSection2sByUser(ctx context.Context, userID string, paginationOptions PaginationOptions) ([]Section2, error) {

	env.logger.Info("Getting all Section 2 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section2{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	sortOrder := "ASC"
	if paginationOptions.SortByNewest {
		sortOrder = "DESC"
	}

	query := fmt.Sprintf("SELECT * FROM sections s WHERE s.user_id = @user_id AND s.section = @section ORDER BY s.created %s", sortOrder)

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@section", Value: 2},
		},
		PageSizeHint: int32(paginationOptions.PerPage),
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section2{}
	currentPage := 0

	for pager.More() {

		if currentPage == paginationOptions.Page {
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

			return sections, nil

		} else {
			_, err := pager.NextPage(ctx)
			if err != nil {
				return []Section2{}, err
			}
			currentPage++
		}

	}

	return sections, nil

}

func (env *env) UpsertSection2(ctx context.Context, section Section2) (Section2, error) {

	env.logger.Info("Upserting section 2")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return section, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return section, err
	}

	return section, nil

}

/*******************************
* SECTION 3
********************************/

type Section3 struct {
	ID            string `json:"id"`
	Nickname      string `json:"nickname"`
	ActivityKind  string `json:"activity_kind"`
	ThingsLearned string `json:"things_learned"`
	Level         string `json:"level"`
	GenericSectionInfo
}

func (env *env) GetSection3ByID(ctx context.Context, userID string, sectionID string) (Section3, error) {

	env.logger.Info("Getting Section 3 by ID")
	section := Section3{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.ReadItem(ctx, partitionKey, sectionID, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	if section.Section != 3 {
		err = &azcore.ResponseError{
			StatusCode: http.StatusNotFound,
		}
		return section, err
	}

	return section, nil

}

func (env *env) GetSection3sByUser(ctx context.Context, userID string, paginationOptions PaginationOptions) ([]Section3, error) {

	env.logger.Info("Getting all Section 3 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section3{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	sortOrder := "ASC"
	if paginationOptions.SortByNewest {
		sortOrder = "DESC"
	}

	query := fmt.Sprintf("SELECT * FROM sections s WHERE s.user_id = @user_id AND s.section = @section ORDER BY s.created %s", sortOrder)

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@section", Value: 3},
		},
		PageSizeHint: int32(paginationOptions.PerPage),
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section3{}
	currentPage := 0

	for pager.More() {

		if currentPage == paginationOptions.Page {
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

			return sections, nil

		} else {
			_, err := pager.NextPage(ctx)
			if err != nil {
				return []Section3{}, err
			}
			currentPage++
		}

	}

	return sections, nil

}

func (env *env) UpsertSection3(ctx context.Context, section Section3) (Section3, error) {

	env.logger.Info("Upserting section 3")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return section, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return section, err
	}

	return section, nil

}

/*******************************
* SECTION 4
********************************/

type Section4 struct {
	ID           string `json:"id"`
	Nickname     string `json:"nickname"`
	ActivityKind string `json:"activity_kind"`
	Scope        string `json:"scope"`
	Level        string `json:"level"`
	GenericSectionInfo
}

func (env *env) GetSection4ByID(ctx context.Context, userID string, sectionID string) (Section4, error) {

	env.logger.Info("Getting Section 4 by ID")
	section := Section4{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.ReadItem(ctx, partitionKey, sectionID, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	if section.Section != 4 {
		err = &azcore.ResponseError{
			StatusCode: http.StatusNotFound,
		}
		return section, err
	}

	return section, nil

}

func (env *env) GetSection4sByUser(ctx context.Context, userID string, paginationOptions PaginationOptions) ([]Section4, error) {

	env.logger.Info("Getting all Section 4 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section4{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	sortOrder := "ASC"
	if paginationOptions.SortByNewest {
		sortOrder = "DESC"
	}

	query := fmt.Sprintf("SELECT * FROM sections s WHERE s.user_id = @user_id AND s.section = @section ORDER BY s.created %s", sortOrder)

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@section", Value: 4},
		},
		PageSizeHint: int32(paginationOptions.PerPage),
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section4{}
	currentPage := 0

	for pager.More() {

		if currentPage == paginationOptions.Page {
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

			return sections, nil

		} else {
			_, err := pager.NextPage(ctx)
			if err != nil {
				return []Section4{}, err
			}
			currentPage++
		}

	}

	return sections, nil

}

func (env *env) UpsertSection4(ctx context.Context, section Section4) (Section4, error) {

	env.logger.Info("Upserting section 4")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return section, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return section, err
	}

	return section, nil

}

/*******************************
* SECTION 5
********************************/

type Section5 struct {
	ID               string `json:"id"`
	Nickname         string `json:"nickname"`
	LeadershipRole   string `json:"leadership_role"`
	HoursSpent       int    `json:"hours_spent"`
	NumPeopleReached int    `json:"num_people_reached"`
	GenericSectionInfo
}

func (env *env) GetSection5ByID(ctx context.Context, userID string, sectionID string) (Section5, error) {

	env.logger.Info("Getting Section 5 by ID")
	section := Section5{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.ReadItem(ctx, partitionKey, sectionID, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	if section.Section != 5 {
		err = &azcore.ResponseError{
			StatusCode: http.StatusNotFound,
		}
		return section, err
	}

	return section, nil

}

func (env *env) GetSection5sByUser(ctx context.Context, userID string, paginationOptions PaginationOptions) ([]Section5, error) {

	env.logger.Info("Getting all Section 5 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section5{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	sortOrder := "ASC"
	if paginationOptions.SortByNewest {
		sortOrder = "DESC"
	}

	query := fmt.Sprintf("SELECT * FROM sections s WHERE s.user_id = @user_id AND s.section = @section ORDER BY s.created %s", sortOrder)

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@section", Value: 5},
		},
		PageSizeHint: int32(paginationOptions.PerPage),
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section5{}
	currentPage := 0

	for pager.More() {

		if currentPage == paginationOptions.Page {
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

			return sections, nil

		} else {
			_, err := pager.NextPage(ctx)
			if err != nil {
				return []Section5{}, err
			}
			currentPage++
		}

	}

	return sections, nil

}

func (env *env) UpsertSection5(ctx context.Context, section Section5) (Section5, error) {

	env.logger.Info("Upserting section 5")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return section, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return section, err
	}

	return section, nil

}

/*******************************
* SECTION 6
********************************/

type Section6 struct {
	ID               string `json:"id"`
	Nickname         string `json:"nickname"`
	OrganizationName string `json:"organization_name"`
	LeadershipRole   string `json:"leadership_role"`
	HoursSpent       int    `json:"hours_spent"`
	NumPeopleReached int    `json:"num_people_reached"`
	GenericSectionInfo
}

func (env *env) GetSection6ByID(ctx context.Context, userID string, sectionID string) (Section6, error) {

	env.logger.Info("Getting Section 6 by ID")
	section := Section6{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.ReadItem(ctx, partitionKey, sectionID, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	if section.Section != 6 {
		err = &azcore.ResponseError{
			StatusCode: http.StatusNotFound,
		}
		return section, err
	}

	return section, nil

}

func (env *env) GetSection6sByUser(ctx context.Context, userID string, paginationOptions PaginationOptions) ([]Section6, error) {

	env.logger.Info("Getting all Section 6 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section6{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	sortOrder := "ASC"
	if paginationOptions.SortByNewest {
		sortOrder = "DESC"
	}

	query := fmt.Sprintf("SELECT * FROM sections s WHERE s.user_id = @user_id AND s.section = @section ORDER BY s.created %s", sortOrder)

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@section", Value: 6},
		},
		PageSizeHint: int32(paginationOptions.PerPage),
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section6{}
	currentPage := 0

	for pager.More() {

		if currentPage == paginationOptions.Page {
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

			return sections, nil

		} else {
			_, err := pager.NextPage(ctx)
			if err != nil {
				return []Section6{}, err
			}
			currentPage++
		}

	}

	return sections, nil

}

func (env *env) UpsertSection6(ctx context.Context, section Section6) (Section6, error) {

	env.logger.Info("Upserting section 6")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return section, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return section, err
	}

	return section, nil

}

/*******************************
* SECTION 7
********************************/

type Section7 struct {
	ID                   string `json:"id"`
	Nickname             string `json:"nickname"`
	ClubMemberActivities string `json:"club_member_activities"`
	HoursSpent           int    `json:"hours_spent"`
	NumPeopleReached     int    `json:"num_people_reached"`
	GenericSectionInfo
}

func (env *env) GetSection7ByID(ctx context.Context, userID string, sectionID string) (Section7, error) {

	env.logger.Info("Getting Section 7 by ID")
	section := Section7{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.ReadItem(ctx, partitionKey, sectionID, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	if section.Section != 7 {
		err = &azcore.ResponseError{
			StatusCode: http.StatusNotFound,
		}
		return section, err
	}

	return section, nil

}

func (env *env) GetSection7sByUser(ctx context.Context, userID string, paginationOptions PaginationOptions) ([]Section7, error) {

	env.logger.Info("Getting all Section 7 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section7{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	sortOrder := "ASC"
	if paginationOptions.SortByNewest {
		sortOrder = "DESC"
	}

	query := fmt.Sprintf("SELECT * FROM sections s WHERE s.user_id = @user_id AND s.section = @section ORDER BY s.created %s", sortOrder)

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@section", Value: 7},
		},
		PageSizeHint: int32(paginationOptions.PerPage),
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section7{}
	currentPage := 0

	for pager.More() {

		if currentPage == paginationOptions.Page {
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

			return sections, nil

		} else {
			_, err := pager.NextPage(ctx)
			if err != nil {
				return []Section7{}, err
			}
			currentPage++
		}

	}

	return sections, nil

}

func (env *env) UpsertSection7(ctx context.Context, section Section7) (Section7, error) {

	env.logger.Info("Upserting section 7")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return section, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return section, err
	}

	return section, nil

}

/*******************************
* SECTION 8
********************************/

type Section8 struct {
	ID                        string `json:"id"`
	Nickname                  string `json:"nickname"`
	IndividualGroupActivities string `json:"individual_group_activities"`
	HoursSpent                int    `json:"hours_spent"`
	NumPeopleReached          int    `json:"num_people_reached"`
	GenericSectionInfo
}

func (env *env) GetSection8ByID(ctx context.Context, userID string, sectionID string) (Section8, error) {

	env.logger.Info("Getting Section 8 by ID")
	section := Section8{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.ReadItem(ctx, partitionKey, sectionID, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	if section.Section != 8 {
		err = &azcore.ResponseError{
			StatusCode: http.StatusNotFound,
		}
		return section, err
	}

	return section, nil

}

func (env *env) GetSection8sByUser(ctx context.Context, userID string, paginationOptions PaginationOptions) ([]Section8, error) {

	env.logger.Info("Getting all Section 8 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section8{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	sortOrder := "ASC"
	if paginationOptions.SortByNewest {
		sortOrder = "DESC"
	}

	query := fmt.Sprintf("SELECT * FROM sections s WHERE s.user_id = @user_id AND s.section = @section ORDER BY s.created %s", sortOrder)

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@section", Value: 8},
		},
		PageSizeHint: int32(paginationOptions.PerPage),
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section8{}
	currentPage := 0

	for pager.More() {

		if currentPage == paginationOptions.Page {
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

			return sections, nil

		} else {
			_, err := pager.NextPage(ctx)
			if err != nil {
				return []Section8{}, err
			}
			currentPage++
		}

	}

	return sections, nil

}

func (env *env) UpsertSection8(ctx context.Context, section Section8) (Section8, error) {

	env.logger.Info("Upserting section 8")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return section, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return section, err
	}

	return section, nil

}

/*******************************
* SECTION 9
********************************/

type Section9 struct {
	ID                string `json:"id"`
	Nickname          string `json:"nickname"`
	CommunicationType string `json:"communication_type"`
	Topic             string `json:"topic"`
	TimesGiven        int    `json:"times_given"`
	Location          string `json:"location"`
	AudienceSize      int    `json:"audience_size"`
	GenericSectionInfo
}

func (env *env) GetSection9ByID(ctx context.Context, userID string, sectionID string) (Section9, error) {

	env.logger.Info("Getting Section 9 by ID")
	section := Section9{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.ReadItem(ctx, partitionKey, sectionID, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	if section.Section != 9 {
		err = &azcore.ResponseError{
			StatusCode: http.StatusNotFound,
		}
		return section, err
	}

	return section, nil

}

func (env *env) GetSection9sByUser(ctx context.Context, userID string, paginationOptions PaginationOptions) ([]Section9, error) {

	env.logger.Info("Getting all Section 9 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section9{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	sortOrder := "ASC"
	if paginationOptions.SortByNewest {
		sortOrder = "DESC"
	}

	query := fmt.Sprintf("SELECT * FROM sections s WHERE s.user_id = @user_id AND s.section = @section ORDER BY s.created %s", sortOrder)

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@section", Value: 9},
		},
		PageSizeHint: int32(paginationOptions.PerPage),
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section9{}
	currentPage := 0

	for pager.More() {

		if currentPage == paginationOptions.Page {
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

			return sections, nil

		} else {
			_, err := pager.NextPage(ctx)
			if err != nil {
				return []Section9{}, err
			}
			currentPage++
		}

	}

	return sections, nil

}

func (env *env) UpsertSection9(ctx context.Context, section Section9) (Section9, error) {

	env.logger.Info("Upserting section 9")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return section, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return section, err
	}

	return section, nil

}

/*******************************
* SECTION 10
********************************/

type Section10 struct {
	ID                string `json:"id"`
	Nickname          string `json:"nickname"`
	CommunicationType string `json:"communication_type"`
	Topic             string `json:"topic"`
	TimesGiven        int    `json:"times_given"`
	Location          string `json:"location"`
	AudienceSize      int    `json:"audience_size"`
	GenericSectionInfo
}

func (env *env) GetSection10ByID(ctx context.Context, userID string, sectionID string) (Section10, error) {

	env.logger.Info("Getting Section 10 by ID")
	section := Section10{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.ReadItem(ctx, partitionKey, sectionID, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	if section.Section != 10 {
		err = &azcore.ResponseError{
			StatusCode: http.StatusNotFound,
		}
		return section, err
	}

	return section, nil

}

func (env *env) GetSection10sByUser(ctx context.Context, userID string, paginationOptions PaginationOptions) ([]Section10, error) {

	env.logger.Info("Getting all Section 10 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section10{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	sortOrder := "ASC"
	if paginationOptions.SortByNewest {
		sortOrder = "DESC"
	}

	query := fmt.Sprintf("SELECT * FROM sections s WHERE s.user_id = @user_id AND s.section = @section ORDER BY s.created %s", sortOrder)

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@section", Value: 10},
		},
		PageSizeHint: int32(paginationOptions.PerPage),
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section10{}
	currentPage := 0

	for pager.More() {

		if currentPage == paginationOptions.Page {
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

			return sections, nil

		} else {
			_, err := pager.NextPage(ctx)
			if err != nil {
				return []Section10{}, err
			}
			currentPage++
		}

	}

	return sections, nil

}

func (env *env) UpsertSection10(ctx context.Context, section Section10) (Section10, error) {

	env.logger.Info("Upserting section 10")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return section, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return section, err
	}

	return section, nil

}

/*******************************
* SECTION 11
********************************/

type Section11 struct {
	ID                 string `json:"id"`
	Nickname           string `json:"nickname"`
	EventAndLevel      string `json:"event_and_level"`
	ExhibitsOrDivision string `json:"exhibits_or_division"`
	RibbonOrPlacings   string `json:"ribbon_or_placings"`
	GenericSectionInfo
}

func (env *env) GetSection11ByID(ctx context.Context, userID string, sectionID string) (Section11, error) {

	env.logger.Info("Getting Section 11 by ID")
	section := Section11{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.ReadItem(ctx, partitionKey, sectionID, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	if section.Section != 11 {
		err = &azcore.ResponseError{
			StatusCode: http.StatusNotFound,
		}
		return section, err
	}

	return section, nil

}

func (env *env) GetSection11sByUser(ctx context.Context, userID string, paginationOptions PaginationOptions) ([]Section11, error) {

	env.logger.Info("Getting all Section 11 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section11{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	sortOrder := "ASC"
	if paginationOptions.SortByNewest {
		sortOrder = "DESC"
	}

	query := fmt.Sprintf("SELECT * FROM sections s WHERE s.user_id = @user_id AND s.section = @section ORDER BY s.created %s", sortOrder)

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@section", Value: 11},
		},
		PageSizeHint: int32(paginationOptions.PerPage),
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section11{}
	currentPage := 0

	for pager.More() {

		if currentPage == paginationOptions.Page {
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

			return sections, nil

		} else {
			_, err := pager.NextPage(ctx)
			if err != nil {
				return []Section11{}, err
			}
			currentPage++
		}

	}

	return sections, nil

}

func (env *env) UpsertSection11(ctx context.Context, section Section11) (Section11, error) {

	env.logger.Info("Upserting section 11")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return section, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return section, err
	}

	return section, nil

}

/*******************************
* SECTION 12
********************************/

type Section12 struct {
	ID                  string `json:"id"`
	Nickname            string `json:"nickname"`
	ContestOrEvent      string `json:"contest_or_event"`
	RecognitionReceived string `json:"recognition_received"`
	Level               string `json:"level"`
	GenericSectionInfo
}

func (env *env) GetSection12ByID(ctx context.Context, userID string, sectionID string) (Section12, error) {

	env.logger.Info("Getting Section 12 by ID")
	section := Section12{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.ReadItem(ctx, partitionKey, sectionID, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	if section.Section != 12 {
		err = &azcore.ResponseError{
			StatusCode: http.StatusNotFound,
		}
		return section, err
	}

	return section, nil

}

func (env *env) GetSection12sByUser(ctx context.Context, userID string, paginationOptions PaginationOptions) ([]Section12, error) {

	env.logger.Info("Getting all Section 12 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section12{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	sortOrder := "ASC"
	if paginationOptions.SortByNewest {
		sortOrder = "DESC"
	}

	query := fmt.Sprintf("SELECT * FROM sections s WHERE s.user_id = @user_id AND s.section = @section ORDER BY s.created %s", sortOrder)

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@section", Value: 12},
		},
		PageSizeHint: int32(paginationOptions.PerPage),
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section12{}
	currentPage := 0

	for pager.More() {

		if currentPage == paginationOptions.Page {
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

			return sections, nil

		} else {
			_, err := pager.NextPage(ctx)
			if err != nil {
				return []Section12{}, err
			}
			currentPage++
		}

	}

	return sections, nil

}

func (env *env) UpsertSection12(ctx context.Context, section Section12) (Section12, error) {

	env.logger.Info("Upserting section 12")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return section, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return section, err
	}

	return section, nil

}

/*******************************
* SECTION 13
********************************/

type Section13 struct {
	ID              string `json:"id"`
	Nickname        string `json:"nickname"`
	RecognitionType string `json:"recognition_type"`
	GenericSectionInfo
}

func (env *env) GetSection13ByID(ctx context.Context, userID string, sectionID string) (Section13, error) {

	env.logger.Info("Getting Section 13 by ID")
	section := Section13{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.ReadItem(ctx, partitionKey, sectionID, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	if section.Section != 13 {
		err = &azcore.ResponseError{
			StatusCode: http.StatusNotFound,
		}
		return section, err
	}

	return section, nil

}

func (env *env) GetSection13sByUser(ctx context.Context, userID string, paginationOptions PaginationOptions) ([]Section13, error) {

	env.logger.Info("Getting all Section 13 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section13{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	sortOrder := "ASC"
	if paginationOptions.SortByNewest {
		sortOrder = "DESC"
	}

	query := fmt.Sprintf("SELECT * FROM sections s WHERE s.user_id = @user_id AND s.section = @section ORDER BY s.created %s", sortOrder)

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@section", Value: 13},
		},
		PageSizeHint: int32(paginationOptions.PerPage),
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section13{}
	currentPage := 0

	for pager.More() {

		if currentPage == paginationOptions.Page {
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

			return sections, nil

		} else {
			_, err := pager.NextPage(ctx)
			if err != nil {
				return []Section13{}, err
			}
			currentPage++
		}

	}

	return sections, nil

}

func (env *env) UpsertSection13(ctx context.Context, section Section13) (Section13, error) {

	env.logger.Info("Upserting section 13")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return section, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return section, err
	}

	return section, nil

}

/*******************************
* SECTION 14
********************************/

type Section14 struct {
	ID              string `json:"id"`
	Nickname        string `json:"nickname"`
	RecognitionType string `json:"recognition_type"`
	GenericSectionInfo
}

func (env *env) GetSection14ByID(ctx context.Context, userID string, sectionID string) (Section14, error) {

	env.logger.Info("Getting Section 14 by ID")
	section := Section14{}

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.ReadItem(ctx, partitionKey, sectionID, nil)
	if err != nil {
		return section, err
	}

	err = json.Unmarshal(response.Value, &section)
	if err != nil {
		return section, err
	}

	if section.Section != 14 {
		err = &azcore.ResponseError{
			StatusCode: http.StatusNotFound,
		}
		return section, err
	}

	return section, nil

}

func (env *env) GetSection14sByUser(ctx context.Context, userID string, paginationOptions PaginationOptions) ([]Section14, error) {

	env.logger.Info("Getting all Section 14 records")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return []Section14{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	sortOrder := "ASC"
	if paginationOptions.SortByNewest {
		sortOrder = "DESC"
	}

	query := fmt.Sprintf("SELECT * FROM sections s WHERE s.user_id = @user_id AND s.section = @section ORDER BY s.created %s", sortOrder)

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@section", Value: 14},
		},
		PageSizeHint: int32(paginationOptions.PerPage),
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	sections := []Section14{}
	currentPage := 0

	for pager.More() {

		if currentPage == paginationOptions.Page {
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

			return sections, nil

		} else {
			_, err := pager.NextPage(ctx)
			if err != nil {
				return []Section14{}, err
			}
			currentPage++
		}

	}

	return sections, nil

}

func (env *env) UpsertSection14(ctx context.Context, section Section14) (Section14, error) {

	env.logger.Info("Upserting section 14")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return section, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(section.UserID)

	marshalled, err := json.Marshal(section)
	if err != nil {
		return section, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return section, err
	}

	return section, nil

}

/*******************************
* DELETING
********************************/

func (env *env) RemoveSection(ctx context.Context, userID string, sectionID string) (interface{}, error) {

	env.logger.Info("Removing section")

	container, err := env.client.NewContainer("sections")
	if err != nil {
		return nil, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.DeleteItem(ctx, partitionKey, sectionID, nil)
	if err != nil {
		return nil, err
	}

	for _, dependent := range env.dependentsMap["sections"] {
		identifiables, err := dependent.GetRelated(ctx, userID, sectionID)
		if err != nil {
			return nil, err
		}
		for _, identifiable := range identifiables {
			_, err := dependent.Delete(ctx, userID, identifiable.GetID())
			if err != nil {
				return nil, err
			}
		}
	}

	return response, nil

}
