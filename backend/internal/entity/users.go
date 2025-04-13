package entity

import (
	"fmt"
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

type UserExt struct {
	User
	Client     Profile `bson:"client"`
	Freelancer Profile `bson:"freelancer"`
}

type Profile struct {
	Rating      float64   `bson:"rating"`
	Description string    `bson:"description"`
	UpdatedAt   time.Time `bson:"updatedAt"`
	Reviews     []Review  `bson:"reviews"`
}

type Review struct {
	AuthorID   primitive.ObjectID `bson:"authorId"`
	AuthorName string             `bson:"authorName"`
	Score      int                `bson:"score"`
	Content    string             `bson:"content,omitempty"`
	CreatedAt  time.Time          `bson:"createdAt"`
}

func DefaultUser(email, password string, createdAt, updatedAt time.Time) UserExt {
	return UserExt{
		User: User{
			DisplayName: fmt.Sprintf("Пользователь_%s", email),
			Email:       email,
			Password:    password,
			SystemRole:  SystemRoleTypeUser,
			Active:      true,
			Balance:     0,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		},
		Client: Profile{
			Rating:      0,
			Description: "Описание профиля заказчика",
			UpdatedAt:   updatedAt,
			Reviews:     []Review{},
		},
		Freelancer: Profile{
			Rating:      0,
			Description: "Описание профиля исполнителя",
			UpdatedAt:   updatedAt,
			Reviews:     []Review{},
		},
	}
}
