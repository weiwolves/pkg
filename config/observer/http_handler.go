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

// +build csall http

package observer

import (
	"io"
	"net/http"

	"github.com/weiwolves/pkg/config"
	"github.com/weiwolves/pkg/net/mw"
)

// HTTPHandlerOptions can set different behaviour to the
// RegisterWithJSONviaHTTP / DeregisterWithJSONviaHTTP. endpoints.
type HTTPHandlerOptions struct {
	// ErrorHandler custom error handler. Default error handler returns a status
	// code and prints the whole stack trace. May leak sensitive information.
	ErrorHandler   mw.ErrorHandler
	MaxRequestSize int64 // Default 10kb
	// StatusCodeOk sets the HTTP status for successful operation, default
	// http.StatusCreated for register and StatusAccepted for deregister.
	StatusCodeOk int
	// StatusCodeError sets the HTTP status for erroneous operation, default
	// http.StatusInternalServerError
	StatusCodeError int
}

// RegisterWithJSONviaHTTP provides an endpoint handler to register validators
// with the concrete type of config.Service
func RegisterWithJSONviaHTTP(or config.ObserverRegisterer, ho HTTPHandlerOptions) http.Handler {
	if ho.StatusCodeOk == 0 {
		ho.StatusCodeOk = http.StatusCreated
	}
	if ho.StatusCodeError == 0 {
		ho.StatusCodeError = http.StatusInternalServerError
	}
	if ho.MaxRequestSize == 0 {
		ho.MaxRequestSize = 1024 * 20 // 20kb
	}
	if ho.ErrorHandler == nil {
		ho.ErrorHandler = mw.ErrorWithStatusCode(ho.StatusCodeError)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if err := RegisterWithJSON(or, io.LimitReader(r.Body, ho.MaxRequestSize)); err != nil {
			ho.ErrorHandler(err).ServeHTTP(w, r)
			return
		}
		// io.Copy(ioutil.Discard, r.Body) // I don't know if that is necessary
		w.WriteHeader(ho.StatusCodeOk)
	})
}

// DeregisterWithJSONviaHTTP provides an endpoint handler to deregister
// validators with the concrete type of config.Service
func DeregisterWithJSONviaHTTP(or config.ObserverRegisterer, ho HTTPHandlerOptions) http.Handler {
	if ho.StatusCodeOk == 0 {
		ho.StatusCodeOk = http.StatusAccepted
	}
	if ho.StatusCodeError == 0 {
		ho.StatusCodeError = http.StatusInternalServerError
	}
	if ho.MaxRequestSize == 0 {
		ho.MaxRequestSize = 1024 * 10 // 10kb
	}
	if ho.ErrorHandler == nil {
		ho.ErrorHandler = mw.ErrorWithStatusCode(ho.StatusCodeError)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if err := DeregisterWithJSON(or, io.LimitReader(r.Body, ho.MaxRequestSize)); err != nil {
			ho.ErrorHandler(err).ServeHTTP(w, r)
			return
		}
		// io.Copy(ioutil.Discard, r.Body) // I don't know if that is necessary
		w.WriteHeader(ho.StatusCodeOk)
	})
}
