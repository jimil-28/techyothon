package models

type User struct {
	PhoneNumber string `json:"phone_number" firestore:"phone_number"`
	Name        string `json:"name" firestore:"name"`
	Rank        string `json:"rank" firestore:"rank"`
	Department  string `json:"department" firestore:"department"`
	IDCardNumber string `json:"id_card_number" firestore:"id_card_number"`
}

type OTPRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type OTPVerifyRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	OTPCode     string `json:"otp_code" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}