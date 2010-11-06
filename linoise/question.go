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
	"os"
	"strconv"
	"strings"

	"github.com/kless/go-term/term"
)


// Values by default
var (
	QuestionPrefix      = " + " // String placed before of questions
	QuestionErrPrefix   = "  "  // String placed before of error messages
	QuestionTrueString  = "y"   // String to represent 'true' by default
	QuestionFalseString = "n"   // Idem but for 'false'

	QuestionFloatFmt  byte = 'g' // Format for float numbers
	QuestionFloatPrec int  = -1  // Precision for float numbers
)

// To pass strings in another languages.
var ExtraBoolString = make(map[string]bool)


// === Type
// ===

type Question struct {
	trueString, falseString string // Strings that represent booleans.
}


// Gets a question type.
func NewQuestion() *Question {
	// === Check the strings that represent a boolean.
	_, err := atob(QuestionTrueString)
	if err != nil {
		panic(fmt.Sprintf("the string %q does not represent a boolean 'true'",
			QuestionTrueString))
	}

	if _, err = atob(QuestionFalseString); err != nil {
		panic(fmt.Sprintf("the string %q does not represent a boolean 'false'",
			QuestionFalseString))
	}

	return &Question{
		strings.ToLower(QuestionTrueString),
		strings.ToLower(QuestionFalseString),
	}
}

// Restores terminal settings.
func (q *Question) RestoreTerm() {
	term.RestoreTerm()
}
// ===


// Gets a line type ready to show questions.
func (q *Question) getLine(prompt, defaultAnswer string, hasDefault bool) *Line {
	prompt = QuestionPrefix + prompt

	// Add the value by default
	if hasDefault {
		prompt = fmt.Sprintf("%s [%s]", prompt, defaultAnswer)
	}

	// Add spaces
	if strings.HasSuffix(prompt, "?") {
		prompt += " "
	} else {
		prompt += ": "
	}

	return NewLinePrompt(prompt, nil) // No history.
}

// Prints the prompt waiting to get a string not empty.
func (q *Question) Read(prompt string) (answer string, err os.Error) {
	line := q.getLine(prompt, "", false)

	for {
		answer, err = line.Read()
		if err != nil {
			return "", err
		}
		if answer != "" {
			return
		}
	}
	return
}

// Base to read strings.
func (q *Question) _baseReadString(prompt, defaultAnswer string, hasDefault bool) (answer string, err os.Error) {
	line := q.getLine(prompt, defaultAnswer, hasDefault)

	for {
		answer, err = line.Read()
		if err != nil {
			return "", err
		}
		if answer != "" {
			// === Check if it is a number
			if _, err := strconv.Atoi(answer); err == nil {
				goto _error
			}
			if _, err := strconv.Atof(answer); err == nil {
				goto _error
			}

			return answer, nil
		}

		if hasDefault {
			return defaultAnswer, nil
		}
		continue

	_error:
		fmt.Fprintf(output, "%s%v: the value has to be a string\r\n",
			QuestionErrPrefix, answer)
	}
	return
}

// Prints the prompt waiting to get a string.
func (q *Question) ReadString(prompt string) (answer string, err os.Error) {
	return q._baseReadString(prompt, "", false)
}

// Prints the prompt waiting to get a string.
// If input is nil then it returns the answer by default.
func (q *Question) ReadStringDefault(prompt, defaultAnswer string) (answer string, err os.Error) {
	return q._baseReadString(prompt, defaultAnswer, true)
}

// Base to read integer numbers.
func (q *Question) _baseReadInt(prompt string, defaultAnswer int, hasDefault bool) (answer int, err os.Error) {
	line := q.getLine(prompt, strconv.Itoa(defaultAnswer), hasDefault)

	for {
		input, err := line.Read()
		if err != nil {
			return 0, err
		}
		if input == "" && hasDefault {
			return defaultAnswer, nil
		}

		answer, err = strconv.Atoi(input)
		if err != nil {
			fmt.Fprintf(output, "%s%q: the value has to be an integer\r\n",
				QuestionErrPrefix, input)
			continue
		} else {
			return answer, nil
		}
	}
	return
}

// Prints the prompt waiting to get an integer number.
func (q *Question) ReadInt(prompt string) (answer int, err os.Error) {
	return q._baseReadInt(prompt, 0, false)
}

