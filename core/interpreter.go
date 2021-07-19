package core

import (
	"errors"
	"fmt"
)

var UNKNOWN_ERROR error = errors.New("UNKNOWN ERROR: something went wrong :(")

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

// main Interpreter
type Interpreter struct {
	scope *Scope
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		scope: MakeScope(),
	}
}

func (self *Interpreter) EnterNewScope() {
	self.scope = self.scope.NewChild()
}

func (self *Interpreter) LeaveScope() {
	self.scope = self.scope.prev
}

func (self *Interpreter) Interpret(program []STATEMENT_NODE) error {
	cb, err := self.EvalStatementList(program)
	if err != nil {
		return err
	}
	if cb != nil {
		return errors.New(fmt.Sprintf("unexpected callback %s in program", cb.FormatCallback()))
	}
	return nil
}

func (self *Interpreter) EvalStatementList(list []STATEMENT_NODE) (CALLBACK, error) {
	for _, statement := range list {

		StmtErr := func(info string) error {
			return errors.New(fmt.Sprintf("%s: %s", info, statement.GetPosition().Format()))
		}

		switch st := statement.(type) {
		case *BREAK_STATEMENT:
			return BREAK_CALLBACK{}, nil
		case *CONTINUE_STATEMENT:
			return CONTINUE_CALLBACK{}, nil
		case *RETURN_STATEMENT:
			val, err := self.EvalExpression(st.Expression)
			if err != nil {
				return nil, Chain(StmtErr("while evaluating RETURN statement"), err)
			}
			return RETURN_CALLBACK{val}, nil
		case *LET_STATEMENT:
			var initial Object
			if st.Expression != nil {
				obj, err := self.EvalExpression(st.Expression)
				if err != nil {
					return nil, Chain(StmtErr("while evaluating LET statement"), err)
				}
				initial = obj
			} else {
				initial = NULL{}
			}
			if err := self.scope.Init(st.Identifier, initial); err != nil {
				return nil, Chain(StmtErr("while evaluating LET statement"), err)
			}
		case *FOR_STATEMENT:
			for {
				val, err := self.EvalExpression(st.Condition)
				if err != nil {
					return nil, Chain(StmtErr("while evaluating FOR condition"), err)
				}
				ok, err := val.ToBoolean()
				if err != nil {
					return nil, Chain(StmtErr("white evaluating FOR condition"), err)
				}
				if !ok {
					break
				}

				self.EnterNewScope()
				cb, err := self.EvalStatementList(st.Body)
				self.LeaveScope()

				if err != nil {
					return nil, Chain(StmtErr("while evaluating FOR body"), err)
				}
				_, ok_br := cb.(BREAK_CALLBACK)
				if ok_br {
					break
				}
				_, ok_rt := cb.(RETURN_CALLBACK)
				if ok_rt {
					return cb, nil
				}
			}
		case *IF_STATEMENT:
			val, err := self.EvalExpression(st.Condition)
			if err != nil {
				return nil, Chain(StmtErr("while evaluating IF condition"), err)
			}
			ok, err := val.ToBoolean()
			if err != nil {
				return nil, Chain(StmtErr("white evaluating IF condition"), err)
			}

			if ok {
				self.EnterNewScope()
				cb, err := self.EvalStatementList(st.Then)
				self.LeaveScope()
				if err != nil {
					return nil, Chain(StmtErr("white evaluating IF then branch"), err)
				}
				if cb != nil {
					return cb, nil
				}
			} else if st.Els != nil {
				self.EnterNewScope()
				cb, err := self.EvalStatementList(st.Els)
				self.LeaveScope()
				if err != nil {
					return nil, Chain(StmtErr("white evaluating IF else branch"), err)
				}
				if cb != nil {
					return cb, nil
				}
			}
		case *SAY_STATEMENT:
			obj, err := self.EvalExpression(st.Expression)
			if err != nil {
				return nil, Chain(StmtErr("while evaluating SAY statement"), err)
			}
			fmt.Println(obj.ToString())
		case *EXPRESSION_STATEMENT:
			_, err := self.EvalExpression(st.Expression)
			if err != nil {
				return nil, Chain(StmtErr("while evaluating EXPRESSION statement"), err)
			}
		}
	}
	return nil, nil
}

