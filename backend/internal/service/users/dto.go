package users

type OperationType string

const (
	OperationTypeWithdraw OperationType = "withdraw"
	OperationTypeDeposit  OperationType = "deposit"
)
