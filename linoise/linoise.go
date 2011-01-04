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


// === Init
// ===

func init() {
	if err := term.CheckIsatty(input.Fd()); err != nil {
		panic(err)
	}
}


// === Type
// ===

// Represents a line.
type Line struct {
	useHistory bool
	ps1Len     int      // Primary prompt size
	ps1        string   // Primary prompt
	ps2        string   // Command continuations
	buf        *buffer  // Text buffer
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

// Prints the primary prompt.
func (ln *Line) prompt() (err os.Error) {
	if _, err = output.Write(delLine_CR); err != nil {
		return outputError(err.String())
	}
	if _, err = fmt.Fprint(output, ln.ps1); err != nil {
		return outputError(err.String())
	}

	ln.buf.pos, ln.buf.size = ln.ps1Len, ln.ps1Len
	return
}


// === Get
// ===

// Reads charactes from input to write them to output, allowing line editing.
// The errors that could return are to indicate if Ctrl-D was pressed, and for
// both input / output errors.
func (ln *Line) Read() (line string, err os.Error) {
	var anotherLine []int  // For lines got from history.
	var isHistoryUsed bool // If the history has been accessed.

	in := bufio.NewReader(input) // Read input.
	seq := make([]byte, 2)       // For escape sequences.
	seq2 := make([]byte, 1)      // Extended escape sequences.

	// Print the primary prompt.
	if err = ln.prompt(); err != nil {
		return "", err
	}

	// === Detect change of window size.
	go term.TrapWinsize()

	go func() {
		for {
			<-term.WinsizeChan // Wait for.

			_, ln.buf.winColumns = term.GetWinsizeInChar()
			ln.buf.refresh()
		}
	}()

	for {
		rune, _, err := in.ReadRune()
		if err != nil {
			return "", inputError(err.String())
		}

		switch rune {
		default:
			if err = ln.buf.insertRune(rune); err != nil {
				return "", err
			}
			continue

		case 13: // enter
			line = ln.buf.toString()

			if ln.useHistory {
				ln.hist.Add(line)
			}
			if _, err = output.Write(_CR_LF); err != nil {
				return "", outputError(err.String())
			}

			return strings.TrimSpace(line), nil

		case 127, 8: // backspace, Ctrl-h
			if err = ln.buf.deletePrev(); err != nil {
				return "", err
			}
			continue

		case 9: // horizontal tab
			// TODO: disabled by now
			continue

		case 3: // Ctrl-c
			if err = ln.buf.insertRunes(ctrlC); err != nil {
				return "", err
			}
			if _, err = output.Write(_CR_LF); err != nil {
				return "", outputError(err.String())
			}
			if err = ln.prompt(); err != nil {
				return "", err
			}

			continue

		case 4: // Ctrl-d
			if err = ln.buf.insertRunes(ctrlD); err != nil {
				return "", err
			}
			if _, err = output.Write(_CR_LF); err != nil {
				return "", outputError(err.String())
			}

			return "", ErrCtrlD

		// Escape sequence
		case 27: // Escape: Ctrl-[ ("033" in octal, "\x1b" in hexadecimal)
			if _, err = in.Read(seq); err != nil {
				return "", inputError(err.String())
			}

			if seq[0] == 79 { // 'O'
				switch seq[1] {
				case 72: // Home: "\x1bOH"
					goto _start
				case 70: // End: "\x1bOF"
					goto _end
				}
			}

			if seq[0] == 91 { // Left square bracket: "["
				switch seq[1] {
				case 68: // "\x1b[D"
					goto _leftArrow
				case 67: // "\x1b[C"
					goto _rightArrow
				case 65, 66: // Up: "\x1b[A"; Down: "\x1b[B"
					goto _upDownArrow
				}

				// Extended escape.
				if seq[1] > 48 && seq[1] < 55 {
					if _, err = in.Read(seq2); err != nil {
						return "", inputError(err.String())
					}

					if seq2[0] == 126 { // '~'
						switch seq[1] {
						//case 50: // Insert: "\x1b[2~"
							
						case 51: // Delete: "\x1b[3~"
							if err = ln.buf.delete(); err != nil {
								return "", err
							}
						//case 53: // RePag: "\x1b[5~"
							
						//case 54: // AvPag: "\x1b[6~"
							
						}
					}
				}
			}
			continue

		case 20: // Ctrl-t, swap actual character by the previous one.
			if err = ln.buf.swap(); err != nil {
				return "", err
			}
			continue

		case 21: // Ctrl+u, delete the whole line.
			if err = ln.buf.deleteLine(); err != nil {
				return "", err
			}
			if err = ln.prompt(); err != nil {
				return "", err
			}
			continue

		case 11: // Ctrl+k, delete from current to end of line.
			if err = ln.buf.deleteRight(); err != nil {
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
			ln.hist.Add(ln.buf.toString())
		}
		isHistoryUsed = true

		ln.buf.grow(len(anotherLine))
		ln.buf.size = len(anotherLine) + ln.buf.promptLen
		copy(ln.buf.data[ln.ps1Len:], anotherLine)

		if err = ln.buf.refresh(); err != nil {
			return "", err
		}
		continue

	_leftArrow:
		if err = ln.buf.backward(); err != nil {
			return "", err
		}
		continue

	_rightArrow:
		if err = ln.buf.forward(); err != nil {
			return "", err
		}
		continue

	_start:
		if err = ln.buf.start(); err != nil {
			return "", err
		}
		continue

	_end:
		if _, err = ln.buf.end(); err != nil {
			return "", err
		}
		continue
	}
	return
}

