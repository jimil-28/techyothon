package handlers

import (
    "log"
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/jimil-28/crowd-monitor/internal/services/firebase"
    "github.com/jimil-28/crowd-monitor/internal/models"
    "github.com/jimil-28/crowd-monitor/internal/utils"  // Add this import
    "fmt"
)

type VideoAnalysisHandler struct {
    firebaseClient *firebase.Client
}

func NewVideoAnalysisHandler(firebaseClient *firebase.Client) *VideoAnalysisHandler {
    return &VideoAnalysisHandler{
        firebaseClient: firebaseClient,
    }
}

func (h *VideoAnalysisHandler) GetAllVideoAnalyses(c *gin.Context) {
    log.Println("Fetching video analyses...")
    analyses, err := h.firebaseClient.GetAllVideoAnalyses(c)
    if (err != nil) {
        log.Printf("Error fetching analyses: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    log.Printf("Found %d analyses", len(analyses))
    c.JSON(http.StatusOK, analyses)
}

func (h *VideoAnalysisHandler) GetVideoAnalysisByID(c *gin.Context) {
    videoID := c.Param("videoId")
    if videoID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "videoId is required"})
        return
    }

    analysis, err := h.firebaseClient.GetVideoAnalysisByID(c, videoID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, analysis)
}

func (h *VideoAnalysisHandler) GetNearbyVideoAnalyses(c *gin.Context) {
    utils.Logger.Printf("Starting nearby video analyses search...")
    
    lat, err := strconv.ParseFloat(c.Query("latitude"), 64)
    if err != nil {
        utils.Logger.Printf("Invalid latitude parameter: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error": "invalid latitude parameter",
            "details": err.Error(),
        })
        return
    }

    lon, err := strconv.ParseFloat(c.Query("longitude"), 64)
    if err != nil {
        utils.Logger.Printf("Invalid longitude parameter: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error": "invalid longitude parameter",
            "details": err.Error(),
        })
        return
    }

    utils.Logger.Printf("Searching for videos near coordinates: %f, %f", lat, lon)
    const radiusKm = 10.0

    analyses, err := h.firebaseClient.GetVideoAnalysesNearby(c, lat, lon, radiusKm)
    if err != nil {
        utils.Logger.Printf("Error fetching nearby analyses: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error": "failed to fetch nearby video analyses",
            "details": err.Error(),
        })
        return
    }

    // Return empty array instead of null when no results found
    if analyses == nil {
        analyses = []models.VideoAnalysis{}
    }

    utils.Logger.Printf("Found %d video analyses within %0.1f km", len(analyses), radiusKm)
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data": analyses,
        "message": fmt.Sprintf("Found %d video analyses within %0.1f km", len(analyses), radiusKm),
    })
}