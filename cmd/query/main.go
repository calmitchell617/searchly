package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"math/rand"
)

// GenerateRandomVector generates a random vector of the specified size
func generateRandomVector(size int) []float32 {
	// Seed the random number generator to get different results each time
	// Create a slice to hold the vector
	vector := make([]float32, size)
	for i := 0; i < size; i++ {
		vector[i] = rand.Float32() // Random float between 0.0 and 1.0
	}
	return vector
}

// SearchCollection searches the collection for points similar to the provided vector
func searchCollection() {
	url := "http://localhost:6333/collections/sample_collection/points/search"

	// Define the search parameters
	searchRequest := map[string]interface{}{
		"vector": generateRandomVector(128), // Use a random vector or specify your own
		"top":    1,                         // Number of closest points to retrieve
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(searchRequest)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// Send request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error searching collection: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Error decoding response: %v", err)
	}

	fmt.Printf("Search results: %+v\n", result)
}

func main() {
	searchCollection() // Call the search function
}
