package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	//	"encoding/json"
	"errors"
	"pdm.com/covip/backend/model"
	"sort"
)

const path = "githubfiles/covid_data"

type CovidData struct {
	Confirmed      float64 `json:"confirmed"`
	Province_state string  `json:"province_state,omitempty"`
	Country_region string  `json:"country_region"`
}

func processData(text string) []CovidData {
	fmt.Println("In processData function")
	reader := csv.NewReader(strings.NewReader(text))
	record, err := reader.Read()
	if err != nil {
		fmt.Println("Error encountered reading first line of csv:", err)
		return []CovidData{}
	}
	var data []CovidData
	for {
		record, err = reader.Read()
		if err == io.EOF {
			fmt.Println("EOF!")
			break
		}
		if err == nil {
			confirmed_int, err := strconv.ParseFloat(record[6], 64)
			if err != nil {
				confirmed_int = 0
			}
			data = append(data, CovidData{
				Confirmed:      confirmed_int,
				Province_state: record[2],
				Country_region: record[3],
			})
		} else {
			fmt.Println("Error encountered in reading csv:", err)
			return data
		}
	}
	return data
}

func createSummary() ([]CovidData, error) {
	fmt.Println("In createSummary function")
	text := ReadRawData()
	if text == "" {
		return []CovidData{}, model.WrapError(errors.New("No data available for this day"), model.ErrNotFound)
	}
	processedData := processData(text)
	if processedData == nil {
		fmt.Println("Processed Data is empty")
		return []CovidData{}, errors.New("Internal error encountered")
	}
	var sortedarr []CovidData
	found := 0
	var arr1 []CovidData
	for i := 0; i < len(processedData); i++ {
		for j := 0; j < len(arr1); j++ {
			if processedData[i].Country_region == arr1[j].Country_region {
				arr1[j].Confirmed += processedData[i].Confirmed
				found = 1
				break
			}
		}
		if found == 0 {
			arr1 = append(arr1, CovidData{
				Country_region: processedData[i].Country_region,
				Confirmed:      processedData[i].Confirmed,
			})
		}
		found = 0
	}
	sortedarr = arr1
	sort.Slice(sortedarr, func(i, j int) bool {
		return sortedarr[i].Confirmed > sortedarr[j].Confirmed
	})
	return sortedarr, nil
}

func GetSummary() ([]CovidData, error) {
	fmt.Println("In GetSummary function")
	summary, err := createSummary()
	return summary, err
}

func GetCountryCases(name string) ([]CovidData, error) {
	fmt.Println("In GetCountryCases function")
	text := ReadRawData()
	if text == "" {
		return []CovidData{}, model.WrapError(errors.New("No data available for this day"), model.ErrNotFound)
	}
	processedData := processData(text)
	if processedData == nil {
		fmt.Println("Processed Data is empty")
		return []CovidData{}, errors.New("Internal error encountered")
	}
	var countrycase []CovidData
	for i := 0; i < len(processedData); i++ {
		if processedData[i].Country_region == name {
			countrycase = append(countrycase, CovidData{
				Country_region: processedData[i].Country_region,
				Province_state: processedData[i].Province_state,
				Confirmed:      processedData[i].Confirmed,
			})
		}
	}
	return countrycase, nil
}
