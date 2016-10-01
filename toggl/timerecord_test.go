package toggl

import (
	"fmt"
	"testing"
	"time"

	"github.com/andreaskoch/togglapi/model"
)

// The model.TimeEntryAPI interface provides functions for fetching and creating time entries.
type mockTimeEntryAPI struct {
	createTimeEntry func(timeEntry model.TimeEntry) (model.TimeEntry, error)
	getTimeEntries  func(start, end time.Time) ([]model.TimeEntry, error)
}

func (timeEntryAPI *mockTimeEntryAPI) CreateTimeEntry(timeEntry model.TimeEntry) (model.TimeEntry, error) {
	return timeEntryAPI.createTimeEntry(timeEntry)
}

func (timeEntryAPI *mockTimeEntryAPI) GetTimeEntries(start, end time.Time) ([]model.TimeEntry, error) {
	return timeEntryAPI.getTimeEntries(start, end)
}

type mockModelConverter struct {
	convertTimeEntryToTimeRecord func(timeEntry model.TimeEntry) (TimeRecord, error)
	convertTimeRecordToTimeEntry func(timeRecord TimeRecord) (model.TimeEntry, error)
}

func (converter *mockModelConverter) ConvertTimeEntryToTimeRecord(timeEntry model.TimeEntry) (TimeRecord, error) {
	return converter.convertTimeEntryToTimeRecord(timeEntry)
}

func (converter *mockModelConverter) ConvertTimeRecordToTimeEntry(timeRecord TimeRecord) (model.TimeEntry, error) {
	return converter.convertTimeRecordToTimeEntry(timeRecord)
}

func Test_NewTimeRecordRepository(t *testing.T) {
	// arrange
	timeEntryAPI := &mockTimeEntryAPI{}
	workspacer := &mockWorkspacer{}
	projecter := &mockProjecter{}
	clienter := &mockClienter{}

	// act
	repository := NewTimeRecordRepository(timeEntryAPI, workspacer, projecter, clienter)

	if repository == nil {
		t.Fail()
		t.Logf("NewTimeRecordRepository should have returned a repositor")
	}
}

func Test_CreateTimeRecord_ConversionFails_ErrorIsReturned(t *testing.T) {
	// arrange
	repository := &TimeRecordRepository{
		modelConverter: &mockModelConverter{
			convertTimeRecordToTimeEntry: func(timeRecord TimeRecord) (model.TimeEntry, error) {
				return model.TimeEntry{}, fmt.Errorf("Conversion failed")
			},
		},
		workspaces: &mockWorkspacer{
			getWorkspaceByName: func(workspaceName string) (Workspace, error) {
				return Workspace{}, nil
			},
		},

		projects: &mockProjecter{
			getProjectByName: func(projectName, workspaceName, clientName string) (Project, error) {
				return Project{}, nil
			},
		},
	}

	start := time.Date(2016, 8, 1, 0, 0, 1, 0, time.UTC)
	stop := time.Date(2016, 8, 31, 23, 59, 59, 0, time.UTC)

	inputTimeRecord := TimeRecord{
		Start:       start,
		Stop:        stop,
		Description: "Yada Yada",
	}

	// act
	err := repository.CreateTimeRecord(inputTimeRecord)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("CreateTimeRecord(%q) should have returned an error", inputTimeRecord)
	}
}

func Test_CreateTimeRecord_CreateFails_ErrorIsReturned(t *testing.T) {
	// arrange
	timeEntryAPI := &mockTimeEntryAPI{
		createTimeEntry: func(timeEntry model.TimeEntry) (model.TimeEntry, error) {
			return model.TimeEntry{}, fmt.Errorf("Fail")
		},
	}

	repository := &TimeRecordRepository{
		timeEntryAPI: timeEntryAPI,
		modelConverter: &mockModelConverter{
			convertTimeRecordToTimeEntry: func(timeRecord TimeRecord) (model.TimeEntry, error) {
				return model.TimeEntry{}, nil
			},
		},
		workspaces: &mockWorkspacer{
			getWorkspaceByName: func(workspaceName string) (Workspace, error) {
				return Workspace{}, nil
			},
		},

		projects: &mockProjecter{
			getProjectByName: func(projectName, workspaceName, clientName string) (Project, error) {
				return Project{}, nil
			},
		},
	}

	start := time.Date(2016, 8, 1, 0, 0, 1, 0, time.UTC)
	stop := time.Date(2016, 8, 31, 23, 59, 59, 0, time.UTC)

	inputTimeRecord := TimeRecord{
		Start:       start,
		Stop:        stop,
		Description: "Yada Yada",
	}

	// act
	err := repository.CreateTimeRecord(inputTimeRecord)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("CreateTimeRecord(%q) should have returned an error", inputTimeRecord)
	}
}

