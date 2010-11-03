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
	"utf8"
)


// Buffer size
var (
	BufferCap = 4096
	BufferLen = 64 // Initial length
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
	return &buffer{0, 0, make([]int, BufferLen, BufferCap)}
}
// ===


// Grows buffer to guarantee space for n more byte.
func (b *buffer) grow(n int) {
	for n > len(b.data) {
		b.data = b.data[:len(b.data)+BufferLen]
	}
}

// Inserts a character in the cursor position.
func (b *buffer) InsertRune(rune int) (useRefresh bool) {
	b.grow(b.size + 1) // Check if there is free space for one more character

	// Avoid a full update of the line.
	if b.cursor == b.size {
		char := make([]byte, utf8.UTFMax)
		utf8.EncodeRune(rune, char)
		output.Write(char)
	} else {
		useRefresh = true
		copy(b.data[b.cursor+1:b.size+1], b.data[b.cursor:b.size])
	}

	b.data[b.cursor] = rune
	b.cursor++
	b.size++

	return
}

// Inserts several characters.
func (b *buffer) InsertRunes(runes []int) (useRefresh bool) {
	for _, r := range runes {
		useRefresh = b.InsertRune(r)
	}
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

// Deletes the previous character from cursor.
func (b *buffer) DeletePrev() bool {
	if b.cursor > 0 {
		copy(b.data[b.cursor-1:], b.data[b.cursor:b.size])
		b.cursor--
		b.size--
		return true
	}
	return false
}

// Swaps the actual character by the previous one. If it is the end of the line
// then it is swapped the 2nd previous by the previous one.
func (b *buffer) Swap() bool {
	if b.cursor <= 0 {
		return false
	}

	if b.cursor < b.size {
		aux := b.data[b.cursor-1]
		b.data[b.cursor-1] = b.data[b.cursor]
		b.data[b.cursor] = aux
		b.cursor++
		// End of line
	} else {
		aux := b.data[b.cursor-2]
		b.data[b.cursor-2] = b.data[b.cursor-1]
		b.data[b.cursor-1] = aux
	}

	return true
}

