package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

// User profile data
type User struct {
	PhoneNumber string `firestore:"phone_number"`
	Name        string `firestore:"name"`
	Rank        string `firestore:"rank"`
	Department  string `firestore:"department"`
	IDCardNumber string `firestore:"id_card_number"`
}

// Location data
type Location struct {
	ID            string   `firestore:"id"`
	Name          string   `firestore:"name"`
	IsOvercrowded bool     `firestore:"is_overcrowded"`
	CameraIDs     []string `firestore:"camera_ids"`
}

// Camera data
type Camera struct {
	ID                      string   `firestore:"id"`
	LocationID              string   `firestore:"location_id"`
	CrowdCount              int      `firestore:"crowd_count"`
	PoliceInterventionNeeded bool     `firestore:"police_intervention_needed"`
	InterventionSuggestions []string `firestore:"intervention_suggestions"`
}

func main() {
	// Initialize seed for randomization
	rand.Seed(time.Now().UnixNano())

	// Initialize Firebase app
	ctx := context.Background()
	opt := option.WithCredentialsFile("./firebase-credentials.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase app: %v", err)
	}

	// Get Firestore client
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	// Create users
	users := []User{
		{
			PhoneNumber: "+919175045787",
			Name:        "Rajesh Kumar",
			Rank:        "ASI",
			Department:  "Madgaon Police Department",
			IDCardNumber: "123",
		},
		{
			PhoneNumber: "+919130232897",
			Name:        "Priya Sharma",
			Rank:        "SI",
			Department:  "Vasco Police Department",
			IDCardNumber: "456",
		},
		{
			PhoneNumber: "+917708122103",
			Name:        "Amit Patel",
			Rank:        "PI",
			Department:  "Panjim Police Department",
			IDCardNumber: "789",
		},
		{
			PhoneNumber: "+919405061349",
			Name:        "Kavita Naik",
			Rank:        "DYSP",
			Department:  "Mapusa Police Department",
			IDCardNumber: "321",
		},
	}

	// Create locations
	locations := []Location{
		{
			ID:            "loc1",
			Name:          "Miramar Beach",
			IsOvercrowded: true,
			CameraIDs:     []string{"cam1", "cam2", "cam3"},
		},
		{
			ID:            "loc2",
			Name:          "Dona Paula",
			IsOvercrowded: true,
			CameraIDs:     []string{"cam4", "cam5", "cam6"},
		},
		{
			ID:            "loc3",
			Name:          "Colva Beach",
			IsOvercrowded: false,
			CameraIDs:     []string{"cam7", "cam8", "cam9"},
		},
		{
			ID:            "loc4",
			Name:          "Bodgini",
			IsOvercrowded: true,
			CameraIDs:     []string{"cam10", "cam11", "cam12"},
		},
		{
			ID:            "loc5",
			Name:          "Calangute Beach",
			IsOvercrowded: true,
			CameraIDs:     []string{"cam13", "cam14", "cam15"},
		},
		{
			ID:            "loc6",
			Name:          "Vasco Damodar Temple",
			IsOvercrowded: false,
			CameraIDs:     []string{"cam16", "cam17", "cam18"},
		},
	}

	// Intervention suggestions
	interventionSuggestions := []string{
		"Immediately stop any further crowd flow towards this congested area.",
		"Open or create multiple emergency exit routes away from the crush point to relieve pressure.",
		"Deploy emergency responders (police, medical) to the immediate area for assistance and upstream to manage flow control and diversions.",
		"Use clear, loud communication (e.g., loudspeakers) to instruct the crowd to stop pushing, remain calm if possible, and direct them towards available exit routes.",
		"Set up barricades to redirect crowd flow and prevent further congestion.",
		"Establish a command post near the affected area to coordinate response efforts.",
		"Initiate evacuation procedures for the most densely packed areas first.",
		"Dispatch additional personnel to manage crowd movement at key entry/exit points.",
		"Activate emergency protocols for rapid response team deployment.",
		"Implement traffic control measures in surrounding areas to facilitate emergency vehicle access.",
	}

	// Create cameras
	var cameras []Camera
	for _, loc := range locations {
		for _, camID := range loc.CameraIDs {
			// Select 3-5 random intervention suggestions
			numSuggestions := rand.Intn(3) + 3 // 3 to 5 suggestions
			suggestions := make([]string, numSuggestions)
			
			// Get random suggestions
			for i := 0; i < numSuggestions; i++ {
				suggestions[i] = interventionSuggestions[rand.Intn(len(interventionSuggestions))]
			}
			
			cameras = append(cameras, Camera{
				ID:                      camID,
				LocationID:              loc.ID,
				CrowdCount:              rand.Intn(400) + 50, // Random crowd between 50-450
				PoliceInterventionNeeded: true,
				InterventionSuggestions:  suggestions,
			})
		}
	}

	// Batch write users
	batch := client.Batch()
	for _, user := range users {
		ref := client.Collection("users").Doc(user.PhoneNumber)
		batch.Set(ref, user)
	}

	// Batch write locations
	for _, location := range locations {
		ref := client.Collection("locations").Doc(location.ID)
		batch.Set(ref, location)
	}

	// Batch write cameras
	for _, camera := range cameras {
		ref := client.Collection("cameras").Doc(camera.ID)
		batch.Set(ref, camera)
	}

	// Commit the batch
	_, err = batch.Commit(ctx)
	if err != nil {
		log.Fatalf("Failed to commit batch: %v", err)
	}

	fmt.Println("Successfully seeded the database!")
}