func Test_CreateTimeRecord_CreateSucceeds_NoErrorIsReturned(t *testing.T) {
	// arrange
	timeEntryAPI := &mockTimeEntryAPI{
		createTimeEntry: func(timeEntry model.TimeEntry) (model.TimeEntry, error) {
			return model.TimeEntry{}, nil
		},
	}

	repository := &TimeRecordRepository{
		timeEntryAPI: timeEntryAPI,
		modelConverter: &mockModelConverter{
			convertTimeRecordToTimeEntry: func(timeRecord TimeRecord) (model.TimeEntry, error) {
				return model.TimeEntry{}, nil
			},
		},
		workspaces: &mockWorkspacer{
			getWorkspaceByName: func(workspaceName string) (Workspace, error) {
				return Workspace{}, nil
			},
		},

		projects: &mockProjecter{
			getProjectByName: func(projectName, workspaceName, clientName string) (Project, error) {
				return Project{}, nil
			},
		},
	}

	start := time.Date(2016, 8, 1, 0, 0, 1, 0, time.UTC)
	stop := time.Date(2016, 8, 31, 23, 59, 59, 0, time.UTC)

	inputTimeRecord := TimeRecord{
		Start:       start,
		Stop:        stop,
		Description: "Yada Yada",
	}

	// act
	err := repository.CreateTimeRecord(inputTimeRecord)

	// assert
	if err != nil {
		t.Fail()
		t.Logf("CreateTimeRecord(%q) should not have returned an error but returned this instead: %s", inputTimeRecord, err)
	}
}

func Test_GetTimeRecords_StartDateAfterStopDate_ErrorIsReturned(t *testing.T) {
	// arrange
	repository := &TimeRecordRepository{}

	stop := time.Date(2016, 8, 1, 0, 0, 1, 0, time.UTC)
	start := time.Date(2016, 8, 31, 23, 59, 59, 0, time.UTC)

	// act
	_, err := repository.GetTimeRecords(start, stop)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetTimeRecords(%q, %q) should have returned an error", start, stop)
	}
}

func Test_GetTimeRecords_NoRangesReturned_EmptySliceIsReturned(t *testing.T) {
	// arrange
	timeRangeProvider := &mockTimeRangeProvider{
		getTimeRanges: func(startDate, endDate time.Time) ([]timeRange, error) {
			return []timeRange{}, nil
		},
	}

	repository := &TimeRecordRepository{
		timeRangeProvider: timeRangeProvider,
	}

	start := time.Date(2016, 8, 1, 0, 0, 1, 0, time.UTC)
	stop := time.Date(2016, 8, 31, 23, 59, 59, 0, time.UTC)

	// act
	records, err := repository.GetTimeRecords(start, stop)

	// assert
	if len(records) > 0 {
		t.Fail()
		t.Logf("GetTimeRecords(%q, %q) should not have returned any records", start, stop)
	}

	if err != nil {
		t.Fail()
		t.Logf("GetTimeRecords(%q, %q) should not have an error but returned this: %s", start, stop, err)
	}
}

func Test_GetTimeRecords_OneTimeRangeReturned_APIReturnsError_ErrorIsReturned(t *testing.T) {
	// arrange
	start := time.Date(2016, 8, 1, 0, 0, 1, 0, time.UTC)
	stop := time.Date(2016, 8, 31, 23, 59, 59, 0, time.UTC)

	timeEntryAPI := &mockTimeEntryAPI{
		getTimeEntries: func(start, end time.Time) ([]model.TimeEntry, error) {
			return nil, fmt.Errorf("Time Entry API error")
		},
	}

	timeRangeProvider := &mockTimeRangeProvider{
		getTimeRanges: func(startDate, endDate time.Time) ([]timeRange, error) {
			return []timeRange{
				timeRange{start, stop},
			}, nil
		},
	}

	repository := &TimeRecordRepository{
		timeEntryAPI:      timeEntryAPI,
		timeRangeProvider: timeRangeProvider,
	}

	// act
	_, err := repository.GetTimeRecords(start, stop)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetTimeRecords(%q, %q) should have returned an error but didn't.", start, stop)
	}
}

