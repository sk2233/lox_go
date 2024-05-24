/*
@author: sk
@date: 2024/3/16
*/
package main

import (
	"os"
)

func main() {
	//runFile("/Users/sk/Documents/go/my_lox/main.lox") // 执行源文件模式
	//runCode("var a = 2233")                                  // 交互模式
	runCode("class A{" +
		"test(){" +
		"print \"AAA\";" +
		"}" +
		"}" +
		"class B < A{" +
		"test(){" +
		"super.test();" +
		"print \"BBB\";" +
		"}" +
		"}" +
		"var b=B();" +
		"b.test();" +
		"var num=22;" +
		"if(num>33){" +
		"print 22;" +
		"}else{" +
		"print 33;" +
		"}" +
		"for(var i=0;i<10;i=i+1){" +
		"print i;" +
		"}")
}

func runFile(path string) {
	bs, err := os.ReadFile(path)
	HandleErr(err)
	runCode(string(bs))
}

func runCode(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()
	parser := NewParser(tokens)
	stmts := parser.Parse()
	InjectNativeFunc() // 注入本地方法
	for _, stmt := range stmts {
		stmt.Exec()
	}
}
