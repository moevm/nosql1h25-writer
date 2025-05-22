package orders

import "go.mongodb.org/mongo-driver/bson/primitive"

type FindOut struct {
	Orders []OrderWithClientData
	Total  int
}

type OrderWithClientData struct {
	ID             string
	Title          string
	Description    string
	CompletionTime int
	Cost           int
	ClientName     string
	Rating         float64
}

type CreateIn struct {
	ClientID       primitive.ObjectID
	Title          string
	Description    string
	CompletionTime int64
	Cost           int
}
