package models

import "time"

type PublicationId int64

type Publication struct {
	Id           int64     `json:"id"`
	AuthorId     UserId    `json:"authorId"`
	CreationTime time.Time `json:"creationTime"`
}

func CreateNewPublication(userId UserId) *Publication {
	return &Publication{
		Id:           -1,
		AuthorId:     userId,
		CreationTime: time.Now().UTC(),
	}
}
