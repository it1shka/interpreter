package core

import (
	"errors"
	"fmt"
	"strconv"
)

func Chain(err1, err2 error) error {
	return fmt.Errorf("%s,\n%s", err1, err2)
}

func Includes(operators []string, operator string) bool {
	for _, op := range operators {
		if operator == op {
			return true
		}
	}
	return false
}

type Parser struct {
	stream *Lexer
	pos    Position
}

func NewParser(code string) *Parser {
	return &Parser{
		stream: NewLexer(code),
		pos:    Position{0, 0},
	}
}

func (self *Parser) SetPosition() {
	current := self.stream.Peek()
	self.pos = Position{current.Line, current.Column}
}

func (self *Parser) Err(info string) error {
	return errors.New(fmt.Sprintf("%s: %s", info, self.pos.Format()))
}

// main functions
func (self *Parser) ParseProgram() ([]STATEMENT_NODE, error) {
	if self.stream.Eof() {
		return []STATEMENT_NODE{}, nil
	}
	statement, err := self.ParseStatement()
	if err != nil {
		return nil, err
	}

	if self.stream.NextIf(";") {
		nextStatements, err := self.ParseProgram()
		if err != nil {
			return nil, err
		}
		return append([]STATEMENT_NODE{statement}, nextStatements...), nil
	}
	if self.stream.Peek().Type != EOF_TOKEN {
		return nil, errors.New(fmt.Sprintf("EOF or \";\" expected, found %s",
			self.stream.Next().Format()))
	}
	return []STATEMENT_NODE{statement}, nil
}

func (self *Parser) ParseStatement() (STATEMENT_NODE, error) {
	self.SetPosition()
	var node STATEMENT_NODE

	if self.stream.NextIf("break") {
		node = &BREAK_STATEMENT{self.pos}
	} else if self.stream.NextIf("continue") {
		node = &CONTINUE_STATEMENT{self.pos}
	} else if self.stream.NextIf("return") {
		expression, err := self.ParseExpression()
		if err != nil {
			return nil, Chain(self.Err("while parsing RETURN statement expression"), err)
		}
		node = &RETURN_STATEMENT{expression, self.pos}
	} else if self.stream.NextIf("let") {
		identifier := self.stream.Next()
		if identifier.Type != ID_TOKEN {
			return nil, self.Err(fmt.Sprintf(
				"while parsing LET statement: expected IDENTIFIER, found %s",
				identifier.Format()))
		}
		var initial EXPRESSION_NODE
		if self.stream.NextIf("=") {
			expression, err := self.ParseExpression()
			if err != nil {
				return nil, Chain(self.Err("while parsing LET statement"), err)
			}
			initial = expression
		}
		node = &LET_STATEMENT{identifier.Literal, initial, self.pos}
	} else if self.stream.NextIf("for") {
		expression, err := self.ParseExpression()
		if err != nil {
			return nil, Chain(self.Err("while parsing FOR statement condition"), err)
		}
		body, err := self.ParseStatementList()
		if err != nil {
			return nil, Chain(self.Err("while parsing FOR statement body"), err)
		}
		node = &FOR_STATEMENT{expression, body, self.pos}
	} else if self.stream.NextIf("if") {
		condition, err := self.ParseExpression()
		if err != nil {
			return nil, Chain(self.Err("while parsing IF statement condition"), err)
		}
		then, err := self.ParseStatementList()
		if err != nil {
			return nil, Chain(self.Err("while parsing IF statement THEN brahcn"), err)
		}
		var els []STATEMENT_NODE
		if self.stream.NextIf("else") {
			statements, err := self.ParseStatementList()
			if err != nil {
				return nil, Chain(self.Err("while parsing IF statement ELSE branch"), err)
			}
			els = statements
		}
		node = &IF_STATEMENT{condition, then, els, self.pos}
	} else if self.stream.NextIf("say") {
		expression, err := self.ParseExpression()
		if err != nil {
			return nil, Chain(self.Err("while parsing SAY statement"), err)
		}
		node = &SAY_STATEMENT{expression, self.pos}
	} else {
		expression, err := self.ParseExpression()
		if err != nil {
			return nil, Chain(self.Err("while parsing EXPRESSION statement"), err)
		}
		node = &EXPRESSION_STATEMENT{expression, self.pos}
	}

	return node, nil
}

func (self *Parser) ParseStatementList() ([]STATEMENT_NODE, error) {
	next_token := self.stream.Next()
	if next_token.Literal != "{" {
		return nil, errors.New(fmt.Sprintf("expected \"{\" while parsing statement list, found %s",
			next_token.Format()))
	}

	return self.ParseStatementList_()
}

