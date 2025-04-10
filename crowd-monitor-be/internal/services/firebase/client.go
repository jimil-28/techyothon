package firebase

import (
	"context"
	"errors"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/db"
	"cloud.google.com/go/firestore"
	"github.com/jimil-28/crowd-monitor/internal/models"
	"google.golang.org/api/option"
)

type Client struct {
	app      *firebase.App
	auth     *auth.Client
	firestore *firestore.Client
	database *db.Client
}

func NewFirebaseClient(credentialsPath, databaseURL string) (*Client, error) {
	if credentialsPath == "" {
		return nil, errors.New("firebase credentials path is required")
	}

	opt := option.WithCredentialsFile(credentialsPath)
	config := &firebase.Config{
		DatabaseURL: databaseURL,
	}

	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return nil, err
	}

	auth, err := app.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	firestore, err := app.Firestore(context.Background())
	if err != nil {
		return nil, err
	}

	database, err := app.Database(context.Background())
	if err != nil {
		log.Printf("Warning: Could not initialize Realtime Database: %v", err)
		// Continue without realtime database
	}

	return &Client{
		app:      app,
		auth:     auth,
		firestore: firestore,
		database: database,
	}, nil
}

func (c *Client) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	query := c.firestore.Collection("users").Where("phone_number", "==", phoneNumber).Limit(1)
	
	iter := query.Documents(ctx)
	doc, err := iter.Next()
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Client) SaveUser(ctx context.Context, user models.User) error {
	_, err := c.firestore.Collection("users").Doc(user.PhoneNumber).Set(ctx, user)
	return err
}

func (c *Client) GetAllLocations(ctx context.Context) ([]models.Location, error) {
	query := c.firestore.Collection("locations")
	iter := query.Documents(ctx)
	
	var locations []models.Location
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		
		var location models.Location
		if err := doc.DataTo(&location); err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}
	
	return locations, nil
}

func (c *Client) GetCamerasByLocationID(ctx context.Context, locationID string) ([]models.Camera, error) {
	query := c.firestore.Collection("cameras").Where("location_id", "==", locationID)
	iter := query.Documents(ctx)
	
	var cameras []models.Camera
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		
		var camera models.Camera
		if err := doc.DataTo(&camera); err != nil {
			return nil, err
		}
		cameras = append(cameras, camera)
	}
	
	return cameras, nil
}

func (c *Client) Close() {
	if c.firestore != nil {
		c.firestore.Close()
	}
}