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

package null

import (
	"bytes"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/weiwolves/pkg/util/assert"
	"github.com/gogo/protobuf/proto"
)

var (
	_ fmt.GoStringer             = (*String)(nil)
	_ json.Marshaler             = (*String)(nil)
	_ json.Unmarshaler           = (*String)(nil)
	_ encoding.BinaryMarshaler   = (*String)(nil)
	_ encoding.BinaryUnmarshaler = (*String)(nil)
	_ encoding.TextMarshaler     = (*String)(nil)
	_ encoding.TextUnmarshaler   = (*String)(nil)
	_ driver.Valuer              = (*String)(nil)
	_ proto.Marshaler            = (*String)(nil)
	_ proto.Unmarshaler          = (*String)(nil)
	_ proto.Sizer                = (*String)(nil)
	_ protoMarshalToer           = (*String)(nil)
)

var (
	stringJSON      = []byte(`"test"`)
	blankStringJSON = []byte(`""`)
	nullStringJSON  = []byte(`{"String":"test","Valid":true}`)

	nullJSON    = []byte(sqlStrNullLC)
	invalidJSON = []byte(`:)`)
)

func TestStringFrom(t *testing.T) {
	str := MakeString("test")
	assertStr(t, str, "MakeString() string")
	assert.Exactly(t, 8, str.Size())

	zero := MakeString("")
	if !zero.Valid {
		t.Error("MakeString(0)", "is invalid, but should be valid")
	}
	assert.Exactly(t, 2, zero.Size())
}

func TestNullString_JsonUnmarshal(t *testing.T) {
	var str String
	maybePanic(json.Unmarshal(stringJSON, &str))
	assertStr(t, str, "string json")

	var ns String
	maybePanic(json.Unmarshal(nullStringJSON, &ns))
	assertStr(t, ns, "sql.String json")

	var blank String
	maybePanic(json.Unmarshal(blankStringJSON, &blank))
	if !blank.Valid {
		t.Error("blank string should be valid")
	}

	var null String
	maybePanic(json.Unmarshal(nullJSON, &null))
	assertNullStr(t, null, "null json")

	var badType String
	err := json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullStr(t, badType, "wrong type json")

	var invalid String
	err = invalid.UnmarshalJSON(invalidJSON)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("expected json.SyntaxError, not %T", err)
	}
	assertNullStr(t, invalid, "invalid json")
}

func TestNullString_TextUnmarshal(t *testing.T) {
	var str String
	err := str.UnmarshalText([]byte("test"))
	maybePanic(err)
	assertStr(t, str, "UnmarshalText() string")

	var null String
	err = null.UnmarshalText([]byte(""))
	maybePanic(err)
	assertEmptyStr(t, null, "UnmarshalText() empty string")

	var null2 String
	err = null2.UnmarshalText(nil)
	maybePanic(err)
	assertNullStr(t, null2, "UnmarshalText() empty string")

	var iv String
	err = iv.UnmarshalText([]byte{0x44, 0xff, 0x01})
	assert.True(t, errors.NotValid.Match(err), "%+v", err)
}

func TestNullString_MarshalText(t *testing.T) {
	str := MakeString("test")
	data, err := json.Marshal(str)
	maybePanic(err)
	assertJSONEquals(t, data, `"test"`, "non-empty json marshal")
	data, err = str.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "test", "non-empty text marshal")

	// empty values should be encoded as an empty string
	zero := MakeString("")
	data, err = json.Marshal(zero)
	maybePanic(err)
	assertJSONEquals(t, data, `""`, "empty json marshal")
	data, err = zero.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "string marshal text")

	zero.Valid = false
	data, err = zero.MarshalText()
	maybePanic(err)
	assert.Nil(t, data)
}

func TestNullString_BinaryEncoding(t *testing.T) {
	runner := func(b String, want []byte) func(*testing.T) {
		return func(t *testing.T) {
			data, err := b.MarshalBinary()
			assert.NoError(t, err)
			assert.Exactly(t, want, data, t.Name()+": MarshalBinary %q", data)
			data, err = b.Marshal()
			assert.NoError(t, err)
			assert.Exactly(t, want, data, t.Name()+": Marshal %q", data)

			var decoded String
			assert.NoError(t, decoded.UnmarshalBinary(data), "UnmarshalBinary")
			assert.Exactly(t, b, decoded)
		}
	}
	t.Run("HelloWorld", runner(MakeString("HelloWorld"), []byte("\n\rHello\uf8ffWorld\x10\x01")))
	t.Run("null", runner(String{}, []byte("")))
}