func (self *Parser) ParseStatementList_() ([]STATEMENT_NODE, error) {
	if self.stream.NextIf("}") {
		return []STATEMENT_NODE{}, nil
	}
	statement, err := self.ParseStatement()
	if err != nil {
		return nil, err
	}
	if self.stream.NextIf(";") {
		next_statements, err := self.ParseStatementList_()
		if err != nil {
			return nil, err
		}
		return append([]STATEMENT_NODE{statement}, next_statements...), nil
	}
	next_tok := self.stream.Next()
	if next_tok.Literal != "}" {
		return nil, errors.New(fmt.Sprintf("closing \"}\" or \";\" expected, found %s", next_tok.Format()))
	}
	return []STATEMENT_NODE{statement}, nil
}

// hope nobody will see it
func (self *Parser) ParseExpression() (EXPRESSION_NODE, error) {
	return self.ParseBinaryAssignExpression(
		[]string{"+=", "-=", "*=", "/=", "%=", "&=", "|=", "="},
		func() (EXPRESSION_NODE, error) {
			return self.ParseBinaryExpression(
				[]string{"|"},
				func() (EXPRESSION_NODE, error) {
					return self.ParseBinaryExpression(
						[]string{"&"},
						func() (EXPRESSION_NODE, error) {
							return self.ParseBinaryExpression(
								[]string{"==", "!="},
								func() (EXPRESSION_NODE, error) {
									return self.ParseBinaryExpression(
										[]string{">", "<", ">=", "<="},
										func() (EXPRESSION_NODE, error) {
											return self.ParseBinaryExpression(
												[]string{"+", "-"},
												func() (EXPRESSION_NODE, error) {
													return self.ParseBinaryExpression(
														[]string{"*", "/", "%"},
														func() (EXPRESSION_NODE, error) {
															return self.ParsePrimaryExpression()
														},
													)
												},
											)
										},
									)
								},
							)
						},
					)
				},
			)
		},
	)
}

func (self *Parser) ParseBinaryAssignExpression(
	operators []string,
	parser func() (EXPRESSION_NODE, error),
) (EXPRESSION_NODE, error) {
	left, err := parser()
	if err != nil {
		return nil, err
	}
	for Includes(operators, self.stream.Peek().Literal) {
		id, ok := left.(*VARIABLE_EXPRESSION)
		if !ok {
			return nil, errors.New("expected identifier in ASSIGN expression")
		}
		operator := self.stream.Next().Literal
		right, err := self.ParseBinaryAssignExpression(operators, parser)
		if err != nil {
			return nil, err
		}
		left = &BINARY_ASSIGN_EXPRESSION{operator, id.Identifier, right}
	}
	return left, nil
}

func (self *Parser) ParseBinaryExpression(
	operators []string,
	parser func() (EXPRESSION_NODE, error),
) (EXPRESSION_NODE, error) {
	left, err := parser()
	if err != nil {
		return nil, err
	}
	for Includes(operators, self.stream.Peek().Literal) {
		operator := self.stream.Next().Literal
		right, err := parser()
		if err != nil {
			return nil, err
		}
		left = &BINARY_EXPRESSION{operator, left, right}
	}
	return left, nil
}

func (self *Parser) ParsePrimaryExpression() (EXPRESSION_NODE, error) {
	return self.ParseUnaryOperatorExpression()
}

func (self *Parser) ParseUnaryOperatorExpression() (EXPRESSION_NODE, error) {
	if Includes([]string{"!", "-", "+"}, self.stream.Peek().Literal) {
		operator := self.stream.Next().Literal
		expression, err := self.ParseUnaryOperatorExpression()
		if err != nil {
			return nil, err
		}
		return &UNARY_OPERATION_EXPRESSION{operator, expression}, nil
	}
	expression, err := self.ParseValueExpression()
	if err != nil {
		return nil, err
	}
	return self.ParsePostExpressionOperator(expression)
}

func (self *Parser) ParsePostExpressionOperator(
	prev EXPRESSION_NODE,
) (EXPRESSION_NODE, error) {
	if self.stream.NextIf("[") {
		index, err := self.ParseExpression()
		if err != nil {
			return nil, err
		}
		next_tok := self.stream.Next()
		if next_tok.Literal != "]" {
			return nil, errors.New(fmt.Sprintf("expected closing \"]\", found %s",
				next_tok.Format()))
		}
		expression := &INDEX_OPERATOR_EXPRESSION{prev, index}
		return self.ParsePostExpressionOperator(expression)
	}

	if self.stream.NextIf("(") {
		arglist, err := self.ParseExpressionList(")")
		if err != nil {
			return nil, err
		}
		expression := &FUNCTION_CALL_EXPRESSION{prev, arglist}
		return self.ParsePostExpressionOperator(expression)
	}

	return prev, nil
}

