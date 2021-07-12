package core

import (
	"strings"
	"unicode"
)

type Lexer struct {
	buffer       *LexerBuffer
	current      *Token
	char         rune
	line, column int
}

// constructor
func NewLexer(code string) *Lexer {
	return &Lexer{
		buffer:  NewLexerBuffer(code),
		current: nil,
		line:    1,
		column:  1,
	}
}

// help functions
func (self *Lexer) NewToken(tokenType TokenType, literal string) *Token {
	return &Token{
		Type:    tokenType,
		Literal: literal,
		Line:    self.line,
		Column:  self.column,
	}
}

func (self *Lexer) ReadWhile(predicate func(rune) bool) string {
	var builder strings.Builder
	for !self.buffer.Eof() && predicate(self.buffer.Peek()) {
		builder.WriteRune(self.buffer.Next())
	}
	return builder.String()
}

func (self *Lexer) ReadString() string {
	str := self.ReadWhile(func(r rune) bool {
		return r != self.char
	})
	self.buffer.Next()
	return str
}

func (self *Lexer) ReadNumber() *Token {
	part := func() string {
		return self.ReadWhile(func(r rune) bool {
			return unicode.IsDigit(r)
		})
	}

	var number strings.Builder
	number.WriteRune(self.char)
	number.WriteString(part())
	if self.buffer.NextIf('.') {
		number.WriteRune('.')
		number.WriteString(part())
		return self.NewToken(FLOAT_TOKEN, number.String())
	}
	return self.NewToken(INT_TOKEN, number.String())
}

func (self *Lexer) ReadWord() *Token {
	word := string(self.char) + self.ReadWhile(func(r rune) bool {
		return unicode.IsDigit(r) || unicode.IsLetter(r) || r == '_' || r == '$'
	})

	// checking if word is a keyword(boolean, null)
	switch word {
	case "true", "false":
		return self.NewToken(BOOL_TOKEN, word)
	case "null":
		return self.NewToken(NULL_TOKEN, word)
	case "let", "break", "continue", "return",
		"for", "if", "else", "fn", "lambda", "say":
		return self.NewToken(KEYWORD_TOKEN, word)
	default:
		return self.NewToken(ID_TOKEN, word)
	}
}

// main function
func (self *Lexer) ReadToken() *Token {
	self.ReadWhile(func(r rune) bool {
		return unicode.IsSpace(r)
	})

	self.line, self.column = self.buffer.Pos()
	self.char = self.buffer.Next()

	switch self.char {
	//end of file
	case rune(0):
		return self.NewToken(EOF_TOKEN, "")
	// one line comment
	case '#':
		self.ReadWhile(func(r rune) bool {
			return r != '\n'
		})
		return self.ReadToken()
	// punctuation
	case '(', ')', '{', '}', ';', ',', ':', '[', ']':
		return self.NewToken(PUNC_TOKEN, string(self.char))
	//operators
	case '+', '-', '*', '/', '%', '=', '!', '>', '<', '&', '|':
		if self.buffer.NextIf('=') {
			return self.NewToken(OP_TOKEN, string([]rune{self.char, '='}))
		}
		return self.NewToken(OP_TOKEN, string(self.char))
	// string literal
	case '"', '\'':
		return self.NewToken(STRING_TOKEN, self.ReadString())
	default:
		//keyword or identifier
		if unicode.IsLetter(self.char) || self.char == '_' || self.char == '$' {
			return self.ReadWord()
		}
		//number
		if unicode.IsDigit(self.char) {
			return self.ReadNumber()
		}
		// unknown, illegal token
		illegalLiteral := string(self.char) + self.ReadWhile(func(r rune) bool {
			return !unicode.IsSpace(r)
		})
		return self.NewToken(ILLEGAL_TOKEN, illegalLiteral)
	}
}

// public interface
func (self *Lexer) Peek() *Token {
	if self.current == nil {
		self.current = self.ReadToken()
	}
	return self.current
}

func (self *Lexer) Next() *Token {
	current := self.Peek()
	self.current = nil
	return current
}

func (self *Lexer) Eof() bool {
	return self.Peek().Type == EOF_TOKEN
}

func (self *Lexer) NextIf(expected string) bool {
	if self.Peek().Literal == expected {
		self.Next()
		return true
	}
	return false
}
