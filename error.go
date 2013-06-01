// Copyright 2010  The "go-linoise" Authors
//
// Use of this source code is governed by the Simplified BSD License
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

package linoise

import (
	"fmt"
)

var (
	ErrCtrlD = fmt.Errorf("Interrumpted (Ctrl-d)")

	ErrEmptyHist  = fmt.Errorf("history: empty")
	ErrNilElement = fmt.Errorf("history: no more elements")
)

// Represents a failure on input.
type inputError string

func (i inputError) Error() string {
	return "could not read from input: " + string(i)
}

// Represents a failure in output.
type outputError string

func (o outputError) Error() string {
	return "could not write to output: " + string(o)
}
