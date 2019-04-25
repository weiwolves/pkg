// package pseudo is the fake data generator, heavily inspired by
// forgery and ffaker Ruby gems.
//
// CoreStore: This package has been refactored to avoid package global variables
// and global settings which are an anti-pattern. In a multi store/language
// environment a package global language limits everything. Also the global PRNG
// got eliminate and reduces a mutex bottle neck. This package can also handle
// max_len values for generated data and supports embedded structs which
// implements sql.Scanner interface.
//
// Most data and methods are ported from forgery/ffaker Ruby gems.
//
// Currently english and russian languages are available.
//
// For the list of available methods please look at
// https://godoc.org/github.com/icrowley/fake.
//
// Fake embeds samples data files unless you call UseExternalData(true) in order
// to be able to work without external files dependencies when compiled, so, if
// you add new data files or make changes to existing ones don't forget to
// regenerate data.go file using github.com/mjibson/esc tool and esc -o data.go
// -pkg fake data command (or you can just use go generate command if you are
// using Go 1.4 or later).
package pseudo

import (
	"encoding"
	"io"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/util/conv"
	"github.com/oklog/ulid"
	"golang.org/x/exp/rand"
)

// Reason for using this package: at the moment well maintained and does not include net/http.
// Run: $ go get -u github.com/shuLhan/go-bindata
//go:generate go-bindata -o bindata.go -pkg pseudo data/...

// Faker allows a type to implement a custom fake data generation. The argument
// fieldName contains the name of the current field for which random/fake data
// should be generated. The return argument hasFakeDataApplied can be set to
// true, if fake data gets generated for the current field. Setting
// hasFakeDataApplied to false, the fake data should be generated by this
// package.
type Faker interface {
	Fake(fieldName string) (hasFakeDataApplied bool, err error)
}

// Supported tags
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	tagName       = "faker"
	tagMaxLenName = "max_len"
	// Skip indicates a struct tag, that the field should be skipped.
	Skip                     = "-"
	DefaultMaxRecursionLevel = 10
)

// FakeFunc generates new specific fake values. The returned interface{} type
// can only contain primitive types and time.Time.
type FakeFunc func(maxLen int) (interface{}, error)

// Options applied for the Service type.
type Options struct {
	Lang            string
	UseExternalData bool
	EnFallback      bool
	TimeLocation    *time.Location // defaults to UTC
	// RespectValidField if enabled allows to observe the `Valid bool` field of
	// a struct. Like sql.NullString, sqlNullInt64 or all null.* types. If Valid
	// is false, the other fields will be reset to their repective default
	// values. Reason: All fields of a struct are getting applied with fake
	// data, if Valid is false and the field gets written to the e.g. DB and
	// later compared.
	RespectValidField bool
	// DisabledFieldNameUse if enabled avoids the usage of the field name
	// instead of the struct tag to find out which kind of random function is
	// needed.
	DisabledFieldNameUse bool
	// MaxLenStringLimit defines the upper bound of the maximal length for a
	// string. If not set, defaults to 512.
	MaxLenStringLimit uint64
	// FloatMaxDecimals limits the float generation to this amount of decimals.
	// Useful for MySQL/MariaDB float column type.
	FloatMaxDecimals int
	// MaxRecursionLevel default see DefaultMaxRecursionLevel
	MaxRecursionLevel int
}

const maxLenStringLimit = 512

type optionFn struct {
	sortOrder int
	fn        func(*Service) error
}

// WithTagFakeFunc extends faker with a new tag (or field name) to generate fake
// data with a specified custom algorithm. It will overwrite a previous set
// function.
func WithTagFakeFunc(tag string, provider FakeFunc) optionFn {
	return optionFn{
		sortOrder: 10,
		fn: func(s *Service) error {
			s.mu.Lock()
			defer s.mu.Unlock()
			s.funcs[tag] = provider
			return nil
		},
	}
}

