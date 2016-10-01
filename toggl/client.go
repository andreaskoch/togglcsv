package toggl

import (
	"fmt"

	"github.com/andreaskoch/togglapi/model"
	"github.com/pkg/errors"
)

// A Client describes a time reporting client
type Client struct {
	ID        int
	Name      string
	Workspace Workspace
}

// A Clienter interface provides read/write access to Toggl clients.
type Clienter interface {
	// CreateClient creates a new client with the given name.
	// Returns an error of the creation failed.
	CreateClient(workspaceID int, name string) (Client, error)

	// GetClients returns all clients.
	GetClients() ([]Client, error)

	// GetClientByID returns the client for the given client id.
	// Returns an error if the client was not found.
	GetClientByID(clientID int) (Client, error)

	// GetClientByName returns the client for the given client name.
	// Returns an error if no matching client was found.
	GetClientByName(workspaceName, clientName string) (Client, error)
}

// NewClientRepository creates a new client repository instance.
func NewClientRepository(clientAPI model.ClientAPI, workspaceProvider Workspacer) Clienter {
	return &ClientRepository{
		clientAPI:  clientAPI,
		workspaces: workspaceProvider,
	}
}

// ClientRepository provides read/write access to Toggl clients.
type ClientRepository struct {
	clientAPI    model.ClientAPI
	clientsCache []Client

	workspaces Workspacer
}

// CreateClient creates a new client with the given name.
// Returns an error of the creation failed.
func (repository *ClientRepository) CreateClient(workspaceID int, name string) (Client, error) {

	workspace, workspaceError := repository.workspaces.GetWorkspaceByID(workspaceID)
	if workspaceError != nil {
		return Client{}, errors.Wrap(workspaceError, fmt.Sprintf("Failed to get workspace with id %d", workspaceID))
	}

	createdClient, err := repository.clientAPI.CreateClient(model.Client{
		Name:        name,
		WorkspaceID: workspace.ID,
	})

	if err != nil {
		return Client{}, err
	}

	// reset the clients cache
	repository.clientsCache = nil

	return Client{
		ID:        createdClient.ID,
		Name:      createdClient.Name,
		Workspace: workspace,
	}, nil
}

// GetClients returns all clients.
func (repository *ClientRepository) GetClients() ([]Client, error) {
	if repository.clientsCache != nil {
		return repository.clientsCache, nil
	}

	var clientModels []Client

	clients, clientsError := repository.clientAPI.GetClients()
	if clientsError != nil {
		return nil, errors.Wrap(clientsError, "Failed to get clients from Toggl")
	}

	for _, client := range clients {

		workspace, workspaceErr := repository.workspaces.GetWorkspaceByID(client.WorkspaceID)
		if workspaceErr != nil {
			return nil, errors.Wrap(workspaceErr, fmt.Sprintf("Failed to get workspace for client %d", client.ID))
		}

		clientModels = append(clientModels, Client{
			ID:        client.ID,
			Name:      client.Name,
			Workspace: workspace,
		})
	}

	// store in cache
	repository.clientsCache = clientModels

	return clientModels, nil
}

// GetClientByID returns the client for the given client id.
// Returns an error if the client was not found.
func (repository *ClientRepository) GetClientByID(clientID int) (Client, error) {
	clients, clientsError := repository.GetClients()
	if clientsError != nil {
		return Client{}, clientsError
	}

	for _, client := range clients {
		if client.ID == clientID {
			return client, nil
		}
	}

	return Client{}, fmt.Errorf("No client found with id %d", clientID)
}

// GetClientByName returns the client for the given client name.
// Returns an error if no matching client was found.
func (repository *ClientRepository) GetClientByName(workspaceName, clientName string) (Client, error) {

	clients, clientsError := repository.GetClients()
	if clientsError != nil {
		return Client{}, clientsError
	}

	for _, client := range clients {
		if client.Workspace.Name == workspaceName && client.Name == clientName {
			return client, nil
		}
	}

	return Client{}, fmt.Errorf("Client %q was not found (Workspace: %q)", clientName, workspaceName)
}
