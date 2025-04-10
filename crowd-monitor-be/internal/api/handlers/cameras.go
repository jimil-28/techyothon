package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jimil-28/crowd-monitor/internal/services/firebase"
)

type CamerasHandler struct {
	firebaseClient *firebase.Client
}

func NewCamerasHandler(firebaseClient *firebase.Client) *CamerasHandler {
	return &CamerasHandler{
		firebaseClient: firebaseClient,
	}
}

func (h *CamerasHandler) GetCamerasByLocation(c *gin.Context) {
	locationID := c.Param("locationId")
	if locationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "locationId is required"})
		return
	}

	cameras, err := h.firebaseClient.GetCamerasByLocationID(c, locationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cameras)
}