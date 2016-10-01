package toggl

import (
	"fmt"
	"testing"
	"time"

	"github.com/andreaskoch/togglapi/model"
)

type mockWorkspacer struct {
	createWorkspace    func(name string) (Workspace, error)
	getWorkspaces      func() ([]Workspace, error)
	getWorkspaceByID   func(workspaceID int) (Workspace, error)
	getWorkspaceByName func(workspaceName string) (Workspace, error)
}

func (workspacer *mockWorkspacer) CreateWorkspace(name string) (Workspace, error) {
	return workspacer.createWorkspace(name)
}

func (workspacer *mockWorkspacer) GetWorkspaces() ([]Workspace, error) {
	return workspacer.getWorkspaces()
}

func (workspacer *mockWorkspacer) GetWorkspaceByID(workspaceID int) (Workspace, error) {
	return workspacer.getWorkspaceByID(workspaceID)
}

func (workspacer *mockWorkspacer) GetWorkspaceByName(workspaceName string) (Workspace, error) {
	return workspacer.getWorkspaceByName(workspaceName)
}

type mockProjecter struct {
	createProject    func(projectName, workspaceName, clientName string) (Project, error)
	getProjects      func() ([]Project, error)
	getProjectByID   func(projectID int) (Project, error)
	getProjectByName func(projectName, workspaceName, clientName string) (Project, error)
}

func (projecter *mockProjecter) CreateProject(projectName, workspaceName, clientName string) (Project, error) {
	return projecter.createProject(projectName, workspaceName, clientName)
}

func (projecter *mockProjecter) GetProjects() ([]Project, error) {
	return projecter.getProjects()
}

func (projecter *mockProjecter) GetProjectByID(projectID int) (Project, error) {
	return projecter.getProjectByID(projectID)
}

func (projecter *mockProjecter) GetProjectByName(projectName, workspaceName, clientName string) (Project, error) {
	return projecter.getProjectByName(projectName, workspaceName, clientName)
}

func Test_ConvertTimeRecordToTimeEntry_WorspaceNotFound_ErrorIsReturned(t *testing.T) {
	// arrange

	modelConverter := &togglModelConverter{
		workspaces: &mockWorkspacer{
			getWorkspaceByName: func(workspaceName string) (Workspace, error) {
				return Workspace{}, fmt.Errorf("Workspace not found")
			},
		},
		projects: &mockProjecter{
			getProjectByName: func(projectName, workspaceName, clientName string) (Project, error) {
				return Project{
					ID: 1,
					Workspace: Workspace{
						ID:   1,
						Name: workspaceName,
					},
					Name: projectName,
				}, nil
			},
		},
	}

	inputTimeRecord := TimeRecord{
		WorkspaceName: "Workspace",
	}

	// act
	_, err := modelConverter.ConvertTimeRecordToTimeEntry(inputTimeRecord)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("ConvertTimeRecordToTimeEntry should have returned an error if the given workspace does not exist")
	}
}

func Test_ConvertTimeRecordToTimeEntry_ProjectNotFound_ErrorIsReturned(t *testing.T) {
	// arrange
	modelConverter := &togglModelConverter{
		workspaces: &mockWorkspacer{
			getWorkspaceByName: func(workspaceName string) (Workspace, error) {
				return Workspace{
					ID:   1,
					Name: workspaceName,
				}, nil
			},
		},
		projects: &mockProjecter{
			getProjectByName: func(projectName, workspaceName, clientName string) (Project, error) {
				return Project{}, fmt.Errorf("Project not found")
			},
		},
	}

	inputTimeRecord := TimeRecord{
		WorkspaceName: "Workspace",
		ProjectName:   "Project",
	}

	// act
	_, err := modelConverter.ConvertTimeRecordToTimeEntry(inputTimeRecord)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("ConvertTimeRecordToTimeEntry should return an error of the project does not exist")
	}

}

