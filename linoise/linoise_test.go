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
	"testing"

	"github.com/kless/go-term/term"
)


func Test(t *testing.T) {
	term.MakeRaw(Input.(*os.File).Fd())
	defer term.RestoreTermios()

	hist, err := NewHistory("/tmp/go-history")
	if err != nil {
		t.Error(err)
	}
	hist.Load()

	ln := NewLine(hist, "matrix> ")
	if err = ln.Run(); err != nil {
		fmt.Println(err)
	} else {
		hist.Save()
	}
}

