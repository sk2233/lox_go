/*
@author: sk
@date: 2024/3/17
*/
package main

import (
	"fmt"
	"strings"
)

type IExpr interface { // 表达式基类
	fmt.Stringer
	GetValue() any
}

type Binary struct { // 二元表达式
	Left, Right IExpr
	Operator    *Token
}

func (b *Binary) GetValue() any {
	left := b.Left.GetValue()
	right := b.Right.GetValue()
	switch b.Operator.Type {
	case GT:
		return left.(float64) > right.(float64)
	case GE:
		return left.(float64) >= right.(float64)
	case LT:
		return left.(float64) < right.(float64)
	case LE:
		return left.(float64) <= right.(float64)
	case DIV:
		return left.(float64) / right.(float64)
	case MUL:
		return left.(float64) * right.(float64)
	case SUB:
		return left.(float64) - right.(float64)
	case ADD: // 字符串也可以相加
		if _, ok := left.(string); ok {
			return left.(string) + right.(string)
		}
		return left.(float64) + right.(float64)
	case NE: // == != 可以应用到 数字 文本 布尔值上
		return left != right
	case EQ:
		return left == right
	default:
		panic(fmt.Sprintf("invalid TokenType %v to left %v right %v", b.Operator.Lexeme, left, right))
	}
}

func (b *Binary) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left, b.Operator.Lexeme, b.Right)
}

func NewBinary(left IExpr, right IExpr, operator *Token) *Binary {
	return &Binary{Left: left, Right: right, Operator: operator}
}

type Group struct { // ( expr )
	Expr IExpr
}

func (g *Group) GetValue() any {
	return g.Expr.GetValue()
}

func (g *Group) String() string {
	return fmt.Sprintf("( %s )", g.Expr)
}

func NewGroup(expr IExpr) *Group {
	return &Group{Expr: expr}
}

type Literal struct { // 基本元素  数字 字符串 true false nil
	Token *Token
}

func (l *Literal) GetValue() any {
	switch l.Token.Type {
	case FALSE:
		return false
	case TRUE:
		return true
	default:
		return l.Token.Value
	}
}

func (l *Literal) String() string {
	return l.Token.Lexeme
}

func NewLiteral(token *Token) *Literal {
	return &Literal{Token: token}
}

type Unary struct { // 一元表达式
	Token *Token
	Expr  IExpr
}

func (u *Unary) GetValue() any {
	val := u.Expr.GetValue()
	switch u.Token.Type {
	case NOT:
		return !val.(bool)
	case SUB:
		return -val.(float64)
	default:
		panic(fmt.Sprintf("invalud TokenType %v apply %v", u.Token.Type, val))
	}
}

func (u *Unary) String() string {
	return fmt.Sprintf("%s%s", u.Token.Lexeme, u.Expr)
}

func NewUnary(token *Token, expr IExpr) *Unary {
	return &Unary{Token: token, Expr: expr}
}

type Variable struct {
	Name *Token
}

func (v *Variable) String() string {
	return fmt.Sprintf("token %v", v.Name)
}

func (v *Variable) GetValue() any {
	return currEnv.Get(v.Name.Lexeme)
}

func NewVariable(name *Token) *Variable {
	return &Variable{Name: name}
}

type Logical struct { // OR  AND  与 NewBinary 类似 但是计算时会短路
	Left, Right IExpr
	Operator    *Token
}

func (l *Logical) String() string {
	return fmt.Sprintf("%s %s %s", l.Left, l.Operator, l.Right)
}

func (l *Logical) GetValue() any {
	val := l.Left.GetValue().(bool)
	switch l.Operator.Type {
	case AND:
		if !val {
			return false
		}
		return l.Right.GetValue()
	case OR:
		if val {
			return true
		}
		return l.Right.GetValue()
	default:
		panic(fmt.Sprintf("invalid Operator %s", l.Operator))
	}
	return nil
}

func NewLogical(left IExpr, right IExpr, operator *Token) *Logical {
	return &Logical{Left: left, Right: right, Operator: operator}
}

type Call struct { // func(args...)
	Caller IExpr   // id 调用变量
	Args   []IExpr // 参数列表
}

func (c *Call) String() string {
	buff := strings.Builder{}
	buff.WriteString(c.Caller.String())
	buff.WriteString("(")
	for _, arg := range c.Args {
		buff.WriteString(arg.String())
		buff.WriteString(",")
	}
	buff.WriteString(")")
	return buff.String()
}

func (c *Call) GetValue() any {
	temp := c.Caller.GetValue()
	if _, ok := temp.(ICall); !ok {
		panic(fmt.Sprintf("%v can't callable", temp))
	}
	caller := temp.(ICall) // 获取调用对象
	args := make([]any, 0, len(c.Args))
	for _, arg := range c.Args {
		args = append(args, arg.GetValue())
	}
	if len(args) != caller.ArgsSize() { // 调用参数校验
		panic(fmt.Sprintf("func %v args not match %d != %d", caller, len(args), caller.ArgsSize()))
	}
	return caller.Call(args) // 进行调用
}

func NewCall(caller IExpr, args []IExpr) *Call {
	return &Call{Caller: caller, Args: args}
}

type Get struct {
	Object IExpr
	Name   *Token
}

func NewGet(object IExpr, name *Token) *Get {
	return &Get{Object: object, Name: name}
}

func (g *Get) String() string {
	return fmt.Sprintf("obj %v . name %v", g.Object, g.Name)
}

func (g *Get) GetValue() any {
	temp := g.Object.GetValue()
	if inst, ok := temp.(IInstance); ok {
		return inst.Get(g.Name.Lexeme)
	}
	panic(fmt.Sprintf("obj %v not a Instance", temp))
}

type This struct {
}

func (t *This) String() string {
	return "this"
}

func (t *This) GetValue() any {
	return currEnv.Get("this")
}

func NewThis() *This {
	return &This{}
}

type Super struct {
	Method *Token // 暂时从父类能拿的只能是方法，没有属性预定义
}

func (s *Super) String() string {
	return fmt.Sprintf("super.%s", s.Method)
}

func (s *Super) GetValue() any {
	obj := currEnv.Get("this").(*BaseInstance)
	return obj.GetSuperMethod(s.Method.Lexeme)
}

func NewSuper(method *Token) *Super {
	return &Super{Method: method}
}
