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

	ans := q.ReadString("What is your name?")
	print(ans)

	ans = q.ReadStringDefault("What color is your hair?", "brown")
	print(ans)

	bAns := q.ReadBool("Do you watch television?", true)
	print(bAns)

	color := []string{"red", "blue", "black"}
	ans = q.ReadChoice("What is you favorite color?", color)
	print(ans)

	iAns := q.ReadIntDefault("What is your age?", 16)
	print(iAns)

	fAns := q.ReadFloat("How tall are you?")
	print(fAns)

	if _, err := atob("not-found"); err == nil {
		t.Error("should return an error")
	}
}

func print(v interface{}) {
	fmt.Printf("  answer: %v\r\n", v)
}