// WithTagFakeFuncAlias sets a new alias. E.g. when WithTagFakeFunc("http",...)
// adds a function to generate HTTP links, then calling
// WithTagFakeFuncAlias("https","http","ftps","ftp") would say that https is an
// alias of the http function and ftps is an alias of ftp.
func WithTagFakeFuncAlias(aliasTag ...string) optionFn {
	return optionFn{
		sortOrder: 100,
		fn: func(s *Service) error {
			s.mu.Lock()
			defer s.mu.Unlock()
			for i := 0; i < len(aliasTag); i = i + 2 {
				// alias name => original name
				s.funcsAliases[aliasTag[i]] = aliasTag[i+1]
			}
			return nil
		},
	}
}

// Service provides a service to generate fake data.
type Service struct {
	r           *rand.Rand
	o           Options
	id          *uint64
	ulidEntropy io.Reader

	mu           sync.RWMutex
	langMapping  map[string]map[string][]string // cat/subcat/lang/samples
	funcs        map[string]FakeFunc
	funcsAliases map[string]string // alias name => original name
}

// MustNewService creates a new Service but panics on error.
func MustNewService(seed uint64, o *Options, opts ...optionFn) *Service {
	s, err := NewService(seed, o, opts...)
	if err != nil {
		panic(err)
	}
	return s
}

