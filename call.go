/*
@author: sk
@date: 2024/3/23
*/
package main

import "time"

func InjectNativeFunc() {
	currEnv.Define("clock", NewClock())
}

type ICall interface { // 可被调用的函数
	Call(args []any) any
	ArgsSize() int
}

type Clock struct { // 本地方法 获取时间
}

func NewClock() *Clock {
	return &Clock{}
}

func (c *Clock) Call(args []any) any { // 数字只支持小数
	return float64(time.Now().Unix())
}

func (c *Clock) ArgsSize() int {
	return 0
}

const (
	RETURN_KEY = "$RETURN_KEY$"
)

type BaseCall struct { // 相当与一种新的类型
	// 定义改函数时的 env 函数调用时不应该使用函数调用时的env
	// 那样会访问到函数使用外的变量 应该使用函数定义时的环境，顺便实现闭包的功能
	DefineEnv *Environment
	Params    []*Token
	Body      IStmt
}

func (b *BaseCall) Call(args []any) any {
	oldEnv := currEnv
	currEnv = NewEnvironmentWithParent(b.DefineEnv) // 添加作用域
	for i := 0; i < len(b.Params); i++ {            // 绑定参数
		currEnv.Define(b.Params[i].Lexeme, args[i])
	}
	currEnv.Define(RETURN_KEY, nil) // 预定义返回值
	b.Body.Exec()                   // 执行函数体
	res := currEnv.Get(RETURN_KEY)  // 获取返回值 必须在移除作用域前
	currEnv = oldEnv                // 移除作用域
	return res
}

func (b *BaseCall) ArgsSize() int {
	return len(b.Params)
}

func (b *BaseCall) BindThis(this *BaseInstance) *BaseCall {
	env := NewEnvironmentWithParent(b.DefineEnv)
	env.Define("this", this) // 创建新的作用域并添加 this 变量
	return NewBaseCall(b.Params, b.Body, env)
}

func NewBaseCall(params []*Token, body IStmt, defineEnv *Environment) *BaseCall {
	return &BaseCall{Params: params, Body: body, DefineEnv: defineEnv}
}
