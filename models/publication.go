package models

import "time"

type PublicationId int64

type Publication struct {
	AuthorId     UserId    `json:"authorId"`
	CreationTime time.Time `json:"creationTime"`
}
