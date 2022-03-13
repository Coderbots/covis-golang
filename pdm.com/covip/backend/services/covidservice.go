package services

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"
	//	"encoding/json"
	"errors"
	"log"
	"pdm.com/covip/backend/model"
	"sort"
)

const path = "githubfiles/covid_data"

type CovidData struct {
	Confirmed     float64 `json:"confirmed"`
	ProvinceState string  `json:"province_state,omitempty"`
	CountryRegion string  `json:"country_region"`
}

func processData(text string) []CovidData {
	log.Println("In processData function")

	reader := csv.NewReader(strings.NewReader(text))
	record, err := reader.Read()
	if err != nil {
		log.Println("Error encountered reading first line of csv:", err)
		return []CovidData{}
	}

	var data []CovidData
	for {
		record, err = reader.Read()
		if err == io.EOF {
			log.Println("EOF!")
			break
		}
		if err == nil {
			confirmedCount, errParse := strconv.ParseFloat(record[6], 64)
			if errParse != nil {
				confirmedCount = 0
			}
			data = append(data, CovidData{
				Confirmed:     confirmedCount,
				ProvinceState: record[2],
				CountryRegion: record[3],
			})
		} else {
			log.Println("Error encountered in reading csv:", err)
			return data
		}
	}
	return data
}

func readProcessData() ([]CovidData, error) {
	log.Println("In readProcessData function")

	text := ReadRawData()
	if text == "" {
		return []CovidData{}, model.WrapError(errors.New("No data available for this day"), model.ErrNotFound)
	}

	processedData := processData(text)
	if processedData == nil {
		log.Println("Processed Data is empty")
		return []CovidData{}, errors.New("Internal error encountered")
	}
	return processedData, nil
}

func createSummary() ([]CovidData, error) {
	log.Println("In createSummary function")

	processedData, errProcessData := readProcessData()
	if errProcessData != nil {
		return processedData, errProcessData
	}

	var sortedArr []CovidData
	found := 0
	var arr1 []CovidData

	for i := 0; i < len(processedData); i++ {
		for j := 0; j < len(arr1); j++ {
			if processedData[i].CountryRegion == arr1[j].CountryRegion {
				arr1[j].Confirmed += processedData[i].Confirmed
				found = 1
				break
			}
		}
		if found == 0 {
			arr1 = append(arr1, CovidData{
				CountryRegion: processedData[i].CountryRegion,
				Confirmed:     processedData[i].Confirmed,
			})
		}
		found = 0
	}
	sortedArr = arr1
	sort.Slice(sortedArr, func(i, j int) bool {
		return sortedArr[i].Confirmed > sortedArr[j].Confirmed
	})
	return sortedArr, nil
}

func GetSummary() ([]CovidData, error) {
	log.Println("In GetSummary function")
	summary, err := createSummary()
	return summary, err
}

func GetCountryCases(name string) ([]CovidData, error) {
	log.Println("In GetCountryCases function")

	processedData, errProcessData := readProcessData()
	if errProcessData != nil {
		return processedData, errProcessData
	}

	var countryCase []CovidData
	for i := 0; i < len(processedData); i++ {
		if processedData[i].CountryRegion == name {
			countryCase = append(countryCase, CovidData{
				CountryRegion: processedData[i].CountryRegion,
				ProvinceState: processedData[i].ProvinceState,
				Confirmed:     processedData[i].Confirmed,
			})
		}
	}
	return countryCase, nil
}
