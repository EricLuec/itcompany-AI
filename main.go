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
		action := rand.Intn(2) // 0 oder 1
		if action == 0 {
			employee, err := employees.GenerateEmployee()
			if err != nil {
				log.Printf("Error generating employee: %v", err)
				continue
			}

			err = employees.PostEmployee(employee)
			if err != nil {
				log.Printf("Error posting employee: %v", err)
				_ = logger.CreateLogEntry("employee", fmt.Sprintf("Error posting employee: %v", err))
			} else {
				msg := fmt.Sprintf("Successfully posted employee: %s %s", employee.FirstName, employee.LastName)
				fmt.Println(msg)
				_ = logger.CreateLogEntry("employee", msg)
			}

		} else {
			ids, err := employees.GetAllEmployeeIDs()
			if err != nil {
				log.Printf("Error retrieving employee IDs: %v", err)
				_ = logger.CreateLogEntry("employee", fmt.Sprintf("Error retrieving employee IDs: %v", err))
				continue
			}

			if len(ids) > 0 {
				id := ids[rand.Intn(len(ids))]
				err := employees.DeleteEmployee(id)
				if err != nil {
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

		time.Sleep(time.Duration(rand.Intn(6)+5) * time.Second)
	}
}
