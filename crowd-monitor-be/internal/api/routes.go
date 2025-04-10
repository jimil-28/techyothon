package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jimil-28/crowd-monitor/internal/api/handlers"
	"github.com/jimil-28/crowd-monitor/internal/api/middleware"
)

func SetupRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	videoAnalysisHandler *handlers.VideoAnalysisHandler,
	userHandler *handlers.UserHandler, // Add this parameter
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
		// Existing routes
		protected.GET("/video-analyses", videoAnalysisHandler.GetAllVideoAnalyses)
		protected.GET("/video-analyses/:videoId", videoAnalysisHandler.GetVideoAnalysisByID)
		protected.GET("/video-analyses/nearby", videoAnalysisHandler.GetNearbyVideoAnalyses)

		// New user routes
		protected.GET("/users", userHandler.GetAllUsers)
		protected.POST("/users", userHandler.AddUser)
	}
}
