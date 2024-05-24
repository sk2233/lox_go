/*
@author: sk
@date: 2024/3/18
*/
package main

import "fmt"

type Parser struct {
	Tokens []*Token
	Index  int
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{
		Tokens: tokens,
		Index:  0,
	}
}

func (p *Parser) Parse() []IStmt {
	res := make([]IStmt, 0)
	for p.Get().Type != EOF {
		res = append(res, p.Declaration())
	}
	return res
}

func (p *Parser) Declaration() IStmt { // Declaration -> ClassDeclaration | FuncDeclaration | VarDeclaration | Statement | Assignment
	if p.Match(CLASS) {
		return p.ClassDeclaration()
	}
	if p.Match(FUNC) {
		return p.FuncDeclaration()
	}
	if p.Match(VAR) {
		return p.VarDeclaration()
	}
	if p.Get().Type == ID {
		return p.Assignment()
	}
	return p.Statement()
}

func (p *Parser) ClassDeclaration() IStmt { // ClassDeclaration -> class ID ( < ID )? { FuncDeclaration* }
	name := p.MustRead(ID)
	var parent *Token
	if p.Match(LT) { // 可选父类
		parent = p.MustRead(ID)
	}
	p.MustMatch(LEFT2)
	methods := make([]*Function, 0)
	for !p.Match(RIGHT2) {
		methods = append(methods, p.FuncDeclaration())
	}
	return NewClass(name, parent, methods)
}

func (p *Parser) FuncDeclaration() *Function { // FuncDeclaration -> func ID( Param? )block
	name := p.MustRead(ID)
	p.MustMatch(LEFT)
	params := make([]*Token, 0)
	if !p.Match(RIGHT) { // Param -> ID ( , ID )*
		params = append(params, p.MustRead(ID))
		for p.Match(COMMA) {
			params = append(params, p.MustRead(ID))
		}
		p.MustMatch(RIGHT)
	}
	p.MustMatch(LEFT2)
	body := p.Block()
	return NewFunction(name, params, body)
}

func (p *Parser) Assignment() IStmt { // Assignment -> ( call . )? leftExpr = Expression;
	leftExpr := p.Expression() // 可能是  leftExpr  或  obj.leftExpr method().leftExpr  或单独的表达式
	if p.Match(SEMI) {         // 只有左边单独的表达式
		return NewExpression(leftExpr)
	}
	p.MustMatch(ASSIGN) // 赋值
	rightExpr := p.Expression()
	p.MustMatch(SEMI)
	switch temp := leftExpr.(type) {
	case *Variable:
		return NewAssign(temp.Name, rightExpr)
	case *Get: // 把解析到的 Get 转换为 Set
		return NewSet(temp.Object, temp.Name, rightExpr)
	default:
		panic(fmt.Sprintf("invalid left obj %v", leftExpr))
	}
}

func (p *Parser) VarDeclaration() IStmt { // VarDeclaration -> var name ( = Expression) ? ;
	name := p.MustRead(ID)
	var expr IExpr
	if p.Match(ASSIGN) {
		expr = p.Expression()
	}
	p.MustMatch(SEMI)
	return NewVar(name, expr)
}

func (p *Parser) Statement() IStmt { // Statement -> ReturnStatement | ForStatement | IfStatement | ExpressionStatement | PrintStatement | Block
	if p.Match(RETURN) {
		return p.ReturnStatement()
	}
	if p.Match(FOR) {
		return p.ForStatement()
	}
	if p.Match(IF) {
		return p.IfStatement()
	}
	if p.Match(PRINT) {
		return p.PrintStatement()
	}
	if p.Match(LEFT2) {
		return p.Block()
	}
	return p.ExpressionStatement()
}

func (p *Parser) ReturnStatement() IStmt { // ReturnStatement -> return expression ;
	var res IExpr
	if !p.Match(SEMI) {
		res = p.Expression()
		p.MustMatch(SEMI)
	}
	return NewReturn(res)
}

func (p *Parser) ForStatement() IStmt { // For -> for (VarDeclaration?;Expression?;Assignment?){Statement?}
	p.MustMatch(LEFT)
	var init IStmt
	if !p.Match(SEMI) {
		p.MustMatch(VAR)
		init = p.VarDeclaration()
	}
	var condition IExpr
	if !p.Match(SEMI) {
		condition = p.Expression()
		p.MustMatch(SEMI)
	}
	var change IStmt
	if !p.Match(RIGHT) { // 这里没有 ; 不能直接使用 Assignment
		name := p.Read()
		p.MustMatch(ASSIGN)
		expr := p.Expression()
		change = NewAssign(name, expr)
		p.MustMatch(RIGHT)
	}
	p.MustMatch(LEFT2)
	body := p.Block() // 直接复用block 会再创建一个 变量作用域还好
	return NewFor(init, condition, change, body)
}

func (p *Parser) IfStatement() IStmt { // If -> if ( Expression ) { Statement } else { Statement }
	p.MustMatch(LEFT)
	condition := p.Expression()
	p.MustMatch(RIGHT)
	p.MustMatch(LEFT2)
	ifBranch := p.Block()
	var elseBranch IStmt
	if p.Match(ELSE) {
		p.MustMatch(LEFT2)
		elseBranch = p.Block()
	}
	return NewIf(condition, ifBranch, elseBranch)
}

