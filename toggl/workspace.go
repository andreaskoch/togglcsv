package toggl

import (
	"fmt"

	"github.com/andreaskoch/togglapi/model"
)

// A Workspace defines a grouping for time reporting entries (e.g. "Company XY")
type Workspace struct {
	ID   int
	Name string
}

// A Workspacer interface provides read/write access to Toggl workspaces.
type Workspacer interface {
	// CreateWorkspace creates a new Toggl workspace.
	CreateWorkspace(name string) (Workspace, error)

	// GetWorkspaces returns all available workspaces.
	GetWorkspaces() ([]Workspace, error)

	// GetWorkspaceByID returns the workspace for the given workspace id.
	// Returns an error if the workspace was not found.
	GetWorkspaceByID(workspaceID int) (Workspace, error)

	// GetWorkspaceByName returns the workspace for the given workspace name.
	// Returns an error if no matching workspace was found.
	GetWorkspaceByName(workspaceName string) (Workspace, error)
}

// NewWorkspaceRepository creates a new workspace provider instance.
func NewWorkspaceRepository(workspaceAPI model.WorkspaceAPI) Workspacer {
	return &WorkspaceRepository{
		workspaceAPI: workspaceAPI,
	}
}

// WorkspaceRepository provides read access to the Toggl workspaces.
// Write is unfortunately not supported by the Toggl API.
type WorkspaceRepository struct {
	workspaceAPI    model.WorkspaceAPI
	workspacesCache []Workspace
}

// CreateWorkspace creates a new Toggl workspace.
func (repository *WorkspaceRepository) CreateWorkspace(name string) (Workspace, error) {
	return Workspace{}, fmt.Errorf("Creating workspaces is unfortunately not supported by the Toggl API. You must create the worspace %q from the Toggl website.", name)
}

// GetWorkspaces returns all available workspaces.
func (repository *WorkspaceRepository) GetWorkspaces() ([]Workspace, error) {
	if repository.workspacesCache != nil {
		return repository.workspacesCache, nil
	}

	workspaces, workspacesError := repository.workspaceAPI.GetWorkspaces()
	if workspacesError != nil {
		return nil, fmt.Errorf("Failed to get workspaces from Toggl: %s", workspacesError.Error())
	}

	var workspaceModels []Workspace
	for _, workspace := range workspaces {
		workspaceModels = append(workspaceModels, Workspace{
			ID:   workspace.ID,
			Name: workspace.Name,
		})
	}

	// store in cache
	repository.workspacesCache = workspaceModels

	return workspaceModels, nil
}

// GetWorkspaceByID returns the workspace for the given workspace id.
// Returns an error if the workspace was not found.
func (repository *WorkspaceRepository) GetWorkspaceByID(workspaceID int) (Workspace, error) {
	workspaces, workspacesError := repository.GetWorkspaces()
	if workspacesError != nil {
		return Workspace{}, workspacesError
	}

	for _, workspace := range workspaces {
		if workspace.ID == workspaceID {
			return workspace, nil
		}
	}

	return Workspace{}, fmt.Errorf("unknown workspace (%d)", workspaceID)
}

// GetWorkspaceByName returns the workspace for the given workspace name.
// Returns an error if no matching workspace was found.
func (repository *WorkspaceRepository) GetWorkspaceByName(workspaceName string) (Workspace, error) {

	workspaces, workspacesError := repository.GetWorkspaces()
	if workspacesError != nil {
		return Workspace{}, workspacesError
	}

	for _, workspace := range workspaces {
		if workspace.Name == workspaceName {
			return workspace, nil
		}
	}

	return Workspace{}, fmt.Errorf("Workspace %q was not found", workspaceName)
}
