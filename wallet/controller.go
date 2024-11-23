package wallet

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	factory dependency
	service *Service
}

func NewController(f dependency) *Controller {
	return &Controller{
		factory: f,
		service: NewService(f),
	}
}

type Request struct {
	Username string
	Amount   float64
}

func (c Controller) Deposit(ctx *gin.Context) {
	var req Request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := req.Username
	amount := req.Amount

	db, _ := c.factory.DB()

	var tx *sql.Tx
	var err error

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	steps := []func(){
		func() {
			tx, err = db.Begin()
		},
		func() {
			_, err = tx.Exec("UPDATE users SET balance = balance + $1 WHERE username = $2", amount, username)
		},
		func() {
			err = c.logTransaction(tx, username, amount, "deposit")
		},
		func() {
			err = tx.Commit()
		},
	}

	for _, step := range steps {
		step()
		if c.handleError(ctx, err) {
			return
		}
	}

	rdb, _ := c.factory.Redis()
	rdb.Del(ctx, username)

	balance, err := c.service.GetBalance(ctx, username)
	if c.handleError(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"balance": balance})
}

func (c Controller) Withdraw(ctx *gin.Context) {
	var req Request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	balance, err := c.service.GetBalance(ctx, req.Username)
	if c.handleError(ctx, err) {
		return
	}
	if balance < req.Amount {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "insufficient balance"})
		return
	}

	var tx *sql.Tx
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	db, _ := c.factory.DB()

	steps := []func(){
		func() {
			tx, err = db.Begin()
		},
		func() {
			_, err = db.Exec("UPDATE users SET balance = balance - $1 WHERE username = $2", req.Amount, req.Username)
		},
		func() {
			err = c.logTransaction(tx, req.Username, -req.Amount, "withdraw")
		},
		func() {
			err = tx.Commit()
		},
	}

	for _, step := range steps {
		step()
		if c.handleError(ctx, err) {
			return
		}
	}

	rdb, _ := c.factory.Redis()

	rdb.Del(ctx, req.Username)

	ctx.JSON(http.StatusOK, gin.H{"balance": balance - req.Amount})
}

type TransferRequest struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

func (c Controller) Transfer(ctx *gin.Context) {
	var req TransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	balance, err := c.service.GetBalance(ctx, req.From)
	if c.handleError(ctx, err) {
		return
	}

	if balance < req.Amount {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "insufficient balance"})
		return
	}

	var tx *sql.Tx

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	steps := []func(){
		// start a transaction
		func() {
			db, _ := c.factory.DB()
			tx, err = db.Begin()
		},
		// lock sender and receiver balance
		func() {
			_, err = tx.Exec("SELECT balance FROM users WHERE username = $1 OR username = $2 FOR UPDATE", req.From, req.To)
		},
		// withdraw from sender
		func() {
			_, err = tx.Exec("UPDATE users SET balance = balance - $1 WHERE username = $2", req.Amount, req.From)
		},
		// deposit to receiver
		func() {
			_, err = tx.Exec("UPDATE users SET balance = balance + $1 WHERE username = $2", req.Amount, req.To)
		},
		// log transactions for both users
		func() { err = c.logTransaction(tx, req.From, -req.Amount, "transfer") },
		func() { err = c.logTransaction(tx, req.To, req.Amount, "transfer") },
		func() { err = tx.Commit() },
	}

	for _, step := range steps {
		step()
		if c.handleError(ctx, err) {
			return
		}
	}

	rdb, _ := c.factory.Redis()

	// del cache for both users
	rdb.Del(ctx, req.From, req.To)

	ctx.Status(http.StatusOK)
}

func (c Controller) GetBalance(ctx *gin.Context) {
	username := ctx.Param("username")

	balance, err := c.service.GetBalance(ctx, username)
	if c.handleError(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"balance": balance})
}

func (c Controller) GetTransactionHistory(ctx *gin.Context) {
	username := ctx.Param("username")

	db, _ := c.factory.DB()

	rows, err := db.Query("SELECT id, user_id, amount, transaction_type, created_at FROM transactions WHERE user_id = (SELECT id FROM users WHERE username = $1)", username)
	if c.handleError(ctx, err) {
		return
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		if err = rows.Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.TransactionType, &transaction.CreatedAt); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		transactions = append(transactions, transaction)
	}

	ctx.JSON(http.StatusOK, transactions)
}

func (c Controller) logTransaction(tx *sql.Tx, username string, amount float64, transactionType string) error {
	logger := c.factory.Logger()

	var userID int
	err := tx.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		logger.Error("error fetching user ID:", err)
		return err
	}

	_, err = tx.Exec("INSERT INTO transactions (user_id, amount, transaction_type) VALUES ($1, $2, $3)", userID, amount, transactionType)
	if err != nil {
		logger.Error("error logging transaction:", err)
		return err
	}
	return nil
}

func (c Controller) RegisterRoutes(router *gin.Engine) {
	router.POST("/deposit", c.Deposit)
	router.POST("/withdraw", c.Withdraw)
	router.POST("/transfer", c.Transfer)
	router.GET("/balance/:username", c.GetBalance)
	router.GET("/transactions/:username", c.GetTransactionHistory)
}

func (c Controller) handleError(ctx *gin.Context, err error) bool {
	l := c.factory.Logger()
	if err != nil {
		l.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return true
	}
	return false
}
