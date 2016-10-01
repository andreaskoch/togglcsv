package toggl

import (
	"fmt"

	"github.com/andreaskoch/togglapi/model"
	"github.com/pkg/errors"
)

type modelConverter interface {
	// ConvertTimeEntryToTimeRecord converts the given TimeEntry model into a TimeRecord model.
	// Returns an error if the given TimeEntry could not be converted.
	ConvertTimeEntryToTimeRecord(timeEntry model.TimeEntry) (TimeRecord, error)

	// ConvertTimeRecordToTimeEntry converts the given TimeRecord model into an TimeEntry model.
	// Returns an error if the given TimeEntry could not be converted.
	ConvertTimeRecordToTimeEntry(timeRecord TimeRecord) (model.TimeEntry, error)
}

// togglModelConverter converts TimeEntry models into TimeRecord models and vice versa
type togglModelConverter struct {
	workspaces Workspacer
	projects   Projecter
	clients    Clienter
}

// ConvertTimeRecordToTimeEntry converts the given TimeRecord model into an TimeEntry model.
// Returns an error if the given TimeEntry could not be converted.
func (converter *togglModelConverter) ConvertTimeRecordToTimeEntry(timeRecord TimeRecord) (model.TimeEntry, error) {

	// lookup the workspace
	workspace, workspaceError := converter.workspaces.GetWorkspaceByName(timeRecord.WorkspaceName)
	if workspaceError != nil {
		return model.TimeEntry{}, errors.Wrap(workspaceError, "Cannot convert time record to time entry.")
	}

	// lookup the project
	project, projectError := converter.projects.GetProjectByName(timeRecord.ProjectName, timeRecord.WorkspaceName, timeRecord.ClientName)
	if projectError != nil {
		return model.TimeEntry{}, errors.Wrap(projectError, "Cannot convert time record to time entry.")
	}

	// create the time entry
	timeEntryModel := model.TimeEntry{
		Wid:         workspace.ID,
		Pid:         project.ID,
		Start:       timeRecord.Start,
		Stop:        timeRecord.Stop,
		Description: timeRecord.Description,
		Tags:        timeRecord.Tags,
	}

	return timeEntryModel, nil
}

// ConvertTimeEntryToTimeRecord converts the given TimeEntry model into a TimeRecord model.
// Returns an error if the given TimeEntry could not be converted.
func (converter *togglModelConverter) ConvertTimeEntryToTimeRecord(timeEntry model.TimeEntry) (TimeRecord, error) {

	// workspace
	workspace, workspaceError := converter.workspaces.GetWorkspaceByID(timeEntry.Wid)
	if workspaceError != nil {
		return TimeRecord{}, fmt.Errorf("No workspace found with ID %d", timeEntry.Wid)
	}

	// project + client (optional)
	var project Project
	var client Client

	projectID := timeEntry.Pid
	if projectID != 0 {

		// load the project
		projectByID, projectByIDError := converter.projects.GetProjectByID(projectID)
		if projectByIDError != nil {
			return TimeRecord{}, errors.Wrap(projectByIDError, fmt.Sprintf("No project found with ID %d", projectID))
		}
		project = projectByID

		// load the client id
		clientID := project.Client.ID
		if clientID != 0 {

			clientByID, clientByIDError := converter.clients.GetClientByID(clientID)
			if clientByIDError != nil {
				return TimeRecord{}, errors.Wrap(clientByIDError, fmt.Sprintf("No client found with ID %d", clientID))
			}

			client = clientByID
		}
	}

	record := TimeRecord{
		WorkspaceName: workspace.Name,
		ProjectName:   project.Name,
		ClientName:    client.Name,

		Start:       timeEntry.Start,
		Stop:        timeEntry.Stop,
		Tags:        timeEntry.Tags,
		Description: timeEntry.Description,
	}

	return record, nil

}
