package core

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type ObjectType int

const (
	INTEGER_TYPE ObjectType = iota
	FLOATING_TYPE
	BOOLEAN_TYPE
	STRING_TYPE
	NULL_TYPE
	ARRAY_TYPE
	FUNCTION_TYPE
)

func Typeof(obj Object) string {
	return obj.Typeof().FormatObjectType()
}

func (s ObjectType) FormatObjectType() string {
	switch s {
	case INTEGER_TYPE:
		return "INTEGER"
	case FLOATING_TYPE:
		return "FLOATING"
	case BOOLEAN_TYPE:
		return "BOOLEAN"
	case STRING_TYPE:
		return "STRING"
	case NULL_TYPE:
		return "NULL"
	case ARRAY_TYPE:
		return "ARRAY"
	case FUNCTION_TYPE:
		return "FUNCTION"
	default:
		return "UNKNOWN"
	}
}

type Object interface {
	Typeof() ObjectType

	ToString() string
	ToInteger() (INT, error)
	ToFloating() (FLOAT, error)
	ToBoolean() (BOOL, error)
}

// primitives
type BOOL bool

func (s BOOL) Typeof() ObjectType {
	return BOOLEAN_TYPE
}
func (s BOOL) ToString() string {
	if s {
		return "true"
	}
	return "false"
}
func (s BOOL) ToBoolean() (BOOL, error) {
	return s, nil
}
func (s BOOL) ToInteger() (INT, error) {
	if s {
		return 1, nil
	} else {
		return 0, nil
	}
}
func (s BOOL) ToFloating() (FLOAT, error) {
	if s {
		return 1.0, nil
	} else {
		return 0, nil
	}
}

type INT int

func (s INT) Typeof() ObjectType {
	return INTEGER_TYPE
}
func (s INT) ToString() string {
	return strconv.Itoa(int(s))
}
func (s INT) ToBoolean() (BOOL, error) {
	return BOOL(s != 0), nil
}
func (s INT) ToInteger() (INT, error) {
	return s, nil
}
func (s INT) ToFloating() (FLOAT, error) {
	return FLOAT(s), nil
}

type FLOAT float64

func (s FLOAT) Typeof() ObjectType {
	return FLOATING_TYPE
}
func (s FLOAT) ToString() string {
	return fmt.Sprintf("%f", s)
}
func (s FLOAT) ToBoolean() (BOOL, error) {
	return BOOL(s != 0), nil
}
func (s FLOAT) ToInteger() (INT, error) {
	return INT(s), nil
}
func (s FLOAT) ToFloating() (FLOAT, error) {
	return s, nil
}

type STRING string

func (s STRING) Typeof() ObjectType {
	return STRING_TYPE
}
func (s STRING) ToString() string {
	return string(s)
}
func (s STRING) ToBoolean() (BOOL, error) {
	return false, errors.New("invalid conversion: STRING to BOOLEAN")
}
func (s STRING) ToInteger() (INT, error) {
	val, ok := strconv.Atoi(string(s))
	if ok != nil {
		return 0, errors.New(fmt.Sprintf("cannot convert STRING %s to INT", s))
	}
	return INT(val), nil
}
func (s STRING) ToFloating() (FLOAT, error) {
	val, ok := strconv.ParseFloat(string(s), 64)
	if ok != nil {
		return 0, errors.New(fmt.Sprintf("cannot convert STRING %s to FLOAT", s))
	}
	return FLOAT(val), nil
}

type NULL struct{}

func (s NULL) Typeof() ObjectType {
	return NULL_TYPE
}
func (s NULL) ToString() string {
	return "null"
}
func (s NULL) ToBoolean() (BOOL, error) {
	return false, nil
}
func (s NULL) ToInteger() (INT, error) {
	return 0, nil
}
func (s NULL) ToFloating() (FLOAT, error) {
	return 0.0, nil
}

type ARRAY []Object

func (s ARRAY) Typeof() ObjectType {
	return ARRAY_TYPE
}
func (s ARRAY) ToString() string {
	strs := make([]string, len(s))
	for i, val := range s {
		strs[i] = val.ToString()
	}
	return "[" + strings.Join(strs, ", ") + "]"
}
func (s ARRAY) ToBoolean() (BOOL, error) {
	return false, errors.New("invalid conversion: ARRAY to BOOLEAN")
}
func (s ARRAY) ToInteger() (INT, error) {
	return 0, errors.New("invalid conversion: ARRAY to INT")
}
func (s ARRAY) ToFloating() (FLOAT, error) {
	return 0.0, errors.New("invalid conversion: ARRAY to FLOAT")
}

type FUNCTION struct {
	Args    []string
	Body    []STATEMENT_NODE
	Context *Scope
}

func (s FUNCTION) Typeof() ObjectType {
	return FUNCTION_TYPE
}
func (s FUNCTION) ToString() string {
	return "function"
}
func (s FUNCTION) ToBoolean() (BOOL, error) {
	return false, errors.New("invalid conversion: FUNCTION to BOOLEAN")
}
func (s FUNCTION) ToInteger() (INT, error) {
	return 0, errors.New("invalid conversion: FUNCTION to INT")
}
func (s FUNCTION) ToFloating() (FLOAT, error) {
	return 0.0, errors.New("invalid conversion: FUNCTION to FLOAT")
}