// NewService creates a new fake service.
func NewService(seed uint64, o *Options, opts ...optionFn) (*Service, error) {
	if seed == 0 {
		seed = uint64(time.Now().UnixNano())
	}
	if o == nil {
		o = &Options{
			EnFallback: true,
		}
	}
	if o.Lang == "" {
		o.Lang = "en"
	}
	if o.TimeLocation == nil {
		o.TimeLocation = time.UTC
	}
	if o.MaxLenStringLimit == 0 {
		o.MaxLenStringLimit = maxLenStringLimit
	}
	if o.MaxRecursionLevel == 0 {
		o.MaxRecursionLevel = DefaultMaxRecursionLevel
	}

	s := &Service{
		langMapping: make(map[string]map[string][]string),
		r:           rand.New(&lockedSource{src: rand.NewSource(seed)}),
		o:           *o,
		id:          new(uint64),
		ulidEntropy: ulid.Monotonic(rand.New(rand.NewSource(seed)), 0),
	}

	s.funcs = map[string]FakeFunc{
		"id": func(maxLen int) (interface{}, error) {
			return s.ID(), nil
		},
		"uuid": func(maxLen int) (interface{}, error) {
			return s.UUID(), nil
		},
		"uuid_string": func(maxLen int) (interface{}, error) {
			return s.UUIDString(), nil
		},
		"ulid": func(maxLen int) (interface{}, error) {
			return s.ULID().String(), nil
		},
		"mac_address": func(maxLen int) (interface{}, error) {
			return s.MacAddress(), nil
		},
		"domain_name": func(maxLen int) (interface{}, error) {
			return s.DomainName(), nil
		},
		"username": func(maxLen int) (interface{}, error) {
			return s.UserName(), nil
		},
		"url": func(maxLen int) (interface{}, error) {
			return s.URL(), nil
		},
		"ipv4": func(maxLen int) (interface{}, error) {
			return s.IPv4(), nil
		},
		"ipv6": func(maxLen int) (interface{}, error) {
			return s.IPv4(), nil
		},
		"password": func(maxLen int) (interface{}, error) {
			return s.Password(6, int(maxLen), true, true, true), nil
		},
		"email": func(maxLen int) (interface{}, error) {
			return s.EmailAddress(), nil
		},
		"lat": func(maxLen int) (interface{}, error) {
			return s.Latitude(), nil
		},
		"long": func(maxLen int) (interface{}, error) {
			return s.Longitude(), nil
		},
		"cc_number": func(maxLen int) (interface{}, error) {
			return s.CreditCardNum(""), nil
		},
		"cc_type": func(maxLen int) (interface{}, error) {
			return s.CreditCardType(), nil
		},
		"phone_number": func(maxLen int) (interface{}, error) {
			return s.Phone(), nil
		},
		"male_first_name": func(maxLen int) (interface{}, error) {
			return s.MaleFirstName(), nil
		},
		"female_first_name": func(maxLen int) (interface{}, error) {
			return s.FemaleFirstName(), nil
		},
		"name": func(maxLen int) (interface{}, error) {
			return s.FullName(), nil
		},
		"last_name": func(maxLen int) (interface{}, error) {
			return s.LastName(), nil
		},
		"first_name": func(maxLen int) (interface{}, error) {
			return s.FirstName(), nil
		},
		"prefix": func(maxLen int) (interface{}, error) {
			return s.Prefix(), nil
		},
		"suffix": func(maxLen int) (interface{}, error) {
			return s.Suffix(), nil
		},
		"date": func(maxLen int) (interface{}, error) {
			return s.Date(), nil
		},
		"clock": func(maxLen int) (interface{}, error) {
			return s.Clock(), nil
		},
		"time": func(maxLen int) (interface{}, error) {
			return s.Time(), nil
		},
		"timestamp": func(maxLen int) (interface{}, error) {
			return s.TimeStamp(), nil
		},
		"dob": func(maxLen int) (interface{}, error) {
			return s.Dob18(), nil
		},
		"timezone": func(maxLen int) (interface{}, error) {
			return s.TimeZone(), nil
		},
		"unix_time": func(maxLen int) (interface{}, error) {
			return s.RandomUnixTime(), nil
		},
		"month_name": func(maxLen int) (interface{}, error) {
			return s.Month(), nil
		},
		"month": func(maxLen int) (interface{}, error) {
			return s.MonthNum(), nil
		},
		"year": func(maxLen int) (interface{}, error) {
			return s.Year(1990, 2025), nil
		},
		"week_day": func(maxLen int) (interface{}, error) {
			return s.WeekDay(), nil
		},

		"sentence": func(maxLen int) (interface{}, error) {
			return s.Sentence(maxLen), nil
		},
		"paragraph": func(maxLen int) (interface{}, error) {
			return s.Paragraph(maxLen), nil
		},
		"currency": func(maxLen int) (interface{}, error) {
			return s.Currency(), nil
		},
		"currency_code": func(maxLen int) (interface{}, error) {
			return s.CurrencyCode(), nil
		},
		"price": func(maxLen int) (interface{}, error) {
			return s.Price(), nil
		},
		"price_currency": func(maxLen int) (interface{}, error) {
			return s.PriceWithCurrency(), nil
		},
		"word": func(maxLen int) (interface{}, error) {
			return s.Word(maxLen), nil
		},
		"city": func(maxLen int) (interface{}, error) {
			return s.City(), nil
		},
		"postcode": func(maxLen int) (interface{}, error) {
			return s.Zip(), nil
		},
		"street": func(maxLen int) (interface{}, error) {
			return s.StreetAddress(), nil
		},
		"company": func(maxLen int) (interface{}, error) {
			return s.CompanyLegal(), nil
		},
		"country": func(maxLen int) (interface{}, error) {
			return s.Country(), nil
		},
		"country_id": func(maxLen int) (interface{}, error) {
			return s.CountryISO2(), nil
		},
		"region": func(maxLen int) (interface{}, error) {
			return s.State(), nil
		},
	}
	s.funcsAliases = map[string]string{
		"firstname":     "first_name",
		"middlename":    "first_name",
		"lastname":      "last_name",
		"password_hash": "password",
		"zip":           "postcode",
		"address":       "street",
		"increment_id":  "ulid",
		"telephone":     "phone_number",
		"fax":           "phone_number",
	}

	sort.Slice(opts, func(i, j int) bool {
		return opts[i].sortOrder < opts[j].sortOrder // ascending 0-9 sorting ;-)
	})

	for _, o := range opts {
		if err := o.fn(s); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	// validate that the alias target exists
	for alias, target := range s.funcsAliases {
		if _, ok := s.funcs[alias]; ok {
			return nil, errors.AlreadyExists.Newf("[pseudo] Alias %q already exists as a fakeFunc", alias)
		}
		if _, ok := s.funcs[target]; !ok {
			return nil, errors.NotImplemented.Newf("[pseudo] Alias %q has an undefined target %q", alias, target)
		}
	}

	return s, nil
}

// GetLangs returns a slice of available languages
func (s *Service) GetLangs() []string {
	lng, _ := AssetDir("data")
	return lng
}

// SetLang sets the language in which the data should be generated
// returns error if passed language is not available
func (s *Service) SetLang(newLang string) error {
	found := false
	for _, l := range s.GetLangs() {
		if newLang == l {
			found = true
			s.o.Lang = newLang
			break
		}
	}
	if !found {
		return errors.NotFound.Newf("[pseudo] The language passed (%s) is not available", newLang)
	}
	return nil
}

func join(parts ...string) string {
	var filtered []string
	for _, part := range parts {
		if part != "" {
			filtered = append(filtered, part)
		}
	}
	return strings.Join(filtered, " ")
}

func (s *Service) generate(lang, cat string, fallback bool) string {
	format := s.lookup(lang, cat+"_format", fallback)
	var result string
	for _, ru := range format {
		if ru != '#' {
			result += string(ru)
		} else {
			result += strconv.Itoa(s.r.Intn(10))
		}
	}
	return result
}

func (s *Service) lookup(lang, cat string, fallback bool) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s._lookup(lang, cat, fallback)
}

