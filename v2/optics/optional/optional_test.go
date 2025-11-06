// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package optional

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"

	"github.com/stretchr/testify/assert"
)

type (
	Phone struct {
		number string
	}

	Employment struct {
		phone *Phone
	}

	Info struct {
		employment *Employment
	}

	Response struct {
		info *Info
	}
)

func (response *Response) GetInfo() *Info {
	return response.info
}

func (response *Response) SetInfo(info *Info) *Response {
	response.info = info
	return response
}

var (
	responseOptional = FromPredicateRef[Response](F.IsNonNil[Info])((*Response).GetInfo, (*Response).SetInfo)

	sampleResponse      = Response{info: &Info{}}
	sampleEmptyResponse = Response{}
)

func TestOptional(t *testing.T) {
	assert.Equal(t, O.Of(sampleResponse.info), responseOptional.GetOption(&sampleResponse))
	assert.Equal(t, O.None[*Info](), responseOptional.GetOption(&sampleEmptyResponse))
}
