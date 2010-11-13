// Copyright 2010  The "go-linoise" Authors
//
// Use of this source code is governed by the Simplified BSD License
// that can be found in the LICENSE file.
//
// This software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied. See the License
// for more details.

/* Important: linoise sets tty in 'raw mode' so there is to use CR+LF (\r\n) at
writing.
*/

package linoise

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"utf8"

	"github.com/kless/go-term/term"
)


// Values by default for prompts.
var (
	PS1 = "linoise$ "
	PS2 = "> "
)

// Input / Output
var (
	input  *os.File = os.Stdin
	output *os.File = os.Stdout
)


// === Type
// ===

// Represents a line.
type Line struct {
	useHistory bool
	ps1Len     int      // Primary prompt size
	ps1        string   // Primary prompt
	ps2        string   // Command continuations
	*buffer             // Text buffer
	hist       *history // History file
}


// Gets a line type using the primary prompt by default. Sets the TTY raw mode.
func NewLine(hist *history) *Line {
	term.MakeRaw()

	buf := newBuffer(len(PS1))
	buf.insertRunes([]int(PS1))

	return &Line{
		hasHistory(hist),
		len(PS1),
		PS1,
		PS2,
		buf,
		hist,
	}
}

// Gets a line type using the given prompt as primary. Sets the TTY raw mode.
// 'ansiLen' is the length of ANSI codes that the prompt could have.
func NewLinePrompt(prompt string, ansiLen int, hist *history) *Line {
	term.MakeRaw()

	buf := newBuffer(len(prompt) - ansiLen)
	buf.insertRunes([]int(prompt))

	return &Line{
		hasHistory(hist),
		len(prompt) - ansiLen,
		prompt,
		PS2,
		buf,
		hist,
	}
}

// Restores terminal settings so it is disabled the raw mode.
func (ln *Line) RestoreTerm() {
	term.RestoreTerm()
}

// Tests if it has an history file.
func hasHistory(h *history) bool {
	if h == nil {
		return false
	}
	return true
}


// === Output
// ===