func (s *Service) _lookup(lang, cat string, fallback bool) string {
	if samples, ok := s.langMapping[lang][cat]; ok {
		return samples[s.r.Intn(len(samples))]
	}

	samples, err := s.populateSamples(lang, cat)
	if err != nil {
		if pe, ok := err.(*os.PathError); lang != "en" && fallback && s.o.EnFallback && ok && pe.Err == os.ErrNotExist {
			return s._lookup("en", cat, false)
		}
		return ""
	}

	return samples[s.r.Intn(len(samples))]
}

func (s *Service) populateSamples(lang, cat string) ([]string, error) {
	data, err := s.readFile(lang, cat)
	if err != nil {
		return nil, err
	}

	if _, ok := s.langMapping[lang]; !ok {
		s.langMapping[lang] = make(map[string][]string)
	}

	samples := strings.Split(strings.TrimSpace(string(data)), "\n")

	s.langMapping[lang][cat] = samples
	return samples, nil
}

func (s *Service) readFile(lang, cat string) (_ []byte, err error) {
	fullPath := filepath.Join("data", lang, cat+".txt")
	data, err := Asset(fullPath)
	return data, err
}

// FakeData is the main function. Will generate a fake data based on your
// struct.  You can use this for automation testing, or anything that need
// automated data. You don't need to Create your own data for your testing.
// Unsupported types are getting ignored.
func (s *Service) FakeData(ptr interface{}) error {

	reflectType := reflect.TypeOf(ptr)

	if ptr == nil || reflectType.Kind() != reflect.Ptr || reflect.ValueOf(ptr).IsNil() {
		return errors.NotSupported.Newf("[pseudo] Nil/Non-pointer values are not supported. Argument ptr should be a pointer.")
	}

	finalValue, err := s.getValue(reflectType.Elem(), 0, 0)
	if err != nil {
		if errors.Exceeded.Match(err) {
			err = nil
		}
		return err
	}

	rVal := reflect.ValueOf(ptr)
	rVal.Elem().Set(finalValue.Convert(reflectType.Elem()))
	return nil
}

type scanner interface {
	Scan(interface{}) error
}

type errRecursionExceeded string

func (e errRecursionExceeded) ErrorKind() errors.Kind {
	return errors.Exceeded
}

func (e errRecursionExceeded) Error() string {
	return string(e)
}