func (p *Parser) Block() IStmt { // Block - > { Declaration* }
	res := make([]IStmt, 0)
	for !p.Match(RIGHT2) {
		res = append(res, p.Declaration())
	}
	return NewBlock(res)
}

func (p *Parser) PrintStatement() IStmt { // PrintStatement -> print Expression ;
	expr := p.Expression()
	p.MustMatch(SEMI)
	return NewPrint(expr)
}

func (p *Parser) ExpressionStatement() IStmt { // ExpressionStatement -> Expression ;
	expr := p.Expression()
	p.MustMatch(SEMI)
	return NewExpression(expr)
}

// 各种语法对应的递归实现，实际就是对语法优先级的定义，最高优先级在最低树上先运算，()可以提高优先级

func (p *Parser) Expression() IExpr { // Expression -> Equality ( ( AND | OR) Equality)*
	left := p.Equality()                            // 这里 and or 视为同等优先级，实际不相等
	for p.Get().Type == AND || p.Get().Type == OR { // 不停合并
		operator := p.Read()
		right := p.Equality()
		left = NewLogical(left, right, operator)
	}
	return left
}

func (p *Parser) Equality() IExpr { // Equality -> Comparison (( != | == )Comparison)*
	left := p.Comparison()
	for p.Get().Type == NE || p.Get().Type == EQ { // 不停合并
		operator := p.Read()
		right := p.Comparison()
		left = NewBinary(left, right, operator)
	}
	return left
}

func (p *Parser) Comparison() IExpr { // Comparison -> Term (( > | >= | < | <= )Term)*
	left := p.Term()
	for p.Get().Type == GT || p.Get().Type == GE || p.Get().Type == LT || p.Get().Type == LE {
		operator := p.Read()
		right := p.Term()
		left = NewBinary(left, right, operator)
	}
	return left
}

func (p *Parser) Term() IExpr { // Term -> Factor (( - | + )Factor)*
	left := p.Factor()
	for p.Get().Type == SUB || p.Get().Type == ADD {
		operator := p.Read()
		right := p.Factor()
		left = NewBinary(left, right, operator)
	}
	return left
}

func (p *Parser) Factor() IExpr { // Factor -> Unary (( / | * )Unary)*
	left := p.Unary()
	for p.Get().Type == DIV || p.Get().Type == MUL {
		operator := p.Read()
		right := p.Unary()
		left = NewBinary(left, right, operator)
	}
	return left
}

func (p *Parser) Unary() IExpr { // Unary -> (( ! | - )Unary)|Call
	if p.Get().Type == NOT || p.Get().Type == SUB {
		operator := p.Read()
		right := p.Unary()
		return NewUnary(operator, right)
	}
	return p.Call()
}

func (p *Parser) Call() IExpr { // Call -> Primary ( ( args? ) | . ID )*    函数的多重调用 属性的多重调用
	expr := p.Primary() // call的对象主要是 id
	for {
		if p.Match(LEFT) {
			expr = p.SingleCall(expr) // 递归函数的单次调用
		} else if p.Match(DOT) {
			name := p.MustRead(ID)
			expr = NewGet(expr, name) // 属性多次点链接
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) SingleCall(expr IExpr) IExpr { // 单次调用
	args := make([]IExpr, 0) // args -> ( Expression ( , Expression )* )
	if !p.Match(RIGHT) {
		args = append(args, p.Expression())
		for p.Match(COMMA) {
			args = append(args, p.Expression())
		}
		p.MustMatch(RIGHT)
	}
	return NewCall(expr, args)
}

func (p *Parser) Primary() IExpr { // Primary -> NUM | STR | true | false | nil | '(' Expression ')' | id
	if p.Get().Type == FALSE || p.Get().Type == TRUE || p.Get().Type == NIL ||
		p.Get().Type == NUM || p.Get().Type == STR {
		return NewLiteral(p.Read())
	}
	if p.Match(THIS) { // 与id类似 不过取固定变量名 this
		return NewThis()
	}
	if p.Match(SUPER) {
		p.MustMatch(DOT)
		method := p.MustRead(ID)
		return NewSuper(method)
	}
	if p.Get().Type == ID {
		return NewVariable(p.Read())
	}
	p.MustMatch(LEFT)
	expr := p.Expression()
	p.MustMatch(RIGHT)
	return NewGroup(expr)
}

func (p *Parser) MustMatch(type0 TokenType) {
	if !p.Match(type0) {
		panic(fmt.Sprintf("invalid token %v", p.Get()))
	}
}

func (p *Parser) Match(type0 TokenType) bool {
	if p.Index >= len(p.Tokens) {
		return false
	}
	if p.Tokens[p.Index].Type != type0 {
		return false
	}
	p.Index++
	return true
}

func (p *Parser) Get() *Token {
	return p.Tokens[p.Index]
}

func (p *Parser) Read() *Token {
	p.Index++
	return p.Tokens[p.Index-1]
}

func (p *Parser) MustRead(type0 TokenType) *Token {
	res := p.Read()
	if res.Type != type0 {
		panic(fmt.Sprintf("invalid token %v need type %v", res, type0))
	}
	return res
}
