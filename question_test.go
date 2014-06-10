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
	"testing"
)

func TestQuest(t *testing.T) {
	// === To handle boolean strings in french
	QuestionTrueString = "oui"
	QuestionFalseString = "non"
	ExtraBoolString = map[string]bool{
		"oui": true,
		"non": false,
		"o":   true,
	}

	fmt.Println("\n == Questions\n")

	q := NewQuestion()
	defer q.RestoreTerm()

	ans, err := q.Read("What is your name?")
	print(ans, err)

	ans, err = q.ReadStringDefault("What color is your hair?", "brown")
	print(ans, err)

	bAns, err := q.ReadBool("Do you watch television?", true)
	print(bAns, err)

	color := []string{"red", "blue", "black"}
	ans, err = q.ReadChoice("What is you favorite color?", color)
	print(ans, err)

	iAns, err := q.ReadIntDefault("What is your age?", 16)
	print(iAns, err)

	fAns, err := q.ReadFloat("How tall are you?")
	print(fAns, err)

	if _, err := atob("not-found"); err == nil {
		t.Error("should return an error")
	}
}

func print(a interface{}, err error) {
	if err == nil {
		fmt.Printf("  answer: %v\r\n", a)
	} else if err != ErrCtrlD {
		fmt.Printf(err.Error() + "\r\n")
	}
}
