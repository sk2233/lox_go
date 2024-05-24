/*
@author: sk
@date: 2024/3/16
*/
package main

import (
	"bytes"
	"fmt"
	"strconv"
)

type Scanner struct {
	Source string
	Index  int
	Line   int
}

func (s *Scanner) ScanTokens() []*Token {
	tokens := make([]*Token, 0)
	s.Index = 0
	s.Line = 1
	for s.Index < len(s.Source) {
		if token := s.ScanToken(); token != nil {
			tokens = append(tokens, token)
		}
	}
	tokens = append(tokens, NewToken(EOF, "", nil, s.Line)) // 添加截断标记
	return tokens
}

func (s *Scanner) ScanToken() *Token {
	ch := s.Read()
	switch ch {
	case '(':
		return NewToken(LEFT, "(", nil, s.Line)
	case ')':
		return NewToken(RIGHT, ")", nil, s.Line)
	case '{':
		return NewToken(LEFT2, "{", nil, s.Line)
	case '}':
		return NewToken(RIGHT2, "}", nil, s.Line)
	case ',':
		return NewToken(COMMA, ",", nil, s.Line)
	case '.':
		return NewToken(DOT, ".", nil, s.Line)
	case ';':
		return NewToken(SEMI, ";", nil, s.Line)
	case '+':
		return NewToken(ADD, "+", nil, s.Line)
	case '-':
		return NewToken(SUB, "-", nil, s.Line)
	case '*':
		return NewToken(MUL, "*", nil, s.Line)
	case '/':
		if s.Match('/') { // 注释  暂时只支持单行注释
			for s.HasMore() && s.Read() != '\n' { // 移除全部注释
			}
			s.Line++
			return nil
		}
		return NewToken(DIV, "/", nil, s.Line)
	case '!':
		if s.Match('=') {
			return NewToken(NE, "!=", nil, s.Line)
		}
		return NewToken(NOT, "!", nil, s.Line)
	case '=':
		if s.Match('=') {
			return NewToken(EQ, "==", nil, s.Line)
		}
		return NewToken(ASSIGN, "=", nil, s.Line)
	case '<':
		if s.Match('=') {
			return NewToken(LE, "<=", nil, s.Line)
		}
		return NewToken(LT, "<", nil, s.Line)
	case '>':
		if s.Match('=') {
			return NewToken(GE, ">=", nil, s.Line)
		}
		return NewToken(GT, ">", nil, s.Line)
	case '"': // 字符串处理 只能处理单行字符串
		buff := bytes.Buffer{}
		for s.HasMore() {
			if ch = s.Read(); ch != '"' {
				buff.WriteByte(ch)
			} else {
				return NewToken(STR, buff.String(), buff.String(), s.Line)
			}
		}
		Error(s.Line, "no end string")
		return nil
	case ' ', '\t', '\r':
		return nil // skip
	case '\n':
		s.Line++
		return nil
	default:
		if IsDigit(ch) {
			buff := bytes.Buffer{}
			buff.WriteByte(ch)
			for s.HasMore() && (IsDigit(s.Get()) || s.Get() == '.') {
				buff.WriteByte(s.Read())
			}
			res, err := strconv.ParseFloat(buff.String(), 64)
			if err != nil {
				Error(s.Line, fmt.Sprintf("err num of %s err = %v", buff.String(), err))
				return nil
			}
			return NewToken(NUM, buff.String(), res, s.Line)
		} else if IsAlpha(ch) {
			buff := bytes.Buffer{}
			buff.WriteByte(ch)
			for s.HasMore() && (IsAlpha(s.Get()) || IsDigit(s.Get())) {
				buff.WriteByte(s.Read())
			}
			str := buff.String()
			if type0, ok := Keywords[str]; ok { // 关键字处理
				return NewToken(type0, str, nil, s.Line)
			}
			return NewToken(ID, buff.String(), nil, s.Line) // 变量处理
		}
		Error(s.Line, fmt.Sprintf("unknown of %c", ch))
		return nil
	}
}

func (s *Scanner) Read() uint8 {
	s.Index++
	return s.Source[s.Index-1]
}

func (s *Scanner) Get() uint8 {
	return s.Source[s.Index]
}

func (s *Scanner) Match(ch uint8) bool {
	if !s.HasMore() {
		return false
	}
	if s.Source[s.Index] != ch {
		return false
	}
	s.Index++ // 匹配的话自动后移
	return true
}

func (s *Scanner) HasMore() bool {
	return s.Index < len(s.Source)
}

func (s *Scanner) UnRead() {
	s.Index--
}

func NewScanner(source string) *Scanner {
	return &Scanner{Source: source}
}
