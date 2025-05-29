package users

import (
	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OperationType string

const (
	OperationTypeWithdraw OperationType = "withdraw"
	OperationTypeDeposit  OperationType = "deposit"
)

type UpdateInput struct {
	RequesterID           primitive.ObjectID
	RequesterRole         entity.SystemRoleType
	UserID                primitive.ObjectID
	DisplayName           *string
	FreelancerDescription *string
	ClientDescription     *string
}
