package users

import "go.mongodb.org/mongo-driver/bson/primitive"

type OperationType string

const (
	OperationTypeWithdraw OperationType = "withdraw"
	OperationTypeDeposit  OperationType = "deposit"
)

type UpdateIn struct {
	UserID                primitive.ObjectID
	DisplayName           *string
	FreelancerDescription *string
	ClientDescription     *string
}
