package toggl

import (
	"fmt"
	"time"

	"github.com/andreaskoch/togglapi/model"
	"github.com/pkg/errors"
)

// TimeRecord represents a single time tracking record
type TimeRecord struct {
	WorkspaceName string
	ProjectName   string
	ClientName    string

	Start       time.Time
	Stop        time.Time
	Description string
	Tags        []string
}

// A TimeRecorder interface provides functions for reading and writing time records.
type TimeRecorder interface {
	// CreateTimeRecord creates a new time record.
	// Returns an error if the creation failed.
	CreateTimeRecord(timeRecord TimeRecord) error

	// GetTimeRecords returns all time records from the given start date until the given stop date.
	// Returns an error of the time records could not be retrieved.
	GetTimeRecords(start, stop time.Time) ([]TimeRecord, error)
}

// NewTimeRecordRepository creates a new time record repository instance.
func NewTimeRecordRepository(
	timeEntryAPI model.TimeEntryAPI,
	workspaceRepository Workspacer,
	projectRepository Projecter,
	clientRepository Clienter) TimeRecorder {

	return &TimeRecordRepository{
		timeEntryAPI:      timeEntryAPI,
		timeRangeProvider: fullMonthTimeRangeProvider{},

		workspaces: workspaceRepository,
		projects:   projectRepository,

		modelConverter: &togglModelConverter{
			workspaces: workspaceRepository,
			projects:   projectRepository,
			clients:    clientRepository,
		},
	}
}

// TimeRecordRepository provides read/write access to the Toggl time records.
type TimeRecordRepository struct {
	timeEntryAPI      model.TimeEntryAPI
	timeRangeProvider timeRangeProvider

	workspaces Workspacer
	projects   Projecter

	modelConverter modelConverter
}

// CreateTimeRecord creates a new time record.
// Returns an error if the creation failed.
func (repository *TimeRecordRepository) CreateTimeRecord(timeRecord TimeRecord) error {

	// create the project if it does not exist
	if _, projectError := repository.projects.GetProjectByName(timeRecord.ProjectName, timeRecord.WorkspaceName, timeRecord.ClientName); projectError != nil {

		if _, createProjectError := repository.projects.CreateProject(timeRecord.ProjectName, timeRecord.WorkspaceName, timeRecord.ClientName); createProjectError != nil {
			return errors.Wrap(createProjectError, fmt.Sprintf("Failed to create project for time record: %#v", timeRecord))
		}
	}

	timeEntry, conversionError := repository.modelConverter.ConvertTimeRecordToTimeEntry(timeRecord)
	if conversionError != nil {
		return errors.Wrap(conversionError, fmt.Sprintf("Failed to convert the given time record (%#v) into a valid time entry", timeRecord))
	}

	_, createError := repository.timeEntryAPI.CreateTimeEntry(timeEntry)
	if createError != nil {
		return errors.Wrap(createError, fmt.Sprintf("Failed to create time record (%v)", timeRecord))
	}

	return nil
}

// GetTimeRecords returns all time records from the given start date until the given stop date.
// Returns an error of the time records could not be retrieved.
func (repository *TimeRecordRepository) GetTimeRecords(start, stop time.Time) ([]TimeRecord, error) {

	if start.After(stop) {
		return nil, fmt.Errorf("The start date cannot be before the stop date")
	}

	var timeRecords []TimeRecord

	// fetch the records in monthly chunks because Toggl has a max length of 1000 entries per request
	rangesSinceStartDate, timeRangeError := repository.timeRangeProvider.GetTimeRanges(start, stop)
	if timeRangeError != nil {
		return nil, timeRangeError
	}

	for _, monthRange := range rangesSinceStartDate {

		// get the records for the current chunk
		records, timeRecordsError := repository.getTimeRecords(monthRange.Start(), monthRange.Stop())
		if timeRecordsError != nil {
			return nil, timeRecordsError
		}

		timeRecords = append(timeRecords, records...)

	}

	return timeRecords, nil
}

// getTimeRecords returns all time records from the given start date until the given stop date.
// Returns an error of the time records could not be retrieved.
//
// Note:
// Don't choose a too large duration because Toggl will only return 1000 entries for each request.
func (repository *TimeRecordRepository) getTimeRecords(start, stop time.Time) ([]TimeRecord, error) {
	timeEntries, err := repository.timeEntryAPI.GetTimeEntries(start, stop)
	if err != nil {
		return nil, err
	}

	var records []TimeRecord
	for _, timeEntry := range timeEntries {

		// skip entries without project
		if noProjectIDSet := timeEntry.Pid == 0; noProjectIDSet {
			continue
		}

		// skip running entries
		if isRunningEntry := timeEntry.Stop.Equal(time.Time{}); isRunningEntry {
			continue
		}

		timeRecord, conversionError := repository.modelConverter.ConvertTimeEntryToTimeRecord(timeEntry)
		if conversionError != nil {
			return nil, errors.Wrap(conversionError, fmt.Sprintf("Failed to convert time entry (%#v)", timeEntry))
		}

		records = append(records, timeRecord)
	}

	return records, nil
}
