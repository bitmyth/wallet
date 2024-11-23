package wallet_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bitmyth/walletserivce/config"
	"github.com/bitmyth/walletserivce/factory"
	"github.com/bitmyth/walletserivce/route"
	"github.com/bitmyth/walletserivce/wallet"
	"github.com/bitmyth/walletserivce/wallet/fixtures"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

var (
	f factory.Factory
	r *gin.Engine
)

func TestMain(m *testing.M) {
	config.SetConfigPath("../")

	var err error
	f, err = factory.New()
	if err != nil {
		return
	}

	db, _ := f.DB()
	if err := db.Migrate(); err != nil {
		log.Fatal(err)
		return
	}

	r = route.Router(f)
	fixtures.PreloadTestingData(f)

	f.RegisterRoutes(r)

	m.Run()
}

func Test_Deposit(t *testing.T) {
	req := wallet.Request{Username: "user1", Amount: 0.123456}
	marshal, _ := json.Marshal(req)
	body := strings.NewReader(string(marshal))
	request := httptest.NewRequest(http.MethodPost, "/deposit", body)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, request)
	respBody, _ := io.ReadAll(resp.Result().Body)
	if resp.Code != http.StatusOK {
		t.Error("code is not 200")
	}
	t.Log(string(respBody))
}

func Test_DepositBadRequest(t *testing.T) {
	body := strings.NewReader(`{a:"2"}`)
	request := httptest.NewRequest(http.MethodPost, "/deposit", body)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, request)
	if resp.Code == http.StatusOK {
		t.Error("expect bad request")
	}
}

func Test_DepositDbError(t *testing.T) {
	tf, _ := factory.NewTesting()
	tr := route.Router(tf)
	tf.RegisterRoutes(tr)

	body := strings.NewReader(`{a:"2"}`)
	request := httptest.NewRequest(http.MethodPost, "/deposit", body)
	resp := httptest.NewRecorder()
	tr.ServeHTTP(resp, request)
	if resp.Code != http.StatusInternalServerError {
		t.Error("expect bad request")
	}
}

func Test_Withdraw(t *testing.T) {
	req := wallet.Request{Username: "user1", Amount: 1.123456}
	marshal, _ := json.Marshal(req)
	body := strings.NewReader(string(marshal))
	request := httptest.NewRequest(http.MethodPost, "/withdraw", body)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, request)
	respBody, _ := io.ReadAll(resp.Result().Body)
	if resp.Code != http.StatusOK {
		t.Error("code is not 200")
	}
	t.Log(string(respBody))
}

func Test_WithdrawBadRequest(t *testing.T) {
	body := strings.NewReader(`{a:"2"}`)
	request := httptest.NewRequest(http.MethodPost, "/withdraw", body)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, request)
	if resp.Code == http.StatusOK {
		t.Error("expect bad request")
	}
}

func Test_WithdrawMoreThanBalance(t *testing.T) {
	req := wallet.Request{Username: "user1", Amount: 1000}
	marshal, _ := json.Marshal(req)
	body := strings.NewReader(string(marshal))
	request := httptest.NewRequest(http.MethodPost, "/withdraw", body)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, request)
	respBody, _ := io.ReadAll(resp.Result().Body)
	if resp.Code != http.StatusBadRequest {
		t.Error("expect bad request response")
	}
	t.Log(string(respBody))
}

func Test_Transfer(t *testing.T) {
	req := wallet.TransferRequest{
		From:   "user1",
		To:     "user2",
		Amount: 1.123456,
	}
	marshal, _ := json.Marshal(req)
	body := strings.NewReader(string(marshal))
	request := httptest.NewRequest(http.MethodPost, "/transfer", body)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, request)
	respBody, _ := io.ReadAll(resp.Result().Body)
	t.Log(resp.Code)
	t.Log(string(respBody))
}

