package core

import (
	"fmt"
	"strconv"
	"strings"
)

type Object interface {
	Typeof() string
	ToString() string
}

// primitives
type BOOL bool

func (s BOOL) Typeof() string {
	return "BOOL"
}
func (s BOOL) ToString() string {
	if s {
		return "true"
	}
	return "false"
}

type INT int

func (s INT) Typeof() string {
	return "INT"
}
func (s INT) ToString() string {
	return strconv.Itoa(int(s))
}

type FLOAT float64

func (s FLOAT) Typeof() string {
	return "FLOAT"
}
func (s FLOAT) ToString() string {
	return fmt.Sprintf("%f", s)
}

type STRING string

func (s STRING) Typeof() string {
	return "STRING"
}
func (s STRING) ToString() string {
	return string(s)
}

type NULL struct{}

func (s NULL) Typeof() string {
	return "NULL"
}
func (s NULL) ToString() string {
	return "null"
}

type ARRAY []Object

func (s ARRAY) Typeof() string {
	return "ARRAY"
}
func (s ARRAY) ToString() string {
	strs := make([]string, len(s))
	for i, val := range s {
		strs[i] = val.ToString()
	}
	return "[" + strings.Join(strs, ", ") + "]"
}

type FUNCTION struct {
	Args    []string
	Body    []STATEMENT_NODE
	Context *Scope
}

func (s FUNCTION) Typeof() string {
	return "FUNCTION"
}
func (s FUNCTION) ToString() string {
	return "function"
}
