package models

import "strings"

type UserId int64

type EmailAddress struct {
	emailAddressAsString string
}

func NewEmailAddress(emailAddressAsString string) *EmailAddress {
	return &EmailAddress{emailAddressAsString: strings.ToLower(emailAddressAsString)}
}

func (emailAddress *EmailAddress) String() string {
	return emailAddress.emailAddressAsString
}
