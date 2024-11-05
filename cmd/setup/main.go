package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

// DeleteCollection deletes the collection if it exists to ensure idempotency
func deleteCollection() {
	url := "http://localhost:6333/collections/sample_collection"
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatalf("Error creating DELETE request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error deleting collection: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("Collection deletion response: %s\n", resp.Status)
}

// CreateCollection creates a new collection in Qdrant
func createCollection() {
	url := "http://localhost:6333/collections/sample_collection"
	collectionConfig := map[string]interface{}{
		"vectors": map[string]interface{}{
			"size":     128,      // Vector size
			"distance": "Cosine", // Distance metric
		},
	}

	// Convert the config to JSON
	jsonData, err := json.Marshal(collectionConfig)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// Send request
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating PUT request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error creating collection: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("Collection created with status: %s\n", resp.Status)
}

// InsertSampleData inserts sample data points into the collection
func insertSampleData() {
	url := "http://localhost:6333/collections/sample_collection/points"

	// Define sample points
	points := map[string]interface{}{
		"points": []map[string]interface{}{
			{
				"id":     1,
				"vector": generateRandomVector(128),
				"payload": map[string]interface{}{
					"title":       "Sample Document 1",
					"description": "This is a sample description for document 1.",
				},
			},
			{
				"id":     2,
				"vector": generateRandomVector(128),
				"payload": map[string]interface{}{
					"title":       "Sample Document 2",
					"description": "This is a sample description for document 2.",
				},
			},
			// Add more points as needed
		},
	}

	// Convert points to JSON
	jsonData, err := json.Marshal(points)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// Send request
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating PUT request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error inserting points: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("Points inserted with status: %s\n", resp.Status)
}

// GenerateRandomVector generates a random vector of specified size
func generateRandomVector(size int) []float32 {
	vector := make([]float32, size)
	for i := 0; i < size; i++ {
		vector[i] = rand.Float32()
	}
	return vector
}

func main() {
	deleteCollection() // Ensure the collection is deleted first
	createCollection() // Create a new collection
	insertSampleData() // Insert sample data into the collection
}