func (s *Service) getValue(t reflect.Type, maxLen uint64, recursionLevel int) (rVal reflect.Value, err error) {

	if recursionLevel > s.o.MaxRecursionLevel {
		return rVal, errRecursionExceeded("Reached max recursion level")
	}

	k := t.Kind()

	if maxLen == 0 {
		maxLen = math.MaxInt8
	}

	switch k {
	case reflect.Ptr:
		v := reflect.New(t.Elem())
		val, err := s.getValue(t.Elem(), maxLen, recursionLevel+1)
		if err != nil {
			return rVal, err
		}
		v.Elem().Set(val.Convert(t.Elem()))
		return v, nil
	case reflect.Struct:

		switch ts := t.String(); ts {
		case "time.Time":
			var neg int64 = 1
			if s.r.Int()%2 == 0 {
				neg = -1
			}
			ft := time.Now().Add(time.Second * time.Duration(neg*s.r.Int63n(3600*24*90)))
			// proper way to get rid of the nano seconds
			ft = time.Unix(ft.Unix(), 0).In(s.o.TimeLocation)
			return reflect.ValueOf(ft), nil
		default:
			v := reflect.New(t).Elem()

			var fkr Faker
			if v.CanInterface() {
				// The _ is important to not cause a panic: "comma ok" pattern with ignored ok
				fkr, _ = v.Addr().Interface().(Faker)
			}
			shouldResetField := false
			typeName := v.Type().String()
			for i := 0; i < v.NumField(); i++ {
				vf := v.Field(i) // value field
				tf := t.Field(i) // type field

				if !vf.CanSet() {
					continue // to avoid panic to set on unexported field in struct
				}

				tag := tf.Tag.Get(tagName)
				if tag == Skip {
					continue
				}

				if maxLenTag := tf.Tag.Get(tagMaxLenName); maxLenTag != "" {
					maxLen, err = strconv.ParseUint(maxLenTag, 10, 64)
					if err != nil {
						return rVal, err
					}
				}

				// Custom functions must have precedence before interface Faker
				// implementation of a type.
				if vf.CanInterface() && vf.CanAddr() {
					sTag := tag
					if sTag == "" {
						sTag = toSnakeCase(tf.Name)
					}

					iFaceVal, ok, err := s.getFuncsValue(typeName, tf.Name, sTag, maxLen)
					if err != nil {
						return rVal, errors.WithStack(err)
					}

					if ok {
						switch tFace := vf.Addr().Interface().(type) {
						case *time.Time:
							iFaceValTime, err := conv.ToTimeE(iFaceVal)
							if err != nil {
								return rVal, errors.WithStack(err)
							}
							vf.Set(reflect.ValueOf(iFaceValTime))
							continue

						case scanner:
							if err := tFace.Scan(iFaceVal); err != nil {
								return rVal, errors.WithStack(err)
							}
							continue

						case encoding.TextUnmarshaler:
							bt, err := conv.ToByteE(iFaceVal)
							if err != nil {
								return rVal, errors.WithStack(err)
							}
							if err := tFace.UnmarshalText(bt); err != nil {
								return rVal, errors.WithStack(err)
							}
							continue
						default:
							val := reflect.ValueOf(iFaceVal)
							if isConvertibleType(vf.Kind()) && isConvertibleType(val.Kind()) {
								val = val.Convert(vf.Type())
								vf.Set(val)
							}
							continue
						}
					}
				}

				if fkr != nil {
					// fieldName is the original field name as written in the struct.
					if skip, err := fkr.Fake(tf.Name); err != nil {
						return rVal, err
					} else if skip {
						continue
					}
				}

				switch {
				case tag == "":
					// The check of isConvertibleType protects from a panic of
					// Convert. Especially useful when an exported field has a
					// func or interface or channel type.
					val, err := s.getValue(vf.Type(), maxLen, recursionLevel+1)
					if err != nil {
						return reflect.Value{}, err
					}
					if isConvertibleType(vf.Kind()) {
						val = val.Convert(vf.Type())
					}
					if val.IsValid() {
						vf.Set(val)
					}
				default:
					err := s.setDataWithTag(vf.Addr(), tag, maxLen, false)
					if err != nil {
						return reflect.Value{}, err
					}
				}

				if !shouldResetField && t.Field(i).Name == "Valid" && vf.Kind() == reflect.Bool && !vf.Bool() {
					shouldResetField = true
				}
			}

			if shouldResetField {
				v.Set(reflect.New(t).Elem())
			}

			return v, nil
		}

	case reflect.String:
		ml := maxLen
		if ml > s.o.MaxLenStringLimit {
			ml = s.o.MaxLenStringLimit
		}
		res := s.randomString(ml)
		return reflect.ValueOf(res), nil
	case reflect.Array, reflect.Slice:
		ml := maxLen
		if ml > s.o.MaxLenStringLimit {
			ml = s.o.MaxLenStringLimit
		}
		slLen := s.r.Uint64n(ml)
		v := reflect.MakeSlice(t, int(slLen), int(slLen))
		for i := 0; i < v.Len(); i++ {
			val, err := s.getValue(t.Elem(), ml, recursionLevel+1)
			if err != nil {
				return rVal, err
			}
			v.Index(i).Set(val)
		}
		return v, nil
	case reflect.Int:
		return reflect.ValueOf(int(s.r.Uint64n(maxLen))), nil
	case reflect.Int8:
		return reflect.ValueOf(int8(s.r.Uint64n(maxLen))), nil
	case reflect.Int16:
		return reflect.ValueOf(int16(s.r.Uint64n(maxLen))), nil
	case reflect.Int32:
		return reflect.ValueOf(int32(s.r.Uint64n(maxLen))), nil
	case reflect.Int64:
		return reflect.ValueOf(int64(s.r.Uint64n(maxLen))), nil

	case reflect.Float32:
		return reflect.ValueOf(s.r.Float32()), nil
	case reflect.Float64:
		f := s.r.Float64()
		if s.o.FloatMaxDecimals > 0 {
			p10 := math.Pow10(s.o.FloatMaxDecimals)
			f = math.Round(f*p10) / p10
		}
		return reflect.ValueOf(f), nil

	case reflect.Bool:
		val := s.r.Uint64n(3) > 0 // create more true values
		return reflect.ValueOf(val), nil

	case reflect.Uint:
		return reflect.ValueOf(uint(s.r.Uint64n(maxLen))), nil

	case reflect.Uint8:
		return reflect.ValueOf(uint8(s.r.Uint64n(maxLen))), nil

	case reflect.Uint16:
		return reflect.ValueOf(uint16(s.r.Uint64n(maxLen))), nil

	case reflect.Uint32:
		return reflect.ValueOf(uint32(s.r.Uint64n(maxLen))), nil

	case reflect.Uint64:
		return reflect.ValueOf(uint64(s.r.Uint64n(maxLen))), nil

	case reflect.Map:
		v := reflect.MakeMap(t)
		randLen := s.r.Uint64n(maxLen)
		var i uint64
		for ; i < randLen; i++ {
			key, err := s.getValue(t.Key(), maxLen, recursionLevel+1)
			if err != nil {
				return rVal, err
			}
			val, err := s.getValue(t.Elem(), maxLen, recursionLevel+1)
			if err != nil {
				return rVal, err
			}
			v.SetMapIndex(key, val)
		}
		return v, nil
	case reflect.Func, reflect.Chan, reflect.Interface, reflect.Complex64, reflect.Complex128:
		// ignore
		return rVal, nil
	}
	return rVal, errors.NotSupported.Newf("[pseudo] Type %+v not supported", t)
}

