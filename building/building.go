package building

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
	"time"
)

var buildingAPI = "http://localhost:8080/buildings"
var buildingAllIDsAPI = "http://localhost:8080/buildings/allIds"

type Building struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	BuildingDate string `json:"buildingDate"`
	Capacity     int    `json:"capacity"`
	City         string `json:"city"`
}

func getRandomBuildingName() string {
	randomBuildingName := logger.OllamaRequester("give me a nice name for a building, answer only with 1 name")
	fmt.Println(randomBuildingName)
	return randomBuildingName
}

func GenerateBuilding() (*Building, error) {
	buildingName := getRandomBuildingName()
	buildingDesc := logger.OllamaRequester("write a description of a few words for a new building answer only with the sentence")
	buildingDate := time.Now().Format("2006-01-02")
	buildingCity := logger.OllamaRequester("give me a random city and its country only the city and country and respond like this City, Country")
	capacity := rand.Intn(1500) + 500

	building := &Building{
		Name:         buildingName,
		Description:  buildingDesc,
		BuildingDate: buildingDate,
		Capacity:     capacity,
		City:         buildingCity,
	}

	return building, nil
}

func PostBuilding(building *Building) error {
	url := "http://localhost:8080/buildings"

	employeeJSON, err := json.Marshal(building)
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
		return fmt.Errorf("failed to post Building: %s, response body: %s", resp.Status, string(body))
	}

	log.Println("Building posted successfully")
	return nil
}

func GetAllBuildingIds() ([]int, error) {
	url := "http://localhost:8080/buildings/allIds"
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

func DeleteBuilding(id int) error {
	url := fmt.Sprintf("http://localhost:8080/buildings/%d", id)
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
	return fmt.Errorf("failed to delete building with ID %d: %d , response body: %s",
		id, resp.StatusCode, string(body))
}

func ExecuteRandomBuildingFunc() {
	action := rand.Intn(2)
	if action == 0 {
		building, err := GenerateBuilding()
		if err != nil {
			log.Printf("Error generating Building: %v", err)
			return
		}

		err = PostBuilding(building)
		if err != nil {
			log.Printf("Error posting building: %v", err)
			_ = logger.CreateLogEntry("building", fmt.Sprintf("Error posting building: %v", err))
		} else {
			msg := fmt.Sprintf("Successfully posted building: %s %s", building.Name, building.City)
			fmt.Println(msg)
			_ = logger.CreateLogEntry("building", msg)
		}

	} else {
		ids, err := GetAllBuildingIds()
		if err != nil {
			log.Printf("Error retrieving building IDs: %v", err)
			_ = logger.CreateLogEntry("building", fmt.Sprintf("Error retrieving building IDs: %v", err))
			return
		}

		if len(ids) > 0 {
			id := ids[rand.Intn(len(ids))]
			err := DeleteBuilding(id)
			if err != nil {
				log.Printf("Error deleting building with ID %d: %v", id, err)
				_ = logger.CreateLogEntry("building", fmt.Sprintf("Error deleting building with ID %d: %v", id, err))
			} else {
				msg := fmt.Sprintf("Successfully deleted building with ID %d", id)
				fmt.Println(msg)
				_ = logger.CreateLogEntry("building", msg)
			}
		} else {
			log.Println("No buildings available to delete.")
			_ = logger.CreateLogEntry("building", "No buildings available to delete.")
		}
	}
}
