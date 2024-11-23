package wallet_test

import (
	"context"
	"github.com/bitmyth/walletserivce/factory"
	"github.com/bitmyth/walletserivce/wallet"
	"testing"
)

func TestService_GetBalance(t *testing.T) {
	s := wallet.NewService(f)
	balance, err := s.GetBalance(context.Background(), "user1")
	if err != nil {
		t.Error(err)
		return
	}
	if balance < 0 {
		t.Error("user1 balance is wrong")
	}

	rdb, _ := f.Redis()
	rdb.Del(context.Background(), "user1")

	balance, err = s.GetBalance(context.Background(), "user1")
	if err != nil {
		t.Error(err)
		return
	}
	if balance < 0 {
		t.Error("user1 balance is wrong")
	}
}

func TestGetBalanceForNotFoundUser(t *testing.T) {
	s := wallet.NewService(f)
	balance, err := s.GetBalance(context.Background(), "notfound")
	if err == nil {
		t.Error("expect return err")
		return
	}
	if balance > 0 {
		t.Error("balance should be 0")
	}
}

func TestGetBalanceDbError(t *testing.T) {
	ft, _ := factory.NewTesting()
	s := wallet.NewService(ft)
	_, err := s.GetBalance(context.Background(), "notfound")
	if err == nil {
		t.Error("expect return err")
		return
	}
}

func BenchmarkGetBalance(b *testing.B) {
	s := wallet.NewService(f)
	for i := 0; i < b.N; i++ {
		_, err := s.GetBalance(context.Background(), "notfound")
		if err == nil {
			b.Error(err)
			return
		}
	}
}
