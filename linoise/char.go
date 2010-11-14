package linoise


// Characters
var (
	_CR    = []byte{13}     // Carriage return (in hexadecimal: '\x0D')
	_CR_LF = []byte{13, 10} // CR+LF is used for a new line in raw mode - (\r\n)
	ctrlC  = []int("^C")
	ctrlD  = []int("^D")
)

// ASCII codes
const (
	_ESC       = 27 // Escape: Ctrl-[ (033 in octal)
	_L_BRACKET = 91 // Left square bracket: [
)

// ANSI terminal escape controls
var (
	// === Cursor control
	cursorUp       = []byte("\033[A") // Up
	cursorDown     = []byte("\033[B") // Down
	cursorForward  = []byte("\033[C") // Forward
	cursorBackward = []byte("\033[D") // Backward

	toNextLine = []byte("\033[E") // To next line
	toPreviousLine = []byte("\033[F") // To previous line

	// === Erasing Text
	//delScreen = []byte("\033[2J") // Erase the screen

	delRight         = []byte("\033[0K")       // Erase to right
	delLine_CR       = []byte("\033[2K\r")     // Erase line; carriage return
	delLine_cursorUp = []byte("\033[2K\033[A") // Erase line; cursor up

	//delChar      = []byte("\033[1X") // Erase character
	delChar      = []byte("\033[P")  // Delete character, from current position
	delBackspace = []byte("\033[D\033[P")

	// === Misc.
	//insertChar  = []byte("\033[@")   // Insert CHaracter
	//setLineWrap = []byte("\033[?7h") // Enable Line Wrap
)

