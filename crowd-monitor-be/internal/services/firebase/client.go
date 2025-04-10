package firebase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/db"
	"github.com/jimil-28/crowd-monitor/internal/models"
	"github.com/jimil-28/crowd-monitor/internal/utils"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Client struct {
	app       *firebase.App
	auth      *auth.Client
	firestore *firestore.Client
	database  *db.Client
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
		app:       app,
		auth:      auth,
		firestore: firestore,
		database:  database,
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

func (c *Client) GetAllVideoAnalyses(ctx context.Context) ([]models.VideoAnalysis, error) {
	log.Println("Starting Firestore query...")
	query := c.firestore.Collection("video-analysis")
	iter := query.Documents(ctx)

	var analyses []models.VideoAnalysis
	docCount := 0

	for {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				log.Printf("Query completed. Found %d documents", docCount)
				break
			}
			return nil, err
		}
		docCount++

		data := doc.Data()
		var analysis models.VideoAnalysis

		// Handle basic fields
		analysis.VideoID = data["video_id"].(string)
		analysis.VideoDuration = data["video_duration"].(float64)

		// Handle timestamps
		if ts, ok := data["timestamp"].(string); ok {
			analysis.Timestamp, _ = time.Parse(time.RFC3339Nano, ts)
		}
		if ts, ok := data["created_at"].(string); ok {
			analysis.CreatedAt, _ = time.Parse(time.RFC3339Nano, ts)
		}

		// Handle location
		if loc, ok := data["location"].(map[string]interface{}); ok {
			analysis.Location = models.Location{
				Latitude:  loc["latitude"].(float64),
				Longitude: loc["longitude"].(float64),
			}
			if ts, ok := loc["timestamp"].(string); ok {
				analysis.Location.Timestamp, _ = time.Parse(time.RFC3339Nano, ts)
			}
		}

		// Handle analysis fields
		if analysisData, ok := data["analysis"].(map[string]interface{}); ok {
			analysis.Analysis = models.Analysis{
				CrowdCount:                 analysisData["crowd_count"].(string),
				CrowdLevel:                 analysisData["crowd_level"].(string),
				CrowdPresent:               analysisData["crowd_present"].(string),
				IsPeakHour:                 analysisData["is_peak_hour"].(string),
				PoliceInterventionRequired: analysisData["police_intervention_required"].(string),
			}
			if suggestions, ok := analysisData["police_intervention_suggestions"].([]interface{}); ok {
				for _, suggestion := range suggestions {
					analysis.Analysis.PoliceInterventionSuggestions = append(
						analysis.Analysis.PoliceInterventionSuggestions,
						suggestion.(string),
					)
				}
			}
		}

		// Handle frame URLs
		if urls, ok := data["frame_urls"].([]interface{}); ok {
			for _, url := range urls {
				analysis.FrameURLs = append(analysis.FrameURLs, url.(string))
			}
		}

		analyses = append(analyses, analysis)
	}

	return analyses, nil
}

func (c *Client) GetVideoAnalysisByID(ctx context.Context, videoID string) (*models.VideoAnalysis, error) {
	query := c.firestore.Collection("video-analysis").Where("video_id", "==", videoID)
	iter := query.Documents(ctx)

	doc, err := iter.Next()
	if err != nil {
		return nil, err
	}

	data := doc.Data()
	analysis := &models.VideoAnalysis{}

	// Use the same parsing logic as GetAllVideoAnalyses
	analysis.VideoID = data["video_id"].(string)
	analysis.VideoDuration = data["video_duration"].(float64)

	if ts, ok := data["timestamp"].(string); ok {
		analysis.Timestamp, _ = time.Parse(time.RFC3339Nano, ts)
	}
	if ts, ok := data["created_at"].(string); ok {
		analysis.CreatedAt, _ = time.Parse(time.RFC3339Nano, ts)
	}

	// ... rest of the parsing logic (same as above)

	return analysis, nil
}

func (c *Client) GetVideoAnalysesNearby(ctx context.Context, lat, lon float64, radiusKm float64) ([]models.VideoAnalysis, error) {
    utils.Logger.Printf("Fetching all video analyses from Firestore...")
    
    // Get collection reference
    collection := c.firestore.Collection("video-analysis")
    if collection == nil {
        utils.Logger.Printf("Error: video-analysis collection not found")
        return nil, fmt.Errorf("video-analysis collection not found")
    }

    // Get all documents
    docs, err := collection.Documents(ctx).GetAll()
    if err != nil {
        utils.Logger.Printf("Error fetching documents: %v", err)
        return []models.VideoAnalysis{}, err
    }

    utils.Logger.Printf("Found %d total documents", len(docs))
    var nearbyAnalyses []models.VideoAnalysis

    for _, doc := range docs {
        var analysis models.VideoAnalysis
        if err := doc.DataTo(&analysis); err != nil {
            utils.Logger.Printf("Error parsing document %s: %v", doc.Ref.ID, err)
            continue
        }

        distance := calculateDistance(
            lat, lon,
            analysis.Location.Latitude,
            analysis.Location.Longitude,
        )

        utils.Logger.Printf("Document %s is %.2f km away", analysis.VideoID, distance)
        
        if distance <= radiusKm {
            nearbyAnalyses = append(nearbyAnalyses, analysis)
            utils.Logger.Printf("Added document %s to results (within range)", analysis.VideoID)
        }
    }

    utils.Logger.Printf("Returning %d nearby analyses", len(nearbyAnalyses))
    return nearbyAnalyses, nil
}

// calculateDistance returns distance in kilometers between two points using Haversine formula
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
    const R = 6371 // Earth's radius in kilometers

    lat1Rad := lat1 * math.Pi / 180
    lat2Rad := lat2 * math.Pi / 180
    lon1Rad := lon1 * math.Pi / 180
    lon2Rad := lon2 * math.Pi / 180

    dlat := lat2Rad - lat1Rad
    dlon := lon2Rad - lon1Rad

    a := math.Sin(dlat/2)*math.Sin(dlat/2) +
        math.Cos(lat1Rad)*math.Cos(lat2Rad)*
            math.Sin(dlon/2)*math.Sin(dlon/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

    return R * c
}

func (c *Client) Close() {
	if c.firestore != nil {
		c.firestore.Close()
	}
}
