package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jimil-28/crowd-monitor/internal/models"
	"github.com/jimil-28/crowd-monitor/internal/services/firebase"
)

type LocationsHandler struct {
	firebaseClient *firebase.Client
}

func NewLocationsHandler(firebaseClient *firebase.Client) *LocationsHandler {
	return &LocationsHandler{
		firebaseClient: firebaseClient,
	}
}

func (h *LocationsHandler) GetAllLocations(c *gin.Context) {
	locations, err := h.firebaseClient.GetAllLocations(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var locationResponses []models.LocationResponse
	for _, location := range locations {
		cameras, err := h.firebaseClient.GetCamerasByLocationID(c, location.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Get intervention suggestions from the first camera (they're the same for all cameras in the same location)
		var suggestions []string
		if len(cameras) > 0 {
			suggestions = cameras[0].InterventionSuggestions
		}

		locationResponse := models.LocationResponse{
			ID:            location.ID,
			Name:          location.Name,
			IsOvercrowded: location.IsOvercrowded,
			Cameras:       cameras,
			Suggestions:   suggestions,
		}
		
		locationResponses = append(locationResponses, locationResponse)
	}

	c.JSON(http.StatusOK, locationResponses)
}