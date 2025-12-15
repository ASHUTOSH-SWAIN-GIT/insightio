package main

import (
	"log"
	"time"

	insightio "github.com/ASHUTOSH-SWAIN-GIT/insightio_sdk"
)

func main() {
	// Create a new client
	// Note: The API key must match one of the keys configured in INSIGHTIO_API_KEY
	// Set the environment variable before starting the server, e.g.:
	// export INSIGHTIO_API_KEY="test-api-key-123"
	client, err := insightio.NewClient(insightio.ClientConfig{
		APIKey:  "test-api-key-123", // Must match a key in INSIGHTIO_API_KEY
		BaseURL: "http://localhost:8080",
	})
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	log.Println("=== InsightIO SDK Test ===")
	log.Println("Make sure the API server is running on http://localhost:8080")
	log.Println("and the gRPC backend is running on localhost:50051")

	// Example 1: Simple event
	log.Println("1. Sending simple page_view event...")
	event1 := insightio.Event{
		Type: "page_view",
	}
	resp, err := client.SendEvent(event1)
	if err != nil {
		log.Printf("   [ERROR] Error: %v\n", err)
	} else {
		log.Printf("   [OK] Success: %s - %s\n", resp.Status, resp.Message)
	}

	time.Sleep(500 * time.Millisecond)

	// Example 2: Event with user ID and value
	log.Println("\n2. Sending purchase event...")
	event2 := insightio.Event{
		Type:   "purchase",
		UserId: "user123",
		Value:  99.99,
		Metadata: map[string]string{
			"product_id": "prod_456",
			"currency":   "USD",
		},
	}
	resp, err = client.SendEvent(event2)
	if err != nil {
		log.Printf("   [ERROR] Error: %v\n", err)
	} else {
		log.Printf("   [OK] Success: %s - %s\n", resp.Status, resp.Message)
	}

	time.Sleep(500 * time.Millisecond)

	// Example 3: Event with metadata
	log.Println("\n3. Sending button_click event...")
	event3 := insightio.Event{
		Type:   "button_click",
		UserId: "user123",
		Metadata: map[string]string{
			"button_id": "checkout",
			"page":      "cart",
			"source":    "web",
		},
	}
	resp, err = client.SendEvent(event3)
	if err != nil {
		log.Printf("   [ERROR] Error: %v\n", err)
	} else {
		log.Printf("   [OK] Success: %s - %s\n", resp.Status, resp.Message)
	}

	time.Sleep(500 * time.Millisecond)

	// Example 4: Multiple events in sequence
	log.Println("\n4. Sending multiple events in sequence...")
	events := []insightio.Event{
		{Type: "login", UserId: "user123"},
		{Type: "search", UserId: "user123", Metadata: map[string]string{"query": "golang"}},
		{Type: "view_item", UserId: "user123", Metadata: map[string]string{"item_id": "item_789"}},
	}

	successCount := 0
	for i, event := range events {
		resp, err := client.SendEvent(event)
		if err != nil {
			log.Printf("   Event %d: [FAIL] Failed - %v\n", i+1, err)
		} else {
			log.Printf("   Event %d: [OK] %s\n", i+1, resp.Status)
			successCount++
		}
		time.Sleep(200 * time.Millisecond)
	}

	log.Printf("\n=== Test Complete: %d/%d events sent successfully ===\n", successCount, len(events))
}
