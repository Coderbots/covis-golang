package services

import (
	"fmt"
	"testing"
)

type readUrlServiceMock struct {
	testData string
}

func (service readUrlServiceMock) ReadRawData() string {
	return service.testData
}

func TestReadProcessData(t *testing.T) {

	type w struct {
		data   []CovidData
		errTxt string
	}

	// Table of tests.
	var tests = []struct {
		testData string
		want     w
	}{
		{testData: "Hello", want: w{[]CovidData{}, "Internal error encountered"}},
		{testData: `90033,Unassigned,New Hampshire,US,2022-03-12 04:20:50,,,7956,26,,,"Unassigned, New Hampshire, US",,0.32679738562091504\n
		34001,Atlantic,New Jersey,US,2022-03-12 04:20:50,39.47538693,-74.65848483,68019,986,,,"Atlantic New Jersey US",25797.019001024008,1.4495949661124097`,
			want: w{[]CovidData{{-74.65848483, "New Jersey", "US"}}, ""}},
	}

	for _, tt := range tests {
		// t.Run enables running "subtests", one for each
		// table entry. These are shown separately
		// when executing `go test -v`.
		testname := fmt.Sprintf("%s", tt.testData)
		t.Run(testname, func(t *testing.T) {

			readUrlSvc = readUrlServiceMock{tt.testData}
			covidstats, err := readProcessData()

			// Fails for non-empty errors if error returned does not match with expected.
			if tt.want.errTxt != "" && err.Error() != tt.want.errTxt {
				t.Errorf("Expected error! Got %+v instead with %+v error", covidstats, err)
			}

			// Fails when no error if 1. Returned array size does not match. 2. Value does not match.
			if tt.want.errTxt == "" {
				if len(tt.want.data) != len(covidstats) {
					fmt.Println(covidstats)

					t.Errorf("Data sizes does not match!")
					return
				}
				for i, val := range tt.want.data {
					if covidstats[i] != val {
						t.Errorf("Result of %d:%+v does not match", i, covidstats[i])
					}
				}
			}
			readUrlSvc = readUrlServiceImpl{}
		})
	}
}
