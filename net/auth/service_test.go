// Copyright (c) 2014 Olivier Poitrey <rs@dailymotion.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package auth_test

import (
	"testing"

	"github.com/corestoreio/errors"
	"github.com/weiwolves/pkg/net/auth"
	"github.com/weiwolves/pkg/util/assert"
)

func withError() auth.Option {
	return func(s *auth.Service) error {
		return errors.NotValid.Newf("Paaaaanic!")
	}
}

func TestMustNew_Default(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.True(t, errors.NotValid.Match(err), "Error: %s", err)
		} else {
			t.Fatal("Expecting a Panic")
		}
	}()
	_ = auth.MustNew(nil, withError())
}

func TestMustNew_Website(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			assert.True(t, errors.NotValid.Match(err), "Error: %s", err)
		} else {
			t.Fatal("Expecting a Panic")
		}
	}()
	_ = auth.MustNew(nil, withError())
}
