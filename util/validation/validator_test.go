// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// The MIT License (MIT)
//
// Copyright (c) 2014 Alex Saskevich
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package validation

import (
	"strings"
	"testing"
	"time"

	"github.com/weiwolves/pkg/sync/bgwork"
)

func TestIsAlpha(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"\n", false},
		{"\r", false},
		{"Ⅸ", false},
		{"", true},
		{"   fooo   ", false},
		{"abc!!!", false},
		{"abc1", false},
		{"abc〩", false},
		{"abc", true},
		{"소주", false},
		{"ABC", true},
		{"FoObAr", true},
		{"소aBC", false},
		{"소", false},
		{"달기&Co.", false},
		{"〩Hours", false},
		{"\ufff0", false},
		{"\u0070", true},  // UTF-8(ASCII): p
		{"\u0026", false}, // UTF-8(ASCII): &
		{"\u0030", false}, // UTF-8(ASCII): 0
		{"123", false},
		{"0123", false},
		{"-00123", false},
		{"0", false},
		{"-0", false},
		{"123.123", false},
		{" ", false},
		{".", false},
		{"-1¾", false},
		{"1¾", false},
		{"〥〩", false},
		{"모자", false},
		{"ix", true},
		{"۳۵۶۰", false},
		{"1--", false},
		{"1-1", false},
		{"-", false},
		{"--", false},
		{"1++", false},
		{"1+1", false},
		{"+", false},
		{"++", false},
		{"+1", false},
	}
	for _, test := range tests {
		actual := IsAlpha(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsAlpha(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsUTFLetter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"\n", false},
		{"\r", false},
		{"Ⅸ", false},
		{"", true},
		{"   fooo   ", false},
		{"abc!!!", false},
		{"abc1", false},
		{"abc〩", false},
		{"", true},
		{"abc", true},
		{"소주", true},
		{"ABC", true},
		{"FoObAr", true},
		{"소aBC", true},
		{"소", true},
		{"달기&Co.", false},
		{"〩Hours", false},
		{"\ufff0", false},
		{"\u0070", true},  // UTF-8(ASCII): p
		{"\u0026", false}, // UTF-8(ASCII): &
		{"\u0030", false}, // UTF-8(ASCII): 0
		{"123", false},
		{"0123", false},
		{"-00123", false},
		{"0", false},
		{"-0", false},
		{"123.123", false},
		{" ", false},
		{".", false},
		{"-1¾", false},
		{"1¾", false},
		{"〥〩", false},
		{"모자", true},
		{"ix", true},
		{"۳۵۶۰", false},
		{"1--", false},
		{"1-1", false},
		{"-", false},
		{"--", false},
		{"1++", false},
		{"1+1", false},
		{"+", false},
		{"++", false},
		{"+1", false},
	}
	for _, test := range tests {
		actual := IsUTFLetter(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsUTFLetter(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsAlphanumeric(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"\n", false},
		{"\r", false},
		{"Ⅸ", false},
		{"", true},
		{"   fooo   ", false},
		{"abc!!!", false},
		{"abc123", true},
		{"ABC111", true},
		{"abc1", true},
		{"abc〩", false},
		{"abc", true},
		{"소주", false},
		{"ABC", true},
		{"FoObAr", true},
		{"소aBC", false},
		{"소", false},
		{"달기&Co.", false},
		{"〩Hours", false},
		{"\ufff0", false},
		{"\u0070", true},  // UTF-8(ASCII): p
		{"\u0026", false}, // UTF-8(ASCII): &
		{"\u0030", true},  // UTF-8(ASCII): 0
		{"123", true},
		{"0123", true},
		{"-00123", false},
		{"0", true},
		{"-0", false},
		{"123.123", false},
		{" ", false},
		{".", false},
		{"-1¾", false},
		{"1¾", false},
		{"〥〩", false},
		{"모자", false},
		{"ix", true},
		{"۳۵۶۰", false},
		{"1--", false},
		{"1-1", false},
		{"-", false},
		{"--", false},
		{"1++", false},
		{"1+1", false},
		{"+", false},
		{"++", false},
		{"+1", false},
	}
	for _, test := range tests {
		actual := IsAlphanumeric(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsAlphanumeric(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsUTFLetterNumeric(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"\n", false},
		{"\r", false},
		{"Ⅸ", true},
		{"", true},
		{"   fooo   ", false},
		{"abc!!!", false},
		{"abc1", true},
		{"abc〩", true},
		{"abc", true},
		{"소주", true},
		{"ABC", true},
		{"FoObAr", true},
		{"소aBC", true},
		{"소", true},
		{"달기&Co.", false},
		{"〩Hours", true},
		{"\ufff0", false},
		{"\u0070", true},  // UTF-8(ASCII): p
		{"\u0026", false}, // UTF-8(ASCII): &
		{"\u0030", true},  // UTF-8(ASCII): 0
		{"123", true},
		{"0123", true},
		{"-00123", false},
		{"0", true},
		{"-0", false},
		{"123.123", false},
		{" ", false},
		{".", false},
		{"-1¾", false},
		{"1¾", true},
		{"〥〩", true},
		{"모자", true},
		{"ix", true},
		{"۳۵۶۰", true},
		{"1--", false},
		{"1-1", false},
		{"-", false},
		{"--", false},
		{"1++", false},
		{"1+1", false},
		{"+", false},
		{"++", false},
		{"+1", false},
	}
	for _, test := range tests {
		actual := IsUTFLetterNumeric(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsUTFLetterNumeric(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsNumeric(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"\n", false},
		{"\r", false},
		{"Ⅸ", false},
		{"", true},
		{"   fooo   ", false},
		{"abc!!!", false},
		{"abc1", false},
		{"abc〩", false},
		{"abc", false},
		{"소주", false},
		{"ABC", false},
		{"FoObAr", false},
		{"소aBC", false},
		{"소", false},
		{"달기&Co.", false},
		{"〩Hours", false},
		{"\ufff0", false},
		{"\u0070", false}, // UTF-8(ASCII): p
		{"\u0026", false}, // UTF-8(ASCII): &
		{"\u0030", true},  // UTF-8(ASCII): 0
		{"123", true},
		{"0123", true},
		{"-00123", false},
		{"+00123", false},
		{"0", true},
		{"-0", false},
		{"123.123", false},
		{" ", false},
		{".", false},
		{"12𐅪3", false},
		{"-1¾", false},
		{"1¾", false},
		{"〥〩", false},
		{"모자", false},
		{"ix", false},
		{"۳۵۶۰", false},
		{"1--", false},
		{"1-1", false},
		{"-", false},
		{"--", false},
		{"1++", false},
		{"1+1", false},
		{"+", false},
		{"++", false},
		{"+1", false},
	}
	for _, test := range tests {
		actual := IsNumeric(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsNumeric(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsUTFNumeric(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"\n", false},
		{"\r", false},
		{"Ⅸ", true},
		{"", true},
		{"   fooo   ", false},
		{"abc!!!", false},
		{"abc1", false},
		{"abc〩", false},
		{"abc", false},
		{"소주", false},
		{"ABC", false},
		{"FoObAr", false},
		{"소aBC", false},
		{"소", false},
		{"달기&Co.", false},
		{"〩Hours", false},
		{"\ufff0", false},
		{"\u0070", false}, // UTF-8(ASCII): p
		{"\u0026", false}, // UTF-8(ASCII): &
		{"\u0030", true},  // UTF-8(ASCII): 0
		{"123", true},
		{"0123", true},
		{"-00123", true},
		{"0", true},
		{"-0", true},
		{"--0", false},
		{"-0-", false},
		{"123.123", false},
		{" ", false},
		{".", false},
		{"12𐅪3", true},
		{"-1¾", true},
		{"1¾", true},
		{"〥〩", true},
		{"모자", false},
		{"ix", false},
		{"۳۵۶۰", true},
		{"1++", false},
		{"1+1", false},
		{"+", false},
		{"++", false},
		{"+1", true},
	}
	for _, test := range tests {
		actual := IsUTFNumeric(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsUTFNumeric(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsUTFDigit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{

		{"\n", false},
		{"\r", false},
		{"Ⅸ", false},
		{"", true},
		{"   fooo   ", false},
		{"abc!!!", false},
		{"abc1", false},
		{"abc〩", false},
		{"abc", false},
		{"소주", false},
		{"ABC", false},
		{"FoObAr", false},
		{"소aBC", false},
		{"소", false},
		{"달기&Co.", false},
		{"〩Hours", false},
		{"\ufff0", false},
		{"\u0070", false}, // UTF-8(ASCII): p
		{"\u0026", false}, // UTF-8(ASCII): &
		{"\u0030", true},  // UTF-8(ASCII): 0
		{"123", true},
		{"0123", true},
		{"-00123", true},
		{"0", true},
		{"-0", true},
		{"--0", false},
		{"-0-", false},
		{"123.123", false},
		{" ", false},
		{".", false},
		{"12𐅪3", false},
		{"1483920", true},
		{"", true},
		{"۳۵۶۰", true},
		{"-29", true},
		{"-1¾", false},
		{"1¾", false},
		{"〥〩", false},
		{"모자", false},
		{"ix", false},
		{"۳۵۶۰", true},
		{"1++", false},
		{"1+1", false},
		{"+", false},
		{"++", false},
		{"+1", true},
	}
	for _, test := range tests {
		actual := IsUTFDigit(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsUTFDigit(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsLowerCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", true},
		{"abc123", true},
		{"abc", true},
		{"a b c", true},
		{"abcß", true},
		{"abcẞ", false},
		{"ABCẞ", false},
		{"tr竪s 端ber", true},
		{"fooBar", false},
		{"123ABC", false},
		{"ABC123", false},
		{"ABC", false},
		{"S T R", false},
		{"fooBar", false},
		{"abacaba123", true},
	}
	for _, test := range tests {
		actual := IsLowerCase(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsLowerCase(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsUpperCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", true},
		{"abc123", false},
		{"abc", false},
		{"a b c", false},
		{"abcß", false},
		{"abcẞ", false},
		{"ABCẞ", true},
		{"tr竪s 端ber", false},
		{"fooBar", false},
		{"123ABC", true},
		{"ABC123", true},
		{"ABC", true},
		{"S T R", true},
		{"fooBar", false},
		{"abacaba123", false},
	}
	for _, test := range tests {
		actual := IsUpperCase(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsUpperCase(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestHasLowerCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", true},
		{"abc123", true},
		{"abc", true},
		{"a b c", true},
		{"abcß", true},
		{"abcẞ", true},
		{"ABCẞ", false},
		{"tr竪s 端ber", true},
		{"fooBar", true},
		{"123ABC", false},
		{"ABC123", false},
		{"ABC", false},
		{"S T R", false},
		{"fooBar", true},
		{"abacaba123", true},
		{"FÒÔBÀŘ", false},
		{"fòôbàř", true},
		{"fÒÔBÀŘ", true},
	}
	for _, test := range tests {
		actual := HasLowerCase(test.param)
		if actual != test.expected {
			t.Errorf("Expected HasLowerCase(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestHasUpperCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", true},
		{"abc123", false},
		{"abc", false},
		{"a b c", false},
		{"abcß", false},
		{"abcẞ", false},
		{"ABCẞ", true},
		{"tr竪s 端ber", false},
		{"fooBar", true},
		{"123ABC", true},
		{"ABC123", true},
		{"ABC", true},
		{"S T R", true},
		{"fooBar", true},
		{"abacaba123", false},
		{"FÒÔBÀŘ", true},
		{"fòôbàř", false},
		{"Fòôbàř", true},
	}
	for _, test := range tests {
		actual := HasUpperCase(test.param)
		if actual != test.expected {
			t.Errorf("Expected HasUpperCase(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"-2147483648", true},          // Signed 32 Bit Min Int
		{"2147483647", true},           // Signed 32 Bit Max Int
		{"-2147483649", true},          // Signed 32 Bit Min Int - 1
		{"2147483648", true},           // Signed 32 Bit Max Int + 1
		{"4294967295", true},           // Unsigned 32 Bit Max Int
		{"4294967296", true},           // Unsigned 32 Bit Max Int + 1
		{"-9223372036854775808", true}, // Signed 64 Bit Min Int
		{"9223372036854775807", true},  // Signed 64 Bit Max Int
		{"-9223372036854775809", true}, // Signed 64 Bit Min Int - 1
		{"9223372036854775808", true},  // Signed 64 Bit Max Int + 1
		{"18446744073709551615", true}, // Unsigned 64 Bit Max Int
		{"18446744073709551616", true}, // Unsigned 64 Bit Max Int + 1
		{"", true},
		{"123", true},
		{"0", true},
		{"-0", true},
		{"+0", true},
		{"01", false},
		{"123.123", false},
		{" ", false},
		{"000", false},
	}
	for _, test := range tests {
		actual := IsInt(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsInt(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsBool(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", true},
		{"0", true},
		{"1", true},
		{"true", true},
		{"truE", false},
		{"TRUE", true},
		{"True", true},
		{"false", true},
		{"f", true},
		{"t", true},
		{"2", false},
		{"OK", false},
		{"YO", false},
	}
	for _, test := range tests {
		actual := IsBool(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsBool(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

var benchmarkIsInt bool

// REGEX
// BenchmarkIsInt-4   	 2000000	       775 ns/op	       0 B/op	       0 allocs/op
// BenchmarkIsInt-4   	 2000000	       769 ns/op	       0 B/op	       0 allocs/op
// BenchmarkIsInt-4   	 2000000	       750 ns/op	       0 B/op	       0 allocs/op
// BenchmarkIsInt-4   	50000000	        33.9 ns/op	       0 B/op	       0 allocs/op
// BenchmarkIsInt-4   	50000000	        33.3 ns/op	       0 B/op	       0 allocs/op
// BenchmarkIsInt-4   	50000000	        33.8 ns/op	       0 B/op	       0 allocs/op
func BenchmarkIsInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkIsInt = IsInt("-9223372036854775809")
	}
	if !benchmarkIsInt {
		b.Fatal("Must be an INT")
	}
}

func TestIsHash(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		algo     string
		expected bool
	}{
		{"3ca25ae354e192b26879f651a51d92aa8a34d8d3", "sha1", true},
		{"3ca25ae354e192b26879f651a51d34d8d3", "sha1", false},
		{"3ca25ae354e192b26879f651a51d92aa8a34d8d3", "Tiger160", true},
		{"3ca25ae354e192b26879f651a51d34d8d3", "ripemd160", false},
		{"579282cfb65ca1f109b78536effaf621b853c9f7079664a3fbe2b519f435898c", "sha256", true},
		{"579282cfb65ca1f109b78536effaf621b853c9f7079664a3fbe2b519f435898casfdsafsadfsdf", "sha256", false},
		{"bf547c3fc5841a377eb1519c2890344dbab15c40ae4150b4b34443d2212e5b04aa9d58865bf03d8ae27840fef430b891", "sha384", true},
		{"579282cfb65ca1f109b78536effaf621b853c9f7079664a3fbe2b519f435898casfdsafsadfsdf", "sha384", false},
		{"45bc5fa8cb45ee408c04b6269e9f1e1c17090c5ce26ffeeda2af097735b29953ce547e40ff3ad0d120e5361cc5f9cee35ea91ecd4077f3f589b4d439168f91b9", "sha512", true},
		{"579282cfb65ca1f109b78536effaf621b853c9f7079664a3fbe2b519f435898casfdsafsadfsdf", "sha512", false},
		{"46fc0125a148788a3ac1d649566fc04eb84a746f1a6e4fa7", "tiger192", true},
		{"46fc0125a148788a3ac1d649566fc04eb84a746f1a6$$%@^", "TIGER192", false},
		{"46fc0125a148788a3ac1d649566fc04eb84a746f1a6$$%@^", "SOMEHASH", false},
	}
	for _, test := range tests {
		actual := IsHash(test.param, test.algo)
		if actual != test.expected {
			t.Errorf("Expected IsHash(%q, %q) to be %v, got %v", test.param, test.algo, test.expected, actual)
		}
	}
}

func TestIsExistingEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"foo@bar.com", true},
		{"foo@bar.com.au", true},
		{"foo+bar@bar.com", true},
		{"foo@bar.coffee", false},
		{"foo@bar.coffee..coffee", false},
		{"invalidemail@", false},
		{"invalid.com", false},
		{"@invalid.com", false},
		{"NathAn.daVIeS@DomaIn.cOM", true},
		{"NATHAN.DAVIES@DOMAIN.CO.UK", true},
	}
	for _, test := range tests {
		actual := IsExistingEmail(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsExistingEmail(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsEmailRegexp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"foo@bar.com", true},
		{"x@x.x", true},
		{"foo@bar.com.au", true},
		{"foo+bar@bar.com", true},
		{"foo@bar.coffee", true},
		{"foo@bar.coffee..coffee", false},
		{"foo@bar.bar.coffee", true},
		{"foo@bar.中文网", true},
		{"invalidemail@", false},
		{"invalid.com", false},
		{"@invalid.com", false},
		{"test|123@m端ller.com", true},
		{"hans@m端ller.com", true},
		{"hans.m端ller@test.com", true},
		{"NathAn.daVIeS@DomaIn.cOM", true},
		{"NATHAN.DAVIES@DOMAIN.CO.UK", true},
	}
	for _, test := range tests {
		actual := IsEmailRegexp(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsEmailRegexp(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsEmailSimple(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"HiFish", false},
		{"foo@bar.com", true},
		{"x@x.x", true},
		{"foo@bar.com.au", true},
		{"foo+bar@bar.com", true},
		{"foo@bar.coffee", true},
		{"foo@bar.coffee..coffee", true},
		{"foo@bar.bar.coffee", true},
		{"foo@bar.中文网", true},
		{"invalidemail@", false},
		{"invalid.com", false},
		{"@invalid.com", false},
		{"test|123@m端ller.com", true},
		{"hans@m端ller.com", true},
		{"hans.m端ller@test.com", true},
		{"NathAn.daVIeS@DomaIn.cOM", true},
		{"NATHAN.DAVIES@DOMAIN.CO.UK", true},
		{"NATHAN.DAVIES@DOMAIN@.UK", false},
	}
	for _, test := range tests {
		actual := IsEmailSimple(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsEmailSimple(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"http://foo.bar#com", true},
		{"http://foobar.com", true},
		{"https://foobar.com", true},
		{"foobar.com", true},
		{"http://foobar.coffee/", true},
		{"http://foobar.中文网/", true},
		{"http://foobar.org/", true},
		{"http://foobar.ORG", true},
		{"http://foobar.org:8080/", true},
		{"ftp://foobar.ru/", true},
		{"ftp.foo.bar", true},
		{"http://user:pass@www.foobar.com/", true},
		{"http://user:pass@www.foobar.com/path/file", true},
		{"http://127.0.0.1/", true},
		{"http://duckduckgo.com/?q=%2F", true},
		{"http://localhost:3000/", true},
		{"http://foobar.com/?foo=bar#baz=qux", true},
		{"http://foobar.com?foo=bar", true},
		{"http://www.xn--froschgrn-x9a.net/", true},
		{"http://foobar.com/a-", true},
		{"http://foobar.پاکستان/", true},
		{"http://foobar.c_o_m", false},
		{"", false},
		{"xyz://foobar.com", false},
		// {"invalid.", false}, is it false like "localhost."?
		{".com", false},
		{"rtmp://foobar.com", false},
		{"http://www.foo_bar.com/", false},
		{"http://localhost:3000/", true},
		{"http://foobar.com#baz=qux", true},
		{"http://foobar.com/t$-_.+!*\\'(),", true},
		{"http://www.foobar.com/~foobar", true},
		{"http://www.-foobar.com/", false},
		{"http://www.foo---bar.com/", false},
		{"http://r6---snnvoxuioq6.googlevideo.com", true},
		{"mailto:someone@example.com", true},
		{"irc://irc.server.org/channel", false},
		{"irc://#channel@network", true},
		{"/abs/test/dir", false},
		{"./rel/test/dir", false},
		{"http://foo^bar.org", false},
		{"http://foo&*bar.org", false},
		{"http://foo&bar.org", false},
		{"http://foo bar.org", false},
		{"http://foo.bar.org", true},
		{"http://www.foo.bar.org", true},
		{"http://www.foo.co.uk", true},
		{"foo", false},
		{"http://.foo.com", false},
		{"http://,foo.com", false},
		{",foo.com", false},
		{"http://myservice.:9093/", true},
		// according to issues #62 #66
		{"https://pbs.twimg.com/profile_images/560826135676588032/j8fWrmYY_normal.jpeg", true},
		// according to #125
		{"http://prometheus-alertmanager.service.q:9093", true},
		{"aio1_alertmanager_container-63376c45:9093", true},
		{"https://www.logn-123-123.url.with.sigle.letter.d:12345/url/path/foo?bar=zzz#user", true},
		{"http://me.example.com", true},
		{"http://www.me.example.com", true},
		{"https://farm6.static.flickr.com", true},
		{"https://zh.wikipedia.org/wiki/Wikipedia:%E9%A6%96%E9%A1%B5", true},
		{"google", false},
		// According to #87
		{"http://hyphenated-host-name.example.co.in", true},
		{"http://cant-end-with-hyphen-.example.com", false},
		{"http://-cant-start-with-hyphen.example.com", false},
		{"http://www.domain-can-have-dashes.com", true},
		{"http://m.abcd.com/test.html", true},
		{"http://m.abcd.com/a/b/c/d/test.html?args=a&b=c", true},
		{"http://[::1]:9093", true},
		{"http://[::1]:909388", false},
		{"1200::AB00:1234::2552:7777:1313", false},
		{"http://[2001:db8:a0b:12f0::1]/index.html", true},
		{"http://[1200:0000:AB00:1234:0000:2552:7777:1313]", true},
		{"http://user:pass@[::1]:9093/a/b/c/?a=v#abc", true},
		{"https://127.0.0.1/a/b/c?a=v&c=11d", true},
		{"https://foo_bar.example.com", true},
		{"http://foo_bar.example.com", true},
		{"http://foo_bar_fizz_buzz.example.com", true},
		{"http://_cant_start_with_underescore", false},
		{"http://cant_end_with_underescore_", false},
		{"foo_bar.example.com", true},
		{"foo_bar_fizz_buzz.example.com", true},
		{"http://hello_world.example.com", true},
		// According to #212
		{"foo_bar-fizz-buzz:1313", true},
		{"foo_bar-fizz-buzz:13:13", false},
		{"foo_bar-fizz-buzz://1313", false},
	}
	for _, test := range tests {
		actual := IsURL(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsURL(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsRequestURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"http://foo.bar/#com", true},
		{"http://foobar.com", true},
		{"https://foobar.com", true},
		{"foobar.com", false},
		{"http://foobar.coffee/", true},
		{"http://foobar.中文网/", true},
		{"http://foobar.org/", true},
		{"http://foobar.org:8080/", true},
		{"ftp://foobar.ru/", true},
		{"http://user:pass@www.foobar.com/", true},
		{"http://127.0.0.1/", true},
		{"http://duckduckgo.com/?q=%2F", true},
		{"http://localhost:3000/", true},
		{"http://foobar.com/?foo=bar#baz=qux", true},
		{"http://foobar.com?foo=bar", true},
		{"http://www.xn--froschgrn-x9a.net/", true},
		{"", false},
		{"xyz://foobar.com", true},
		{"invalid.", false},
		{".com", false},
		{"rtmp://foobar.com", true},
		{"http://www.foo_bar.com/", true},
		{"http://localhost:3000/", true},
		{"http://foobar.com/#baz=qux", true},
		{"http://foobar.com/t$-_.+!*\\'(),", true},
		{"http://www.foobar.com/~foobar", true},
		{"http://www.-foobar.com/", true},
		{"http://www.foo---bar.com/", true},
		{"mailto:someone@example.com", true},
		{"irc://irc.server.org/channel", true},
		{"/abs/test/dir", false},
		{"./rel/test/dir", false},
	}
	for _, test := range tests {
		actual := IsRequestURL(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsRequestURL(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsRequestURI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"http://foo.bar/#com", true},
		{"http://foobar.com", true},
		{"https://foobar.com", true},
		{"foobar.com", false},
		{"http://foobar.coffee/", true},
		{"http://foobar.中文网/", true},
		{"http://foobar.org/", true},
		{"http://foobar.org:8080/", true},
		{"ftp://foobar.ru/", true},
		{"http://user:pass@www.foobar.com/", true},
		{"http://127.0.0.1/", true},
		{"http://duckduckgo.com/?q=%2F", true},
		{"http://localhost:3000/", true},
		{"http://foobar.com/?foo=bar#baz=qux", true},
		{"http://foobar.com?foo=bar", true},
		{"http://www.xn--froschgrn-x9a.net/", true},
		{"xyz://foobar.com", true},
		{"invalid.", false},
		{".com", false},
		{"rtmp://foobar.com", true},
		{"http://www.foo_bar.com/", true},
		{"http://localhost:3000/", true},
		{"http://foobar.com/#baz=qux", true},
		{"http://foobar.com/t$-_.+!*\\'(),", true},
		{"http://www.foobar.com/~foobar", true},
		{"http://www.-foobar.com/", true},
		{"http://www.foo---bar.com/", true},
		{"mailto:someone@example.com", true},
		{"irc://irc.server.org/channel", true},
		{"/abs/test/dir", true},
		{"./rel/test/dir", false},
	}
	for _, test := range tests {
		actual := IsRequestURI(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsRequestURI(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsFloat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"  ", false},
		{"-.123", true},
		{"abacaba", false},
		{"1f", false},
		{"-1f", false},
		{"+1f", false},
		{"123", true},
		{"123.", true},
		{"123.123", true},
		{"-123.123", true},
		{"+123.123", true},
		{"0.123", true},
		{"-0.123", true},
		{"+0.123", true},
		{".0", true},
		{"01.123", true},
		{"-0.22250738585072011e-307", true},
		{"+0.22250738585072011e-307", true},
	}
	for _, test := range tests {
		actual := IsFloat(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsFloat(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsHexadecimal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"abcdefg", false},
		{"", false},
		{"..", false},
		{"deadBEEF", true},
		{"ff0044", true},
	}
	for _, test := range tests {
		actual := IsHexadecimal(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsHexadecimal(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsHexcolor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"#ff", false},
		{"fff0", false},
		{"#ff12FG", false},
		{"CCccCC", true},
		{"fff", true},
		{"#f00", true},
	}
	for _, test := range tests {
		actual := IsHexcolor(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsHexcolor(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsRGBcolor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"rgb(0,31,255)", true},
		{"rgb(1,349,275)", false},
		{"rgb(01,31,255)", false},
		{"rgb(0.6,31,255)", false},
		{"rgba(0,31,255)", false},
		{"rgb(0,  31, 255)", true},
	}
	for _, test := range tests {
		actual := IsRGBcolor(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsRGBcolor(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsNull(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"abacaba", false},
		{"", true},
	}
	for _, test := range tests {
		actual := IsNull(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsNull(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsNotEmpty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"abacaba", true},
		{"", false},
	}
	for _, test := range tests {
		actual := IsNotEmpty(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsNull(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsNotEmptyTrimSpace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"abacaba ", true},
		{" abacaba", true},
		{"\n abacaba", true},
		{"\n abacaba\t", true},
		{"", false},
		{" ", false},
		{"\t", false},
		{"\t \t", false},
		{"\t \t\r", false},
		{"\t \t\r\n", false},
		{"\n", false},
	}
	for _, test := range tests {
		actual := IsNotEmptyTrimSpace(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsNull(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsDivisibleBy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param1   string
		param2   string
		expected bool
	}{
		{"4", "2", true},
		{"100", "10", true},
		{"", "1", true},
		{"123", "foo", false},
		{"123", "0", false},
	}
	for _, test := range tests {
		actual := IsDivisibleBy(test.param1, test.param2)
		if actual != test.expected {
			t.Errorf("Expected IsDivisibleBy(%q, %q) to be %v, got %v", test.param1, test.param2, test.expected, actual)
		}
	}
}

// This small example illustrate how to work with IsDivisibleBy function.
func ExampleIsDivisibleBy() {
	println("1024 is divisible by 64: ", IsDivisibleBy("1024", "64"))
}

func TestIsByteLength(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param1   string
		param2   int
		param3   int
		expected bool
	}{
		{"abacaba", 100, -1, false},
		{"abacaba", 1, 3, false},
		{"abacaba", 1, 7, true},
		{"abacaba", 0, 8, true},
		{"\ufff0", 1, 1, false},
	}
	for _, test := range tests {
		actual := IsByteLength(test.param1, test.param2, test.param3)
		if actual != test.expected {
			t.Errorf("Expected IsByteLength(%q, %q, %q) to be %v, got %v", test.param1, test.param2, test.param3, test.expected, actual)
		}
	}
}

func TestIsMultibyte(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"abc", false},
		{"123", false},
		{"<>@;.-=", false},
		{"ひらがな・カタカナ、．漢字", true},
		{"あいうえお foobar", true},
		{"test＠example.com", true},
		{"test＠example.com", true},
		{"1234abcDEｘｙｚ", true},
		{"ｶﾀｶﾅ", true},
		{"", true},
	}
	for _, test := range tests {
		actual := IsMultibyte(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsMultibyte(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsASCII(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", true},
		{"ｆｏｏbar", false},
		{"ｘｙｚ０９８", false},
		{"１２３456", false},
		{"ｶﾀｶﾅ", false},
		{"foobar", true},
		{"0987654321", true},
		{"test@example.com", true},
		{"1234abcDEF", true},
		{"", true},
	}
	for _, test := range tests {
		actual := IsASCII(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsASCII(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsPrintableASCII(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", true},
		{"ｆｏｏbar", false},
		{"ｘｙｚ０９８", false},
		{"１２３456", false},
		{"ｶﾀｶﾅ", false},
		{"foobar", true},
		{"0987654321", true},
		{"test@example.com", true},
		{"1234abcDEF", true},
		{"newline\n", false},
		{"\x19test\x7F", false},
	}
	for _, test := range tests {
		actual := IsPrintableASCII(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsPrintableASCII(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsFullWidth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", true},
		{"abc", false},
		{"abc123", false},
		{"!\"#$%&()<>/+=-_? ~^|.,@`{}[]", false},
		{"ひらがな・カタカナ、．漢字", true},
		{"３ー０　ａ＠ｃｏｍ", true},
		{"Ｆｶﾀｶﾅﾞﾬ", true},
		{"Good＝Parts", true},
		{"", true},
	}
	for _, test := range tests {
		actual := IsFullWidth(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsFullWidth(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsHalfWidth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", true},
		{"あいうえお", false},
		{"００１１", false},
		{"!\"#$%&()<>/+=-_? ~^|.,@`{}[]", true},
		{"l-btn_02--active", true},
		{"abc123い", true},
		{"ｶﾀｶﾅﾞﾬ￩", true},
		{"", true},
	}
	for _, test := range tests {
		actual := IsHalfWidth(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsHalfWidth(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsVariableWidth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", true},
		{"ひらがなカタカナ漢字ABCDE", true},
		{"３ー０123", true},
		{"Ｆｶﾀｶﾅﾞﾬ", true},
		{"", true},
		{"Good＝Parts", true},
		{"abc", false},
		{"abc123", false},
		{"!\"#$%&()<>/+=-_? ~^|.,@`{}[]", false},
		{"ひらがな・カタカナ、．漢字", false},
		{"１２３４５６", false},
		{"ｶﾀｶﾅﾞﾬ", false},
	}
	for _, test := range tests {
		actual := IsVariableWidth(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsVariableWidth(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsUUID(t *testing.T) {
	t.Parallel()

	// Tests without version
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"xxxa987fbc9-4bed-3078-cf07-9141ba07c9f3", false},
		{"a987fbc9-4bed-3078-cf07-9141ba07c9f3xxx", false},
		{"a987fbc94bed3078cf079141ba07c9f3", false},
		{"934859", false},
		{"987fbc9-4bed-3078-cf07a-9141ba07c9f3", false},
		{"aaaaaaaa-1111-1111-aaag-111111111111", false},
		{"a987fbc9-4bed-3078-cf07-9141ba07c9f3", true},
	}
	for _, test := range tests {
		actual := IsUUID(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsUUID(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}

	// UUID ver. 3
	tests = []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"412452646", false},
		{"xxxa987fbc9-4bed-3078-cf07-9141ba07c9f3", false},
		{"fbc9-4bed-3078-cf07-9141ba07c9f3", false},
		{"a987fbc9-4bed-4078-8f07-9141ba07c9f3", false},
		{"a987fbc9-4bed-3078-cf07-9141ba07c9f3", true},
		{"a987fbc9-4bed-3078-cf07-9141ba07c9g3", false},
	}
	for _, test := range tests {
		actual := IsUUIDv3(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsUUIDv3(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}

	// UUID ver. 4
	tests = []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"xxxa987fbc9-4bed-3078-cf07-9141ba07c9f3", false},
		{"a987fbc9-4bed-5078-af07-9141ba07c9f3", false},
		{"934859", false},
		{"57b73598-8764-4ad0-a76a-679bb6640eb1", true},
		{"625e63f3-58f5-40b7-83a1-a72ad31acffb", true},
	}
	for _, test := range tests {
		actual := IsUUIDv4(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsUUIDv4(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}

	// UUID ver. 5
	tests = []struct {
		param    string
		expected bool
	}{

		{"", false},
		{"xxxa987fbc9-4bed-3078-cf07-9141ba07c9f3", false},
		{"9c858901-8a57-4791-81fe-4c455b099bc9", false},
		{"a987fbc9-4bed-3078-cf07-9141ba07c9f3", false},
		{"987fbc97-4bed-5078-af07-9141ba07c9f3", true},
		{"987fbc97-4bed-5078-9f07-9141ba07c9f3", true},
	}
	for _, test := range tests {
		actual := IsUUIDv5(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsUUIDv5(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsCreditCard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"foo", false},
		{"5398228707871528", false},
		{"375556917985515", true},
		{"36050234196908", true},
		{"4716461583322103", true},
		{"4716-2210-5188-5662", true},
		{"4929 7226 5379 7141", true},
		{"5398228707871527", true},
	}
	for _, test := range tests {
		actual := IsCreditCard(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsCreditCard(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsISBN(t *testing.T) {
	t.Parallel()

	// Without version
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"foo", false},
		{"3836221195", true},
		{"1-61729-085-8", true},
		{"3 423 21412 0", true},
		{"3 401 01319 X", true},
		{"9784873113685", true},
		{"978-4-87311-368-5", true},
		{"978 3401013190", true},
		{"978-3-8362-2119-1", true},
	}
	for _, test := range tests {
		actual := IsISBN(test.param, -1)
		if actual != test.expected {
			t.Errorf("Expected IsISBN(%q, -1) to be %v, got %v", test.param, test.expected, actual)
		}
	}

	// ISBN 10
	tests = []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"foo", false},
		{"3423214121", false},
		{"978-3836221191", false},
		{"3-423-21412-1", false},
		{"3 423 21412 1", false},
		{"3836221195", true},
		{"1-61729-085-8", true},
		{"3 423 21412 0", true},
		{"3 401 01319 X", true},
	}
	for _, test := range tests {
		actual := IsISBN10(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsISBN10(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}

	// ISBN 13
	tests = []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"foo", false},
		{"3-8362-2119-5", false},
		{"01234567890ab", false},
		{"978 3 8362 2119 0", false},
		{"9784873113685", true},
		{"978-4-87311-368-5", true},
		{"978 3401013190", true},
		{"978-3-8362-2119-1", true},
	}
	for _, test := range tests {
		actual := IsISBN13(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsISBN13(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsDataURI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"data:image/png;base64,TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdC4=", true},
		{"data:text/plain;base64,Vml2YW11cyBmZXJtZW50dW0gc2VtcGVyIHBvcnRhLg==", true},
		{"image/gif;base64,U3VzcGVuZGlzc2UgbGVjdHVzIGxlbw==", false},
		{"data:image/gif;base64,MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuMPNS1Ufof9EW/M98FNw" +
			"UAKrwflsqVxaxQjBQnHQmiI7Vac40t8x7pIb8gLGV6wL7sBTJiPovJ0V7y7oc0Ye" +
			"rhKh0Rm4skP2z/jHwwZICgGzBvA0rH8xlhUiTvcwDCJ0kc+fh35hNt8srZQM4619" +
			"FTgB66Xmp4EtVyhpQV+t02g6NzK72oZI0vnAvqhpkxLeLiMCyrI416wHm5Tkukhx" +
			"QmcL2a6hNOyu0ixX/x2kSFXApEnVrJ+/IxGyfyw8kf4N2IZpW5nEP847lpfj0SZZ" +
			"Fwrd1mnfnDbYohX2zRptLy2ZUn06Qo9pkG5ntvFEPo9bfZeULtjYzIl6K8gJ2uGZ" + "HQIDAQAB", true},
		{"data:image/png;base64,12345", false},
		{"", false},
		{"data:text,:;base85,U3VzcGVuZGlzc2UgbGVjdHVzIGxlbw==", false},
	}
	for _, test := range tests {
		actual := IsDataURI(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsDataURI(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsBase64(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdC4=", true},
		{"TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdC4===", false},
		{"TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdC4", false},
		{"Vml2YW11cyBmZXJtZW50dW0gc2VtcGVyIHBvcnRhLg==", true},
		{"U3VzcGVuZGlzc2UgbGVjdHVzIGxlbw==", true},
		{"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuMPNS1Ufof9EW/M98FNw" +
			"UAKrwflsqVxaxQjBQnHQmiI7Vac40t8x7pIb8gLGV6wL7sBTJiPovJ0V7y7oc0Ye" +
			"rhKh0Rm4skP2z/jHwwZICgGzBvA0rH8xlhUiTvcwDCJ0kc+fh35hNt8srZQM4619" +
			"FTgB66Xmp4EtVyhpQV+t02g6NzK72oZI0vnAvqhpkxLeLiMCyrI416wHm5Tkukhx" +
			"QmcL2a6hNOyu0ixX/x2kSFXApEnVrJ+/IxGyfyw8kf4N2IZpW5nEP847lpfj0SZZ" +
			"Fwrd1mnfnDbYohX2zRptLy2ZUn06Qo9pkG5ntvFEPo9bfZeULtjYzIl6K8gJ2uGZ" + "HQIDAQAB", true},
		{"12345", false},
		{"", false},
		{"Vml2YW11cyBmZXJtZtesting123", false},
	}
	for _, test := range tests {
		actual := IsBase64(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsBase64(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsISO3166Alpha2(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"ABCD", false},
		{"A", false},
		{"AC", false},
		{"AP", false},
		{"GER", false},
		{"NU", true},
		{"DE", true},
		{"JP", true},
		{"JPN", false},
		{"ZWE", false},
		{"GER", false},
		{"DEU", false},
	}
	for _, test := range tests {
		actual := IsISO3166Alpha2(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsISO3166Alpha2(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsISO3166Alpha3(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"ABCD", false},
		{"A", false},
		{"AC", false},
		{"AP", false},
		{"NU", false},
		{"DE", false},
		{"JP", false},
		{"ZWE", true},
		{"JPN", true},
		{"GER", false},
		{"DEU", true},
	}
	for _, test := range tests {
		actual := IsISO3166Alpha3(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsISO3166Alpha3(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsISO693Alpha2(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"abcd", false},
		{"a", false},
		{"ac", false},
		{"ap", false},
		{"de", true},
		{"DE", false},
		{"mk", true},
		{"mac", false},
		{"sw", true},
		{"SW", false},
		{"ger", false},
		{"deu", false},
	}
	for _, test := range tests {
		actual := IsISO693Alpha2(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsISO693Alpha2(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsISO693Alpha3b(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"abcd", false},
		{"a", false},
		{"ac", false},
		{"ap", false},
		{"de", false},
		{"DE", false},
		{"mkd", false},
		{"mac", true},
		{"sw", false},
		{"SW", false},
		{"ger", true},
		{"deu", false},
	}
	for _, test := range tests {
		actual := IsISO693Alpha3b(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsISO693Alpha3b(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsLocale(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"1", false},
		{"-", false},
		{"a", false},
		{"abcd", false},
		{"de_DE", true},
		{"de_de", true},
		{"dede", true},
		{"dech", true},
		{"deCH", true},
		{"de-CH", true},
		{"zhhanscn", true},
		{"xhhanscn", false},
		{"", false},
	}
	// Detect race conditions
	bgwork.Wait(len(tests), func(idx int) {
		test := tests[idx]
		actual := IsLocale(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsLocale(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	})
}

func TestIsIP(t *testing.T) {
	t.Parallel()

	// Without version
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"127.0.0.1", true},
		{"0.0.0.0", true},
		{"255.255.255.255", true},
		{"1.2.3.4", true},
		{"::1", true},
		{"2001:db8:0000:1:1:1:1:1", true},
		{"300.0.0.0", false},
	}
	for _, test := range tests {
		actual := IsIP(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsIP(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}

	// IPv4
	tests = []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"127.0.0.1", true},
		{"0.0.0.0", true},
		{"255.255.255.255", true},
		{"1.2.3.4", true},
		{"::1", false},
		{"2001:db8:0000:1:1:1:1:1", false},
		{"300.0.0.0", false},
	}
	for _, test := range tests {
		actual := IsIPv4(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsIPv4(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}

	// IPv6
	tests = []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"127.0.0.1", false},
		{"0.0.0.0", false},
		{"255.255.255.255", false},
		{"1.2.3.4", false},
		{"::1", true},
		{"2001:db8:0000:1:1:1:1:1", true},
		{"300.0.0.0", false},
	}
	for _, test := range tests {
		actual := IsIPv6(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsIPv6(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsPort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"1", true},
		{"65535", true},
		{"0", false},
		{"65536", false},
		{"65538", false},
	}

	for _, test := range tests {
		actual := IsPort(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsPort(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsDNSName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"localhost", true},
		{"a.bc", true},
		{"a.b.", true},
		{"a.b..", false},
		{"localhost.local", true},
		{"localhost.localdomain.intern", true},
		{"l.local.intern", true},
		{"ru.link.n.svpncloud.com", true},
		{"-localhost", false},
		{"localhost.-localdomain", false},
		{"localhost.localdomain.-int", false},
		{"_localhost", true},
		{"localhost._localdomain", true},
		{"localhost.localdomain._int", true},
		{"lÖcalhost", false},
		{"localhost.lÖcaldomain", false},
		{"localhost.localdomain.üntern", false},
		{"__", true},
		{"localhost/", false},
		{"127.0.0.1", false},
		{"[::1]", false},
		{"50.50.50.50", false},
		{"localhost.localdomain.intern:65535", false},
		{"漢字汉字", false},
		{"www.jubfvq1v3p38i51622y0dvmdk1mymowjyeu26gbtw9andgynj1gg8z3msb1kl5z6906k846pj3sulm4kiyk82ln5teqj9nsht59opr0cs5ssltx78lfyvml19lfq1wp4usbl0o36cmiykch1vywbttcus1p9yu0669h8fj4ll7a6bmop505908s1m83q2ec2qr9nbvql2589adma3xsq2o38os2z3dmfh2tth4is4ixyfasasasefqwe4t2ub2fz1rme.de", false},
	}

	for _, test := range tests {
		actual := IsDNSName(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsDNS(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsHost(t *testing.T) {
	t.Parallel()
	tests := []struct {
		param    string
		expected bool
	}{
		{"localhost", true},
		{"localhost.localdomain", true},
		{"2001:db8:0000:1:1:1:1:1", true},
		{"::1", true},
		{"play.golang.org", true},
		{"localhost.localdomain.intern:65535", false},
		{"-[::1]", false},
		{"-localhost", false},
		{".localhost", false},
	}
	for _, test := range tests {
		actual := IsHost(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsHost(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsDialString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"localhost.local:1", true},
		{"localhost.localdomain:9090", true},
		{"localhost.localdomain.intern:65535", true},
		{"127.0.0.1:30000", true},
		{"[::1]:80", true},
		{"[1200::AB00:1234::2552:7777:1313]:22", false},
		{"-localhost:1", false},
		{"localhost.-localdomain:9090", false},
		{"localhost.localdomain.-int:65535", false},
		{"localhost.loc:100000", false},
		{"漢字汉字:2", false},
		{"www.jubfvq1v3p38i51622y0dvmdk1mymowjyeu26gbtw9andgynj1gg8z3msb1kl5z6906k846pj3sulm4kiyk82ln5teqj9nsht59opr0cs5ssltx78lfyvml19lfq1wp4usbl0o36cmiykch1vywbttcus1p9yu0669h8fj4ll7a6bmop505908s1m83q2ec2qr9nbvql2589adma3xsq2o38os2z3dmfh2tth4is4ixyfasasasefqwe4t2ub2fz1rme.de:20000", false},
	}

	for _, test := range tests {
		actual := IsDialString(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsDialString(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsMAC(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"3D:F2:C9:A6:B3:4F", true},
		{"3D-F2-C9-A6-B3:4F", false},
		{"123", false},
		{"", false},
		{"abacaba", false},
	}
	for _, test := range tests {
		actual := IsMAC(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsMAC(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestFilePath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
		osType   int
	}{
		{"c:\\" + strings.Repeat("a", 32767), true, Win}, // See http://msdn.microsoft.com/en-us/library/aa365247(VS.85).aspx#maxpath
		{"c:\\" + strings.Repeat("a", 32768), false, Win},
		{"c:\\path\\file (x86)\bar", true, Win},
		{"c:\\path\\file", true, Win},
		{"c:\\path\\file:exe", false, Unknown},
		{"C:\\", true, Win},
		{"c:\\path\\file\\", true, Win},
		{"c:/path/file/", false, Unknown},
		{"/path/file/", true, Unix},
		{"/path/file:SAMPLE/", true, Unix},
		{"/path/file:/.txt", true, Unix},
		{"/path", true, Unix},
		{"/path/__bc/file.txt", true, Unix},
		{"/path/a--ac/file.txt", true, Unix},
		{"/_path/file.txt", true, Unix},
		{"/path/__bc/file.txt", true, Unix},
		{"/path/a--ac/file.txt", true, Unix},
		{"/__path/--file.txt", true, Unix},
		{"/path/a bc", true, Unix},
	}
	for _, test := range tests {
		actual, osType := IsFilePath(test.param)
		if actual != test.expected || osType != test.osType {
			t.Errorf("Expected IsFilePath(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsLatitude(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"-90.000", true},
		{"+90", true},
		{"47.1231231", true},
		{"+99.9", false},
		{"108", false},
	}
	for _, test := range tests {
		actual := IsLatitude(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsLatitude(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsLongitude(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"-180.000", true},
		{"180.1", false},
		{"+73.234", true},
		{"+382.3811", false},
		{"23.11111111", true},
	}
	for _, test := range tests {
		actual := IsLongitude(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsLongitude(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsSSN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"00-90-8787", false},
		{"66690-76", false},
		{"191 60 2869", true},
		{"191-60-2869", true},
	}
	for _, test := range tests {
		actual := IsSSN(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsSSN(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsMongoID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"507f1f77bcf86cd799439011", true},
		{"507f1f77bcf86cd7994390", false},
		{"507f1f77bcf86cd79943901z", false},
		{"507f1f77bcf86cd799439011 ", false},
		{"", false},
	}
	for _, test := range tests {
		actual := IsMongoID(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsMongoID(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsSemver(t *testing.T) {
	t.Parallel()
	tests := []struct {
		param    string
		expected bool
	}{
		{"v1.0.0", true},
		{"1.0.0", true},
		{"1.1.01", false},
		{"1.01.0", false},
		{"01.1.0", false},
		{"v1.1.01", false},
		{"v1.01.0", false},
		{"v01.1.0", false},
		{"1.0.0-alpha", true},
		{"1.0.0-alpha.1", true},
		{"1.0.0-0.3.7", true},
		{"1.0.0-0.03.7", false},
		{"1.0.0-00.3.7", false},
		{"1.0.0-x.7.z.92", true},
		{"1.0.0-alpha+001", true},
		{"1.0.0+20130313144700", true},
		{"1.0.0-beta+exp.sha.5114f85", true},
		{"1.0.0-beta+exp.sha.05114f85", true},
		{"1.0.0-+beta", false},
		{"1.0.0-b+-9+eta", false},
		{"v+1.8.0-b+-9+eta", false},
	}
	for _, test := range tests {
		actual := IsSemver(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsSemver(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsTime(t *testing.T) {
	t.Parallel()
	tests := []struct {
		param    string
		format   string
		expected bool
	}{
		{"2016-12-31 11:00", time.RFC3339, false},
		{"2016-12-31 11:00:00", time.RFC3339, false},
		{"2016-12-31T11:00", time.RFC3339, false},
		{"2016-12-31T11:00:00", time.RFC3339, false},
		{"2016-12-31T11:00:00Z", time.RFC3339, true},
		{"2016-12-31T11:00:00+01:00", time.RFC3339, true},
		{"2016-12-31T11:00:00-01:00", time.RFC3339, true},
		{"2016-12-31T11:00:00.05Z", time.RFC3339, true},
		{"2016-12-31T11:00:00.05-01:00", time.RFC3339, true},
		{"2016-12-31T11:00:00.05+01:00", time.RFC3339, true},
		{"2016-12-31T11:00:00", RF3339WithoutZone, true},
		{"2016-12-31T11:00:00Z", RF3339WithoutZone, false},
		{"2016-12-31T11:00:00+01:00", RF3339WithoutZone, false},
		{"2016-12-31T11:00:00-01:00", RF3339WithoutZone, false},
		{"2016-12-31T11:00:00.05Z", RF3339WithoutZone, false},
		{"2016-12-31T11:00:00.05-01:00", RF3339WithoutZone, false},
		{"2016-12-31T11:00:00.05+01:00", RF3339WithoutZone, false},
	}
	for _, test := range tests {
		actual := IsTime(test.param, test.format)
		if actual != test.expected {
			t.Errorf("Expected IsTime(%q, time.RFC3339) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsRFC3339(t *testing.T) {
	t.Parallel()
	tests := []struct {
		param    string
		expected bool
	}{
		{"2016-12-31 11:00", false},
		{"2016-12-31 11:00:00", false},
		{"2016-12-31T11:00", false},
		{"2016-12-31T11:00:00", false},
		{"2016-12-31T11:00:00Z", true},
		{"2016-12-31T11:00:00+01:00", true},
		{"2016-12-31T11:00:00-01:00", true},
		{"2016-12-31T11:00:00.05Z", true},
		{"2016-12-31T11:00:00.05-01:00", true},
		{"2016-12-31T11:00:00.05+01:00", true},
	}
	for _, test := range tests {
		actual := IsRFC3339(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsRFC3339(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsISO4217(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"ABCD", false},
		{"A", false},
		{"ZZZ", false},
		{"usd", false},
		{"USD", true},
	}
	for _, test := range tests {
		actual := IsISO4217(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsISO4217(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestByteLength(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value    string
		min      int
		max      int
		expected bool
	}{
		{"123456", 0, 100, true},
		{"1239999", 0, 0, false},
		{"1239asdfasf99", 100, 200, false},
		{"1239999asdff29", 10, 30, true},
		{"你", 0, 1, false},
	}
	for _, test := range tests {
		actual := ByteLength([]byte(test.value), test.min, test.max)
		if actual != test.expected {
			t.Errorf("Expected ByteLength(%s, %d, %d) to be %v, got %v", test.value, test.min, test.max, test.expected, actual)
		}
	}
}

func TestStringLength(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value    string
		min      int
		max      int
		expected bool
	}{
		{"123456", 0, 100, true},
		{"1239999", 0, 0, false},
		{"1239asdfasf99", 100, 200, false},
		{"1239999asdff29", 10, 30, true},
		{"あいうえお", 0, 5, true},
		{"あいうえおか", 0, 5, false},
		{"あいうえお", 0, 0, false},
		{"あいうえ", 5, 10, false},
	}
	for _, test := range tests {
		actual := StringLength(test.value, test.min, test.max)
		if actual != test.expected {
			t.Errorf("Expected StringLength(%s, %d, %d) to be %v, got %v", test.value, test.min, test.max, test.expected, actual)
		}
	}
}

func TestIsIn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value    string
		params   []string
		expected bool
	}{
		{"PRESENT", []string{"PRESENT"}, true},
		{"PRESENT", []string{"PRESENT", "PRÉSENTE", "NOTABSENT"}, true},
		{"PRÉSENTE", []string{"PRESENT", "PRÉSENTE", "NOTABSENT"}, true},
		{"PRESENT", []string{}, false},
		{"PRESENT", nil, false},
		{"ABSENT", []string{"PRESENT", "PRÉSENTE", "NOTABSENT"}, false},
		{"", []string{"PRESENT", "PRÉSENTE", "NOTABSENT"}, false},
	}
	for _, test := range tests {
		actual := IsIn(test.value, test.params...)
		if actual != test.expected {
			t.Errorf("Expected IsIn(%s, %v) to be %v, got %v", test.value, test.params, test.expected, actual)
		}
	}
}

func TestIsCIDR(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"193.168.3.20/7", true},
		{"2001:db8::/32", true},
		{"2001:0db8:85a3:0000:0000:8a2e:0370:7334/64", true},
		{"193.138.3.20/60", false},
		{"500.323.2.23/43", false},
		{"", false},
	}
	for _, test := range tests {
		actual := IsCIDR(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsCIDR(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestIsRsaPublicKey(t *testing.T) {
	tests := []struct {
		rsastr   string
		keylen   int
		expected bool
	}{
		{`fubar`, 2048, false},
		{`MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvncDCeibmEkabJLmFec7x9y86RP6dIvkVxxbQoOJo06E+p7tH6vCmiGHKnuu
XwKYLq0DKUE3t/HHsNdowfD9+NH8caLzmXqGBx45/Dzxnwqz0qYq7idK+Qff34qrk/YFoU7498U1Ee7PkKb7/VE9BmMEcI3uoKbeXCbJRI
HoTp8bUXOpNTSUfwUNwJzbm2nsHo2xu6virKtAZLTsJFzTUmRd11MrWCvj59lWzt1/eIMN+ekjH8aXeLOOl54CL+kWp48C+V9BchyKCShZ
B7ucimFvjHTtuxziXZQRO7HlcsBOa0WwvDJnRnskdyoD31s4F4jpKEYBJNWTo63v6lUvbQIDAQAB`, 2048, true},
		{`MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvncDCeibmEkabJLmFec7x9y86RP6dIvkVxxbQoOJo06E+p7tH6vCmiGHKnuu
XwKYLq0DKUE3t/HHsNdowfD9+NH8caLzmXqGBx45/Dzxnwqz0qYq7idK+Qff34qrk/YFoU7498U1Ee7PkKb7/VE9BmMEcI3uoKbeXCbJRI
HoTp8bUXOpNTSUfwUNwJzbm2nsHo2xu6virKtAZLTsJFzTUmRd11MrWCvj59lWzt1/eIMN+ekjH8aXeLOOl54CL+kWp48C+V9BchyKCShZ
B7ucimFvjHTtuxziXZQRO7HlcsBOa0WwvDJnRnskdyoD31s4F4jpKEYBJNWTo63v6lUvbQIDAQAB`, 1024, false},
		{`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvncDCeibmEkabJLmFec7
x9y86RP6dIvkVxxbQoOJo06E+p7tH6vCmiGHKnuuXwKYLq0DKUE3t/HHsNdowfD9
+NH8caLzmXqGBx45/Dzxnwqz0qYq7idK+Qff34qrk/YFoU7498U1Ee7PkKb7/VE9
BmMEcI3uoKbeXCbJRIHoTp8bUXOpNTSUfwUNwJzbm2nsHo2xu6virKtAZLTsJFzT
UmRd11MrWCvj59lWzt1/eIMN+ekjH8aXeLOOl54CL+kWp48C+V9BchyKCShZB7uc
imFvjHTtuxziXZQRO7HlcsBOa0WwvDJnRnskdyoD31s4F4jpKEYBJNWTo63v6lUv
bQIDAQAB
-----END PUBLIC KEY-----`, 2048, true},
		{`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvncDCeibmEkabJLmFec7
x9y86RP6dIvkVxxbQoOJo06E+p7tH6vCmiGHKnuuXwKYLq0DKUE3t/HHsNdowfD9
+NH8caLzmXqGBx45/Dzxnwqz0qYq7idK+Qff34qrk/YFoU7498U1Ee7PkKb7/VE9
BmMEcI3uoKbeXCbJRIHoTp8bUXOpNTSUfwUNwJzbm2nsHo2xu6virKtAZLTsJFzT
UmRd11MrWCvj59lWzt1/eIMN+ekjH8aXeLOOl54CL+kWp48C+V9BchyKCShZB7uc
imFvjHTtuxziXZQRO7HlcsBOa0WwvDJnRnskdyoD31s4F4jpKEYBJNWTo63v6lUv
bQIDAQAB
-----END PUBLIC KEY-----`, 4096, false},
	}
	for i, test := range tests {
		actual := IsRsaPublicKey(test.rsastr, test.keylen)
		if actual != test.expected {
			t.Errorf("Expected TestIsRsaPublicKey(%d, %d) to be %v, got %v", i, test.keylen, test.expected, actual)
		}
	}
}
