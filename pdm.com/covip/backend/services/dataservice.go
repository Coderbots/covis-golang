package services

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func download(url string, filepath string) error {

	resp, err := http.Get(url)
	//fmt.Println("Response from download", resp)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func Readrawdata() string {
	url := "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_daily_reports/05-31-2021.csv"
	pathdir := "githubfiles"

	errmkdir := os.MkdirAll(pathdir, 0777)
	if errmkdir != nil {
		fmt.Println("Err encountered:", errmkdir)
	}

	//path := "githubfiles/covid_data"
	today := time.Now().Format("01-02-2006")
	path_current := filepath.Join(pathdir, "covid_data_"+today+".csv")
	//path_current := path + "_" + today + ".csv"
	fmt.Println("In Readrawdata function")
	//filehandler, err := os.Open(path_current)

	err := download(url, path_current)
	if err != nil {
		fmt.Println("Was unable to download file due to:", err)
		return ""
	}

	filehandler, err := os.Open(path_current)
	if err != nil {
		fmt.Println("Error encountered in opening csv file:", err)
	}
	scanner := bufio.NewScanner(filehandler)
	var filedata string
	for scanner.Scan() {
		filedata += scanner.Text() + "\n"
	}
	//fmt.Println("filedata from Readrawdata ready", filedata)
	return filedata
}
