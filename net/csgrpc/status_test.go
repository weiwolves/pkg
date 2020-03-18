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

package csgrpc_test

import (
	"testing"

	"github.com/weiwolves/pkg/net/csgrpc"
	"github.com/weiwolves/pkg/util/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNewStatusBadRequestError(t *testing.T) {
	err := csgrpc.NewStatusBadRequestError(codes.Aborted, "error message", "field1", "desc1", "field2", "desc2")
	assert.EqualError(t, err, "rpc error: code = Aborted desc = error message")

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Exactly(t, codes.Aborted, st.Code())

	assert.EqualError(t,
		st.Details()[0].(error),
		"any: message type \"google.rpc.BadRequest\" isn't linked in",
	)
}
