package toggl

import (
	"fmt"
	"testing"

	"github.com/andreaskoch/togglapi/model"
)

type mockClientAPI struct {
	createClient func(client model.Client) (model.Client, error)
	getClients   func() ([]model.Client, error)
}

func (clientAPI *mockClientAPI) CreateClient(client model.Client) (model.Client, error) {
	return clientAPI.createClient(client)
}

func (clientAPI *mockClientAPI) GetClients() ([]model.Client, error) {
	return clientAPI.getClients()
}

func Test_CreateClient_CreateSucceeds_ClientIsReturned(t *testing.T) {
	// arrange
	clientAPI := &mockClientAPI{
		createClient: func(client model.Client) (model.Client, error) {
			return model.Client{Name: "Sample Client", WorkspaceID: 1}, nil
		},
	}

	clientRepository := ClientRepository{
		clientAPI: clientAPI,
		workspaces: &mockWorkspacer{
			getWorkspaceByID: func(workspaceID int) (Workspace, error) {
				return Workspace{ID: workspaceID}, nil
			},
		},
	}

	// act
	client, err := clientRepository.CreateClient(1, "Sample Client")

	// assert
	if err != nil {
		t.Logf("CreateClient should not return an error if the create succeeded: %s", err.Error())
	}

	if client.Name != "Sample Client" || client.Workspace.ID != 1 {
		t.Fail()
		t.Logf("CreateClient have returned a client model with the given parameters but returned this instead: %#v", client)
	}
}

func Test_CreateClient_CreateFails_ErrorIsReturned(t *testing.T) {
	// arrange
	clientAPI := &mockClientAPI{
		createClient: func(client model.Client) (model.Client, error) {
			return model.Client{}, fmt.Errorf("Some error")
		},
	}

	clientRepository := ClientRepository{
		clientAPI: clientAPI,
		workspaces: &mockWorkspacer{
			getWorkspaceByID: func(workspaceID int) (Workspace, error) {
				return Workspace{ID: workspaceID}, nil
			},
		},
	}

	// act
	_, err := clientRepository.CreateClient(1, "Sample Client")

	// assert
	if err == nil {
		t.Fail()
		t.Logf("CreateClient return an error if the API returns one")
	}
}

func Test_GetClients_APIReturnsError_ErrorIsReturned(t *testing.T) {
	// arrange
	clientAPI := &mockClientAPI{
		getClients: func() ([]model.Client, error) {
			return nil, fmt.Errorf("Some error")
		},
	}

	workspaceProvider := &mockWorkspacer{
		getWorkspaceByID: func(workspaceID int) (Workspace, error) {
			return Workspace{
				ID:   workspaceID,
				Name: "Some workspace",
			}, nil
		},
	}

	clientRepository := ClientRepository{
		clientAPI:  clientAPI,
		workspaces: workspaceProvider,
	}

	// act
	_, err := clientRepository.GetClients()

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetClients should return an error if the API returned one")
	}
}

func Test_GetClients_AllClientsReturnedByTheAPIAreReturned(t *testing.T) {
	// arrange
	clientAPI := &mockClientAPI{
		getClients: func() ([]model.Client, error) {
			return []model.Client{
				model.Client{},
				model.Client{},
			}, nil
		},
	}

	workspaceProvider := &mockWorkspacer{
		getWorkspaceByID: func(workspaceID int) (Workspace, error) {
			return Workspace{
				ID:   workspaceID,
				Name: "Some workspace",
			}, nil
		},
	}

	clientRepository := ClientRepository{
		clientAPI:  clientAPI,
		workspaces: workspaceProvider,
	}

	// act
	clients, err := clientRepository.GetClients()

	// assert
	if err != nil {
		t.Logf("GetClients should not return an error if the API does not return one: %s", err.Error())
	}

	if len(clients) != 2 {
		t.Fail()
		t.Logf("GetClients should have returned two clients")
	}
}

func Test_GetClients_SecondCallUsesTheCache(t *testing.T) {
	// arrange
	numberOfCalls := 0
	clientAPI := &mockClientAPI{
		getClients: func() ([]model.Client, error) {
			numberOfCalls++

			return []model.Client{
				model.Client{},
				model.Client{},
			}, nil
		},
	}

	workspaceProvider := &mockWorkspacer{
		getWorkspaceByID: func(workspaceID int) (Workspace, error) {
			return Workspace{
				ID:   workspaceID,
				Name: "Some workspace",
			}, nil
		},
	}

	clientRepository := ClientRepository{
		clientAPI:  clientAPI,
		workspaces: workspaceProvider,
	}

	// act
	clientRepository.GetClients() // first call
	clientRepository.GetClients() // second call

	// assert
	if numberOfCalls != 1 {
		t.Fail()
		t.Logf("The seconds call to GetClients should not have called the API again")
	}
}

func Test_GetClientByName_NoClientsAvailabe_ErrorIsReturned(t *testing.T) {
	// arrange
	clientAPI := &mockClientAPI{
		getClients: func() ([]model.Client, error) {
			return []model.Client{}, nil
		},
	}

	workspaceProvider := &mockWorkspacer{
		getWorkspaceByID: func(workspaceID int) (Workspace, error) {
			return Workspace{
				ID:   workspaceID,
				Name: "Some workspace",
			}, nil
		},
	}

	clientRepository := ClientRepository{
		clientAPI:  clientAPI,
		workspaces: workspaceProvider,
	}

	// act
	_, err := clientRepository.GetClientByName("A Workspace", "Client A")

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetClientByName should return an error if no matching client was found")
	}
}

