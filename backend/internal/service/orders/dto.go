package orders

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateIn struct {
	ClientID       primitive.ObjectID
	Title          string
	Description    string
	CompletionTime int64
	Cost           int
}