func Test_GetTimeRecords_OneTimeRangeReturned_APIReturnsNoTimeEntries_EmptySliceIsReturned(t *testing.T) {
	// arrange
	start := time.Date(2016, 8, 1, 0, 0, 1, 0, time.UTC)
	stop := time.Date(2016, 8, 31, 23, 59, 59, 0, time.UTC)

	timeEntryAPI := &mockTimeEntryAPI{
		getTimeEntries: func(start, end time.Time) ([]model.TimeEntry, error) {
			return []model.TimeEntry{}, nil
		},
	}

	timeRangeProvider := &mockTimeRangeProvider{
		getTimeRanges: func(startDate, endDate time.Time) ([]timeRange, error) {
			return []timeRange{
				timeRange{start, stop},
			}, nil
		},
	}

	repository := &TimeRecordRepository{
		timeEntryAPI:      timeEntryAPI,
		timeRangeProvider: timeRangeProvider,
	}

	// act
	records, err := repository.GetTimeRecords(start, stop)

	// assert
	if len(records) > 0 {
		t.Fail()
		t.Logf("GetTimeRecords(%q, %q) should not have returned any records", start, stop)
	}

	if err != nil {
		t.Fail()
		t.Logf("GetTimeRecords(%q, %q) should not have an error but returned this: %s", start, stop, err)
	}
}

func Test_GetTimeRecords_APIReturnsRunningTimeEntry_NoTimeRecordIsReturned(t *testing.T) {
	// arrange
	start := time.Date(2016, 8, 1, 0, 0, 1, 0, time.UTC)
	stop := time.Date(2016, 8, 31, 23, 59, 59, 0, time.UTC)

	runningTimeEntry := model.TimeEntry{
		Pid:   1,
		Wid:   1,
		Start: start,
	}

	timeEntryAPI := &mockTimeEntryAPI{
		getTimeEntries: func(start, end time.Time) ([]model.TimeEntry, error) {
			return []model.TimeEntry{
				runningTimeEntry,
			}, nil
		},
	}

	timeRangeProvider := &mockTimeRangeProvider{
		getTimeRanges: func(startDate, endDate time.Time) ([]timeRange, error) {
			return []timeRange{
				timeRange{start, stop},
			}, nil
		},
	}

	repository := &TimeRecordRepository{
		timeEntryAPI:      timeEntryAPI,
		timeRangeProvider: timeRangeProvider,
		modelConverter: &mockModelConverter{
			convertTimeEntryToTimeRecord: func(timeEntry model.TimeEntry) (TimeRecord, error) {
				return TimeRecord{}, nil
			},
		},
	}

	// act
	records, err := repository.GetTimeRecords(start, stop)

	// assert
	if len(records) != 0 {
		t.Fail()
		t.Logf("GetTimeRecords(%q, %q) should not have returned a time record", start, stop)
	}

	if err != nil {
		t.Fail()
		t.Logf("GetTimeRecords(%q, %q) should not have an error but returned this: %s", start, stop, err)
	}
}

func Test_GetTimeRecords_APIReturnsTimeEntryWithoutProject_NoTimeRecordIsReturned(t *testing.T) {
	// arrange
	start := time.Date(2016, 8, 1, 0, 0, 1, 0, time.UTC)
	stop := time.Date(2016, 8, 31, 23, 59, 59, 0, time.UTC)

	timeEntryWithoutProject := model.TimeEntry{
		Pid:   0,
		Wid:   1,
		Start: start,
		Stop:  stop,
	}

	timeEntryAPI := &mockTimeEntryAPI{
		getTimeEntries: func(start, end time.Time) ([]model.TimeEntry, error) {
			return []model.TimeEntry{
				timeEntryWithoutProject,
			}, nil
		},
	}

	timeRangeProvider := &mockTimeRangeProvider{
		getTimeRanges: func(startDate, endDate time.Time) ([]timeRange, error) {
			return []timeRange{
				timeRange{start, stop},
			}, nil
		},
	}

	repository := &TimeRecordRepository{
		timeEntryAPI:      timeEntryAPI,
		timeRangeProvider: timeRangeProvider,
		modelConverter: &mockModelConverter{
			convertTimeEntryToTimeRecord: func(timeEntry model.TimeEntry) (TimeRecord, error) {
				return TimeRecord{}, nil
			},
		},
	}

	// act
	records, err := repository.GetTimeRecords(start, stop)

	// assert
	if len(records) != 0 {
		t.Fail()
		t.Logf("GetTimeRecords(%q, %q) should not have returned a time record", start, stop)
	}

	if err != nil {
		t.Fail()
		t.Logf("GetTimeRecords(%q, %q) should not have an error but returned this: %s", start, stop, err)
	}
}

