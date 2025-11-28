package hardcodedrepository

import (
	"errors"
	"math/rand"
	"quote-service/internal/repository"
)

var (
	ErrNotFound = errors.New("quote not found")
)

type HardcodedRepository struct {
	quotes map[int]repository.Quote
}

// NewHardcodedRepository creates a new hardcoded repository instance with data from quotes.csv
func NewHardcodedRepository() *HardcodedRepository {
	quotes := map[int]repository.Quote{
		1:  {ID: 1, Message: "Age is an issue of mind over matter. If you don't mind, it doesn't matter.", AuthorID: 1},
		2:  {ID: 2, Message: "Anyone who stops learning is old, whether at twenty or eighty. Anyone who keeps learning stays young. The greatest thing in life is to keep your mind young.", AuthorID: 2},
		3:  {ID: 3, Message: "Wrinkles should merely indicate where smiles have been.", AuthorID: 3},
		4:  {ID: 4, Message: "True terror is to wake up one morning and discover that your high school class is running the country.", AuthorID: 4},
		5:  {ID: 5, Message: "A diplomat is a man who always remembers a woman's birthday but never remembers her age.", AuthorID: 5},
		6:  {ID: 6, Message: "As I grow older, I pay less attention to what men say. I just watch what they do.", AuthorID: 6},
		7:  {ID: 7, Message: "How incessant and great are the ills with which a prolonged old age is replete.", AuthorID: 7},
		8:  {ID: 8, Message: "Old age, believe me, is a good and pleasant thing. It is true you are gently shouldered off the stage, but then you are given such a comfortable front stall as spectator.", AuthorID: 8},
		9:  {ID: 9, Message: "Old age has deformities enough of its own. It should never add to them the deformity of vice.", AuthorID: 9},
		10: {ID: 10, Message: "Nobody grows old merely by living a number of years. We grow old by deserting our ideals. Years may wrinkle the skin, but to give up enthusiasm wrinkles the soul.", AuthorID: 10},
		11: {ID: 11, Message: "An archaeologist is the best husband a woman can have. The older she gets the more interested he is in her.", AuthorID: 11},
		12: {ID: 12, Message: "All diseases run into one, old age.", AuthorID: 12},
		13: {ID: 13, Message: "Bashfulness is an ornament to youth, but a reproach to old age.", AuthorID: 13},
		14: {ID: 14, Message: "Like everyone else who makes the mistake of getting older, I begin each day with coffee and obituaries.", AuthorID: 14},
		15: {ID: 15, Message: "Age appears to be best in four things old wood best to burn, old wine to drink, old friends to trust, and old authors to read.", AuthorID: 15},
		16: {ID: 16, Message: "None are so old as those who have outlived enthusiasm.", AuthorID: 16},
		17: {ID: 17, Message: "Every man over forty is a scoundrel.", AuthorID: 17},
		18: {ID: 18, Message: "Forty is the old age of youth fifty the youth of old age.", AuthorID: 18},
		19: {ID: 19, Message: "You can't help getting older, but you don't have to get old.", AuthorID: 19},
		20: {ID: 20, Message: "Alas, after a certain age every man is responsible for his face.", AuthorID: 20},
	}

	return &HardcodedRepository{
		quotes: quotes,
	}
}

// GetQuoteByID returns a quote by its ID
func (r *HardcodedRepository) GetQuoteByID(id int) (*repository.Quote, error) {
	quote, ok := r.quotes[id]
	if !ok {
		return nil, ErrNotFound
	}

	return &quote, nil
}

// GetRandomQuote returns a random quote from the collection
func (r *HardcodedRepository) GetRandomQuote() (*repository.Quote, error) {
	if len(r.quotes) == 0 {
		return nil, ErrNotFound
	}

	// Get a random ID from 1 to 20
	randomID := rand.Intn(len(r.quotes)) + 1
	quote := r.quotes[randomID]

	return &quote, nil
}
