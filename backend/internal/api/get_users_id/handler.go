package get_users_id

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

type Request struct {
	ID       primitive.ObjectID `param:"id" validate:"required"`
	Profiles []string           `query:"profile" validate:"dive,oneof=client freelancer"`
}

type Response struct {
	ID          primitive.ObjectID    `json:"id"`
	DisplayName string                `json:"displayName"`
	Email       string                `json:"email"`
	SystemRole  entity.SystemRoleType `json:"systemRole"`
	CreatedAt   time.Time             `json:"createdAt"`
	UpdatedAt   time.Time             `json:"updatedAt"`
	Balance     int                   `json:"balance"`
	Client      *Profile              `json:"client,omitempty"`
	Freelancer  *Profile              `json:"freelancer,omitempty"`
}

type Profile struct {
	Rating      float64   `json:"rating"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
