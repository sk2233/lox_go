/*
@author: sk
@date: 2024/3/18
*/
package main

import "fmt"

var (
	currEnv = NewEnvironment()
)

type IStmt interface {
	Exec()
}

type Expression struct { // expression ;
	Expression IExpr
}

func (e *Expression) Exec() {
	e.Expression.GetValue() // 简单执行一下
}

func NewExpression(expression IExpr) *Expression {
	return &Expression{Expression: expression}
}

type Print struct { // print expression ;
	Expression IExpr
}

func (p *Print) Exec() {
	val := p.Expression.GetValue()
	fmt.Println(val)
}

func NewPrint(expression IExpr) *Print {
	return &Print{Expression: expression}
}

type Var struct { // var name = expr ;
	Name *Token
	Expr IExpr // 初始值
}

func (v *Var) Exec() {
	var val any
	if v.Expr != nil {
		val = v.Expr.GetValue()
	}
	currEnv.Define(v.Name.Lexeme, val) // 定义变量
}

func NewVar(name *Token, expr IExpr) *Var {
	return &Var{Name: name, Expr: expr}
}

type Assign struct { // name = expr ;
	Name *Token
	Expr IExpr
}

func (a *Assign) Exec() {
	val := a.Expr.GetValue()
	currEnv.Assign(a.Name.Lexeme, val)
}

func NewAssign(name *Token, expr IExpr) *Assign {
	return &Assign{Name: name, Expr: expr}
}

type Set struct { // Object.Name=Expr
	Object IExpr
	Name   *Token
	Expr   IExpr
}

func (s *Set) Exec() {
	temp := s.Object.GetValue()
	if inst, ok := temp.(IInstance); ok {
		val := s.Expr.GetValue()
		inst.Set(s.Name.Lexeme, val)
		return
	}
	panic(fmt.Sprintf("obj %v not a Instance", temp))
}

func NewSet(object IExpr, name *Token, expr IExpr) *Set {
	return &Set{Object: object, Name: name, Expr: expr}
}

type Block struct { // { Declaration*  }
	Statements []IStmt
}

func NewBlock(statements []IStmt) *Block {
	return &Block{Statements: statements}
}

func (b *Block) Exec() {
	oldEnv := currEnv
	currEnv = NewEnvironmentWithParent(currEnv) // 添加作用域
	for _, stmt := range b.Statements {
		stmt.Exec()
	}
	currEnv = oldEnv // 移除作用域
}

type If struct { // if ( IExpr ) { IfBranch } else { ElseBranch }
	Condition            IExpr
	IfBranch, ElseBranch IStmt
}

func (i *If) Exec() {
	val := i.Condition.GetValue()
	if val.(bool) {
		i.IfBranch.Exec()
	} else if i.ElseBranch != nil {
		i.ElseBranch.Exec()
	}
}

func NewIf(condition IExpr, ifBranch IStmt, elseBranch IStmt) *If {
	return &If{Condition: condition, IfBranch: ifBranch, ElseBranch: elseBranch}
}

type For struct { // for(Init?;Condition?;Change?){Body?} ;不能省略
	Init         IStmt
	Condition    IExpr
	Change, Body IStmt
}

func (f *For) Exec() {
	oldEnv := currEnv
	currEnv = NewEnvironmentWithParent(currEnv) // 添加作用域
	if f.Init != nil {
		f.Init.Exec()
	}
	for f.Condition == nil || f.Condition.GetValue().(bool) { // 条件为空视为 true
		f.Body.Exec()
		if f.Change != nil { // 执行变更
			f.Change.Exec()
		}
	}
	currEnv = oldEnv // 移除作用域
}

func NewFor(init IStmt, condition IExpr, change IStmt, body IStmt) *For {
	return &For{Init: init, Condition: condition, Change: change, Body: body}
}

type Function struct {
	Name   *Token
	Params []*Token
	Body   IStmt
}

func (f *Function) Exec() {
	currEnv.Define(f.Name.Lexeme, NewBaseCall(f.Params, f.Body, currEnv))
}

func NewFunction(name *Token, params []*Token, body IStmt) *Function {
	return &Function{Name: name, Params: params, Body: body}
}

type Return struct {
	Expr IExpr
}

// return 并不会终止函数执行 还会继续执行后面的
// 可以直接panic，并在panic信息中包含返回只在BaseCall.Call中捕获异常，或对Exec添加异常返回
func (r *Return) Exec() {
	var val any
	if r.Expr != nil {
		val = r.Expr.GetValue()
	}
	currEnv.Assign(RETURN_KEY, val) // 必须在函数中使用 return 只有函数中会预定义RETURN_KEY
}

func NewReturn(expr IExpr) *Return {
	return &Return{Expr: expr}
}

type Class struct {
	Name, Parent *Token
	Methods      []*Function
}

func (c *Class) Exec() {
	var parent *BaseClass
	if c.Parent != nil { // 试图获取父类定义
		parent = currEnv.Get(c.Parent.Lexeme).(*BaseClass)
	}
	methods := make(map[string]*BaseCall, len(c.Methods))
	for _, method := range c.Methods {
		methods[method.Name.Lexeme] = NewBaseCall(method.Params, method.Body, currEnv)
	}
	currEnv.Define(c.Name.Lexeme, NewBaseClass(c.Name.Lexeme, parent, methods))
}

func NewClass(name *Token, parent *Token, methods []*Function) *Class {
	return &Class{Name: name, Parent: parent, Methods: methods}
}
