package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSearchParams struct {
	Offset        int64
	Limit         int64
	ProfileFilter []string
}

type Review struct {
	AuthorID   primitive.ObjectID `bson:"authorId"`
	AuthorName string             `bson:"authorName"`
	Score      int                `bson:"score"`
	Content    string             `bson:"content,omitempty"`
	CreatedAt  time.Time          `bson:"createdAt"`
}

type Profile struct {
	Role        string    `bson:"role"`
	Description string    `bson:"description,omitempty"`
	Rating      float64   `bson:"rating"`
	Reviews     []Review  `bson:"reviews"`
	CreatedAt   time.Time `bson:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt"`
}

type UserExt struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	DisplayName string             `bson:"displayName"`
	Email       string             `bson:"email"`
	Password    string             `bson:"password"`
	Balance     int64              `bson:"balance"`
	SystemRole  string             `bson:"systemRole,omitempty"`
	Profiles    []Profile          `bson:"profiles"`
	Active      bool               `bson:"active"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
}
