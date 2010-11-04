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

	q := NewQuestion()
	defer q.RestoreTerm()

	ans := q.ReadString("What is your name?")
	fmt.Printf("  answer: %q\n", ans)

	ans = q.ReadStringDefault("What color is your hair?", "brown")
	fmt.Printf("  answer: %q\n", ans)

	bAns := q.ReadBool("Do you watch television?", true)
	fmt.Printf("  answer: %t\n", bAns)

	color := []string{"red", "blue", "black"}
	ans = q.ReadChoice("What is you favorite color?", color)
	fmt.Printf("  answer: %q\n", ans)

	iAns := q.ReadIntDefault("What is your age?", 16)
	fmt.Printf("  answer: %d\n", iAns)

	fAns := q.ReadFloat("How tall are you?")
	fmt.Printf("  answer: %f\n", fAns)
}

