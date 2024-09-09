package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type LlamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func Llama3dot1(question string) string {
	// === CONFIGURING BODY ===
	body := map[string]string{
		"model":  "llama3.1",
		"prompt": question,
	}

	// === JSONIFY BODY ===
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Fatal("Failed to Marshal body", err)
	}

	// === MAKE REQUEST ===
	apiURL := "http://localhost:11434/api/generate"
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatal("Failed to make request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Request failed with status code: %d", resp.StatusCode)
	}

	// === READ AND PROCESS RESPONSE STREAM ===
	var fullResponse string
	decoder := json.NewDecoder(resp.Body)

	for {
		var part LlamaResponse
		if err := decoder.Decode(&part); err == io.EOF {
			break // end of the response stream
		} else if err != nil {
			log.Fatal("Failed to decode response", err)
		}
		fullResponse += part.Response
		if part.Done {
			break
		}
	}

	return fullResponse
}

func main() {
	var question string = `Can octopuses regrow their limbs?`
	llama3dot1Answer := Llama3dot1(question)
	fmt.Println(llama3dot1Answer)
}
