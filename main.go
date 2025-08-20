package main

import (
	"itCompany-AI/building"
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
			building.ExecuteRandomBuildingFunc()
		}

		time.Sleep(time.Duration(rand.Intn(6)+5) * time.Second)
	}

	/*
		antwort := logger.OllamaRequester("hallo")
		fmt.Println(antwort)

	*/
}