func isConvertibleType(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return true

	case reflect.Float32, reflect.Float64:
		return true

	case reflect.Complex64, reflect.Complex128:
		return true

	case reflect.String:
		return true

	case reflect.Slice:
		return true
	}
	return false
}

func (s *Service) getFuncsValue(typeName, fieldName, tag string, maxLen uint64) (interface{}, bool, error) {
	// typeName = testdata.CustomerAddressEntity
	// fieldName = ParentID
	// tag = parent_id

	// fmt.Printf("%s.%s Tag:%s\n", typeName, fieldName, tag)

	// testdata.CustomerAddressEntity.CountryID
	if fn, ok := s.funcs[typeName+"."+fieldName]; ok {
		iFaceVal, err := fn(int(maxLen))
		return iFaceVal, true, err
	}

	// lookup an alias for parent_id
	if fnAlias, ok := s.funcsAliases[tag]; ok && fnAlias != "" {
		tag = fnAlias
	}
	fn, ok := s.funcs[tag]
	if !ok {
		return nil, false, nil
	}
	iFaceVal, err := fn(int(maxLen))
	return iFaceVal, true, err
}

func (s *Service) setDataWithTag(v reflect.Value, tag string, maxLen uint64, tagIsFieldName bool) error {
	if v.Kind() != reflect.Ptr {
		return errors.NotSupported.Newf("[pseudo] Non-pointer values are not supported. Argument ptr should be a pointer.")
	}

	// TODO check if the map access causes a race condition.
	if fnAlias, ok := s.funcsAliases[tag]; ok && fnAlias != "" {
		tag = fnAlias
	}
	fn, ok := s.funcs[tag]
	if !ok && tagIsFieldName {
		return nil
	}
	if !ok && !tagIsFieldName {
		return errors.NotFound.Newf("[pseudo] Tag %q not found in map", tag)
	}
	iFaceVal, err := fn(int(maxLen))
	if err != nil {
		return errors.WithStack(err)
	}

	v = reflect.Indirect(v)
	switch k := v.Kind(); k {
	case reflect.Float32, reflect.Float64:
		val, err := conv.ToFloat64E(iFaceVal)
		if err != nil {
			return errors.WithStack(err)
		}
		v.SetFloat(val)
	case reflect.Bool:
		val, err := conv.ToBoolE(iFaceVal)
		if err != nil {
			return errors.WithStack(err)
		}
		v.SetBool(val)
	case reflect.String:
		val, err := conv.ToStringE(iFaceVal)
		if err != nil {
			return errors.WithStack(err)
		}
		v.SetString(val)
	case reflect.Slice: // TODO must be improved to detect []byte
		val, err := conv.ToByteE(iFaceVal)
		if err != nil {
			return errors.WithStack(err)
		}
		v.SetBytes(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := conv.ToInt64E(iFaceVal)
		if err != nil {
			return errors.Wrapf(err, "[pseudo] For Tag %q", tag)
		}
		v.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := conv.ToUintE(iFaceVal)
		if err != nil {
			return errors.WithStack(err)
		}
		v.SetUint(uint64(val))
	default:
		return errors.NotSupported.Newf("[pseudo] Kind %q not supported in setDataWithTag", k)
	}

	return nil
}

func (s *Service) randomString(n uint64) string {
	b := make([]byte, n)
	for i, cache, remain := int64(n-1), s.r.Uint64(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = s.r.Uint64(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func (s *Service) randomElementFromSliceString(sl []string) string {
	return sl[s.r.Int()%len(sl)]
}

type caseCache struct {
	mu    sync.RWMutex
	cache map[string]string
}

var (
	matchFirstCap    = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap      = regexp.MustCompile("([a-z0-9])([A-Z])")
	caseCacheService = &caseCache{cache: map[string]string{
		"CountryID":    "country_id",
		"EntityID":     "entity_id",
		"PasswordHash": "password_hash",
	}}
)

func toSnakeCase(str string) string {
	caseCacheService.mu.RLock()
	caStr := caseCacheService.cache[str]
	caseCacheService.mu.RUnlock()
	if caStr != "" {
		return caStr
	}
	caseCacheService.mu.Lock()
	defer caseCacheService.mu.Unlock()

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	lc := strings.ToLower(snake)
	caseCacheService.cache[str] = lc
	return lc
}
