package backend_test

import (
	"testing"

	"github.com/weiwolves/pkg/backend"
	"github.com/weiwolves/pkg/config/cfgmock"
	"github.com/weiwolves/pkg/config/cfgmodel"
	"github.com/weiwolves/pkg/config/source"
)

// benchmarkGlobalStruct trick the compiler to not optimize anything
var benchmarkGlobalStruct bool

func Benchmark_StructGlobal(b *testing.B) {

	sg := cfgmock.NewService().NewScoped(1, 1)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkGlobalStruct, err = backend.Backend.DevCSSMinifyFiles.Get(sg) // any random struct field
		if err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_StructSpecific(b *testing.B) {

	sg := cfgmock.NewService().NewScoped(1, 1)

	mb := cfgmodel.NewBool("aa/bb/cc", cfgmodel.WithFieldFromSectionSlice(backend.ConfigStructure), cfgmodel.WithSource(source.YesNo))

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var err error
		benchmarkGlobalStruct, err = mb.Get(sg)
		if err != nil {
			b.Error(err)
		}
	}
}
