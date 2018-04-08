package models

import "strings"

// UserId is a type that represents the key of the user in the database.
type UserId int64

// EmailAddress is a wrapper of the string class to make sure that emails
// are always formatted properly when passed around in the system.
type EmailAddress struct {
	emailAddressAsString string
}

// NewEmailAddress provides some guranteess that it is a properly formated email
// address.
func NewEmailAddress(emailAddressAsString string) *EmailAddress {
	return &EmailAddress{emailAddressAsString: strings.ToLower(emailAddressAsString)}
}

// String returns the email address as a string
func (emailAddress *EmailAddress) String() string {
	return emailAddress.emailAddressAsString
}
