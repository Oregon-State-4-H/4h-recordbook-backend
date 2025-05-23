package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type Event struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Location    string `json:"location"`
	Description string `json:"description"`
	UserID      string `json:"user_id"`
	GenericDatabaseInfo
}

type EventSection struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	EventID       string `json:"event_id"`
	SectionNumber int    `json:"section_number"`
	SectionID     string `json:"section_id"`
	GenericDatabaseInfo
}

func (es EventSection) GetID() string {
	return es.ID
}

/*******************************
* FULL EVENTS
********************************/

func (env *env) GetEventsByUser(ctx context.Context, userID string, paginationOptions PaginationOptions) ([]Event, error) {

	env.logger.Info("Getting events")

	container, err := env.client.NewContainer("events")
	if err != nil {
		return []Event{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	sortOrder := "ASC"
	if paginationOptions.SortByNewest {
		sortOrder = "DESC"
	}

	query := fmt.Sprintf("SELECT * FROM events e WHERE e.user_id = @user_id ORDER BY e.created %s", sortOrder)

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
		},
		PageSizeHint: int32(paginationOptions.PerPage),
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	events := []Event{}
	currentPage := 0

	for pager.More() {

		if currentPage == paginationOptions.Page {
			response, err := pager.NextPage(ctx)
			if err != nil {
				return []Event{}, err
			}

			for _, bytes := range response.Items {
				event := Event{}
				err := json.Unmarshal(bytes, &event)
				if err != nil {
					return []Event{}, err
				}
				events = append(events, event)
			}

			return events, nil

		} else {
			_, err := pager.NextPage(ctx)
			if err != nil {
				return []Event{}, err
			}
			currentPage++
		}

	}

	return events, nil

}

func (env *env) GetEventByID(ctx context.Context, userID string, eventID string) (Event, error) {

	env.logger.Info("Getting event by ID")
	event := Event{}

	container, err := env.client.NewContainer("events")
	if err != nil {
		return event, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.ReadItem(ctx, partitionKey, eventID, nil)
	if err != nil {
		return event, err
	}

	err = json.Unmarshal(response.Value, &event)
	if err != nil {
		return event, err
	}

	return event, nil

}

func (env *env) UpsertEvent(ctx context.Context, event Event) (Event, error) {

	env.logger.Info("Upserting event")

	container, err := env.client.NewContainer("events")
	if err != nil {
		return event, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(event.UserID)

	marshalled, err := json.Marshal(event)
	if err != nil {
		return event, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return event, err
	}

	return event, nil

}

func (env *env) RemoveEvent(ctx context.Context, userID string, eventID string) (interface{}, error) {

	env.logger.Info("Removing event")

	container, err := env.client.NewContainer("events")
	if err != nil {
		return nil, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.DeleteItem(ctx, partitionKey, eventID, nil)
	if err != nil {
		return nil, err
	}

	for _, dependent := range env.dependentsMap["events"] {
		identifiables, err := dependent.GetRelated(ctx, userID, eventID)
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

/*******************************
* EVENT SECTIONS
********************************/

func (env *env) GetEventSectionByIDs(ctx context.Context, userID string, eventID string, sectionID string) (EventSection, error) {

	env.logger.Info("Getting event section")

	container, err := env.client.NewContainer("eventsections")
	if err != nil {
		return EventSection{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	query := "SELECT * FROM eventsections es WHERE es.user_id = @user_id AND es.event_id = @event_id AND es.section_id = @section_id"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@event_id", Value: eventID},
			{Name: "@section_id", Value: sectionID},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	eventSections := []EventSection{}

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return EventSection{}, err
		}

		for _, bytes := range response.Items {
			eventSection := EventSection{}
			err := json.Unmarshal(bytes, &eventSection)
			if err != nil {
				return EventSection{}, err
			}
			eventSections = append(eventSections, eventSection)
		}
	}

	if len(eventSections) == 0 {
		return EventSection{}, err
	}

	return eventSections[0], nil

}

func (env *env) GetEventSectionsByEvent(ctx context.Context, userID string, eventID string) ([]EventSection, error) {

	env.logger.Info("Getting sections by event")

	container, err := env.client.NewContainer("eventsections")
	if err != nil {
		return []EventSection{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	query := "SELECT * FROM eventsections es WHERE es.user_id = @user_id AND es.event_id = @event_id"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@event_id", Value: eventID},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	eventSections := []EventSection{}

	for pager.More() {
		response, err := pager.NextPage(ctx)
		if err != nil {
			return []EventSection{}, err
		}

		for _, bytes := range response.Items {
			eventSection := EventSection{}
			err := json.Unmarshal(bytes, &eventSection)
			if err != nil {
				return []EventSection{}, err
			}
			eventSections = append(eventSections, eventSection)
		}
	}

	return eventSections, nil

}

func (env *env) GetEventDependentEventSections(ctx context.Context, userID string, eventID string) ([]Identifiable, error) {

	env.logger.Info("Getting event dependent event sections")

	container, err := env.client.NewContainer("eventsections")
	if err != nil {
		return []Identifiable{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	query := "SELECT * FROM eventsections es WHERE es.user_id = @user_id AND es.event_id = @event_id"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@event_id", Value: eventID},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	eventSections := []EventSection{}

	for pager.More() {

		response, err := pager.NextPage(ctx)
		if err != nil {
			return []Identifiable{}, err
		}

		for _, bytes := range response.Items {
			eventSection := EventSection{}
			err := json.Unmarshal(bytes, &eventSection)
			if err != nil {
				return []Identifiable{}, err
			}
			eventSections = append(eventSections, eventSection)
		}

	}

	identifiables := []Identifiable{}

	for _, es := range eventSections {
		identifiables = append(identifiables, es)
	}

	return identifiables, nil

}

func (env *env) GetSectionDependentEventSections(ctx context.Context, userID string, sectionID string) ([]Identifiable, error) {

	env.logger.Info("Getting section dependent event sections")

	container, err := env.client.NewContainer("eventsections")
	if err != nil {
		return []Identifiable{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	query := "SELECT * FROM eventsections es WHERE es.user_id = @user_id AND es.section_id = @section_id"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
			{Name: "@section_id", Value: sectionID},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	eventSections := []EventSection{}

	for pager.More() {

		response, err := pager.NextPage(ctx)
		if err != nil {
			return []Identifiable{}, err
		}

		for _, bytes := range response.Items {
			eventSection := EventSection{}
			err := json.Unmarshal(bytes, &eventSection)
			if err != nil {
				return []Identifiable{}, err
			}
			eventSections = append(eventSections, eventSection)
		}

	}

	identifiables := []Identifiable{}

	for _, es := range eventSections {
		identifiables = append(identifiables, es)
	}

	return identifiables, nil

}

func (env *env) UpsertEventSection(ctx context.Context, eventSection EventSection) (EventSection, error) {

	env.logger.Info("Upserting event section")

	container, err := env.client.NewContainer("eventsections")
	if err != nil {
		return eventSection, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(eventSection.UserID)

	marshalled, err := json.Marshal(eventSection)
	if err != nil {
		return eventSection, err
	}

	_, err = container.UpsertItem(ctx, partitionKey, marshalled, nil)
	if err != nil {
		return eventSection, err
	}

	return eventSection, nil

}

func (env *env) RemoveEventSection(ctx context.Context, userID string, eventSectionID string) (interface{}, error) {

	env.logger.Info("Removing event section")

	container, err := env.client.NewContainer("eventsections")
	if err != nil {
		return nil, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	response, err := container.DeleteItem(ctx, partitionKey, eventSectionID, nil)
	if err != nil {
		return nil, err
	}

	return response, nil

}
