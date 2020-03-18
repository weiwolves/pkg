package pseudo

import (
	"testing"

	"github.com/weiwolves/pkg/util/assert"
)

func TestCreditCards(t *testing.T) {
	s := MustNewService(0, nil)
	for _, lang := range s.GetLangs() {
		assert.NoError(t, s.SetLang(lang))
		assert.NotEmpty(t, s.CreditCardType(), "s.CreditCardType()")
		assert.NotEmpty(t, s.CreditCardNum(""), "s.CreditCardNum()")
		assert.NotEmpty(t, s.CreditCardNum("visa"), "s.CreditCardNum(visa)")
	}
}
