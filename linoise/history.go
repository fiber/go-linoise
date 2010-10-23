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
	r        *ring.Ring
}


// Base to create an history file.
func _baseHistory(fname string, size int) (*history, os.Error) {
	file, err := os.Open(fname, os.O_CREATE|os.O_RDWR, FilePerm)
	if err != nil {
		return nil, err
	}

	_history := new(history)
	_history.filename = fname
	_history.file = file
	_history.Cap = size
	_history.r = ring.New(size)

	return _history, nil
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
func (self *history) Add(line string) {
	if strings.HasPrefix(line, " ") {
		return
	}

	// Check the last line.
	_line := strings.TrimSpace(line)
	if _line == "" || _line == self.r.Prev().Value {
		return
	}

	self.r.Value = _line
	self.r = self.r.Next()

	if self.Len < self.Cap {
		self.Len++
	}
}

// Loads the history from the file.
func (self *history) Load() {
	bufin := bufio.NewReader(self.file)

	for {
		line, err := bufin.ReadString('\n')
		if err == os.EOF {
			break
		}

		self.r.Value = strings.TrimRight(line, "\n")
		self.r = self.r.Next()
		self.Len++
	}
}

// Saves to text file.
func (self *history) Save() (err os.Error) {
	bufout := bufio.NewWriter(self.file)

	for v := range self.r.Iter() {
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

	self.closeFile()
	return
}

// Closes the file descriptor.
func (self *history) closeFile() {
	self.file.Close()
}

// Opens the file.
/*func (self *history) openFile() {
	file, err := os.Open(fname, os.O_CREATE|os.O_RDWR, FilePerm)
	if err != nil {
		log.Println("history.openFile:", err)
		return
	}

	self.file = file
}*/

