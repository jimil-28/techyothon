package twilio

import (
	"errors"
	"fmt"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

type Client struct {
	twilioClient *twilio.RestClient
	serviceSid   string
}

func NewTwilioClient(accountSid, authToken, serviceSid string) (*Client, error) {
	if accountSid == "" || authToken == "" || serviceSid == "" {
		return nil, errors.New("missing Twilio credentials")
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	return &Client{
		twilioClient: client,
		serviceSid:   serviceSid,
	}, nil
}

func (c *Client) SendOTP(phoneNumber string) error {
	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(phoneNumber)
	params.SetChannel("sms")

	_, err := c.twilioClient.VerifyV2.CreateVerification(c.serviceSid, params)
	if err != nil {
		return fmt.Errorf("failed to send OTP: %v", err)
	}

	return nil
}

func (c *Client) VerifyOTP(phoneNumber, code string) (bool, error) {
	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo(phoneNumber)
	params.SetCode(code)

	resp, err := c.twilioClient.VerifyV2.CreateVerificationCheck(c.serviceSid, params)
	if err != nil {
		return false, fmt.Errorf("failed to verify OTP: %v", err)
	}

	status := *resp.Status
	return status == "approved", nil
}