func Test_GetTimeRecords_APIReturnsValidTimeEntry_ConversionFails_ErrorIsReturned(t *testing.T) {
	// arrange
	start := time.Date(2016, 8, 1, 0, 0, 1, 0, time.UTC)
	stop := time.Date(2016, 8, 31, 23, 59, 59, 0, time.UTC)

	validTimeEntry := model.TimeEntry{
		Pid:   1,
		Wid:   1,
		Start: start,
		Stop:  stop,
	}

	timeEntryAPI := &mockTimeEntryAPI{
		getTimeEntries: func(start, end time.Time) ([]model.TimeEntry, error) {
			return []model.TimeEntry{
				validTimeEntry,
			}, nil
		},
	}

	timeRangeProvider := &mockTimeRangeProvider{
		getTimeRanges: func(startDate, endDate time.Time) ([]timeRange, error) {
			return []timeRange{
				timeRange{start, stop},
			}, nil
		},
	}

	repository := &TimeRecordRepository{
		timeEntryAPI:      timeEntryAPI,
		timeRangeProvider: timeRangeProvider,
		modelConverter: &mockModelConverter{
			convertTimeEntryToTimeRecord: func(timeEntry model.TimeEntry) (TimeRecord, error) {
				return TimeRecord{}, fmt.Errorf("Some conversion error")
			},
		},
	}

	// act
	records, err := repository.GetTimeRecords(start, stop)

	// assert
	if len(records) != 0 {
		t.Fail()
		t.Logf("GetTimeRecords(%q, %q) should not have returned a time record", start, stop)
	}

	if err == nil {
		t.Fail()
		t.Logf("GetTimeRecords(%q, %q) should have an error but didn't", start, stop)
	}
}

func Test_GetTimeRecords_APIReturnsValidTimeEntry_ConversionSucceeds_TimeRecordIsReturned(t *testing.T) {
	// arrange
	start := time.Date(2016, 8, 1, 0, 0, 1, 0, time.UTC)
	stop := time.Date(2016, 8, 31, 23, 59, 59, 0, time.UTC)

	validTimeEntry := model.TimeEntry{
		Pid:   1,
		Wid:   1,
		Start: start,
		Stop:  stop,
	}

	timeEntryAPI := &mockTimeEntryAPI{
		getTimeEntries: func(start, end time.Time) ([]model.TimeEntry, error) {
			return []model.TimeEntry{
				validTimeEntry,
			}, nil
		},
	}

	timeRangeProvider := &mockTimeRangeProvider{
		getTimeRanges: func(startDate, endDate time.Time) ([]timeRange, error) {
			return []timeRange{
				timeRange{start, stop},
			}, nil
		},
	}

	repository := &TimeRecordRepository{
		timeEntryAPI:      timeEntryAPI,
		timeRangeProvider: timeRangeProvider,
		modelConverter: &mockModelConverter{
			convertTimeEntryToTimeRecord: func(timeEntry model.TimeEntry) (TimeRecord, error) {
				return TimeRecord{}, nil
			},
		},
	}

	// act
	records, err := repository.GetTimeRecords(start, stop)

	// assert
	if len(records) != 1 {
		t.Fail()
		t.Logf("GetTimeRecords(%q, %q) should have returned one time record", start, stop)
	}

	if err != nil {
		t.Fail()
		t.Logf("GetTimeRecords(%q, %q) should not have an error but returned this: %s", start, stop, err)
	}
}
