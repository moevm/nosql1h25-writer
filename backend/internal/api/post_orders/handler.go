package post_orders

import "go.mongodb.org/mongo-driver/bson/primitive"

type Request struct {
	Title          string `json:"title" validate:"required,min=3,max=32"`
	Description    string `json:"description" validate:"required,min=16,max=4096"`
	CompletionTime int    `json:"comletionTime" validate:"required,gte=3600000000000"`
	Cost           *int   `json:"cost" validate:"gte=1"`
}

type Response struct {
	ID primitive.ObjectID `json:"id"`
}
