/*
@author: sk
@date: 2024/3/24
*/
package main

import "fmt"

type BaseClass struct { // 基础类信息  类存储方法信息
	Name    string
	Parent  *BaseClass
	Methods map[string]*BaseCall
}

func (b *BaseClass) Call(args []any) any { // 把类当作方法调用就是创建对象
	res := NewBaseInstance(b)
	if init, ok := b.Methods["init"]; ok { // 若存在初始化方法调用初始化方法，参数透传，也就是参数必须与init方法参数一致
		init.BindThis(res).Call(args)
	}
	return res
}

func (b *BaseClass) ArgsSize() int {
	if init, ok := b.Methods["init"]; ok { // 若存在初始化方法，把改类当方法调用时参数数量必须一致
		return init.ArgsSize()
	}
	return 0
}

func (b *BaseClass) GetMethod(name string) *BaseCall {
	if val, ok := b.Methods[name]; ok {
		return val
	}
	if b.Parent != nil {
		return b.Parent.GetMethod(name)
	}
	panic(fmt.Sprintf("no method name %s", name))
}

func NewBaseClass(name string, parent *BaseClass, methods map[string]*BaseCall) *BaseClass {
	return &BaseClass{Name: name, Parent: parent, Methods: methods}
}

type IInstance interface {
	Get(name string) any
	Set(name string, val any)
}

type BaseInstance struct { // 实例存储字段信息
	Class  *BaseClass
	Fields map[string]any
}

func (b *BaseInstance) Set(name string, val any) {
	b.Fields[name] = val // 存在覆盖，不存在创建
}

func (b *BaseInstance) Get(name string) any {
	if val, ok := b.Fields[name]; ok { // 先找字段
		return val
	}
	if method := b.Class.GetMethod(name); method != nil { // 再找方法
		return method.BindThis(b)
	}
	panic(fmt.Sprintf("field %v not Define", name))
}

func (b *BaseInstance) GetSuperMethod(method string) any {
	if b.Class.Parent == nil {
		panic(fmt.Sprintf("class %v no parent", b.Class))
	}
	return b.Class.Parent.GetMethod(method)
}

func NewBaseInstance(class *BaseClass) *BaseInstance {
	return &BaseInstance{Class: class, Fields: make(map[string]any)}
}
