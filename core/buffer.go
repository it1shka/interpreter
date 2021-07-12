package core

type LexerBuffer struct {
	input                  []rune
	line, column, position int
}

func NewLexerBuffer(input string) *LexerBuffer {
	return &LexerBuffer{
		input:    []rune(input),
		line:     1,
		column:   1,
		position: 0,
	}
}

func (self *LexerBuffer) Peek() rune {
	if self.position >= len(self.input) {
		return rune(0)
	}

	return self.input[self.position]
}

func (self *LexerBuffer) Next() rune {
	if self.position >= len(self.input) {
		return rune(0)
	}

	char := self.input[self.position]
	self.position++
	if char == '\n' {
		self.line++
		self.column = 1
	} else {
		self.column++
	}
	return char
}

func (self *LexerBuffer) NextIf(expected rune) bool {
	if self.Peek() == expected {
		self.Next()
		return true
	}
	return false
}

func (self *LexerBuffer) Eof() bool {
	return self.Peek() == rune(0)
}

func (self *LexerBuffer) Pos() (int, int) {
	return self.line, self.column
}