func Test_ConvertTimeRecordToTimeEntry_ProjectAndWorkspaceExist_TimeRecordIsReturned(t *testing.T) {
	// arrange
	modelConverter := &togglModelConverter{
		workspaces: &mockWorkspacer{
			getWorkspaceByName: func(workspaceName string) (Workspace, error) {
				return Workspace{
					ID:   1,
					Name: workspaceName,
				}, nil
			},
		},
		projects: &mockProjecter{
			getProjectByName: func(projectName, workspaceName, clientName string) (Project, error) {
				return Project{
					Workspace: Workspace{
						ID:   1,
						Name: workspaceName,
					},
					ID:   1,
					Name: fmt.Sprintf("Project %d", 1),
				}, nil
			},
		},
	}

	start := time.Date(2016, 8, 1, 9, 0, 1, 0, time.UTC)
	stop := time.Date(2016, 8, 1, 9, 30, 2, 0, time.UTC)

	inputTimeRecord := TimeRecord{
		WorkspaceName: "Workspace",
		ProjectName:   "Project",
		Start:         start,
		Stop:          stop,
		Description:   "Yada Yada",
	}

	// act
	resultTimeEntry, err := modelConverter.ConvertTimeRecordToTimeEntry(inputTimeRecord)

	// assert
	if resultTimeEntry.Description != "Yada Yada" {
		t.Fail()
		t.Logf("ConvertTimeRecordToTimeEntry(%#v) have returned a time record with the description %q but returned this instead: %#v", inputTimeRecord, inputTimeRecord.Description, resultTimeEntry)
	}

	if err != nil {
		t.Fail()
		t.Logf("ConvertTimeRecordToTimeEntry(%#v) should not have returned an error but returned this: %s", inputTimeRecord, err)
	}

}

func Test_ConvertTimeEntryToTimeRecord_WorspaceNotFound_ErrorIsReturned(t *testing.T) {
	// arrange
	modelConverter := &togglModelConverter{
		workspaces: &mockWorkspacer{
			getWorkspaceByID: func(workspaceID int) (Workspace, error) {
				return Workspace{}, fmt.Errorf("Workspace not found")
			},
		},
		projects: &mockProjecter{},
	}

	timeEntry := model.TimeEntry{
		Wid: 189,
	}

	// act
	_, err := modelConverter.ConvertTimeEntryToTimeRecord(timeEntry)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("ConvertTimeEntryToTimeRecord(%#v) should have returned an error", timeEntry)
	}

}

func Test_ConvertTimeEntryToTimeRecord_ProjectNotFound_ErrorIsReturned(t *testing.T) {
	// arrange
	modelConverter := &togglModelConverter{
		workspaces: &mockWorkspacer{
			getWorkspaceByID: func(workspaceID int) (Workspace, error) {
				return Workspace{
					ID:   workspaceID,
					Name: "Some workspace",
				}, nil
			},
		},
		projects: &mockProjecter{
			getProjectByID: func(projectID int) (Project, error) {
				return Project{}, fmt.Errorf("Project not found")
			},
		},
		clients: &mockClienter{
			getClientByID: func(clientID int) (Client, error) {
				return Client{
					ID: clientID,
				}, nil
			},
		},
	}

	timeEntry := model.TimeEntry{
		Wid: 1,
		Pid: 12893,
	}

	// act
	_, err := modelConverter.ConvertTimeEntryToTimeRecord(timeEntry)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("ConvertTimeEntryToTimeRecord(%#v) should have returned an error", timeEntry)
	}

}

func Test_ConvertTimeEntryToTimeRecord_ProjectAndWorkspaceExist_TimeRecordIsReturned(t *testing.T) {
	// arrange
	modelConverter := &togglModelConverter{
		workspaces: &mockWorkspacer{
			getWorkspaceByID: func(workspaceID int) (Workspace, error) {
				return Workspace{
					ID:   workspaceID,
					Name: "Some workspace",
				}, nil
			},
		},
		projects: &mockProjecter{
			getProjectByID: func(projectID int) (Project, error) {
				return Project{
					ID:   projectID,
					Name: "Some project",
				}, nil
			},
		},
		clients: &mockClienter{
			getClientByID: func(clientID int) (Client, error) {
				return Client{
					ID: clientID,
				}, nil
			},
		},
	}

	start := time.Date(2016, 8, 1, 9, 0, 1, 0, time.UTC)
	stop := time.Date(2016, 8, 1, 9, 30, 2, 0, time.UTC)

	timeEntry := model.TimeEntry{
		Wid:         1,
		Pid:         12893,
		Start:       start,
		Stop:        stop,
		Description: "Yada Yada",
	}

	// act
	timeRecord, err := modelConverter.ConvertTimeEntryToTimeRecord(timeEntry)

	// assert
	if timeRecord.Description != "Yada Yada" {
		t.Fail()
		t.Logf("ConvertTimeEntryToTimeRecord(%#v) have returned a time record with the description %q but returned this instead: %#v", timeEntry, timeEntry.Description, timeRecord)
	}

	if err != nil {
		t.Fail()
		t.Logf("ConvertTimeEntryToTimeRecord(%#v) should not have returned an error but returned this: %s", timeEntry, err)
	}

}
