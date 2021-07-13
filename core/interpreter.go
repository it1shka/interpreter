package core

import (
	"errors"
	"fmt"
)

// Callbacks: break, continue, return
type CALLBACK interface {
	callback()
	FormatCallback() string
}

type BREAK_CALLBACK struct{}

func (s BREAK_CALLBACK) callback() {}
func (s BREAK_CALLBACK) FormatCallback() string {
	return "BREAK"
}

type CONTINUE_CALLBACK struct{}

func (s CONTINUE_CALLBACK) callback() {}
func (s CONTINUE_CALLBACK) FormatCallback() string {
	return "CONTINUE"
}

type RETURN_CALLBACK struct {
	value Object
}

func (s RETURN_CALLBACK) callback() {}
func (s RETURN_CALLBACK) FormatCallback() string {
	return fmt.Sprintf("RETURN %s", s.value.ToString())
}

func StmtErr(info string, stmt STATEMENT_NODE) error {
	return errors.New(fmt.Sprintf("%s: %s", info, stmt.GetPosition().Format()))
}

// main Interpreter
type Interpreter struct {
	scope *Scope
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		scope: MakeScope(),
	}
}

func (self *Interpreter) Interpret(program []STATEMENT_NODE) error {
	cb, err := self.EvalStatementList(program)
	if err != nil {
		return err
	}
	if cb != nil {
		return errors.New(fmt.Sprintf("unexpected callback %s", cb.FormatCallback()))
	}
	return nil
}

func (self *Interpreter) EvalStatementList(list []STATEMENT_NODE) (CALLBACK, error) {
	for _, statement := range list {
		switch st := statement.(type) {
		case *BREAK_STATEMENT:
			return BREAK_CALLBACK{}, nil
		case *CONTINUE_STATEMENT:
			return CONTINUE_CALLBACK{}, nil
		case *RETURN_STATEMENT:
			val, err := self.EvalExpression(st.Expression)
			if err != nil {
				return nil, Chain(StmtErr(), err)
			}
		}
	}
	return nil, nil
}

func (self *Interpreter) EvalExpression(expression EXPRESSION_NODE) (Object, error) {

}
