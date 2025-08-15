package logger

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func OllamaRequester(promptProp string) string {
	prompt := promptProp

	body, _ := json.Marshal(map[string]string{
		"model":  "gemma2:2b",
		"prompt": prompt,
	})

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	var fullAnswer string

	for scanner.Scan() {
		var chunk map[string]interface{}
		err := json.Unmarshal(scanner.Bytes(), &chunk)
		if err != nil {
			continue
		}

		if val, ok := chunk["response"].(string); ok {
			fullAnswer += val
		}

		if done, ok := chunk["done"].(bool); ok && done {
			break
		}
	}

	fmt.Println("Antwort vom Modell:\n", fullAnswer)
	return fullAnswer
}
