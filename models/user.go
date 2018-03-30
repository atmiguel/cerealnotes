package models

import "strings"

type UserId int64

type EmailAddress struct {
	emailAsString string
}

func NewEmailAddress(emailAddr string) *EmailAddress {
	return &EmailAddress{emailAsString: strings.ToLower(emailAddr)}
}

func (emailAddress *EmailAddress) String() string {
	return emailAddress.emailAsString
}
