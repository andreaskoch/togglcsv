package toggl

import (
	"fmt"
	"testing"

	"github.com/andreaskoch/togglapi/model"
)

type mockProjectAPI struct {
	createProject func(project model.Project) (model.Project, error)
	getProjects   func(workspaceID int) ([]model.Project, error)
}

func (projectAPI *mockProjectAPI) CreateProject(project model.Project) (model.Project, error) {
	return projectAPI.createProject(project)
}

func (projectAPI *mockProjectAPI) GetProjects(workspaceID int) ([]model.Project, error) {
	return projectAPI.getProjects(workspaceID)
}

type mockClienter struct {
	createClient    func(workspaceID int, name string) (Client, error)
	getClients      func() ([]Client, error)
	getClientByID   func(clientID int) (Client, error)
	getClientByName func(workspaceName, clientName string) (Client, error)
}

func (clienter *mockClienter) CreateClient(workspaceID int, name string) (Client, error) {
	return clienter.createClient(workspaceID, name)
}

func (clienter *mockClienter) GetClients() ([]Client, error) {
	return clienter.getClients()
}

func (clienter *mockClienter) GetClientByID(clientID int) (Client, error) {
	return clienter.getClientByID(clientID)
}

func (clienter *mockClienter) GetClientByName(workspaceName, clientName string) (Client, error) {
	return clienter.getClientByName(workspaceName, clientName)
}

func Test_CreateProject_CreateSucceeds_ProjectIsReturned(t *testing.T) {
	// arrange
	projectAPI := &mockProjectAPI{
		createProject: func(project model.Project) (model.Project, error) {
			return model.Project{Name: "Sample Project", WorkspaceID: 1}, nil
		},
	}

	projectRepository := ProjectRepository{
		projectAPI: projectAPI,
		workspaces: &mockWorkspacer{
			getWorkspaceByName: func(workspaceName string) (Workspace, error) {
				return Workspace{
					ID:   1,
					Name: workspaceName,
				}, nil
			},
		},
		clients: &mockClienter{
			getClientByName: func(workspaceName, clientName string) (Client, error) {
				return Client{
					ID:   1,
					Name: clientName,
				}, nil
			},
		},
	}

	// act
	project, err := projectRepository.CreateProject("Sample Project", "A workspace", "A client")

	// assert
	if err != nil {
		t.Logf("CreateProject should not return an error if the create succeeded: %s", err.Error())
	}

	if project.Name != "Sample Project" || project.Workspace.ID != 1 {
		t.Fail()
		t.Logf("CreateProject have returned a project model with the given parameters but returned this instead: %#v", project)
	}
}

func Test_CreateProject_NonExistingClient_ClientIsCreated(t *testing.T) {
	// arrange
	projectAPI := &mockProjectAPI{
		createProject: func(project model.Project) (model.Project, error) {
			return model.Project{Name: "Sample Project", WorkspaceID: 1}, nil
		},
	}

	clientIsCreated := false

	projectRepository := ProjectRepository{
		projectAPI: projectAPI,
		workspaces: &mockWorkspacer{
			getWorkspaceByName: func(workspaceName string) (Workspace, error) {
				return Workspace{
					ID:   1,
					Name: workspaceName,
				}, nil
			},
		},
		clients: &mockClienter{
			getClientByName: func(workspaceName, clientName string) (Client, error) {
				return Client{}, fmt.Errorf("Client does not exist")
			},
			createClient: func(workspaceID int, name string) (Client, error) {

				clientIsCreated = true

				return Client{
					ID:   13,
					Name: name,
				}, nil
			},
		},
	}

	// act
	project, _ := projectRepository.CreateProject("Sample Project", "A workspace", "A client")

	// assert
	if !clientIsCreated {
		t.Fail()
		t.Logf("CreateProject should create the client if it does not exist")
	}

	if project.Client.ID != 13 {
		t.Fail()
		t.Logf("CreateProject should use the created client model when the client did not exist before")
	}
}

func Test_CreateProject_CreateFails_ErrorIsReturned(t *testing.T) {
	// arrange
	projectAPI := &mockProjectAPI{
		createProject: func(project model.Project) (model.Project, error) {
			return model.Project{}, fmt.Errorf("Some error")
		},
	}

	projectRepository := ProjectRepository{
		projectAPI: projectAPI,
		workspaces: &mockWorkspacer{
			getWorkspaceByName: func(workspaceName string) (Workspace, error) {
				return Workspace{
					ID:   1,
					Name: workspaceName,
				}, nil
			},
		},
		clients: &mockClienter{
			getClientByName: func(workspaceName, clientName string) (Client, error) {
				return Client{
					ID:   1,
					Name: clientName,
				}, nil
			},
		},
	}

	// act
	_, err := projectRepository.CreateProject("Sample Project", "A workspace", "A client")

	// assert
	if err == nil {
		t.Fail()
		t.Logf("CreateProject return an error if the API returns one")
	}
}

