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
	"fmt"
	"io"
	"os"
)


// ASCII codes
const (
	_ESC       = 27 // Escape: Ctrl-[ (033 in octal)
	_L_BRACKET = 91 // Left square bracket: [
)

// ANSI terminal escape controls
var (
	cursorToleft   = []byte("\033[0G")        // Cursor to left edge.
	delScreen      = []byte("\033[2J")        // Erase the screen.
	toleftDelRight = []byte("\033[0K\033[0G") // Cursor to left; erase to right.
)

// Characters to print in case of interrumption
var (
	ctrl_c = []byte("^C")
	ctrl_d = []byte("^D")
)


// === Type
// ===

// Represents a line.
type line struct {
//	cursor    int      // Location pointer into buffer
	promptLen int      // Prompt size
	prompt    string   // Primary prompt
	//data      *buffer  // Text buffer
	*buffer  // Text buffer
	hist      *history // History file
	input     io.Reader
	output    io.Writer
}

func NewLine(in io.Reader, out io.Writer, hist *history, prompt string) *line {
	return &line{
		len(prompt),
		prompt,
		newBuffer(),
		hist,
		in,
		out,
	}
}


// === Output
// ===

// Returns a slice of the contents of the buffer.
func (ln *line) Bytes() []byte { return ln.data }

// Returns the contents of the buffer as a string.
func (ln *line) String() string { return string(ln.data[:ln.size]) }

// Refreshes the line.
func (ln *line) Refresh() (err os.Error) {
	if _, err = ln.output.Write(cursorToleft); err != nil {
		return
	}
	if _, err = fmt.Fprint(ln.output, ln.prompt); err != nil {
		return
	}
	if _, err = ln.output.Write(ln.Bytes()[:ln.size]); err != nil {
		return
	}
	if _, err = ln.output.Write(toleftDelRight); err != nil {
		return
	}
	// Move cursor to original position.
	_, err = fmt.Fprintf(ln.output, "\033[%dC", ln.promptLen+ln.cursor)

	return
}


// === Get
// ===

func (ln *line) Run() (err os.Error) {
	// Print the primary prompt.
	if _, err = fmt.Fprint(ln.output, ln.prompt); err != nil {
		return
	}

	in := bufio.NewReader(ln.input) // Read input.
	buf := newBuffer()  // Store characters entered.
	seq := make([]byte, 2)  // For escape sequences.
	seq2 := make([]byte, 2) // Extended escape sequences.

	for {
		rune, runeSize, err := in.ReadRune()
		if err != nil {
			return
		}

		switch rune {
		default:
			if err = buf.InsertRune(rune, runeSize); err != nil {
				return
			}
			continue

		case 9: // horizontal tab
			continue

		case 13: // enter
			ln.hist.Add(ln.String())
			goto deleteLine

		case 127, 8: // backspace, Ctrl-h
			if buf.DeleteLast() {
				goto refresh
			}
			continue

		case 20: // Ctrl-t //!!! add
			continue

		case 2: // Ctrl-b
			goto leftArrow

		case 6: // Ctrl-f
			goto rightArrow

		case 16: // Ctrl-p
			seq[1] = 65
			goto upDownArrow

		case 14: // Ctrl-n
			seq[1] = 66
			goto upDownArrow

		case 3, 4: // Ctrl-c, Ctrl-d
			if rune == 3 {
				buf.InsertByte(ctrl_c)
				err = ErrCtrlC
			} else {
				buf.InsertByte(ctrl_d)
				err = ErrCtrlD
			}

			ln.Refresh()
			return

		// Escape sequence
		case _ESC:
			if _, err = in.Read(seq); err != nil {
				return
			}

			if seq[0] == _L_BRACKET {
				switch seq[1] {
				case 68:
					goto leftArrow
				case 67:
					goto rightArrow
				case 65, 66:
					goto upDownArrow
				}

				// Extended escape.
				if seq[1] > 48 && seq[1] < 55 {
					if _, err = in.Read(seq2); err != nil {
						return
					}
//!!! doesn't works
					if seq[1] == 51 && seq2[0] == 126 { // Delete
						if buf.DeleteNext() {
							goto refresh
						}
						continue
					}
				}
			}

		case 21: // Ctrl+u, delete the whole line.
			goto deleteLine

		case 11: // Ctrl+k, delete from current to end of line.
			buf.size = buf.cursor
			goto refresh

		case 1: // Ctrl+a, go to the start of the line.
			buf.cursor = 0
			goto refresh

		case 5: // Ctrl+e, go to the end of the line.
			buf.cursor = buf.size
			goto refresh
		}

	leftArrow:
		if buf.CursorToleft() {
			goto refresh
		}
		continue

	rightArrow:
		if buf.CursorToright() {
			goto refresh
		}
		continue

	upDownArrow:
		if ln.hist.Len > 1 {
			// Update the current history entry before to overwrite it with tne
			// next one.
			
		}
		continue

	deleteLine:
		buf.cursor, buf.size = 0, 0
		//goto refresh

	refresh:
		if err = ln.Refresh(); err != nil {
			return
		}
		continue
	}

	return
}

