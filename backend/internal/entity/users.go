package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	DisplayName string             `bson:"displayName"`
	Email       string             `bson:"email"`
	Password    string             `bson:"password"`
	SystemRole  SystemRoleType     `bson:"systemRole"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
	Balance     int                `bson:"balance"`
	Active      bool               `bson:"active"`
}

type SystemRoleType string

const (
	SystemRoleTypeAdmin SystemRoleType = "admin"
	SystemRoleTypeUser  SystemRoleType = "user"
)
