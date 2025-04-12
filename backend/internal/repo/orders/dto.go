package orders

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FindOrdersOut struct {
	Name         string
	ID           primitive.ObjectID `json:"id" example:"64a1b2c3d4e5f6a7b8c9d0e2"`
	Title        string             `json:"title" example:"Разработка сайта"`
	Description  string             `json:"description" example:"Создание лендинга"`
	Budget       float64            `json:"budget" example:"1000"`
	Active       bool               `json:"active" example:"true"`
	CreatedAt    time.Time          `json:"createdAt" example:"2023-10-01T10:00:00Z"`
	FreelancerID string             `json:"freelancerId,omitempty" example:"64a1b2c3d4e5f6a7b8c9d0e3"`
	ClientID     string             `json:"clientId,omitempty" example:"64a1b2c3d4e5f6a7b8c9d0e1"`
}
