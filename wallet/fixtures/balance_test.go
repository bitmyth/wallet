package fixtures

import (
	"github.com/bitmyth/walletserivce/config"
	"github.com/bitmyth/walletserivce/factory"
	"testing"
)

func TestPreloadTestingData(_ *testing.T) {
	config.SetConfigPath("../../")

	f, _ := factory.New()
	PreloadTestingData(f)
}
