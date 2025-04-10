package models

import "time"

type Location struct {
	Latitude  float64 `firestore:"latitude"`
	Longitude float64 `firestore:"longitude"`
	Timestamp string  `firestore:"timestamp"` // Change to string
}

type Analysis struct {
	CrowdCount                    string   `json:"crowd_count"`
	CrowdLevel                    string   `json:"crowd_level"`
	CrowdPresent                  string   `json:"crowd_present"`
	IsPeakHour                    string   `json:"is_peak_hour"`
	PoliceInterventionRequired    string   `json:"police_intervention_required"`
	PoliceInterventionSuggestions []string `json:"police_intervention_suggestions"`
}

type VideoAnalysis struct {
	VideoID       string    `firestore:"video_id"`
	VideoDuration float64   `json:"video_duration"`
	Timestamp     time.Time `firestore:"timestamp"` // Change to time.Time
	CreatedAt     time.Time `json:"created_at"`     // Change to time.Time
	Location      Location  `firestore:"location"`
	Analysis      Analysis  `firestore:"analysis"`
	FrameURLs     []string  `firestore:"frame_urls"`
}
