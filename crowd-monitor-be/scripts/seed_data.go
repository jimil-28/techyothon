package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func init() {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatal("Failed to create logs directory:", err)
	}

	// Set up logging to file
	logFile := fmt.Sprintf("logs/seed_data_%s.log", time.Now().Format("2006-01-02_15-04-05"))
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	// Set log output to both file and console
	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("=== Seed Data Script Started ===")
}

// User profile data
type User struct {
	PhoneNumber  string `firestore:"phone_number"`
	Name         string `firestore:"name"`
	Rank         string `firestore:"rank"`
	Department   string `firestore:"department"`
	IDCardNumber string `firestore:"id_card_number"`
}

func main() {
	defer log.Println("=== Seed Data Script Ended ===")

	// Initialize Firebase app
	log.Println("Initializing Firebase app...")
	ctx := context.Background()
	opt := option.WithCredentialsFile("./firebase-credentials.json")

	// Log Firebase configuration details
	log.Printf("Using credentials file: %s", "./firebase-credentials.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Printf("ERROR: Failed to initialize Firebase app: %v", err)
		fmt.Printf("Status: %d - Internal Server Error\n", http.StatusInternalServerError)
		return
	}
	log.Println("Firebase app initialized successfully")

	// Get Firestore client
	log.Println("Creating Firestore client...")
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Printf("ERROR: Failed to create Firestore client: %v", err)
		fmt.Printf("Status: %d - Internal Server Error\n", http.StatusInternalServerError)
		return
	}
	defer client.Close()
	log.Println("Firestore client created successfully")

	// Create users
	users := []User{
		{
			PhoneNumber:  "+919175045787",
			Name:         "Rajesh Kumar",
			Rank:         "ASI",
			Department:   "Madgaon Police Department",
			IDCardNumber: "123",
		},
		{
			PhoneNumber:  "+919130232897",
			Name:         "Priya Sharma",
			Rank:         "SI",
			Department:   "Vasco Police Department",
			IDCardNumber: "456",
		},
		{
			PhoneNumber:  "+917708122103",
			Name:         "Amit Patel",
			Rank:         "PI",
			Department:   "Panjim Police Department",
			IDCardNumber: "789",
		},
		{
			PhoneNumber:  "+919405061349",
			Name:         "Kavita Naik",
			Rank:         "DYSP",
			Department:   "Mapusa Police Department",
			IDCardNumber: "321",
		},
	}
	log.Printf("Preparing to seed %d users...", len(users))

	// Batch write users
	log.Println("Creating batch write operation...")
	batch := client.Batch()
	for _, user := range users {
		ref := client.Collection("users").Doc(user.PhoneNumber)
		batch.Set(ref, user)
		log.Printf("Added user to batch: %s (%s)", user.Name, user.PhoneNumber)
	}

	// Commit the batch and check response
	log.Println("Committing batch to Firestore...")
	result, err := batch.Commit(ctx)
	if err != nil {
		log.Printf("ERROR: Failed to seed data: %v", err)
		fmt.Printf("Status: %d - Internal Server Error\n", http.StatusInternalServerError)
		return
	}

	if len(result) == len(users) {
		log.Printf("SUCCESS: Successfully seeded all %d users", len(users))
		fmt.Printf("Status: %d - Success\n", http.StatusCreated)
		fmt.Printf("Successfully seeded %d users to the database!\n", len(users))
	} else {
		log.Printf("WARNING: Partial success - Only %d out of %d users were seeded", len(result), len(users))
		fmt.Printf("Status: %d - Partial Success\n", http.StatusPartialContent)
		fmt.Printf("Only %d out of %d users were seeded\n", len(result), len(users))
	}
}
