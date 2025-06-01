package users

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

type FindIn struct {
	Offset              int
	Limit               int
	NameSearch          *string
	EmailSearch         *string
	Roles               []string
	MinFreelancerRating *float64
	MaxFreelancerRating *float64
	MinClientRating     *float64
	MaxClientRating     *float64
	MinCreatedAt        *time.Time
	MaxCreatedAt        *time.Time
	MaxBalance          *int
	MinBalance          *int
	SortBy              *string
}

type FindOut struct {
	Users []entity.UserExt
	Total int
}

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
