package toggl

import (
	"fmt"
	"testing"

	"github.com/andreaskoch/togglapi/model"
)

type mockWorkspaceAPI struct {
	createWorkspace          func(workspace model.Workspace) (model.Workspace, error)
	getWorkspaces            func() ([]model.Workspace, error)
	getWorkspacesByWorkspace func(workspaceID int) ([]model.Workspace, error)
}

func (workspaceAPI *mockWorkspaceAPI) CreateWorkspace(workspace model.Workspace) (model.Workspace, error) {
	return workspaceAPI.createWorkspace(workspace)
}

func (workspaceAPI *mockWorkspaceAPI) GetWorkspaces() ([]model.Workspace, error) {
	return workspaceAPI.getWorkspaces()
}

func (workspaceAPI *mockWorkspaceAPI) GetWorkspacesByWorkspace(workspaceID int) ([]model.Workspace, error) {
	return workspaceAPI.getWorkspacesByWorkspace(workspaceID)
}

func Test_CreateWorkspace_ErrorIsReturned(t *testing.T) {
	// arrange
	workspaceAPI := &mockWorkspaceAPI{
		createWorkspace: func(workspace model.Workspace) (model.Workspace, error) {
			return model.Workspace{Name: "Sample model.Workspace", ID: 1}, nil
		},
	}

	workspaceRepository := WorkspaceRepository{
		workspaceAPI: workspaceAPI,
	}

	// act
	_, err := workspaceRepository.CreateWorkspace("Sample model.Workspace")

	// assert
	if err == nil {
		t.Logf("CreateWorkspace should return an error because creating workspaces is not supported")
	}
}

func Test_GetWorkspaces_APIReturnsError_ErrorIsReturned(t *testing.T) {
	// arrange
	workspaceAPI := &mockWorkspaceAPI{
		getWorkspaces: func() ([]model.Workspace, error) {
			return nil, fmt.Errorf("Some error")
		},
	}

	workspaceRepository := WorkspaceRepository{
		workspaceAPI: workspaceAPI,
	}

	// act
	_, err := workspaceRepository.GetWorkspaces()

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetWorkspaces should return an error if the API returned one")
	}
}

func Test_GetWorkspaces_AllWorkspacesReturnedByTheAPIAreReturned(t *testing.T) {
	// arrange
	workspaceAPI := &mockWorkspaceAPI{
		getWorkspaces: func() ([]model.Workspace, error) {
			return []model.Workspace{
				model.Workspace{},
				model.Workspace{},
			}, nil
		},
	}

	workspaceRepository := WorkspaceRepository{
		workspaceAPI: workspaceAPI,
	}

	// act
	workspaces, err := workspaceRepository.GetWorkspaces()

	// assert
	if err != nil {
		t.Logf("GetWorkspaces should not return an error if the API does not return one: %s", err.Error())
	}

	if len(workspaces) != 2 {
		t.Fail()
		t.Logf("GetWorkspaces should have returned two workspaces")
	}
}

func Test_GetWorkspaces_SecondCallUsesTheCache(t *testing.T) {
	// arrange
	numberOfCalls := 0
	workspaceAPI := &mockWorkspaceAPI{
		getWorkspaces: func() ([]model.Workspace, error) {
			numberOfCalls++

			return []model.Workspace{
				model.Workspace{},
				model.Workspace{},
			}, nil
		},
	}

	workspaceRepository := WorkspaceRepository{
		workspaceAPI: workspaceAPI,
	}

	// act
	workspaceRepository.GetWorkspaces() // first call
	workspaceRepository.GetWorkspaces() // second call

	// assert
	if numberOfCalls != 1 {
		t.Fail()
		t.Logf("The seconds call to GetWorkspaces should not have called the API again")
	}
}

func Test_GetWorkspaceByName_NoWorkspacesAvailabe_ErrorIsReturned(t *testing.T) {
	// arrange
	workspaceAPI := &mockWorkspaceAPI{
		getWorkspaces: func() ([]model.Workspace, error) {
			return []model.Workspace{}, nil
		},
	}

	workspaceRepository := WorkspaceRepository{
		workspaceAPI: workspaceAPI,
	}

	// act
	_, err := workspaceRepository.GetWorkspaceByName("Workspace A")

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetWorkspaceByName should return an error if no matching workspace was found")
	}
}

func Test_GetWorkspaceByName_WorkspaceNameDoesNotMatch_ErrorIsReturned(t *testing.T) {
	// arrange
	workspaceAPI := &mockWorkspaceAPI{
		getWorkspaces: func() ([]model.Workspace, error) {
			return []model.Workspace{
				model.Workspace{
					Name: "Workspace A",
				},
				model.Workspace{
					Name: "Workspace B",
				},
			}, nil
		},
	}

	workspaceRepository := WorkspaceRepository{
		workspaceAPI: workspaceAPI,
	}

	// act
	_, err := workspaceRepository.GetWorkspaceByName("Workspace C")

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetWorkspaceByName should return an error if no matching workspace was found")
	}
}

func Test_GetWorkspaceByID_NoWorkspacesAvailabe_ErrorIsReturned(t *testing.T) {
	// arrange
	workspaceAPI := &mockWorkspaceAPI{
		getWorkspaces: func() ([]model.Workspace, error) {
			return []model.Workspace{}, nil
		},
	}

	workspaceRepository := WorkspaceRepository{
		workspaceAPI: workspaceAPI,
	}

	// act
	_, err := workspaceRepository.GetWorkspaceByID(1)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetWorkspaceByID should return an error if no matching workspace was found")
	}
}

func Test_GetWorkspaceByID_WorkspaceIDDoesNotMatch_ErrorIsReturned(t *testing.T) {
	// arrange
	workspaceAPI := &mockWorkspaceAPI{
		getWorkspaces: func() ([]model.Workspace, error) {
			return []model.Workspace{
				model.Workspace{
					ID: 1,
				},
				model.Workspace{
					ID: 2,
				},
			}, nil
		},
	}

	workspaceRepository := WorkspaceRepository{
		workspaceAPI: workspaceAPI,
	}

	// act
	_, err := workspaceRepository.GetWorkspaceByID(3)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetWorkspaceByID should return an error if no matching workspace was found")
	}
}
