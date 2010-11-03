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
	"testing"

	"github.com/kless/go-term/term"
)


func TestUtility(t *testing.T) {
	fmt.Println("\n\n == Terminal")

	// === Read single key
	term.KeyPress()
	rune, _ := ReadKey("\n + Mode on single character: ")

	fmt.Printf("\n  pressed: %q", string(rune))
	term.RestoreTerm()

	// === Echo
	term.Echo(false)
	fmt.Print("\n + Echo disabled. Write and press Enter: ")
	line, _ := ReadString("")

	fmt.Printf("\n  pressed: %q\n", line)
	fmt.Println(" + Echo enabled")
	term.Echo(true)
}

