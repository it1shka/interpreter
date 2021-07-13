package core

import (
	"fmt"
)

type Position struct {
	Line, Column int
}

func (self Position) Format() string {
	return fmt.Sprintf("at line %d, at coolumn %d", self.Line, self.Column)
}

type STATEMENT_NODE interface {
	statementNode()
	GetPosition() Position
}

type EXPRESSION_NODE interface {
	expressionNode()
}

// STATEMENTS:

type BREAK_STATEMENT struct {
	Position
}

func (s *BREAK_STATEMENT) statementNode()        {}
func (s *BREAK_STATEMENT) GetPosition() Position { return s.Position }

type CONTINUE_STATEMENT struct {
	Position
}

func (s *CONTINUE_STATEMENT) statementNode()        {}
func (s *CONTINUE_STATEMENT) GetPosition() Position { return s.Position }

type RETURN_STATEMENT struct {
	Expression EXPRESSION_NODE
	Position
}

func (s *RETURN_STATEMENT) statementNode()        {}
func (s *RETURN_STATEMENT) GetPosition() Position { return s.Position }

type LET_STATEMENT struct {
	Identifier string
	Expression EXPRESSION_NODE
	Position
}

func (s *LET_STATEMENT) statementNode()        {}
func (s *LET_STATEMENT) GetPosition() Position { return s.Position }

type FOR_STATEMENT struct {
	Condition EXPRESSION_NODE
	Body      []STATEMENT_NODE
	Position
}

func (s *FOR_STATEMENT) statementNode()        {}
func (s *FOR_STATEMENT) GetPosition() Position { return s.Position }

type IF_STATEMENT struct {
	Condition EXPRESSION_NODE
	Then, Els []STATEMENT_NODE
	Position
}

func (s *IF_STATEMENT) statementNode()        {}
func (s *IF_STATEMENT) GetPosition() Position { return s.Position }

type SAY_STATEMENT struct {
	Expression EXPRESSION_NODE
	Position
}

func (s *SAY_STATEMENT) statementNode()        {}
func (s *SAY_STATEMENT) GetPosition() Position { return s.Position }

type EXPRESSION_STATEMENT struct {
	Expression EXPRESSION_NODE
	Position
}

func (s *EXPRESSION_STATEMENT) statementNode()        {}
func (s *EXPRESSION_STATEMENT) GetPosition() Position { return s.Position }

// EXPRESSIONS:

type BINARY_EXPRESSION struct {
	Operator    string
	Left, Right EXPRESSION_NODE
}

func (s *BINARY_EXPRESSION) expressionNode() {}

type BINARY_ASSIGN_EXPRESSION struct {
	Operator, Left string
	Right          EXPRESSION_NODE
}

func (s *BINARY_ASSIGN_EXPRESSION) expressionNode() {}

type UNARY_OPERATION_EXPRESSION struct {
	Operator   string
	Expression EXPRESSION_NODE
}

func (s *UNARY_OPERATION_EXPRESSION) expressionNode() {}

type VARIABLE_EXPRESSION struct {
	Identifier string
}

func (s *VARIABLE_EXPRESSION) expressionNode() {}

type PRIMITIVE_LITERAL_EXPRESSION struct {
	Value interface{}
}

func (s *PRIMITIVE_LITERAL_EXPRESSION) expressionNode() {}

type NULL_EXPRESSION struct{}

func (s *NULL_EXPRESSION) expressionNode() {}

type FUNCTIONAL_EXPRESSION struct {
	Identifier string
	Args       []string
	Body       []STATEMENT_NODE
}

func (s *FUNCTIONAL_EXPRESSION) expressionNode() {}

type FUNCTION_CALL_EXPRESSION struct {
	Callable EXPRESSION_NODE
	Args     []EXPRESSION_NODE
}

func (s *FUNCTION_CALL_EXPRESSION) expressionNode() {}

type ARRAY_EXPRESSION struct {
	Expressions []EXPRESSION_NODE
}

func (s *ARRAY_EXPRESSION) expressionNode() {}

type INDEX_OPERATOR_EXPRESSION struct {
	Array EXPRESSION_NODE
	Index EXPRESSION_NODE
}

func (s *INDEX_OPERATOR_EXPRESSION) expressionNode() {}

/*
type LAMBDA_EXPRESSION struct {
	Args []string
	Body EXPRESSION_NODE
}

func (s *LAMBDA_EXPRESSION) expressionNode() {}
*/
