package orders

//go:generate go tool mockgen -destination mocks/mock_$GOFILE -package=mocks . Service
type Service interface{}
