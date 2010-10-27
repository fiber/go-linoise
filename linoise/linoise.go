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
	"os"
	"utf8"
)


// Values by default
var (
	Input  *os.File = os.Stdin
	Output *os.File = os.Stdout
	PS1    = "linoise$ "
	PS2    = "> "
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


// === Type
// ===

// Represents a line.
type line struct {
	ps1Len  int      // Primary prompt size
	ps1     string   // Primary prompt
	ps2     string   // Command continuations
	*buffer          // Text buffer
	hist    *history // History file
}

// Gets a line type using the given prompt as primary.
func NewLinePrompt(hist *history, prompt string) *line {
	return &line{
		len(prompt),
		prompt,
		PS2,
		newBuffer(),
		hist,
	}
}

// Gets a line type using the primary prompt by default.
func NewLine(hist *history) *line {
	return &line{
		len(PS1),
		PS1,
		PS2,
		newBuffer(),
		hist,
	}
}


// === Output
// ===

// Returns a slice of the contents of the buffer.
func (ln *line) Bytes() []byte {
	chars := make([]byte, ln.size*utf8.UTFMax)
	var end, runeLen int

	// === Each character (as integer) is encoded to []byte
	for i := 0; i < ln.size; i++ {
		if i != 0 {
			runeLen = utf8.EncodeRune(ln.data[i], chars[end:])
			end += runeLen
		} else {
			runeLen = utf8.EncodeRune(ln.data[i], chars)
			end = runeLen
		}
	}
	return chars[:end]
}

// Returns the contents of the buffer as a string.
func (ln *line) String() string { return string(ln.data[:ln.size]) }

// Refreshes the line.
func (ln *line) Refresh() {
	_, err := Output.Write(cursorToleft)
	if err != nil {
		goto _error
	}
	if _, err = fmt.Fprint(Output, ln.ps1); err != nil {
		goto _error
	}
	if _, err = Output.Write(ln.Bytes()); err != nil {
		goto _error
	}
	if _, err = Output.Write(toleftDelRight); err != nil {
		goto _error
	}
	// Move cursor to original position.
	if _, err = fmt.Fprintf(Output, "\033[%dC", ln.ps1Len+ln.cursor); err != nil {
		goto _error
	}

	return

	_error:
		reportPanic("line.Refresh", err)
}

// Reports an panic message, printing the function name `f`.
func reportPanic(f string, err os.Error) {
	fmt.Println()
	Output.Write(cursorToleft)
	panic(fmt.Sprintf("linoise: %s: %s", f, err.String()))
}


// === Get
// ===

func (ln *line) Run() {
	var anotherLine []int  // For lines got from history.
	var isHistoryUsed bool // If the history has been accessed.

	in := bufio.NewReader(Input) // Read input.
	seq := make([]byte, 2)       // For escape sequences.
	seq2 := make([]byte, 2)      // Extended escape sequences.

	// Print the primary prompt.
	_, err := fmt.Fprint(Output, ln.ps1)
	if err != nil {
		reportPanic("line.Run", err)
	}

	for {
		rune, _, err := in.ReadRune()
		if err != nil {
			reportPanic("line.Run", err)
		}

		switch rune {
		default:
			useRefresh, err := ln.Insert(rune)
			if err != nil {
				reportPanic("line.Run", err)
			}

			if useRefresh {
				ln.Refresh()
			}
			continue

		case 13: // enter
			ln.hist.Add(ln.String())
			goto _deleteLine

		case 127, 8: // backspace, Ctrl-h
			if ln.DeletePrev() {
				goto _refresh
			}
			continue

		case 9: // horizontal tab
			//!!! disabled by now
			continue

		case 3: // Ctrl-c
			ln.Insert('^')
			ln.Insert('C')

			if _, err = fmt.Fprint(Output, "\n"); err != nil {
				reportPanic("line.Run", err)
			}
			goto _deleteLine

		case 4: // Ctrl-d
			ln.Insert('^')
			ln.Insert('D')
			ln.Refresh()
			return

		// Escape sequence
		case _ESC:
			if _, err = in.Read(seq); err != nil {
				reportPanic("line.Run", err)
			}
			//fmt.Print(" >", seq) //!!! For DEBUG

			if seq[0] == _L_BRACKET {
				switch seq[1] {
				case 68:
					goto _leftArrow
				case 67:
					goto _rightArrow
				case 65, 66: // Up, Down
					goto _upDownArrow
				}

				// Extended escape.
				if seq[1] > 48 && seq[1] < 55 {
					if _, err = in.Read(seq2); err != nil {
						reportPanic("line.Run", err)
					}
					//fmt.Print(" >>", seq2) //!!! For DEBUG

					//!!! doesn't works
					if seq[1] == 51 && seq2[0] == 126 { // Delete
						if ln.Delete() {
							goto _refresh
						}
					}
				}
			}
			if seq[0] == 79 {
				switch seq[1] {
				case 72: // Home
					goto _start
				case 70: // End
					goto _end
				}
			}
			continue

		case 20: // Ctrl-t, swap actual character by the previous one.
			if ln.Swap() {
				goto _refresh
			}
			continue

		case 21: // Ctrl+u, delete the whole line.
			goto _deleteLine

		case 11: // Ctrl+k, delete from current to end of line.
			ln.size = ln.cursor
			goto _refresh

		case 1: // Ctrl+a, go to the start of the line.
			goto _start

		case 5: // Ctrl+e, go to the end of the line.
			goto _end

		case 2: // Ctrl-b
			goto _leftArrow

		case 6: // Ctrl-f
			goto _rightArrow

		case 16: // Ctrl-p
			seq[1] = 65
			goto _upDownArrow

		case 14: // Ctrl-n
			seq[1] = 66
			goto _upDownArrow
		}
			continue // To be safe.

	_upDownArrow: // Up and down arrow: history
		// Up
		if seq[1] == 65 {
			anotherLine, err = ln.hist.Prev()
			// Down
		} else {
			anotherLine, err = ln.hist.Next()
		}
		if err != nil {
			continue
		}

		// Update the current history entry before to overwrite it with
		// the next one.
		//!!! it has to be removed before of to be saved the history
		if !isHistoryUsed {
			ln.hist.Add(ln.String())
		}
		isHistoryUsed = true

		ln.grow(len(anotherLine))
		ln.size = len(anotherLine)
		ln.cursor = len(anotherLine)
		copy(ln.data[0:], anotherLine)
		goto _refresh

	_leftArrow:
		if ln.Left() {
			goto _refresh
		}
		continue

	_rightArrow:
		if ln.Right() {
			goto _refresh
		}
		continue

	_start:
		ln.cursor = 0
		goto _refresh

	_end:
		ln.cursor = ln.size
		goto _refresh

	_deleteLine:
		ln.cursor, ln.size = 0, 0
		//goto _refresh

	_refresh:
		ln.Refresh()
	}
}

