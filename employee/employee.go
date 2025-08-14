package employees

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Employee struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Salary    int    `json:"salary"`
	HireDate  string `json:"hireDate"`
	Manager   string `json:"manager"`
}

// CreateEmployee erstellt einen neuen Mitarbeiter über die API
func CreateEmployee(employee Employee) (*Employee, error) {
	// Erstelle die POST-Daten
	data, err := json.Marshal(employee)
	if err != nil {
		return nil, fmt.Errorf("error marshalling employee data: %v", err)
	}

	// Sende die POST-Anfrage
	resp, err := http.Post("http://localhost:8080/employees", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error creating employee: %v", err)
	}
	defer resp.Body.Close()

	// Überprüfe, ob der Status OK ist (ändere dies auf 200 OK)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create employee, server responded with status: %s", resp.Status)
	}

	var createdEmployee Employee
	err = json.NewDecoder(resp.Body).Decode(&createdEmployee)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &createdEmployee, nil
}

func GetEmployeeIDs() ([]int, error) {
	resp, err := http.Get("http://localhost:8080/employees/allIds")
	if err != nil {
		return nil, fmt.Errorf("error fetching employee IDs: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch employee IDs, server responded with status: %s", resp.Status)
	}

	var ids []int
	err = json.NewDecoder(resp.Body).Decode(&ids)
	if err != nil {
		return nil, fmt.Errorf("error decoding employee IDs response: %v", err)
	}

	return ids, nil
}

func DeleteEmployee(employeeID int) error {
	url := fmt.Sprintf("http://localhost:8080/employees/%d", employeeID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("error creating DELETE request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending DELETE request: %v", err)
	}
	defer resp.Body.Close()

	// Überprüfe, ob der Server eine erfolgreiche Antwort gibt
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete employee, server responded with status: %s", resp.Status)
	}

	log.Printf("Employee with ID %d was successfully deleted.", employeeID)
	return nil
}

func RandomAction() error {
	employee := Employee{
		FirstName: fmt.Sprintf("John%d", rand.Intn(100)),
		LastName:  fmt.Sprintf("Doe%d", rand.Intn(100)),
		Email:     fmt.Sprintf("john.doe%d@example.com", rand.Intn(100)),
		Salary:    rand.Intn(10000) + 5000,
		HireDate:  time.Now().Format("2006-01-02"),
		Manager:   "",
	}

	// Erstelle den Mitarbeiter
	createdEmployee, err := CreateEmployee(employee)
	if err != nil {
		return fmt.Errorf("Error creating employee: %v", err)
	}

	log.Printf("Created employee: %+v", createdEmployee)

	ids, err := GetEmployeeIDs()
	if err != nil {
		return fmt.Errorf("Error fetching employee IDs: %v", err)
	}

	if len(ids) == 0 {
		return fmt.Errorf("No employees found to delete.")
	}

	randomID := ids[rand.Intn(len(ids))]

	err = DeleteEmployee(randomID)
	if err != nil {
		return fmt.Errorf("Error deleting employee: %v", err)
	}

	waitTime := time.Duration(rand.Intn(21)+10) * time.Second
	time.Sleep(waitTime)

	return nil
}
