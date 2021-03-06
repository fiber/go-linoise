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
	"os"
	"path"
	"testing"
	"syscall"
	"github.com/kless/term"
)

func init() {
	term.Input = os.Stderr
	term.InputFD = syscall.Stderr
	Init()
}

var linoiseFile = path.Join(os.TempDir(), "go_linoise")

func TestLinoise(t *testing.T) {
	fmt.Println("Press ^D to exit\n")

	hist, err := NewHistory(linoiseFile)
	if err != nil {
		t.Error(err)
	}
	hist.Load()

	ln := NewLine(hist)
	defer ln.RestoreTerm()

	for {
		if _, err = ln.Read(); err != nil {
			if err == ErrCtrlD {
				hist.Save()
			} else {
				fmt.Fprintf(os.Stderr, err.Error())
			}

			break
		}
	}

	//os.Remove(linoiseFile)
}