func TestNullString_MarshalTo(t *testing.T) {
	str := MakeString("HelloWorld")
	var buf4 [17]byte
	n, err := str.MarshalTo(buf4[:])
	maybePanic(err)
	assert.Exactly(t, 17, n)
	assert.Exactly(t, []byte("\n\rHello\uf8ffWorld\x10\x01"), buf4[:])

	bufFit := make([]byte, str.Size())
	n, err = str.MarshalTo(bufFit)
	maybePanic(err)
	assert.Exactly(t, 17, n)
	assert.Exactly(t, []byte("\n\rHello\uf8ffWorld\x10\x01"), bufFit)
}

func TestStringPointer(t *testing.T) {
	str := MakeString("test")
	ptr := str.Ptr()
	if *ptr != "test" {
		t.Errorf("bad %s string: %#v ≠ %s\n", "pointer", ptr, "test")
	}

	null := String{}
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s string: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestStringIsZero(t *testing.T) {
	str := MakeString("test")
	if str.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	blank := MakeString("")
	if blank.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	empty := MakeString("")
	if empty.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestStringSetValid(t *testing.T) {
	change := MakeString("").SetNull()
	assertNullStr(t, change, "SetValid()")
	assertStr(t, change.SetValid("test"), "SetValid()")
}

func TestNullString_Scan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var nv String
		assert.NoError(t, nv.Scan(nil))
		assert.Exactly(t, String{}, nv)
	})
	t.Run("[]byte", func(t *testing.T) {
		var nv String
		assert.NoError(t, nv.Scan([]byte(`12345678910`)))
		assert.Exactly(t, MakeString(`12345678910`), nv)
	})
	t.Run("string", func(t *testing.T) {
		var nv String
		assert.NoError(t, nv.Scan(`12345678910`))
		assert.Exactly(t, MakeString(`12345678910`), nv)
	})
	t.Run("[]rune unsupported", func(t *testing.T) {
		var nv String
		err := nv.Scan([]rune(`1234567`))
		assert.True(t, errors.MatchKind(err, errors.NotSupported), "Error behaviour should be errors.NotSupported")
		assert.Exactly(t, String{}, nv)
	})
}

func TestString_GoString(t *testing.T) {
	s := MakeString("test")
	assert.Exactly(t, "null.MakeString(`test`)", s.GoString())

	s = MakeString("test").SetNull()
	assert.Exactly(t, "null.String{}", s.GoString())

	s = MakeString("te`st")
	gsWant := []byte("null.MakeString(`te`+\"`\"+`st`)")
	if !bytes.Equal(gsWant, []byte(s.GoString())) {
		t.Errorf("Have: %#v Want: %v", s.GoString(), string(gsWant))
	}
}

func assertStr(t *testing.T, s String, from string) {
	t.Helper()
	if s.Data != "test" {
		t.Errorf("bad %q string: %q ≠ %q\n", from, s.Data, "test")
	}
	if !s.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertEmptyStr(t *testing.T, s String, from string) {
	t.Helper()
	if s.Data != "" {
		t.Errorf("bad %q string: %q ≠ %q\n", from, s.Data, "test")
	}
	if !s.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullStr(t *testing.T, s String, from string) {
	t.Helper()
	if s.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertJSONEquals(t *testing.T, data []byte, cmp string, from string) {
	t.Helper()
	if string(data) != cmp {
		t.Errorf("bad %s data: %s ≠ %s\n", from, data, cmp)
	}
}

func TestNullStringFrom(t *testing.T) {
	assert.Equal(t, "product", MakeString("product").Data)
	assert.True(t, MakeString("product").Valid)
	// assert.False(t, NullStringFromPtr(nil).Valid)
	assert.True(t, MakeString("").Valid)
	v, err := MakeString("product").Value()
	assert.NoError(t, err)
	assert.Equal(t, "product", v)
}
