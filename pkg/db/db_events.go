package db

import (
	"context"
	"encoding/json"

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

/*******************************
* FULL EVENTS
********************************/

func (env *env) GetEventsByUser(ctx context.Context, userID string) ([]Event, error) {

	env.logger.Info("Getting events")

	container, err := env.client.NewContainer("events")
	if err != nil {
		return []Event{}, err
	}

	partitionKey := azcosmos.NewPartitionKeyString(userID)

	query := "SELECT * FROM events e WHERE e.user_id = @user_id"

	queryOptions := azcosmos.QueryOptions{
		QueryParameters: []azcosmos.QueryParameter{
			{Name: "@user_id", Value: userID},
		},
	}

	pager := container.NewQueryItemsPager(query, partitionKey, &queryOptions)

	events := []Event{}

	for pager.More() {
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
