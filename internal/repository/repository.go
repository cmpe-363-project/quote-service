package repository

type Quote struct {
	ID       int
	Message  string
	AuthorID int
}

type Repository interface {
	GetQuoteByID(id int) (*Quote, error)
	GetRandomQuote() (*Quote, error)
}
