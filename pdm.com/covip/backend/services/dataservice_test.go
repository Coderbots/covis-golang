package services

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"pdm.com/covip/backend/helpers"
	"testing"
)

type mockClient struct {
	msg        string
	statusCode int
	errorTxt   string
}

func init() {
	// Initializing AppConfig object since main.go is not loaded.
	if nil == helpers.AppConfig {
		helpers.AppConfig = &helpers.Config{helpers.CovidRepo{Url: "dummyurl"}}
	}
}

// Mocks http.Get method.
func (m mockClient) Get(url string) (*http.Response, error) {
	// create a new reader with msg.
	r := io.NopCloser(bytes.NewReader([]byte(m.msg)))
	var err error
	if "" != m.errorTxt {
		err = errors.New(m.errorTxt)
	}
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       r,
	}, err
}

func setUpTestDir() {
	if _, errDir := os.Stat(downloadDirPath); os.IsNotExist(errDir) {
		errMkdir := os.MkdirAll(downloadDirPath, 0777)
		if errMkdir != nil {
			log.Println("Err encountered:", errMkdir)
		}
	}

}

func tearDownTestDir() {
	err := os.RemoveAll(downloadDirPath)
	if err != nil {
		log.Fatal(err)
	}
}

func TestReadRawData(t *testing.T) {
	// Table of tests.
	var tests = []struct {
		respBody, errTxt string
		code             int
		want             string
	}{
		{"404!Not Found", "", 404, ""},
		{"200!Found", "", 200, "200!Found\n"},
		{"404!Not Found 2", "", 404, ""},
	}

	// Define download directory for tests.
	downloadDirName = "datafiles_test"
	downloadDirPath = filepath.Join(os.TempDir(), appName, downloadDirName)

	for _, tt := range tests {

		// Set up test directory.
		setUpTestDir()

		testname := fmt.Sprintf("%s,%s,%d", tt.respBody, tt.errTxt, tt.code)
		t.Run(testname, func(t *testing.T) {
			// Mock HTTP Client.
			client = mockClient{tt.respBody, tt.code, tt.errTxt}
			fileData := readUrlSvc.ReadRawData()

			// Fail test if ReadRawData function returns unexpected data.
			if fileData != tt.want {
				t.Errorf("Responses %+v does not match with %+v test!", fileData, testname)
			}
		})

		// Delete test download directory.
		tearDownTestDir()
	}

	// Set the download directory variables back to original.
	downloadDirName = "datafiles"
	downloadDirPath = filepath.Join(os.TempDir(), appName, downloadDirName)
}
