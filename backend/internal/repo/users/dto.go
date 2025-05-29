package users

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateIn struct {
	DisplayName string
	Email       string
	Password    string
}

type UpdateIn struct {
	UserID                primitive.ObjectID
	DisplayName           *string
	FreelancerDescription *string
	ClientDescription     *string
}
