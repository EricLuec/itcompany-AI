package main

import (
	"fmt"
	employees "itCompany-AI/employee"
	"itCompany-AI/logger"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	for {
		actionType := rand.Intn(2)

		var err error
		var actionDesc string
		var logCategory string

		if actionType == 0 {
			err = employees.RandomAction()
			actionDesc = "Created a new employee"
			logCategory = "employees"
		} else {
			ids, fetchErr := employees.GetEmployeeIDs()
			if fetchErr != nil || len(ids) == 0 {
				log.Printf("No employees to delete: %v", fetchErr)
				continue
			}

			randomID := ids[rand.Intn(len(ids))]
			err = employees.DeleteEmployee(randomID)
			actionDesc = fmt.Sprintf("Deleted employee with ID %d", randomID)
			logCategory = "employees"
		}

		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Println("Action completed successfully.")
		}

		if err == nil {
			err := logger.CreateLogEntry(logCategory, actionDesc, 0)
			if err != nil {
				log.Printf("Error creating log entry: %v", err)
			}
		}

		waitTime := time.Duration(rand.Intn(21)+10) * time.Second
		time.Sleep(waitTime)
	}
}