func TestTransferConcurrent(t *testing.T) {
	a2b := func(group *sync.WaitGroup) {
		req := wallet.TransferRequest{From: "user1", To: "user2", Amount: 0.12}
		marshal, _ := json.Marshal(req)
		body := strings.NewReader(string(marshal))
		request := httptest.NewRequest(http.MethodPost, "/transfer", body)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, request)
		respBody, _ := io.ReadAll(resp.Result().Body)
		t.Log(resp.Code)
		t.Log(string(respBody))
		group.Done()
	}
	b2a := func(group *sync.WaitGroup) {
		req := wallet.TransferRequest{From: "user2", To: "user1", Amount: 0.12}
		marshal, _ := json.Marshal(req)
		body := strings.NewReader(string(marshal))
		request := httptest.NewRequest(http.MethodPost, "/transfer", body)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, request)
		respBody, _ := io.ReadAll(resp.Result().Body)
		t.Log(resp.Code)
		t.Log(string(respBody))
		group.Done()
	}

	svc := wallet.NewService(f)
	balanceBefore, _ := svc.GetBalance(context.Background(), "user1")

	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 50; i++ {
		go a2b(&wg)
		go b2a(&wg)
	}
	wg.Wait()

	balanceAfter, _ := svc.GetBalance(context.Background(), "user1")
	t.Log(balanceBefore, balanceAfter)
	if balanceAfter != balanceBefore {
		t.Error("balance should not change")
	}
}

func Test_TransferBadRequest(t *testing.T) {
	body := strings.NewReader(`{a:"2"}`)
	request := httptest.NewRequest(http.MethodPost, "/transfer", body)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, request)
	if resp.Code == http.StatusOK {
		t.Error("expect bad request")
	}
}

func Test_TransferFailed(t *testing.T) {
	req := wallet.TransferRequest{
		From:   "user1",
		To:     "user2",
		Amount: 1000,
	}
	marshal, _ := json.Marshal(req)
	body := strings.NewReader(string(marshal))
	request := httptest.NewRequest(http.MethodPost, "/transfer", body)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, request)
	respBody, _ := io.ReadAll(resp.Result().Body)
	if resp.Code != http.StatusBadRequest {
		t.Error("expect bad request response")
	}
	t.Log(string(respBody))
}

func Test_TransactionHistory(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/transactions/%s", "user1"), nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, request)
	respBody, _ := io.ReadAll(resp.Result().Body)
	t.Log(resp.Code)
	t.Log(string(respBody))
}

func Test_TransactionHistoryBadRequest(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/transactions/%s", "notfound"), nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, request)
	respBody, _ := io.ReadAll(resp.Result().Body)
	t.Log(resp.Code, string(respBody))
}

func Test_Balance(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/balance/%s", "user1"), nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, request)
	respBody, _ := io.ReadAll(resp.Result().Body)
	if resp.Code != http.StatusOK {
		t.Error("should be ok")
	}
	t.Log(string(respBody))
}

func TestNewController(t *testing.T) {
	controller := wallet.NewController(f)
	if controller == nil {
		t.Error("expect controller not nil")
	}
}

func BenchmarkWithdraw(b *testing.B) {
	req := wallet.Request{Username: "user1", Amount: 0.003456}
	marshal, _ := json.Marshal(req)

	for i := 0; i < b.N; i++ {
		body := strings.NewReader(string(marshal))
		request := httptest.NewRequest(http.MethodPost, "/withdraw", body)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, request)
		respBody, _ := io.ReadAll(resp.Result().Body)
		if resp.Code != http.StatusOK {
			b.Error("code is not 200")
			b.Log(respBody)
		}
	}
}
func BenchmarkDeposit(b *testing.B) {
	req := wallet.Request{Username: "user1", Amount: 0.123456}
	marshal, _ := json.Marshal(req)

	for i := 0; i < b.N; i++ {
		body := strings.NewReader(string(marshal))
		request := httptest.NewRequest(http.MethodPost, "/deposit", body)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, request)
		_, _ = io.ReadAll(resp.Result().Body)
		if resp.Code != http.StatusOK {
			b.Error("code is not 200")
		}
	}
}

func BenchmarkTransfer(b *testing.B) {
	req := wallet.TransferRequest{
		From:   "user1",
		To:     "user2",
		Amount: 0.000123456,
	}
	marshal, _ := json.Marshal(req)
	for i := 0; i < b.N; i++ {
		body := strings.NewReader(string(marshal))
		request := httptest.NewRequest(http.MethodPost, "/transfer", body)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, request)
		_, _ = io.ReadAll(resp.Result().Body)
	}
}