func (self *Parser) ParseExpressionList(end string) ([]EXPRESSION_NODE, error) {
	if self.stream.NextIf(end) {
		return []EXPRESSION_NODE{}, nil
	}
	expression, err := self.ParseExpression()
	if err != nil {
		return nil, err
	}

	if self.stream.NextIf(",") {
		next_expressions, err := self.ParseExpressionList(end)
		if err != nil {
			return nil, err
		}
		return append([]EXPRESSION_NODE{expression}, next_expressions...), nil
	}

	next_tok := self.stream.Next()
	if next_tok.Literal != end {
		return nil, errors.New(fmt.Sprintf("expected closing \"%s\" or \",\", found %s",
			end, next_tok.Format()))
	}
	return []EXPRESSION_NODE{expression}, nil
}

func (self *Parser) ParseFunctionArgsList(end string) ([]string, error) {
	if self.stream.Peek().Literal == end {
		return []string{}, nil
	}
	next_tok := self.stream.Next()
	if next_tok.Type != ID_TOKEN {
		return nil, errors.New(fmt.Sprintf("identifier expected, found %s",
			next_tok.Format()))
	}
	id := next_tok.Literal

	if self.stream.NextIf(",") {
		next_ids, err := self.ParseFunctionArgsList(end)
		if err != nil {
			return nil, err
		}
		return append([]string{id}, next_ids...), nil
	}

	if self.stream.Peek().Literal != end {
		return nil, errors.New(fmt.Sprintf("expected closing \"%s\" or \",\", found %s",
			end, self.stream.Peek().Format()))
	}
	return []string{id}, nil
}

func (self *Parser) ParseValueExpression() (EXPRESSION_NODE, error) {
	if self.stream.NextIf("(") {
		expression, err := self.ParseExpression()
		if err != nil {
			return nil, err
		}
		next_tok := self.stream.Next()
		if next_tok.Literal != ")" {
			return nil, errors.New(fmt.Sprintf("expected closing \")\", found %s",
				next_tok.Format()))
		}
		return expression, nil
	}

	if self.stream.NextIf("[") {
		exprlist, err := self.ParseExpressionList("]")
		if err != nil {
			return nil, err
		}
		return &ARRAY_EXPRESSION{exprlist}, nil
	}

	if self.stream.NextIf("fn") {
		name := ""
		if self.stream.Peek().Type == ID_TOKEN {
			name = self.stream.Next().Literal
		}
		next_tok := self.stream.Next()
		if next_tok.Literal != ":" {
			return nil, errors.New(fmt.Sprintf("\":\" expected, found %s",
				next_tok.Format()))
		}
		args, err := self.ParseFunctionArgsList("{")
		if err != nil {
			return nil, err
		}
		body, err := self.ParseStatementList()
		if err != nil {
			return nil, err
		}

		return &FUNCTIONAL_EXPRESSION{name, args, body}, nil
	}

	if self.stream.NextIf("lambda") {
		args, err := self.ParseFunctionArgsList(":")
		next_tok := self.stream.Next()
		if next_tok.Literal != ":" {
			return nil, errors.New(fmt.Sprintf("\":\" expected, found %s",
				next_tok.Format()))
		}
		if err != nil {
			return nil, err
		}
		body, err := self.ParseExpression()
		if err != nil {
			return nil, err
		}
		return &LAMBDA_EXPRESSION{args, body}, nil
	}

	switch next_token := self.stream.Next(); next_token.Type {
	case ID_TOKEN:
		return &VARIABLE_EXPRESSION{next_token.Literal}, nil
	case INT_TOKEN:
		val, err := strconv.Atoi(next_token.Literal)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to parse %s to INT",
				next_token.Format()))
		}
		return &PRIMITIVE_LITERAL_EXPRESSION{val}, nil
	case FLOAT_TOKEN:
		val, err := strconv.ParseFloat(next_token.Literal, 64)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to parse %s to FLOAT",
				next_token.Format()))
		}
		return &PRIMITIVE_LITERAL_EXPRESSION{val}, nil
	case STRING_TOKEN:
		return &PRIMITIVE_LITERAL_EXPRESSION{next_token.Literal}, nil
	case BOOL_TOKEN:
		if next_token.Literal == "true" {
			return &PRIMITIVE_LITERAL_EXPRESSION{true}, nil
		} else {
			return &PRIMITIVE_LITERAL_EXPRESSION{false}, nil
		}
	case NULL_TOKEN:
		return &NULL_EXPRESSION{}, nil
	default:
		return nil, errors.New(fmt.Sprintf("unexpected %s", next_token.Format()))
	}
}
