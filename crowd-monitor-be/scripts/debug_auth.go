package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8080/api/v1"

type OTPRequest struct {
	PhoneNumber string `json:"phone_number"`
}

type OTPVerifyRequest struct {
	PhoneNumber string `json:"phone_number"`
	OTPCode     string `json:"otp_code"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

func main() {
	fmt.Println("Debugging Auth Flow...")
	
	// Test phone number (use one from your dummy data)
	phoneNumber := "+919175045787"
	
	// Step 1: Send OTP (in development mode this might be bypassed)
	fmt.Println("\n1. Sending OTP to", phoneNumber)
	otpReq := OTPRequest{PhoneNumber: phoneNumber}
	resp, err := postJSON(baseURL+"/auth/send-otp", otpReq)
	if err != nil {
		fmt.Printf("Error sending OTP: %v\n", err)
		return
	}
	fmt.Println("Response:", string(resp))
	
	// For testing, we can use any OTP code since we've bypassed actual verification
	otpCode := "123456"
	fmt.Printf("\nUsing test OTP code: %s\n", otpCode)
	
	// Step 2: Verify OTP
	fmt.Println("\n2. Verifying OTP")
	verifyReq := OTPVerifyRequest{
		PhoneNumber: phoneNumber,
		OTPCode:     otpCode,
	}
	resp, err = postJSON(baseURL+"/auth/verify-otp", verifyReq)
	if err != nil {
		fmt.Printf("Error verifying OTP: %v\n", err)
		return
	}
	
	fmt.Println("Raw verify response:", string(resp))
	
	// Parse the token
	var authResp AuthResponse
	err = json.Unmarshal(resp, &authResp)
	if err != nil {
		fmt.Printf("Error parsing auth response: %v\n", err)
		return
	}
	
	token := authResp.Token
	fmt.Println("\nReceived Token:", token)
	
	// Step 3: Test protected endpoint
	fmt.Println("\n3. Testing protected endpoint with token")
	resp, err = getWithAuth(baseURL+"/locations", token)
	if err != nil {
		fmt.Printf("Error accessing protected endpoint: %v\n", err)
		
		// Additional debug: Try sending the raw request
		fmt.Println("\nAttempting raw request with curl-like headers...")
		req, _ := http.NewRequest("GET", baseURL+"/locations", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Raw request failed: %v\n", err)
			return
		}
		
		fmt.Printf("Raw response status: %s\n", resp.Status)
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Raw response body: %s\n", string(body))
		
		return
	}
	
	fmt.Println("Success! Protected endpoint response:", string(resp))
}

// Helper function to make a POST request with JSON data
func postJSON(url string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s (status code %d)", body, resp.StatusCode)
	}
	
	return body, nil
}

// Helper function to make a GET request with authentication
func getWithAuth(url, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer "+token)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s (status code %d)", body, resp.StatusCode)
	}
	
	return body, nil
}