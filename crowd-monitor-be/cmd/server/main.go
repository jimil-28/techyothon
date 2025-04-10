package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jimil-28/crowd-monitor/config"
	"github.com/jimil-28/crowd-monitor/internal/api"
	"github.com/jimil-28/crowd-monitor/internal/api/handlers"
	"github.com/jimil-28/crowd-monitor/internal/services/auth"
	"github.com/jimil-28/crowd-monitor/internal/services/firebase"
	"github.com/jimil-28/crowd-monitor/internal/services/twilio"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Twilio client
	twilioClient, err := twilio.NewTwilioClient(
		cfg.TwilioAccountSid,
		cfg.TwilioAuthToken,
		cfg.TwilioServiceSid,
	)
	if err != nil {
		log.Fatalf("Failed to initialize Twilio client: %v", err)
	}

	// Initialize Firebase client
	firebaseClient, err := firebase.NewFirebaseClient(
		cfg.FirebaseCredPath,
		cfg.FirebaseDatabaseURL,
	)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase client: %v", err)
	}
	defer firebaseClient.Close()

	// Initialize authentication service
	authService := auth.NewAuthService(twilioClient, firebaseClient)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	locationsHandler := handlers.NewLocationsHandler(firebaseClient)
	camerasHandler := handlers.NewCamerasHandler(firebaseClient)

	// Setup Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	})

	// Setup routes
	api.SetupRoutes(router, authHandler, locationsHandler, camerasHandler)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		fmt.Printf("Server is running on port %s...\n", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}