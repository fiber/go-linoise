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
	"utf8"
)


// Buffer size
var (
	Capacity = 4096
	Length   = 64 // Initial length
)


// === Type
// ===

// Represents the line buffer.
type buffer struct {
	size   int    // Amount of characters added
	cursor int    // Location pointer into buffer
	data   []byte // Text buffer
}

func newBuffer() *buffer {
	return &buffer{0, 0, make([]byte, Length, Capacity)}
}
// ===


// Grows buffer to guarantee space for n more bytes.
func (b *buffer) grow(n int) {
	if b.size+n > len(b.data) {
		b.data = b.data[:len(b.data)+Length]
	}
}

// Base to insert characters immediately after the cursor position.
func (b *buffer) _Insert(chars []byte) (useRefresh bool, err os.Error) {
	b.grow(len(chars)) // Check the free space.

	// Avoid a full update of the line.
	if b.cursor == b.size {
		if _, err = Output.Write(chars); err != nil {
			return
		}
	} else {
		useRefresh = true
		copy(b.data[b.cursor+len(chars):b.size+len(chars)],
			b.data[b.cursor:b.size])
	}

	for _, v := range chars {
		b.data[b.cursor] = v
		b.cursor ++
		b.size ++
	}
	return
}

// Inserts a byte after of cursor.
func (b *buffer) InsertByte(chars []byte) (useRefresh bool, err os.Error) {
	useRefresh, err = b._Insert(chars)
	return
}

// Inserts an unicode character after of cursor.
func (b *buffer) InsertRune(rune, runeSize int) (useRefresh bool, err os.Error) {
	runeEncoded := make([]byte, runeSize)
	utf8.EncodeRune(rune, runeEncoded)

	useRefresh, err = b._Insert(runeEncoded)
	return
}

// Moves the cursor one character backward.
func (b *buffer) CursorToleft() bool {
	if b.cursor > 0 {
		b.cursor--
		return true
	}
	return false
}

// Moves the cursor one character forward.
func (b *buffer) CursorToright() bool {
	if b.cursor < b.size {
		b.cursor++
		return true
	}
	return false
}

// Deletes the last character from cursor.
func (b *buffer) DeleteLast() bool {
	if b.cursor > 0 {
		copy(b.data[b.cursor-1:], b.data[b.cursor:b.size])
		b.cursor--
		b.size--
		return true
	}
	return false
}

// Deletes the next character from cursor.
func (b *buffer) DeleteNext() bool {
	if b.size > 0 && b.cursor < b.size {
		copy(b.data[b.cursor:], b.data[b.cursor+1:b.size])
		b.size--
		return true
	}
	return false
}

