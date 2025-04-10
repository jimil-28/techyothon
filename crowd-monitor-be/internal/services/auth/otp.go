package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jimil-28/crowd-monitor/internal/models"
	"github.com/jimil-28/crowd-monitor/internal/services/firebase"
	"github.com/jimil-28/crowd-monitor/internal/services/twilio"
)

// Make sure this matches exactly the secret used in middleware/auth.go
var jwtSecret = []byte("08c06a7edfa777158e6024e7c2de3a1851fe29a50e48bfa1c3a126e5eb75a0fc")

type Service struct {
	twilioClient   *twilio.Client
	firebaseClient *firebase.Client
}

func NewAuthService(twilioClient *twilio.Client, firebaseClient *firebase.Client) *Service {
	return &Service{
		twilioClient:   twilioClient,
		firebaseClient: firebaseClient,
	}
}

func (s *Service) SendOTP(phoneNumber string) error {
	return s.twilioClient.SendOTP(phoneNumber)
}

func (s *Service) VerifyOTP(ctx context.Context, phoneNumber string, otpCode string) (*models.AuthResponse, error) {
	// For development: Skip actual OTP verification
	// In production, uncomment the verification code below
	
	verified, err := s.twilioClient.VerifyOTP(phoneNumber, otpCode)
	if err != nil {
		return nil, err
	}

	if !verified {
		return nil, fmt.Errorf("invalid OTP code")
	}
	

	// Fetch user from Firebase
	user, err := s.firebaseClient.GetUserByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	// Generate JWT token
	token, err := generateJWT(user.PhoneNumber)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func generateJWT(phoneNumber string) (string, error) {
	// Create a new token with claims
	token := jwt.New(jwt.SigningMethodHS256)
	
	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["phone_number"] = phoneNumber
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Expires in 24 hours
	
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}