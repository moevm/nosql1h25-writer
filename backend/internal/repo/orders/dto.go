package orders

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

type FindIn struct {
	Limit   int
	Offset  int
	MinCost *int
	MaxCost *int
	MinTime *int64
	MaxTime *int64
	Search  *string
	SortBy  *string
}

type FindOut struct {
	Orders []entity.Order
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
	FreelancerID   *primitive.ObjectID
}
