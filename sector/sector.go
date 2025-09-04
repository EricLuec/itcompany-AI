package sector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"itCompany-AI/logger"
	"net/http"
)

var sectorAPI = "http://localhost:8080/sectors"
var sectorAllIDsAPI = "http://localhost:8080/sectors/allIds"

type Sector struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SalaryClass string `json:"salaryClass"`
}

func getRandomSector() string {
	randomSectorName := logger.OllamaRequester("gib mir einen sektor der in einem informatiklastigen Unternehmen existiert. nur 1 wort")
	fmt.Println(randomSectorName)
	return randomSectorName
}

func GenerateSector() (*Sector, error) {
	sectorName := getRandomSector()
	sectorDesc := logger.OllamaRequester("write a description of a few words for" + sectorName + " answer only with the sentence")
	salaryClass := logger.OllamaRequester("choose a random letter from A to D and answer only with the letter")

	sector := &Sector{
		Name:        sectorName,
		Description: sectorDesc,
		SalaryClass: salaryClass,
	}

	return sector, nil
}

func PostSector(sector *Sector) error {
	data, err := json.Marshal(sector)
	if err != nil {
		return err
	}

	resp, err := http.Post(sectorAPI, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to post Sector: %s, body: %s", resp.Status, string(body))
	}
	return nil
}

func GetAllSectorIds() ([]int, error) {
	resp, err := http.Get(sectorAllIDsAPI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ids []int
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return nil, err
	}
	return ids, nil
}

func DeleteSector(id int) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%d", sectorAPI, id), nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete Sector: %s, body: %s", resp.Status, string(body))
	}
	return nil
}

func GetOneSector() (*Sector, error) {
	resp, err := http.Get(sectorAPI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sectors []Sector
	if err := json.NewDecoder(resp.Body).Decode(&sectors); err != nil {
		return nil, err
	}
	if len(sectors) == 0 {
		return nil, fmt.Errorf("no sectors found")
	}
	return &sectors[0], nil
}