func Test_GetProjects_APIReturnsError_ErrorIsReturned(t *testing.T) {
	// arrange
	projectAPI := &mockProjectAPI{
		getProjects: func(workspaceID int) ([]model.Project, error) {
			return nil, fmt.Errorf("Some error")
		},
	}

	workspaceProvider := &mockWorkspacer{
		getWorkspaces: func() ([]Workspace, error) {
			return []Workspace{
				Workspace{},
			}, nil
		},
	}

	projectRepository := ProjectRepository{
		projectAPI: projectAPI,
		workspaces: workspaceProvider,
	}

	// act
	_, err := projectRepository.GetProjects()

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetProjects should return an error if the API returned one")
	}
}

func Test_GetProjects_AllProjectsReturnedByTheAPIAreReturned(t *testing.T) {
	// arrange
	projectAPI := &mockProjectAPI{
		getProjects: func(workspaceID int) ([]model.Project, error) {
			return []model.Project{
				model.Project{},
				model.Project{},
			}, nil
		},
	}

	workspaceProvider := &mockWorkspacer{
		getWorkspaces: func() ([]Workspace, error) {
			return []Workspace{
				Workspace{},
			}, nil
		},
	}

	projectRepository := ProjectRepository{
		projectAPI: projectAPI,
		workspaces: workspaceProvider,
	}

	// act
	projects, err := projectRepository.GetProjects()

	// assert
	if err != nil {
		t.Logf("GetProjects should not return an error if the API does not return one: %s", err.Error())
	}

	if len(projects) != 2 {
		t.Fail()
		t.Logf("GetProjects should have returned two projects")
	}
}

func Test_GetProjects_SecondCallUsesTheCache(t *testing.T) {
	// arrange
	numberOfCalls := 0
	projectAPI := &mockProjectAPI{
		getProjects: func(workspaceID int) ([]model.Project, error) {
			numberOfCalls++

			return []model.Project{
				model.Project{},
				model.Project{},
			}, nil
		},
	}

	workspaceProvider := &mockWorkspacer{
		getWorkspaces: func() ([]Workspace, error) {
			return []Workspace{
				Workspace{},
			}, nil
		},
	}

	projectRepository := ProjectRepository{
		projectAPI: projectAPI,
		workspaces: workspaceProvider,
	}

	// act
	projectRepository.GetProjects() // first call
	projectRepository.GetProjects() // second call

	// assert
	if numberOfCalls != 1 {
		t.Fail()
		t.Logf("The seconds call to GetProjects should not have called the API again")
	}
}

func Test_GetProjectByName_NoProjectsAvailabe_ErrorIsReturned(t *testing.T) {
	// arrange
	projectAPI := &mockProjectAPI{
		getProjects: func(workspaceID int) ([]model.Project, error) {
			return []model.Project{}, nil
		},
	}

	workspaceProvider := &mockWorkspacer{
		getWorkspaces: func() ([]Workspace, error) {
			return []Workspace{
				Workspace{},
			}, nil
		},
	}

	projectRepository := ProjectRepository{
		projectAPI: projectAPI,
		workspaces: workspaceProvider,
	}

	// act
	_, err := projectRepository.GetProjectByName("Project A", "A Workspace", "Client A")

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetProjectByName should return an error if no matching project was found")
	}
}

func Test_GetProjectByName_WorkspaceNameDoesNotMatch_ErrorIsReturned(t *testing.T) {
	// arrange
	projectAPI := &mockProjectAPI{
		getProjects: func(workspaceID int) ([]model.Project, error) {
			return []model.Project{
				model.Project{
					WorkspaceID: 2,
					Name:        "Project A",
				},
				model.Project{
					WorkspaceID: 2,
					Name:        "Project B",
				},
			}, nil
		},
	}

	workspaceProvider := &mockWorkspacer{
		getWorkspaces: func() ([]Workspace, error) {
			return []Workspace{
				Workspace{ID: 1, Name: "Workspace A"},
				Workspace{ID: 2, Name: "Workspace B"},
			}, nil
		},
	}

	projectRepository := ProjectRepository{
		projectAPI: projectAPI,
		workspaces: workspaceProvider,
	}

	// act
	_, err := projectRepository.GetProjectByName("Project A", "A Workspace", "Client A")

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetProjectByName should return an error if no matching project was found")
	}
}

func Test_GetProjectByName_ProjectNameDoesNotMatch_ErrorIsReturned(t *testing.T) {
	// arrange
	projectAPI := &mockProjectAPI{
		getProjects: func(workspaceID int) ([]model.Project, error) {
			return []model.Project{
				model.Project{
					WorkspaceID: 1,
					Name:        "Project A",
				},
				model.Project{
					WorkspaceID: 1,
					Name:        "Project B",
				},
			}, nil
		},
	}

	workspaceProvider := &mockWorkspacer{
		getWorkspaces: func() ([]Workspace, error) {
			return []Workspace{
				Workspace{ID: 1},
			}, nil
		},
	}

	projectRepository := ProjectRepository{
		projectAPI: projectAPI,
		workspaces: workspaceProvider,
	}

	// act
	_, err := projectRepository.GetProjectByName("Project A", "A Workspace", "Client A")

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetProjectByName should return an error if no matching project was found")
	}
}

