package pseudo

import (
	"testing"

	"github.com/weiwolves/pkg/util/assert"
)

func TestCurrencies(t *testing.T) {
	s := MustNewService(0, nil)
	for _, lang := range s.GetLangs() {
		assert.NoError(t, s.SetLang(lang))
		assert.NotEmpty(t, s.Currency(), "Currency()")
		assert.NotEmpty(t, s.CurrencyCode(), "CurrencyCode()")
	}
}