func (self *Interpreter) EvalExpression(expression EXPRESSION_NODE) (Object, error) {
	switch ex := expression.(type) {
	case *NULL_EXPRESSION:
		return NULL{}, nil
	case *VARIABLE_EXPRESSION:
		return self.scope.Get(ex.Identifier)
	case *PRIMITIVE_LITERAL_EXPRESSION:
		switch t := ex.Value.(type) {
		case int:
			return INT(t), nil
		case float64:
			return FLOAT(t), nil
		case string:
			return STRING(t), nil
		case bool:
			return BOOL(t), nil
		default:
			return nil, UNKNOWN_ERROR
		}
	case *INDEX_OPERATOR_EXPRESSION:
		arrv, err := self.EvalExpression(ex.Array)
		if err != nil {
			return nil, err
		}
		arr, ok := arrv.(ARRAY)
		if !ok {
			return nil, errors.New("cannot get index of non-array object")
		}
		indv, err := self.EvalExpression(ex.Index)
		if err != nil {
			return nil, err
		}
		ind, ok := indv.(INT)
		if !ok {
			return nil, errors.New("index must be typeof INT")
		}

		if int(ind) >= len(arr) {
			return NULL{}, nil
		} else {
			return arr[int(ind)], nil
		}
	case *ARRAY_EXPRESSION:
		objarr := make([]Object, len(ex.Expressions))
		for i, v := range ex.Expressions {
			obj, err := self.EvalExpression(v)
			if err != nil {
				return nil, err
			}
			objarr[i] = obj
		}
		return ARRAY(objarr), nil
	case *FUNCTION_CALL_EXPRESSION:
		fv, err := self.EvalExpression(ex.Callable)
		if err != nil {
			return nil, err
		}
		fn, ok := fv.(FUNCTION)
		if !ok {
			return nil, errors.New("cannot call non-callable object")
		}

		if len(ex.Args) != len(fn.Args) {
			return nil, errors.New(fmt.Sprintf("expected %d args, found %d args",
				len(fn.Args), len(ex.Args)))
		}

		before_call := self.scope
		self.scope = fn.Context.NewChild()
		for i := 0; i < len(ex.Args); i++ {
			obj, err := self.EvalExpression(ex.Args[i])
			if err != nil {
				return nil, err
			}
			self.scope.Init(fn.Args[i], obj)
		}
		cb, err := self.EvalStatementList(fn.Body)
		self.scope = before_call
		if err != nil {
			return nil, err
		}
		switch cbv := cb.(type) {
		case BREAK_CALLBACK:
			return nil, errors.New("unexpected BREAK callback in function call")
		case CONTINUE_CALLBACK:
			return nil, errors.New("unexpected CONTINUE callback in function call")
		case RETURN_CALLBACK:
			return cbv.value, nil
		}
		return NULL{}, nil
	case *FUNCTIONAL_EXPRESSION:
		val := FUNCTION{ex.Args, ex.Body, self.scope}
		if len(ex.Identifier) > 0 {
			self.scope.Init(ex.Identifier, val)
		}
		return val, nil
	case *UNARY_OPERATION_EXPRESSION:
		obj, err := self.EvalExpression(ex.Expression)
		if err != nil {
			return nil, err
		}
		return ApplyUnaryOperator(ex.Operator, obj)
	case *BINARY_EXPRESSION:
		left, err := self.EvalExpression(ex.Left)
		if err != nil {
			return nil, err
		}
		right, err := self.EvalExpression(ex.Right)
		if err != nil {
			return nil, err
		}
		return ApplyBinaryOperator(ex.Operator, left, right)
	case *BINARY_ASSIGN_EXPRESSION:
		right, err := self.EvalExpression(ex.Right)
		if err != nil {
			return nil, err
		}
		if len(ex.Operator) == 2 {
			mainop := string([]rune(ex.Operator)[1])
			cur_val, err := self.scope.Get(ex.Left)
			if err != nil {
				return nil, err
			}
			obj, err := ApplyBinaryOperator(mainop, cur_val, right)
			if err != nil {
				return nil, err
			}
			self.scope.Set(ex.Left, obj)
			return obj, nil
		} else { // =
			self.scope.Set(ex.Left, right)
			return right, nil
		}
	default:
		return nil, UNKNOWN_ERROR
	}
}

func ApplyBinaryOperator(operator string, left, right Object) (Object, error) {
	switch operator {
	case "+":
		if left.Typeof() == INTEGER_TYPE && right.Typeof() == INTEGER_TYPE {
			return INT(left.(INT) + right.(INT)), nil
		}
		if left.Typeof() == INTEGER_TYPE && right.Typeof() == FLOATING_TYPE {
			return FLOAT(FLOAT(left.(INT)) + right.(FLOAT)), nil
		}
		if left.Typeof() == FLOATING_TYPE && right.Typeof() == INTEGER_TYPE {
			return FLOAT(left.(FLOAT) + FLOAT(right.(INT))), nil
		}

		return nil, UNKNOWN_ERROR
	default:
		return nil, UNKNOWN_ERROR
	}
}

func ApplyUnaryOperator(operator string, obj Object) (Object, error) {
	switch operator {
	case "-":
		switch t := obj.(type) {
		case INT:
			return -t, nil
		case FLOAT:
			return -t, nil
		default:
			return nil, errors.New(fmt.Sprintf("cannot apply \"-\" operator for %s",
				Typeof(obj)))
		}
	case "!":
		switch t := obj.(type) {
		case BOOL:
			return !t, nil
		default:
			return nil, errors.New(fmt.Sprintf("cannot apply \"!\" operator for %s",
				Typeof(obj)))
		}
	default:
		return nil, UNKNOWN_ERROR
	}
}
