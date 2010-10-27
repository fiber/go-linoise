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
	"strings"
)


// Reads the key pressed. The argument `prompt` is written to standard output,
// if any.
func ReadKey(prompt string) (rune int, err os.Error) {
	in, err := bufio.NewReaderSize(os.Stdin, 1)
	if err != nil {
		return 0, err
	}

	if prompt != "" {
		fmt.Print(prompt)
	}

	rune, _, err = in.ReadRune()
	if err != nil {
		return 0, err
	}

	return rune, nil
}

// Reads a line from input until Return is pressed (stripping a trailing
// newline), and returns that. The argument `prompt` is written to standard
// output, if any.
func ReadString(prompt string) (line string, err os.Error) {
	in := bufio.NewReader(os.Stdin)

	if prompt != "" {
		fmt.Print(prompt)
	}

	line, err = in.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimRight(line, "\n"), nil
}

