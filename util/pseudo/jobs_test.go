package pseudo

import (
	"testing"

	"github.com/weiwolves/pkg/util/assert"
)

func TestJobs(t *testing.T) {
	s := MustNewService(0, nil)
	for _, lang := range s.GetLangs() {
		assert.NoError(t, s.SetLang(lang))

		v := s.Company()
		t.Logf("Company %q %q", lang, v)
		assert.LenBetween(t, v, 1, 40, "Company failed with lang %s", lang)

		v = s.CompanyLegal()
		t.Logf("CompanyLegal %q %q", lang, v)
		assert.LenBetween(t, v, 1, 70, "CompanyLegal failed with lang %s", lang)

		v = s.JobTitle()
		t.Logf("JobTitle %q %q", lang, v)
		assert.LenBetween(t, v, 1, 40, "JobTitle failed with lang %s", lang)

		v = s.Industry()
		t.Logf("Industry %q %q", lang, v)
		assert.LenBetween(t, v, 1, 50, "Industry failed with lang %s", lang)

	}
}
