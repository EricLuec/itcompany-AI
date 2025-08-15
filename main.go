package main

import (
	"fmt"
	employees "itCompany-AI/employee"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	for {
		action := rand.Intn(2)
		if action == 0 {
			employees.ExecuteRandomEmployeeFunc()
		} else {
			fmt.Println("no work")
		}

		time.Sleep(time.Duration(rand.Intn(6)+5) * time.Second)
	}
}
