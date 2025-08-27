package utils

import (
	"fmt"
)

func SendVerificationEmail(email, token string) error {
	// For now, we'll just print the email to the console
	fmt.Printf("Sending verification email to %s with token %s\n", email, token)
	return nil
}