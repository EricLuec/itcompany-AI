package employees

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

type Employee struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	HireDate  string `json:"hireDate"`
	Manager   string `json:"manager"`
	Salary    int    `json:"salary"`
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

	var data RandomUserResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", "", err
	}

	firstName := data.Results[0].Name.First
	lastName := data.Results[0].Name.Last

	return firstName, lastName, nil
}

func generateRandomSalary() int {
	return rand.Intn(1501) + 500
}

func GenerateEmployee() (*Employee, error) {
	firstName, lastName, err := getRandomUser()
	if err != nil {
		return nil, err
	}

	hireDate := time.Now().Format("2006-01-02")

	salary := generateRandomSalary()

	email := fmt.Sprintf("%s.%s@gmail.com", firstName, lastName)

	employee := &Employee{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		HireDate:  hireDate,
		Manager:   "",
		Salary:    salary,
	}

	return employee, nil
}

func PostEmployee(employee *Employee) error {
	url := "http://localhost:8080/employees"
	jsonData, err := json.Marshal(employee)
	if err != nil {
		return err
	}
	fmt.Printf("Sending request with body: %s\n", string(jsonData))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("Server response body:", string(body))
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to post employee: %s, response body: %s", resp.Status, string(body))
	}
	return nil
}
