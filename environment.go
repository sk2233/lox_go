/*
@author: sk
@date: 2024/3/19
*/
package main

import "fmt"

type Environment struct {
	Parent *Environment
	Values map[string]any
}

func NewEnvironmentWithParent(parent *Environment) *Environment {
	return &Environment{Parent: parent, Values: make(map[string]any)}
}

func NewEnvironment() *Environment {
	return &Environment{Values: make(map[string]any)}
}

func (e *Environment) Define(key string, val any) {
	if _, ok := e.Values[key]; ok {
		panic(fmt.Sprintf("repeat Define %v", key))
	}
	e.Values[key] = val
}

func (e *Environment) GetOrNil(key string) any {
	if val, ok := e.Values[key]; ok {
		return val
	}
	if e.Parent != nil {
		return e.Parent.GetOrNil(key)
	}
	return nil
}

func (e *Environment) Get(key string) any {
	if val, ok := e.Values[key]; ok {
		return val
	}
	if e.Parent != nil {
		return e.Parent.Get(key)
	}
	panic(fmt.Sprintf("no Define var %s", key))
}

func (e *Environment) Assign(key string, val any) {
	if _, ok := e.Values[key]; ok {
		e.Values[key] = val
		return
	}
	if e.Parent != nil {
		e.Parent.Assign(key, val)
		return
	}
	panic(fmt.Sprintf("Assign no Define %v", key))
}
