package models

type Location struct {
	ID            string   `json:"id" firestore:"id"`
	Name          string   `json:"name" firestore:"name"`
	IsOvercrowded bool     `json:"is_overcrowded" firestore:"is_overcrowded"`
	CameraIDs     []string `json:"camera_ids" firestore:"camera_ids"`
}

type LocationResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	IsOvercrowded bool      `json:"is_overcrowded"`
	Cameras       []Camera  `json:"cameras"`
	Suggestions   []string  `json:"intervention_suggestions"`
}