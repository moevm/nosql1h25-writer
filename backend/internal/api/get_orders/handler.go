package get_orders

type Request struct {
	Offset *int `query:"offset" validate:"gte=0" example:"0"`
	Limit  *int `query:"limit" validate:"gte=1,lte=100" example:"10"`
}

type Response struct {
	Orders []Order
	Total  int
}

type Order struct {
	Title          string
	Description    string
	CompletionTime int
	Cost           *int
	ClientName     string
	Rating         float64
}
