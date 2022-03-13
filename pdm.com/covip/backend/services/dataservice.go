package services

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"pdm.com/covip/backend/helpers"
	"time"
)

var (
	appName         = "covis-golang"
	downloadDirName = "datafiles"
	downloadDirPath = filepath.Join(os.TempDir(), appName, downloadDirName)
)

func init() {
	if _, errDir := os.Stat(downloadDirPath); os.IsNotExist(errDir) {
		errMkdir := os.MkdirAll(downloadDirPath, 0777)
		if errMkdir != nil {
			fmt.Println("Err encountered:", errMkdir)
		}
	}
}

func download(url string, filePath string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, errIORead := io.ReadAll(resp.Body)
		if errIORead == nil {
			errHttpResponse := errors.New(string(b))
			return errHttpResponse
		}
		return errIORead
	}

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)

	return err
}

//ReadRawData reads from downloaded file and returns covid statistics.
func ReadRawData() string {
	fmt.Println("In Readrawdata function")

	url := helpers.AppConfig.CovidRepo.Url
	var downloadFilePath string

	//Check if data is available for last 4 days. Return the last available data.
	for i := 0; i <= 3; i++ {
		today := time.Now().AddDate(0, 0, -i).Format("01-02-2006")
		downloadUrl := url + today + ".csv"
		downloadFileName := "covid_data_" + today + ".csv"
		downloadFilePath = filepath.Join(downloadDirPath, downloadFileName)

		if _, errFile := os.Stat(downloadFilePath); os.IsNotExist(errFile) {
			fmt.Printf("File not present!Attempting to download %d time ...\n", i+1)
			errDownload := download(downloadUrl, downloadFilePath)
			if errDownload == nil {
				break
			} else if errDownload != nil && i == 3 {
				fmt.Println("Was unable to download file due to:", errDownload)
				return ""
			}
		} else {
			break
		}

	}

	fileHandler, err := os.Open(downloadFilePath)
	if err != nil {
		fmt.Println("Error encountered in opening csv file:", err)
	}

	scanner := bufio.NewScanner(fileHandler)
	var fileData string
	for scanner.Scan() {
		fileData += scanner.Text() + "\n"
	}
	return fileData
}
