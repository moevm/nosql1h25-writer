package orders

type FindOut struct {
	/// тут короче надо будет завести поля, которые получатся после запроса с lookup, где ты сможешь получить имя
	/// и рейтинг заказчика. и сюда запишется как инфа о заказе, так и вот эти дополнительные поля.
	/// предполагаю что будет выглядеть так, но не уверен
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