func Test_GetClientByName_ClientNameDoesNotMatch_ErrorIsReturned(t *testing.T) {
	// arrange
	clientAPI := &mockClientAPI{
		getClients: func() ([]model.Client, error) {
			return []model.Client{
				model.Client{
					WorkspaceID: 1,
					Name:        "Client A",
				},
				model.Client{
					WorkspaceID: 1,
					Name:        "Client B",
				},
			}, nil
		},
	}

	workspaceProvider := &mockWorkspacer{
		getWorkspaceByID: func(workspaceID int) (Workspace, error) {
			return Workspace{
				ID:   workspaceID,
				Name: "Some workspace",
			}, nil
		},
	}

	clientRepository := ClientRepository{
		clientAPI:  clientAPI,
		workspaces: workspaceProvider,
	}

	// act
	_, err := clientRepository.GetClientByName("A workspace", "Client C")

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetClientByName should return an error if no matching client was found")
	}
}

func Test_GetClientByName_WorkspaceNameMatches_ClientIsReturned(t *testing.T) {
	// arrange
	clientAPI := &mockClientAPI{
		getClients: func() ([]model.Client, error) {
			return []model.Client{
				model.Client{
					WorkspaceID: 1,
					Name:        "Client A",
				},
				model.Client{
					WorkspaceID: 1,
					Name:        "Client B",
				},
			}, nil
		},
	}

	workspaceProvider := &mockWorkspacer{
		getWorkspaceByID: func(workspaceID int) (Workspace, error) {
			return Workspace{
				ID:   workspaceID,
				Name: "A Workspace",
			}, nil
		},
	}

	clientRepository := ClientRepository{
		clientAPI:  clientAPI,
		workspaces: workspaceProvider,
	}

	// act
	client, err := clientRepository.GetClientByName("A Workspace", "Client A")

	// assert
	if err != nil {
		t.Logf("GetClientByName should not have returned an error but returned this: %s", err)
	}

	if client.Name != "Client A" || client.Workspace.ID != 1 {
		t.Fail()
		t.Logf("GetClientByName should have returned a matching client")
	}
}

func Test_GetClientByID_NoClientsAvailabe_ErrorIsReturned(t *testing.T) {
	// arrange
	clientAPI := &mockClientAPI{
		getClients: func() ([]model.Client, error) {
			return []model.Client{}, nil
		},
	}

	workspaceProvider := &mockWorkspacer{
		getWorkspaceByID: func(workspaceID int) (Workspace, error) {
			return Workspace{
				ID:   workspaceID,
				Name: "Some workspace",
			}, nil
		},
	}

	clientRepository := ClientRepository{
		clientAPI:  clientAPI,
		workspaces: workspaceProvider,
	}

	// act
	_, err := clientRepository.GetClientByID(1)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetClientByID should return an error if no matching client was found")
	}
}

func Test_GetClientByID_WorkspacerReturnsError_ErrorIsReturned(t *testing.T) {
	// arrange
	clientAPI := &mockClientAPI{
		getClients: func() ([]model.Client, error) {
			return []model.Client{}, nil
		},
	}

	workspaceProvider := &mockWorkspacer{
		getWorkspaces: func() ([]Workspace, error) {
			return nil, fmt.Errorf("Error")
		},
	}

	clientRepository := ClientRepository{
		clientAPI:  clientAPI,
		workspaces: workspaceProvider,
	}

	// act
	_, err := clientRepository.GetClientByID(1)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetClientByID should return an error if no matching client was found")
	}
}

func Test_GetClientByID_ClientNameDoesNotMatch_ErrorIsReturned(t *testing.T) {
	// arrange
	clientAPI := &mockClientAPI{
		getClients: func() ([]model.Client, error) {
			return []model.Client{
				model.Client{
					WorkspaceID: 1,
					ID:          1,
				},
				model.Client{
					WorkspaceID: 1,
					ID:          2,
				},
			}, nil
		},
	}

	workspaceProvider := &mockWorkspacer{
		getWorkspaceByID: func(workspaceID int) (Workspace, error) {
			return Workspace{
				ID:   workspaceID,
				Name: "Some workspace",
			}, nil
		},
	}

	clientRepository := ClientRepository{
		clientAPI:  clientAPI,
		workspaces: workspaceProvider,
	}

	// act
	_, err := clientRepository.GetClientByID(3)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetClientByID should return an error if no matching client was found")
	}
}

func Test_GetClientByID_ClientIDIsFound_ClientIsReturned(t *testing.T) {
	// arrange
	clientAPI := &mockClientAPI{
		getClients: func() ([]model.Client, error) {
			return []model.Client{
				model.Client{
					WorkspaceID: 1,
					ID:          1,
				},
				model.Client{
					WorkspaceID: 1,
					ID:          2,
				},
			}, nil
		},
	}

	workspaceProvider := &mockWorkspacer{
		getWorkspaceByID: func(workspaceID int) (Workspace, error) {
			return Workspace{
				ID:   workspaceID,
				Name: "Some Workspace",
			}, nil
		},
	}

	clientRepository := ClientRepository{
		clientAPI:  clientAPI,
		workspaces: workspaceProvider,
	}

	// act
	client, err := clientRepository.GetClientByID(1)

	// assert
	if err != nil {
		t.Logf("GetClientByID should not have returned an error but returned this: %s", err)
	}

	if client.ID != 1 || client.Workspace.ID != 1 {
		t.Fail()
		t.Logf("GetClientByID should have returned a matching client")
	}
}