// Returns a slice of the contents of the buffer.
func (ln *Line) toBytes() []byte {
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
func (ln *Line) toString() string { return string(ln.data[:ln.size]) }

// Prints the primary prompt.
func (ln *Line) prompt() (err os.Error) {
	if lines, err = ln.end(); err != nil {
		return err
	}

	for lines > 0 {
		if _, err = output.Write(delLine_cursorUp); err != nil {
			return OutputError(err.String())
		}
		lines--
	}

	if _, err = output.Write(delLine_CR); err != nil {
		return OutputError(err.String())
	}
	if _, err = fmt.Fprint(output, ln.ps1); err != nil {
		return OutputError(err.String())
	}

	ln.pos, ln.size = ln.ps1Len, ln.ps1Len
	return
}

// Refreshes the line.
func (ln *Line) refresh() (err os.Error) {
	lines, err := ln.end()
	if err != nil {
		return err
	}

	lineToRefresh := lines
	for lineToRefresh > 0 {
		if _, err = output.Write(delLine_cursorUp); err != nil {
			return OutputError(err.String())
		}
		lineToRefresh--
	}

	if _, err = output.Write(delLine_CR); err != nil {
		return OutputError(err.String())
	}
	if _, err = output.Write(ln.toBytes()); err != nil {
		return OutputError(err.String())
	}

	// === Move cursor to original position.
	for lineToRefresh < lines {
		if _, err = output.Write(cursorDown); err != nil {
			return OutputError(err.String())
		}
		lineToRefresh++
	}
	if _, err = fmt.Fprintf(output, "\r\033[%dC", ln.pos); err != nil {
		return OutputError(err.String())
	}

	return nil
}


// === Get
// ===

// Reads charactes from input to write them to output, allowing line editing.
// The errors that could return are to indicate if Ctrl-D was pressed, and for
// both input / output errors.
func (ln *Line) Read() (line string, err os.Error) {
	var anotherLine []int  // For lines got from history.
	var isHistoryUsed bool // If the history has been accessed.
	var useRefresh bool

	in := bufio.NewReader(input) // Read input.
	seq := make([]byte, 2)       // For escape sequences.
	seq2 := make([]byte, 2)      // Extended escape sequences.

	for {
		rune, _, err := in.ReadRune()
		if err != nil {
			return "", InputError(err.String())
		}

		switch rune {
		default:
			useRefresh, err = ln.insertRune(rune)
			if err != nil {
				return "", err
			}
			if useRefresh {
				if err = ln.refresh(); err != nil {
					return "", err
				}
			}

			continue

		case 13: // enter
			line = ln.toString()

			if ln.useHistory {
				ln.hist.Add(line)
			}

			if _, err = output.Write(_CR_LF); err != nil {
				return "", OutputError(err.String())
			}

			return strings.TrimSpace(line), nil

		case 127, 8: // backspace, Ctrl-h
		useRefresh, err = ln.deletePrev()
		if err != nil {
			return "", err
		}
		if useRefresh {
			if err = ln.refresh(); err != nil {
				return "", err
			}
		}

		continue

		case 9: // horizontal tab
			// TODO: disabled by now
			continue

		case 3: // Ctrl-c
			useRefresh, err := ln.insertRunes(ctrlC)
			if err != nil {
				return "", err
			}
			if useRefresh {
				if err = ln.refresh(); err != nil {
					return "", err
				}
			}

			if _, err = output.Write(_CR_LF); err != nil {
				return "", OutputError(err.String())
			}
			if err = ln.prompt(); err != nil {
				return "", err
			}

			continue

		case 4: // Ctrl-d
			useRefresh, err := ln.insertRunes(ctrlD)
			if err != nil {
				return "", err
			}
			if useRefresh {
				if err = ln.refresh(); err != nil {
					return "", err
				}
			}

			if _, err = output.Write(_CR_LF); err != nil {
				return "", OutputError(err.String())
			}

			return "", ErrCtrlD

		// Escape sequence
		case _ESC:
			if _, err = in.Read(seq); err != nil {
				return "", InputError(err.String())
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
						return "", InputError(err.String())
					}
					//fmt.Print(" >>", seq2) //!!! For DEBUG

					// TODO: doesn't works
					if seq[1] == 51 && seq2[0] == 126 { // Delete
						useRefresh, err = ln.delete()
						if err != nil {
							return "", err
						}
						if useRefresh {
							if err = ln.refresh(); err != nil {
								return "", err
							}
						}
					}
				}
				continue
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
			if ln.swap() {
				goto _refresh
			}
			continue

		case 21: // Ctrl+u, delete the whole line.
			goto _deleteLine

		case 11: // Ctrl+k, delete from current to end of line.
			if err = ln.deleteRight(); err != nil {
				return "", err
			}
			continue

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

	_upDownArrow: // Up and down arrow: history
		if !ln.useHistory {
			continue
		}

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
		// TODO: it has to be removed before of to be saved the history
		if !isHistoryUsed {
			ln.hist.Add(ln.toString())
		}
		isHistoryUsed = true

		ln.grow(len(anotherLine))
		ln.size = len(anotherLine)
		copy(ln.data[0:], anotherLine)
		goto _refresh

	_leftArrow:
		if err = ln.backward(); err != nil {
			return "", err
		}
		continue

	_rightArrow:
		if err = ln.forward(); err != nil {
			return "", err
		}
		continue

	_start:
		if err = ln.start(); err != nil {
			return "", err
		}
		continue

	_end:
		if _, err = ln.end(); err != nil {
			return "", err
		}
		continue

	_deleteLine:
		if err = ln.prompt(); err != nil {
			return "", err
		}
		continue

	_refresh:
		if err = ln.refresh(); err != nil {
			return "", err
		}
	}
	return
}

