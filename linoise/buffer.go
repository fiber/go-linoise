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
	"utf8"

	"github.com/kless/go-term/term"
)


// Buffer size
var (
	BufferCap = 4096
	BufferLen = 64 // Initial length
)


// === Init
// ===

var lines, columns int

func init() {
	lines, columns = term.GetWinsizeInChar()
}


// === Type
// ===

// Represents the line buffer.
type buffer struct {
	promptLen int
	pos       int   // Pointer position into buffer
	size      int   // Amount of characters added
	data      []int // Text buffer
}

func newBuffer(promptLen int) *buffer {
	b := new(buffer)
	b.promptLen = promptLen
	b.data = make([]int, BufferLen, BufferCap)

	return b
}
// ===


// Inserts a character in the cursor position.
func (b *buffer) insertRune(rune int) (useRefresh bool, err os.Error) {
	b.grow(b.size + 1) // Check if there is free space for one more character

	// Avoid a full update of the line.
	if b.pos == b.size {
		char := make([]byte, utf8.UTFMax)
		utf8.EncodeRune(rune, char)

		if _, err = output.Write(char); err != nil {
			return false, OutputError(err.String())
		}
	} else {
		useRefresh = true
		copy(b.data[b.pos+1:b.size+1], b.data[b.pos:b.size])
	}

	b.data[b.pos] = rune
	b.pos++
	b.size++
	return
}

// Inserts several characters.
func (b *buffer) insertRunes(runes []int) (useRefresh bool, err os.Error) {
	for _, r := range runes {
		useRefresh, err = b.insertRune(r)
		if err != nil {
			return
		}
	}
	return
}

// Moves the cursor at the start.
func (b *buffer) start() (err os.Error) {
	if b.pos == b.promptLen {
		return
	}

	for ln, _ := b.pos2xy(b.pos); ln > 0; ln-- {
		if _, err = output.Write(cursorUp); err != nil {
			return OutputError(err.String())
		}
	}

	if _, err = fmt.Fprintf(output, "\r\033[%dC", b.promptLen); err != nil {
		return OutputError(err.String())
	}

	b.pos = b.promptLen
	return
}

// Moves the cursor at the end. Returns the number of lines that fill in the data.
func (b *buffer) end() (lines int, err os.Error) {
	if b.pos == b.size {
		return
	}

	lastLine, lastColumn := b.pos2xy(b.size)

	for ln, _ := b.pos2xy(b.pos); ln < lastLine; ln++ {
		if _, err = output.Write(cursorDown); err != nil {
			return 0, OutputError(err.String())
		}
	}

	if _, err = fmt.Fprintf(output, "\r\033[%dC", lastColumn); err != nil {
		return 0, OutputError(err.String())
	}

	b.pos = b.size
	return lastLine, nil
}

// Moves the cursor one character backward.
func (b *buffer) backward() (err os.Error) {
	if b.pos == b.promptLen {
		return
	}

	b.pos--

	// If position is on the same line.
	if _, col := b.pos2xy(b.pos); col != 0 {
		if _, err = output.Write(cursorBackward); err != nil {
			return OutputError(err.String())
		}
	} else {
		if _, err = output.Write(cursorUp); err != nil {
			return OutputError(err.String())
		}
		if _, err = fmt.Fprintf(output, "\033[%dC", columns); err != nil {
			return OutputError(err.String())
		}
	}

	return
}

// Moves the cursor one character forward.
func (b *buffer) forward() (err os.Error) {
	if b.pos == b.size {
		return
	}

	b.pos++

	if _, col := b.pos2xy(b.pos); col != 0 {
		if _, err = output.Write(cursorForward); err != nil {
			return OutputError(err.String())
		}
	} else {
		if _, err = output.Write(toNextLine); err != nil {
			return OutputError(err.String())
		}
	}

	return
}

// Deletes the character in cursor.
func (b *buffer) delete() (useRefresh bool, err os.Error) {
	if b.pos == b.size {
		return
	}

	copy(b.data[b.pos:], b.data[b.pos+1:b.size])
	b.size--
	//	ln, _ := b.pos2xy(b.pos)

	if _, err = output.Write(delChar); err != nil {
		return false, OutputError(err.String())
	}

	return
}

// Deletes the previous character from cursor.
func (b *buffer) deletePrev() (useRefresh bool, err os.Error) {
	if b.pos == b.promptLen {
		return
	}

	copy(b.data[b.pos-1:], b.data[b.pos:b.size])
	b.pos--
	ln, col := b.pos2xy(b.pos)

	if col != 0 {
		if _, err = output.Write(delBackspace); err != nil {
			return false, OutputError(err.String())
		}
	} else {
		if _, err = output.Write(cursorUp); err != nil {
			return false, OutputError(err.String())
		}
		if _, err = fmt.Fprintf(output, "\033[%dC", columns); err != nil {
			return false, OutputError(err.String())
		}
	}

	b.size--
	lastLine, _ := b.pos2xy(b.size)
	if ln != lastLine {
		return true, nil
	}

	return
}

// Deletes from current position until to end of line.
func (b *buffer) deleteRight() (err os.Error) {
	if b.pos == b.size {
		return
	}

	lastLine, _ := b.pos2xy(b.size)
	posLine, _ := b.pos2xy(b.pos)

	// To the last line.
	for ln := posLine; ln < lastLine; ln++ {
		if _, err = output.Write(cursorDown); err != nil {
			return OutputError(err.String())
		}
	}

	// Delete all lines until the cursor position.
	for ln := lastLine; ln > posLine; ln-- {
		if _, err = output.Write(delLine_cursorUp); err != nil {
			return OutputError(err.String())
		}
	}

	if _, err = output.Write(delRight); err != nil {
		return OutputError(err.String())
	}

	b.size = b.pos
	return nil
}

// Swaps the actual character by the previous one. If it is the end of the line
// then it is swapped the 2nd previous by the previous one.
func (b *buffer) swap() bool {
	if b.pos == b.promptLen {
		return false
	}

	if b.pos < b.size {
		aux := b.data[b.pos-1]
		b.data[b.pos-1] = b.data[b.pos]
		b.data[b.pos] = aux
		b.pos++
		// End of line
	} else {
		aux := b.data[b.pos-2]
		b.data[b.pos-2] = b.data[b.pos-1]
		b.data[b.pos-1] = aux
	}

	return true
}


// === Utility
// ===

// Grows buffer to guarantee space for n more byte.
func (b *buffer) grow(n int) {
	for n > len(b.data) {
		b.data = b.data[:len(b.data)+BufferLen]
	}
}

// Returns the coordinates of a position for a line of size given in `columns`.
func (b *buffer) pos2xy(pos int) (line, column int) {
	if pos < columns {
		return 0, pos
	}

	// pos--
	line = pos / columns
	column = pos - (line * columns)
	return
}

