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
	"os"
)


var (
	ErrCtrlD = os.NewError("Interrumpted (Ctrl-d)")

	ErrEmptyHist  = os.NewError("history: empty")
	ErrNilElement = os.NewError("history: no more elements")
)


// Represents a failure on input.
type InputError string

func (i InputError) String() string {
	return "could not read from input: " + string(i)
}


// Represents a failure in output.
type OutputError string

func (o OutputError) String() string {
	return "could not write to output: " + string(o)
}

