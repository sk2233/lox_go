/*
@author: sk
@date: 2024/3/16
*/
package main

import "fmt"

type TokenType uint64

const (
	EOF TokenType = iota // 文件结尾
	// Single-character tokens.
	LEFT   // (
	RIGHT  // )
	LEFT2  // {
	RIGHT2 // }
	ADD    // +
	SUB    // -
	MUL    // *
	DIV    // /
	COMMA  // ,
	DOT    // .
	SEMI   // ;
	// One or two character tokens.
	NOT    // !
	NE     // !=
	ASSIGN // =
	EQ     // ==
	GT     // >
	GE     // >=
	LT     // <
	LE     // <=
	// Literals.
	ID  // var
	STR // string
	NUM // int float
	// Keywords.
	AND
	CLASS
	ELSE
	FALSE
	FUNC
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
)

var (
	Keywords = map[string]TokenType{
		"and":    AND,
		"class":  CLASS,
		"else":   ELSE,
		"false":  FALSE,
		"for":    FOR,
		"func":   FUNC,
		"if":     IF,
		"nil":    NIL,
		"or":     OR,
		"print":  PRINT,
		"return": RETURN,
		"super":  SUPER,
		"this":   THIS,
		"true":   TRUE,
		"var":    VAR,
	}
)

type Token struct {
	Type   TokenType
	Lexeme string // 语义
	Value  any
	Line   int
}

func (t *Token) String() string {
	return fmt.Sprintf("type:%d,lexeme:%s,value:%v,line:%d", t.Type, t.Lexeme, t.Value, t.Line)
}

func NewToken(type0 TokenType, lexeme string, value any, line int) *Token {
	return &Token{Type: type0, Lexeme: lexeme, Value: value, Line: line}
}
