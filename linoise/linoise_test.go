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

	"github.com/kless/go-term/term"
)


var linoiseFile = path.Join(os.TempDir(), "go_linoise")


func TestLinoise(t *testing.T) {
	term.MakeRaw()
	defer term.RestoreTermios()

	hist, err := NewHistory(linoiseFile)
	if err != nil {
		t.Error(err)
	}
	hist.Load()

	ln := NewLine(hist)

	for {
		_, err = ln.Read()
		if err != nil {
			break
		}
	}

	if err == ErrCtrlD {
		hist.Save()
	} else {
		fmt.Println(err)
	}

	//os.Remove(linoiseFile)
}

