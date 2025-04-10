package models

type Camera struct {
	ID                    string   `json:"id" firestore:"id"`
	LocationID            string   `json:"location_id" firestore:"location_id"`
	CrowdCount            int      `json:"crowd_count" firestore:"crowd_count"`
	PoliceInterventionNeeded bool      `json:"police_intervention_needed" firestore:"police_intervention_needed"`
	InterventionSuggestions []string `json:"intervention_suggestions" firestore:"intervention_suggestions"`
}