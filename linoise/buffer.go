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
	size   int   // Amount of characters added
	cursor int   // Location pointer into buffer
	data   []int // Text buffer
}

func newBuffer() *buffer {
	return &buffer{0, 0, make([]int, Length, Capacity)}
}
// ===


// Inserts a character in the cursor position.
func (b *buffer) Insert(rune int) (useRefresh bool, err os.Error) {
	// Grow buffer to guarantee space for one more byte.
	if b.size+1 > len(b.data) {
		b.data = b.data[:len(b.data)+Length]
	}

	// Avoid a full update of the line.
	if b.cursor == b.size {
		char := make([]byte, utf8.UTFMax)
		utf8.EncodeRune(rune, char)
		if _, err = Output.Write(char); err != nil {
			return
		}
	} else {
		useRefresh = true
		copy(b.data[b.cursor+1:b.size+1], b.data[b.cursor:b.size])
	}

	b.data[b.cursor] = rune
	b.cursor ++
	b.size ++

	return
}

// Moves the cursor one character backward.
func (b *buffer) Left() bool {
	if b.cursor > 0 {
		b.cursor--
		return true
	}
	return false
}

// Moves the cursor one character forward.
func (b *buffer) Right() bool {
	if b.cursor < b.size {
		b.cursor++
		return true
	}
	return false
}

// Deletes the character in cursor.
func (b *buffer) Delete() bool {
	if b.size > 0 && b.cursor < b.size {
		copy(b.data[b.cursor:], b.data[b.cursor+1:b.size])
		b.size--
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

