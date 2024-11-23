package wallet

type User struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Balance  float64 `json:"balance"`
}

type Transaction struct {
	ID              int     `json:"id"`
	UserID          int     `json:"user_id"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
	CreatedAt       string  `json:"created_at"`
}
