package models

import "time"

type Location struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
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
	VideoID       string    `json:"video_id"`
	VideoDuration float64   `json:"video_duration"`
	Timestamp     time.Time `json:"timestamp"`
	CreatedAt     time.Time `json:"created_at"`
	Location      Location  `json:"location"`
	Analysis      Analysis  `json:"analysis"`
	FrameURLs     []string  `json:"frame_urls"`
}
