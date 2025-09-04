package employees

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"itCompany-AI/logger"
	sector2 "itCompany-AI/sector"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	employeeAPI       = "http://localhost:8080/employees"
	employeeAllIDsAPI = "http://localhost:8080/employees/allIds"

	getRandomUserFunc = getRandomUser
	getSectorFunc     = sector2.GetOneSector
)

type Sector struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SalaryClass string `json:"salaryClass"`
}

type Employee struct {
	Email     string         `json:"email"`
	FirstName string         `json:"firstName"`
	LastName  string         `json:"lastName"`
	HireDate  string         `json:"hireDate"`
	Manager   string         `json:"manager"`
	Salary    int            `json:"salary"`
	Sector    sector2.Sector `json:"sector"`
}

type RandomUserResponse struct {
	Results []struct {
		Name struct {
			First string `json:"first"`
			Last  string `json:"last"`
		} `json:"name"`
	} `json:"results"`
}

func getRandomUser() (string, string, error) {
	url := "https://randomuser.me/api?exc=dob,nat,email,location,login,gender,phone,cell,id,info,picture"
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	var userResponse RandomUserResponse
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return "", "", err
	}
	if len(userResponse.Results) > 0 {
		return userResponse.Results[0].Name.First, userResponse.Results[0].Name.Last, nil
	}
	return "", "", fmt.Errorf("no user data available")
}

func GenerateEmployee() (*Employee, error) {
	firstName, lastName, err := getRandomUserFunc()
	if err != nil {
		return nil, err
	}
	hireDate := time.Now().Format("2006-01-02")
	salary := rand.Intn(15000) + 500
	email := fmt.Sprintf("%s.%s@gmail.com", firstName, lastName)
	sector, err := getSectorFunc()
	if err != nil {
		return nil, err
	}
	return &Employee{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		HireDate:  hireDate,
		Manager:   "",
		Salary:    salary,
		Sector:    *sector,
	}, nil
}

func PostEmployee(employee *Employee) error {
	employeeJSON, err := json.Marshal(employee)
	if err != nil {
		return err
	}
	resp, err := http.Post(employeeAPI, "application/json", bytes.NewBuffer(employeeJSON))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to post employee: %s, response body: %s", resp.Status, string(body))
	}
	log.Println("Employee posted successfully")
	return nil
}

func GetAllEmployeeIDs() ([]int, error) {
	resp, err := http.Get(employeeAllIDsAPI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var ids []int
	if err := json.Unmarshal(body, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

func DeleteEmployee(id int) error {
	url := fmt.Sprintf("%s/%d", employeeAPI, id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("failed to delete employee with ID %d: %d, response body: %s", id, resp.StatusCode, string(body))
}

func ExecuteRandomEmployeeFunc() {
	if rand.Intn(2) == 0 {
		emp, err := GenerateEmployee()
		if err != nil {
			log.Printf("Error generating employee: %v", err)
			return
		}
		if err := PostEmployee(emp); err != nil {
			log.Printf("Error posting employee: %v", err)
			_ = logger.CreateLogEntry("employee", fmt.Sprintf("Error posting employee: %v", err))
		} else {
			msg := fmt.Sprintf("Successfully posted employee: %s %s", emp.FirstName, emp.LastName)
			fmt.Println(msg)
			_ = logger.CreateLogEntry("employee", msg)
		}
	} else {
		ids, err := GetAllEmployeeIDs()
		if err != nil {
			log.Printf("Error retrieving employee IDs: %v", err)
			_ = logger.CreateLogEntry("employee", fmt.Sprintf("Error retrieving employee IDs: %v", err))
			return
		}
		if len(ids) > 0 {
			id := ids[rand.Intn(len(ids))]
			if err := DeleteEmployee(id); err != nil {
				log.Printf("Error deleting employee with ID %d: %v", id, err)
				_ = logger.CreateLogEntry("employee", fmt.Sprintf("Error deleting employee with ID %d: %v", id, err))
			} else {
				msg := fmt.Sprintf("Successfully deleted employee with ID %d", id)
				fmt.Println(msg)
				_ = logger.CreateLogEntry("employee", msg)
			}
		} else {
			log.Println("No employees available to delete.")
			_ = logger.CreateLogEntry("employee", "No employees available to delete.")
		}
	}
}
