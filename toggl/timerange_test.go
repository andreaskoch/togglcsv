package toggl

import (
	"testing"
	"time"
)

type mockTimeRangeProvider struct {
	getTimeRanges func(startDate, endDate time.Time) ([]timeRange, error)
}

func (timeRangeProvider *mockTimeRangeProvider) GetTimeRanges(startDate, endDate time.Time) ([]timeRange, error) {
	return timeRangeProvider.getTimeRanges(startDate, endDate)
}

func Test_newTimeRange_StopBeforeStart_ErrorIsReturned(t *testing.T) {
	// arrange
	start := time.Date(2016, 8, 30, 0, 0, 1, 0, time.UTC)
	stop := time.Date(2016, 1, 1, 0, 0, 1, 0, time.UTC)

	// act
	_, err := newTimeRange(start, stop)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("newTimeRange(%q, %q) should have returned an error", start, stop)
	}
}

func Test_newTimeRange_StopEqualsStart_ErrorIsReturned(t *testing.T) {
	// arrange
	start := time.Date(2016, 8, 30, 0, 0, 1, 0, time.UTC)
	stop := time.Date(2016, 8, 30, 0, 0, 1, 0, time.UTC)

	// act
	_, err := newTimeRange(start, stop)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("newTimeRange(%q, %q) should have returned an error", start, stop)
	}
}

func Test_GetTimeRanges_ExpectedRangesAreReturned(t *testing.T) {
	// arrange
	timeRangeProvider := &fullMonthTimeRangeProvider{}
	inputs := []struct {
		Start                  time.Time
		Stop                   time.Time
		ExpectedNumberOfRanges int
	}{
		// 8 month, 8 ranges are returned (beginning of January till end of August)
		{
			Start: time.Date(2016, 1, 1, 0, 0, 1, 0, time.UTC),
			Stop:  time.Date(2016, 8, 30, 0, 0, 1, 0, time.UTC),
			ExpectedNumberOfRanges: 8,
		},

		// 8 month, 8 ranges are returned (beginning of January till beginning of August)
		{
			Start: time.Date(2016, 1, 1, 0, 0, 1, 0, time.UTC),
			Stop:  time.Date(2016, 8, 1, 0, 0, 1, 0, time.UTC),
			ExpectedNumberOfRanges: 8,
		},

		// 8 month, 8 ranges are returned (end of January till end of August)
		{
			Start: time.Date(2016, 1, 31, 0, 0, 1, 0, time.UTC),
			Stop:  time.Date(2016, 8, 30, 0, 0, 1, 0, time.UTC),
			ExpectedNumberOfRanges: 8,
		},

		// 1 month, 1 range is returned (beginning of August till mid of August)
		{
			Start: time.Date(2016, 8, 1, 0, 0, 1, 0, time.UTC),
			Stop:  time.Date(2016, 8, 12, 0, 0, 1, 0, time.UTC),
			ExpectedNumberOfRanges: 1,
		},

		// 2 month, 2 ranges are returned (mid of July till mid of August)
		{
			Start: time.Date(2016, 7, 12, 0, 0, 1, 0, time.UTC),
			Stop:  time.Date(2016, 8, 12, 0, 0, 1, 0, time.UTC),
			ExpectedNumberOfRanges: 2,
		},
	}

	for _, input := range inputs {

		// act
		ranges, _ := timeRangeProvider.GetTimeRanges(input.Start, input.Stop)

		// assert
		if len(ranges) != input.ExpectedNumberOfRanges {
			t.Fail()
			t.Logf("GetTimeRanges(%q, %q) should have returned %d ranges but returned %d instead",
				input.Start,
				input.Stop,
				input.ExpectedNumberOfRanges,
				len(ranges))
		}

	}

}

func Test_GetTimeRanges_StopDateBeforeStartDate_ErrorIsReturned(t *testing.T) {
	// arrange
	stop := time.Date(2016, 1, 1, 0, 0, 1, 0, time.UTC)
	start := time.Date(2016, 8, 1, 0, 0, 1, 0, time.UTC)
	timeRangeProvider := &fullMonthTimeRangeProvider{}

	// act
	_, err := timeRangeProvider.GetTimeRanges(start, stop)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetTimeRanges(%q, %q) should have returned an error", start, stop)
	}

}
