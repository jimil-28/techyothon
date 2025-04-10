package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jimil-28/crowd-monitor/internal/api/handlers"
	"github.com/jimil-28/crowd-monitor/internal/api/middleware"
)

func SetupRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	locationsHandler *handlers.LocationsHandler,
	camerasHandler *handlers.CamerasHandler,
) {
	// Public routes
	public := router.Group("/api/v1")
	{
		public.POST("/auth/send-otp", authHandler.SendOTP)
		public.POST("/auth/verify-otp", authHandler.VerifyOTP)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/locations", locationsHandler.GetAllLocations)
		protected.GET("/locations/:locationId/cameras", camerasHandler.GetCamerasByLocation)
	}
}