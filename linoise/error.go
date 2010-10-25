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
	"strconv"
)


// The error representing an invalid history size.
type HistSizeError int

func (self HistSizeError) String() string {
	return "history: bad size " + strconv.Itoa(int(self))
}

// ===

type error struct {
	os.ErrorString
}

var (
	ErrCtrlC      = &error{"Interrumpted: End of Text (Ctrl-c)"}
	ErrCtrlD      = &error{"Interrumpted: End of Transmission (Ctrl-d)"}
	ErrEmptyHist  = &error{"history: empty"}
	ErrNilElement = &error{"history: no more elements"}
)

