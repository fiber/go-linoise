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
	"path"
	"strings"
	"testing"
)


var (
	historyFile = path.Join(os.TempDir(), "go_history")
	historyLen  int
)


func TestSave(t *testing.T) {
	hist, err := NewHistorySize(historyFile, 10)
	if err != nil {
		t.Error("could not create history", err)
	}

	if hist.Cap != hist.rng.Len() {
		t.Error("bad capacity size")
	}

	hist.Add("line with trailing spaces ")
	hist.Add("line without trailing spaces")
	hist.Add("line without trailing spaces")
	hist.Add("with trailing tabulator\t")
	hist.Add("with trailing new line\n")
	hist.Add(" ")
	hist.Add(" leading space")
	hist.Add("")
	hist.Add("line without trailing spaces")
	hist.Add("line number 6")

	hist.Save()

	historyLen = hist.Len
}

func TestLoad(t *testing.T) {
	hist, err := NewHistorySize(historyFile, 10)
	if err != nil {
		t.Error("could not load history", err)
	}

	hist.Load()

	for v := range hist.rng.Iter() {
		if v != nil {
			line := v.(string)
			if strings.HasSuffix(line, "\n") || strings.HasSuffix(line, "\t") ||
				strings.HasSuffix(line, " ") {
				t.Error("line saved with any trailing Unicode space")
			}
		}
	}

	if hist.Len != historyLen {
		t.Error("length doesn't match with values saved")
	}

	os.Remove(historyFile)
}

