package orders

import (
	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FindOut struct {
	Orders []OrderWithClientData
	Total  int
}

type OrderWithClientData struct {
	ID             primitive.ObjectID
	ClientID       primitive.ObjectID
	Title          string
	Description    string
	CompletionTime int64
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

type UpdateIn struct {
	OrderID        primitive.ObjectID
	Title          *string
	Description    *string
	CompletionTime *int64
	Cost           *int
	Status         *entity.StatusType
}
