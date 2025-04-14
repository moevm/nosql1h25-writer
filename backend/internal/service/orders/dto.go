package orders

type FindOut struct {
	Orders []OrderWithClientData
}

type OrderWithClientData struct {
	Title          string
	Description    string
	CompletionTime int
	Cost           int
	ClientName     string
	Rating         float64
}
