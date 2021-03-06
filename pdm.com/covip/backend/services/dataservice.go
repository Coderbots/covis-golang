package services

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"pdm.com/covip/backend/helpers"
	"time"
)

type readUrlService interface {
	ReadRawData() string
}

type readUrlServiceImpl struct{}

type httpClient interface {
	Get(url string) (*http.Response, error)
}

var (
	appName                        = "covis-golang"
	downloadDirName                = "datafiles"
	downloadDirPath                = filepath.Join(os.TempDir(), appName, downloadDirName)
	readUrlSvc      readUrlService = readUrlServiceImpl{}
	client          httpClient
)

func init() {
	// Initialize HTTP Client.
	client = &http.Client{}
	//Create download directory for data files.
	if _, errDir := os.Stat(downloadDirPath); os.IsNotExist(errDir) {
		errMkdir := os.MkdirAll(downloadDirPath, 0777)
		if errMkdir != nil {
			log.Println("Err encountered:", errMkdir)
		}
	}
}

// download calls public covid repo and returns data.
func download(url string, filePath string) error {
	resp, err := client.Get(url)
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
func (service readUrlServiceImpl) ReadRawData() string {
	log.Println("In Readrawdata function")

	url := helpers.AppConfig.CovidRepo.Url

	var (
		downloadFilePath string
		errDownload      error
		fileDownloaded   = false
	)

	//Check if data is available for last 4 days. Return the last available data.
	for i := 0; i <= 3; i++ {
		today := time.Now().AddDate(0, 0, -i).Format("01-02-2006")
		downloadUrl := url + today + ".csv"
		downloadFileName := "covid_data_" + today + ".csv"
		downloadFilePath = filepath.Join(downloadDirPath, downloadFileName)
		if _, errFile := os.Stat(downloadFilePath); os.IsNotExist(errFile) {
			log.Printf("File not present!Attempting to download %d time ...\n", i+1)
			errDownload = download(downloadUrl, downloadFilePath)
			if errDownload == nil {
				fileDownloaded = true
				break
			}
		} else {
			fileDownloaded = true
			break
		}

	}

	if !fileDownloaded {
		log.Println("Was unable to download file due to:", errDownload)
		return ""
	}

	fileHandler, err := os.Open(downloadFilePath)
	if err != nil {
		log.Println("Error encountered in opening csv file:", err)
		return ""
	}

	scanner := bufio.NewScanner(fileHandler)
	var fileData string
	for scanner.Scan() {
		// Add back new line removed by scanner.
		fileData += scanner.Text() + "\n"
	}
	return fileData
}
