package core

import (
	"fmt"
)

type TokenType int

const (
	// primitive types
	INT_TOKEN TokenType = iota
	FLOAT_TOKEN
	STRING_TOKEN
	BOOL_TOKEN
	NULL_TOKEN

	// terminal tokens
	ID_TOKEN
	KEYWORD_TOKEN
	PUNC_TOKEN
	OP_TOKEN

	// special tokens
	EOF_TOKEN
	ILLEGAL_TOKEN
)

func (self TokenType) Format() string {
	switch self {
	case INT_TOKEN:
		return "INT"
	case FLOAT_TOKEN:
		return "FLOAT"
	case STRING_TOKEN:
		return "STRING"
	case BOOL_TOKEN:
		return "BOOL"
	case ID_TOKEN:
		return "INDENTIFIER"
	case KEYWORD_TOKEN:
		return "KEYWORD"
	case PUNC_TOKEN:
		return "PUNCTUATION"
	case OP_TOKEN:
		return "OPERATOR"
	case EOF_TOKEN:
		return "END OF FILE"
	case ILLEGAL_TOKEN:
		return "ILLEGAL TOKEN"
	case NULL_TOKEN:
		return "NULL"
	default:
		return "UNKNOWN TYPE"
	}
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

func (self *Token) Format() string {
	return fmt.Sprintf("\"%s\" of type %s: (%d, %d)",
		self.Literal, self.Type.Format(), self.Line, self.Column)
}
