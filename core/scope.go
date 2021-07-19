package core

import (
	"errors"
	"fmt"
)

type Scope struct {
	current map[string]Object
	prev    *Scope
}

func MakeScope() *Scope {
	return &Scope{
		current: make(map[string]Object),
		prev:    nil,
	}
}

func (self *Scope) Init(identifier string, value Object) error {
	if _, ok := self.current[identifier]; ok {
		return errors.New(fmt.Sprintf("trying to initialize \"%s\" for the second time", identifier))
	}
	self.current[identifier] = value
	return nil
}

func (self *Scope) Set(identifier string, value Object) (Object, error) {
	if self.current[identifier] != nil {
		fmt.Println("Not Nil :)")
		fmt.Printf("last val: %s", self.current[identifier].ToString())
		self.current[identifier] = value
		return value, nil
	}
	if self.prev != nil {
		return self.prev.Set(identifier, value)
	}
	return nil, errors.New(fmt.Sprintf("trying to set uninitialized variable \"%s\"", identifier))
}

func (self *Scope) Get(identifier string) (Object, error) {
	if value, ok := self.current[identifier]; ok {
		return value, nil
	}
	if self.prev != nil {
		return self.prev.Get(identifier)
	}
	return nil, errors.New(fmt.Sprintf("trying to get uninitialized variable \"%s\"", identifier))
}

func (self *Scope) NewChild() *Scope {
	return &Scope{
		current: make(map[string]Object),
		prev:    self,
	}
}
