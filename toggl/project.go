package toggl

import (
	"fmt"

	"github.com/andreaskoch/togglapi/model"
	"github.com/pkg/errors"
)

// A Project describes a time reporting project (e.g. "Manhattan Project")
type Project struct {
	ID        int
	Name      string
	Client    Client
	Workspace Workspace
}

// A Projecter interface provides read/write access to Toggl projects.
type Projecter interface {
	// CreateProject creates a new project with the given name.
	// Returns an error of the creation failed.
	CreateProject(projectName, workspaceName, clientName string) (Project, error)

	// GetProjects returns all projects.
	GetProjects() ([]Project, error)

	// GetProjectByID returns the project for the given project id.
	// Returns an error if the project was not found.
	GetProjectByID(projectID int) (Project, error)

	// GetProjectByName returns the project for the given project name.
	// Returns an error if no matching project was found.
	GetProjectByName(projectName, workspaceName, clientName string) (Project, error)
}

// NewProjectRepository creates a new project repository instance.
func NewProjectRepository(projectAPI model.ProjectAPI, workspaceProvider Workspacer, clientProvider Clienter) Projecter {
	return &ProjectRepository{
		projectAPI: projectAPI,
		workspaces: workspaceProvider,
		clients:    clientProvider,
	}
}

// ProjectRepository provides read/write access to Toggl projects.
type ProjectRepository struct {
	projectAPI    model.ProjectAPI
	projectsCache []Project

	workspaces Workspacer
	clients    Clienter
}

// CreateProject creates a new project with the given name.
// Returns an error of the creation failed.
func (repository *ProjectRepository) CreateProject(projectName, workspaceName, clientName string) (Project, error) {

	workspace, workspaceError := repository.workspaces.GetWorkspaceByName(workspaceName)
	if workspaceError != nil {
		return Project{}, errors.Wrap(workspaceError, fmt.Sprintf("Failed to get workspace %q", workspaceName))
	}

	var client Client
	if clientName != "" {

		existingClient, existingClientError := repository.clients.GetClientByName(workspaceName, clientName)
		if existingClientError != nil {

			createdClient, createClientError := repository.clients.CreateClient(workspace.ID, clientName)

			if createClientError != nil {
				return Project{}, errors.Wrap(createClientError, fmt.Sprintf("Failed to create client: %s", clientName))
			}

			client = createdClient
		} else {
			client = existingClient
		}

	}

	createdProject, createClientError := repository.projectAPI.CreateProject(model.Project{
		Name:        projectName,
		WorkspaceID: workspace.ID,
		ClientID:    client.ID,
	})

	if createClientError != nil {
		return Project{}, createClientError
	}

	// reset the projects cache
	repository.projectsCache = nil

	return Project{
		ID:        createdProject.ID,
		Name:      createdProject.Name,
		Workspace: workspace,
		Client:    client,
	}, nil
}

// GetProjects returns all projects.
func (repository *ProjectRepository) GetProjects() ([]Project, error) {
	if repository.projectsCache != nil {
		return repository.projectsCache, nil
	}

	workspaces, workspacesError := repository.workspaces.GetWorkspaces()
	if workspacesError != nil {
		return nil, errors.Wrap(workspacesError, "Failed to retrieve workspaces")
	}

	var projects []Project
	for _, workspace := range workspaces {

		projectsByWorkspace, projectsByWorkspaceError := repository.projectAPI.GetProjects(workspace.ID)
		if projectsByWorkspaceError != nil {
			return nil, errors.Wrap(projectsByWorkspaceError, "Failed to get projects from Toggl")
		}

		for _, projectModel := range projectsByWorkspace {

			var client Client
			clientID := projectModel.ClientID
			if clientID != 0 {

				clientByID, clientByIDError := repository.clients.GetClientByID(clientID)
				if clientByIDError != nil {
					return nil, errors.Wrap(clientByIDError, fmt.Sprintf("Failed to get client %d", clientID))
				}

				client = clientByID
			}

			projects = append(projects, Project{
				ID:        projectModel.ID,
				Name:      projectModel.Name,
				Workspace: workspace,
				Client:    client,
			})
		}
	}

	// store in cache
	repository.projectsCache = projects

	return projects, nil
}

// GetProjectByID returns the project for the given project id.
// Returns an error if the project was not found.
func (repository *ProjectRepository) GetProjectByID(projectID int) (Project, error) {
	projects, projectsError := repository.GetProjects()
	if projectsError != nil {
		return Project{}, projectsError
	}

	for _, project := range projects {
		if project.ID == projectID {
			return project, nil
		}
	}

	return Project{}, fmt.Errorf("No project found with id %d", projectID)
}

// GetProjectByName returns the project for the given project name.
// Returns an error if no matching project was found.
func (repository *ProjectRepository) GetProjectByName(projectName, workspaceName, clientName string) (Project, error) {

	projects, projectsError := repository.GetProjects()
	if projectsError != nil {
		return Project{}, projectsError
	}

	for _, project := range projects {
		if project.Name == projectName && project.Workspace.Name == workspaceName && project.Client.Name == clientName {
			return project, nil
		}
	}

	return Project{}, fmt.Errorf("Project %q was not found (Workspace: %q, Client: %q)", projectName, workspaceName, clientName)
}
