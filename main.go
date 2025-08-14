package main

import (
	"fmt"
	employees "itCompany-AI/employee"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	for {
		employee, err := employees.GenerateEmployee()
		if err != nil {
			log.Printf("Error generating employee: %v", err)
			continue
		}

		err = employees.PostEmployee(employee)
		if err != nil {
			log.Printf("Error posting employee: %v", err)
		} else {
			fmt.Printf("Successfully posted employee: %v\n", employee)
		}

		time.Sleep(time.Duration(rand.Intn(6)+5) * time.Second)
	}
}
