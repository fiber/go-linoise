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
	"bufio"
	"container/ring"
	"log"
	"os"
	"strings"
)


// Values by default
var (
	FilePerm   uint32 = 0600 // History file permission
	HistoryCap = 500         // Capacity
)


// === Type
// ===

type history struct {
	Cap, Len int
	filename string
	file     *os.File
	rng      *ring.Ring
}


// Base to create an history file.
func _baseHistory(fname string, size int) (*history, os.Error) {
	file, err := os.Open(fname, os.O_CREATE|os.O_RDWR, FilePerm)
	if err != nil {
		return nil, err
	}

	h := new(history)
	h.filename = fname
	h.file = file
	h.Cap = size
	h.rng = ring.New(size)

	return h, nil
}

// Creates a new history using the maximum length by default.
func NewHistory(filename string) (*history, os.Error) {
	return _baseHistory(filename, HistoryCap)
}

// Creates a new history whose buffer has the specified size, which must be
// greater than zero.
func NewHistorySize(filename string, size int) (*history, os.Error) {
	if size <= 0 {
		return nil, HistSizeError(size)
	}

	return _baseHistory(filename, size)
}
// ===


// Adds a new line, except when:
// + it starts with some space
// + it's the same line than the previous one
func (h *history) Add(line string) {
	if strings.HasPrefix(line, " ") {
		return
	}

	// Check the last line.
	_line := strings.TrimSpace(line)
	if _line == "" || _line == h.rng.Prev().Value {
		return
	}

	h.rng.Value = _line
	h.rng = h.rng.Next()

	if h.Len < h.Cap {
		h.Len++
	}
}

// Loads the history from the file.
func (h *history) Load() {
	bufin := bufio.NewReader(h.file)

	for {
		line, err := bufin.ReadString('\n')
		if err == os.EOF {
			break
		}

		h.rng.Value = strings.TrimRight(line, "\n")
		h.rng = h.rng.Next()
		h.Len++
	}
}

// Saves to text file.
func (h *history) Save() (err os.Error) {
	if _, err = h.file.Seek(0, 0); err != nil {
		return
	}

	bufout := bufio.NewWriter(h.file)

	for v := range h.rng.Iter() {
		if v != nil {
			if _, err = bufout.WriteString(v.(string) + "\n"); err != nil {
				log.Println("history.Save:", err)
				break
			}
		}
	}

	if err = bufout.Flush(); err != nil {
		log.Println("history.Save:", err)
	}

	h.closeFile()
	return
}

// Closes the file descriptor.
func (h *history) closeFile() {
	h.file.Close()
}

// Opens the file.
/*func (h *history) openFile() {
	file, err := os.Open(fname, os.O_CREATE|os.O_RDWR, FilePerm)
	if err != nil {
		log.Println("history.openFile:", err)
		return
	}

	h.file = file
}*/

