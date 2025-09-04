package main

import (
	employees "itCompany-AI/employee"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	/*
		for {
			action := rand.Intn(3)
			if action == 0 {
				employees.ExecuteRandomEmployeeFunc()
			} else if action == 1 {
				building.ExecuteRandomBuildingFunc()
			} else if action == 2 {
				sector.ExecuteRandomSectorFunc()
			}

			time.Sleep(time.Duration(rand.Intn(6)+5) * time.Second)
		}


	*/
	/*
		antwort := logger.OllamaRequester("hallo")
		fmt.Println(antwort)

	*/
	for {
		employees.ExecuteRandomEmployeeFunc()
		time.Sleep(time.Duration(rand.Intn(6)+5) * time.Second)

	}
}

// ollama run gemma2:2b
