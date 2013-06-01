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

func TestHistSave(t *testing.T) {
	hist, err := NewHistorySize(historyFile, 10)
	if err != nil {
		t.Error("could not create history", err)
	}

	if hist.li.Len() > hist.Cap {
		t.Error("bad capacity size")
	}

	hist.Add("1 line with trailing spaces ")
	hist.Add("2 line without trailing spaces")
	hist.Add("3 line without trailing spaces")
	hist.Add("4 with trailing tabulator\t")
	hist.Add("5 with trailing new line\n")
	hist.Add(" ")              // Not saved to file
	hist.Add(" leading space") // Idem
	hist.Add("")               // Idem
	hist.Add("9 line without trailing spaces")
	hist.Add("10 line number 6")
	hist.Save()

	historyLen = hist.li.Len() - 3 // 3 lines should not be saved
}

func TestHistLoad(t *testing.T) {
	hist, err := NewHistorySize(historyFile, 10)
	if err != nil {
		t.Error("could not load history", err)
	}

	hist.Load()
	e := hist.li.Front()

	for i := 0; i < hist.li.Len(); i++ {
		line := e.Value.(string)

		if strings.HasSuffix(line, "\n") || strings.HasSuffix(line, "\t") ||
			strings.HasSuffix(line, " ") {
			t.Error("line saved with any trailing Unicode space")
		}
	}

	if hist.li.Len() != historyLen {
		t.Error("length doesn't match with values saved")
	}

	os.Remove(historyFile)
}
