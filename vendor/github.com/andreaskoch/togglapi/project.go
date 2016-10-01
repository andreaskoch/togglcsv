package togglapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andreaskoch/togglapi/model"
	"github.com/pkg/errors"
)

// NewProjectAPI create a new client for the Toggl project API.
func NewProjectAPI(baseURL, token string) model.ProjectAPI {
	return &ProjectAPI{
		restClient: &togglRESTAPIClient{
			baseURL: baseURL,
			token:   token,
		},
	}
}

// ProjectAPI provides functions for interacting with Toggls' project API.
type ProjectAPI struct {
	restClient RESTRequester
}

// CreateProject creates a new project.
func (repository *ProjectAPI) CreateProject(project model.Project) (model.Project, error) {

	projectRequest := struct {
		Project model.Project `json:"project"`
	}{
		Project: project,
	}

	jsonBody, marshalError := json.Marshal(projectRequest)
	if marshalError != nil {
		return model.Project{}, errors.Wrap(marshalError, "Failed to serialize the project")
	}

	content, err := repository.restClient.Request(http.MethodPost, "projects", bytes.NewBuffer(jsonBody))
	if err != nil {
		return model.Project{}, errors.Wrap(err, "Failed to create project")
	}

	var projectResponse struct {
		Project model.Project `json:"data"`
	}

	if unmarshalError := json.Unmarshal(content, &projectResponse); unmarshalError != nil {
		return model.Project{}, errors.Wrap(unmarshalError, "Failed to deserialize the created project")
	}

	return projectResponse.Project, nil
}

// GetProjects returns all projects for the given workspace.
func (repository *ProjectAPI) GetProjects(workspaceID int) ([]model.Project, error) {

	route := fmt.Sprintf(
		"workspaces/%d/projects",
		workspaceID,
	)

	content, err := repository.restClient.Request(http.MethodGet, route, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve projects")
	}

	var projects []model.Project
	if unmarshalError := json.Unmarshal(content, &projects); unmarshalError != nil {
		return nil, errors.Wrap(unmarshalError, "Failed to deserialize the projects")
	}

	return projects, nil
}
