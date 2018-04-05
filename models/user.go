package models

import "strings"

// UserId is used to symbolize that this int64 is a specifically a UserId
type UserId int64

// EmailAddress is a wrapper struct of the private string emailAddress to make
// sure that emails are always formatted properly when passed around in the system.
type EmailAddress struct {
	emailAddressAsString string
}

// NewEmailAddress constructs an EmailAddress objects and provides some guranteess
// that it is a properly formated email address.
func NewEmailAddress(emailAddressAsString string) *EmailAddress {
	return &EmailAddress{emailAddressAsString: strings.ToLower(emailAddressAsString)}
}

// String returns the email address as a string
func (emailAddress *EmailAddress) String() string {
	return emailAddress.emailAddressAsString
}