func Test_GetProjectByName_WorkspaceNameMatches_ProjectIsReturned(t *testing.T) {
	// arrange
	projectAPI := &mockProjectAPI{
		getProjects: func(workspaceID int) ([]model.Project, error) {
			return []model.Project{
				model.Project{
					WorkspaceID: 1,
					ClientID:    1,
					Name:        "Project A",
				},
				model.Project{
					WorkspaceID: 1,
					ClientID:    1,
					Name:        "Project B",
				},
			}, nil
		},
	}

	workspaceProvider := &mockWorkspacer{
		getWorkspaces: func() ([]Workspace, error) {
			return []Workspace{
				Workspace{ID: 1, Name: "A Workspace"},
			}, nil
		},
	}

	clientProvider := &mockClienter{
		getClientByID: func(clientID int) (Client, error) {
			return Client{ID: clientID, Name: "Client A"}, nil
		},
	}

	projectRepository := ProjectRepository{
		projectAPI: projectAPI,
		clients:    clientProvider,
		workspaces: workspaceProvider,
	}

	// act
	project, err := projectRepository.GetProjectByName("Project A", "A Workspace", "Client A")

	// assert
	if err != nil {
		t.Logf("GetProjectByName should not have returned an error but returned this: %s", err)
	}

	if project.Name != "Project A" || project.Workspace.ID != 1 {
		t.Fail()
		t.Logf("GetProjectByName should have returned a matching project")
	}
}

func Test_GetProjectByID_NoProjectsAvailabe_ErrorIsReturned(t *testing.T) {
	// arrange
	projectAPI := &mockProjectAPI{
		getProjects: func(workspaceID int) ([]model.Project, error) {
			return []model.Project{}, nil
		},
	}

	workspaceProvider := &mockWorkspacer{
		getWorkspaces: func() ([]Workspace, error) {
			return []Workspace{
				Workspace{ID: 1},
			}, nil
		},
	}

	projectRepository := ProjectRepository{
		projectAPI: projectAPI,
		workspaces: workspaceProvider,
	}

	// act
	_, err := projectRepository.GetProjectByID(1)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetProjectByID should return an error if no matching project was found")
	}
}

func Test_GetProjectByID_WorkspacerReturnsError_ErrorIsReturned(t *testing.T) {
	// arrange
	projectAPI := &mockProjectAPI{
		getProjects: func(workspaceID int) ([]model.Project, error) {
			return []model.Project{}, nil
		},
	}

	workspaceProvider := &mockWorkspacer{
		getWorkspaces: func() ([]Workspace, error) {
			return nil, fmt.Errorf("Error")
		},
	}

	projectRepository := ProjectRepository{
		projectAPI: projectAPI,
		workspaces: workspaceProvider,
	}

	// act
	_, err := projectRepository.GetProjectByID(1)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetProjectByID should return an error if no matching project was found")
	}
}

func Test_GetProjectByID_ProjectNameDoesNotMatch_ErrorIsReturned(t *testing.T) {
	// arrange
	projectAPI := &mockProjectAPI{
		getProjects: func(workspaceID int) ([]model.Project, error) {
			return []model.Project{
				model.Project{
					WorkspaceID: 1,
					ID:          1,
				},
				model.Project{
					WorkspaceID: 1,
					ID:          2,
				},
			}, nil
		},
	}

	workspaceProvider := &mockWorkspacer{
		getWorkspaces: func() ([]Workspace, error) {
			return []Workspace{
				Workspace{ID: 1},
			}, nil
		},
	}

	projectRepository := ProjectRepository{
		projectAPI: projectAPI,
		workspaces: workspaceProvider,
	}

	// act
	_, err := projectRepository.GetProjectByID(3)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetProjectByID should return an error if no matching project was found")
	}
}

func Test_GetProjectByID_ProjectIDIsFound_ProjectIsReturned(t *testing.T) {
	// arrange
	projectAPI := &mockProjectAPI{
		getProjects: func(workspaceID int) ([]model.Project, error) {
			return []model.Project{
				model.Project{
					WorkspaceID: 1,
					ID:          1,
				},
				model.Project{
					WorkspaceID: 1,
					ID:          2,
				},
			}, nil
		},
	}

	workspaceProvider := &mockWorkspacer{
		getWorkspaces: func() ([]Workspace, error) {
			return []Workspace{
				Workspace{ID: 1, Name: "A workspace"},
			}, nil
		},
	}

	projectRepository := ProjectRepository{
		projectAPI: projectAPI,
		workspaces: workspaceProvider,
	}

	// act
	project, err := projectRepository.GetProjectByID(1)

	// assert
	if err != nil {
		t.Logf("GetProjectByID should not have returned an error but returned this: %s", err)
	}

	if project.ID != 1 || project.Workspace.ID != 1 {
		t.Fail()
		t.Logf("GetProjectByID should have returned a matching project")
	}
}