// Prints the prompt waiting to get an integer number.
// If input is nil then it returns the answer by default.
func (q *Question) ReadIntDefault(prompt string, defaultAnswer int) (answer int, err os.Error) {
	return q._baseReadInt(prompt, defaultAnswer, true)
}

// Base to read float numbers.
func (q *Question) _baseReadFloat(prompt string, defaultAnswer float, hasDefault bool) (answer float, err os.Error) {
	line := q.getLine(
		prompt,
		strconv.Ftoa(defaultAnswer, QuestionFloatFmt, QuestionFloatPrec),
		hasDefault,
	)

	for {
		input, err := line.Read()
		if err != nil {
			return 0.0, err
		}
		if input == "" && hasDefault {
			return defaultAnswer, nil
		}

		answer, err = strconv.Atof(input)
		if err != nil {
			fmt.Fprintf(output, "%s%q: the value has to be a float\r\n",
				QuestionErrPrefix, input)
			continue
		} else {
			return answer, nil
		}
	}
	return
}

// Prints the prompt waiting to get a float number.
func (q *Question) ReadFloat(prompt string) (answer float, err os.Error) {
	return q._baseReadFloat(prompt, 0.0, false)
}

// Prints the prompt waiting to get a float number.
// If input is nil then it returns the answer by default.
func (q *Question) ReadFloatDefault(prompt string, defaultAnswer float) (answer float, err os.Error) {
	return q._baseReadFloat(prompt, defaultAnswer, true)
}

// Prints the prompt waiting to get a string that represents a boolean.
func (q *Question) ReadBool(prompt string, defaultAnswer bool) (answer bool, err os.Error) {
	var options string

	if defaultAnswer {
		options = fmt.Sprintf("%s/%s", strings.ToUpper(q.trueString), q.falseString)
	} else {
		options = fmt.Sprintf("%s/%s", q.trueString, strings.ToUpper(q.falseString))
	}

	line := q.getLine(prompt, options, true)

	for {
		input, err := line.Read()
		if err != nil {
			return false, err
		}
		if input == "" {
			return defaultAnswer, nil
		}

		answer, err = atob(input)
		if err != nil {
			fmt.Fprintf(output, "%s%s: the value does not represent a boolean\r\n",
				QuestionErrPrefix, input)
			continue
		} else {
			return answer, nil
		}
	}
	return
}

// Base to read strings from a set.
func (q *Question) _baseReadChoice(prompt string, a []string, defaultAnswer uint) (answer string, err os.Error) {
	prompt = fmt.Sprintf("%s (%s)", prompt, strings.Join(a, ","))
	line := q.getLine(prompt, a[defaultAnswer], true)

	for {
		answer, err = line.Read()
		if err != nil {
			return "", err
		}
		if answer == "" {
			return a[defaultAnswer], nil
		}

		for _, v := range a {
			if answer == v {
				return answer, nil
			}
		}
	}
	return
}

// Prints the prompt waiting to get an element from array `a`.
// If input is nil then it returns the first element of `a`.
func (q *Question) ReadChoice(prompt string, a []string) (answer string, err os.Error) {
	return q._baseReadChoice(prompt, a, 0)
}

// Prints the prompt waiting to get an element from array `a`.
// If input is nil then it returns the answer by default which is the position
// inner array.
func (q *Question) ReadChoiceDefault(prompt string, a []string, defaultAnswer uint) (answer string, err os.Error) {
	if defaultAnswer >= uint(len(a)) {
		panic(fmt.Sprintf("ReadChoiceDefault: element %d is not in array",
			defaultAnswer))
	}
	return q._baseReadChoice(prompt, a, defaultAnswer)
}


// === Utility
// ===

// Returns the boolean value represented by the string.
// It accepts "y, Y, yes, YES, Yes, n, N, no, NO, No". And values in
// 'strconv.Atob', and 'ExtraBoolString'. Any other value returns an error.
func atob(str string) (value bool, err os.Error) {
	v, err := strconv.Atob(str)
	if err == nil {
		return v, nil
	}

	switch str {
	case "y", "Y", "yes", "YES", "Yes":
		return true, nil
	case "n", "N", "no", "NO", "No":
		return false, nil
	}

	// Check extra characters, if any.
	boolExtra, ok := ExtraBoolString[strings.ToLower(str)]
	if ok {
		return boolExtra, nil
	}

	return false, err // Return error of 'strconv.Atob'
}

