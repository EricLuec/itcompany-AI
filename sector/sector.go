package sector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"itCompany-AI/logger"
	"log"
	"math/rand"
	"net/http"
)

type Sector struct {
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
	url := "http://localhost:8080/sectors"

	employeeJSON, err := json.Marshal(sector)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(employeeJSON))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to post Sector: %s, response body: %s", resp.Status, string(body))
	}

	log.Println("Sector posted successfully")
	return nil
}

func GetAllSectorIds() ([]int, error) {
	url := "http://localhost:8080/sectors/allIds"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ids []int
	err = json.Unmarshal(body, &ids)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func DeleteSector(id int) error {
	url := fmt.Sprintf("http://localhost:8080/sectors/%d", id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("failed to delete Sector with ID %d: %d , response body: %s",
		id, resp.StatusCode, string(body))
}

func GetOneSector() {
	resp, err := http.Get("http://localhost:8080/sectors")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var sectors []map[string]interface{}
	if err := json.Unmarshal(body, &sectors); err != nil {
		panic(err)
	}

	if len(sectors) == 0 {
		panic("no sectors found")
	}

	out, _ := json.MarshalIndent(sectors[0], "", "  ")
	fmt.Println(string(out))
}

func ExecuteRandomSectorFunc() {
	action := rand.Intn(2)
	if action == 0 {
		sector, err := GenerateSector()
		if err != nil {
			log.Printf("Error generating Sector: %v", err)
			return
		}

		err = PostSector(sector)
		if err != nil {
			log.Printf("Error posting Sector: %v", err)
			_ = logger.CreateLogEntry("Sector", fmt.Sprintf("Error posting Sector: %v", err))
		} else {
			msg := fmt.Sprintf("Successfully posted Sector: %s %s", sector.Name, sector.SalaryClass)
			fmt.Println(msg)
			_ = logger.CreateLogEntry("Sector", msg)
		}

	} else {
		ids, err := GetAllSectorIds()
		if err != nil {
			log.Printf("Error retrieving Sector IDs: %v", err)
			_ = logger.CreateLogEntry("sector", fmt.Sprintf("Error retrieving Sector IDs: %v", err))
			return
		}

		if len(ids) > 0 {
			id := ids[rand.Intn(len(ids))]
			err := DeleteSector(id)
			if err != nil {
				log.Printf("Error deleting Sector with ID %d: %v", id, err)
				_ = logger.CreateLogEntry("Sector", fmt.Sprintf("Error deleting Sector with ID %d: %v", id, err))
			} else {
				msg := fmt.Sprintf("Successfully deleted Sector with ID %d", id)
				fmt.Println(msg)
				_ = logger.CreateLogEntry("Sector", msg)
			}
		} else {
			log.Println("No Sector available to delete.")
			_ = logger.CreateLogEntry("Sector", "No Sector available to delete.")
		}
	}
}
