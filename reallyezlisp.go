package main

import "fmt"
import "strings"
import "strconv"
import "regexp"

type ReallyEzType int
var fnMap map[string]string=make(map[string]string)

func fnRewrite(from string, to string){
	fnMap[from]=to
}

const (
	Number ReallyEzType = 1
	List ReallyEzType = 2
	Stmt ReallyEzType = 3
	Reference ReallyEzType = 4
)

type Argument struct{
	//number is int[1] so
	//todo:find some way to make it smaller because slice has
	//     length and capacity whitch is not needed for just number
	data []Argument
	value int
	stmt Statement
	argType ReallyEzType
	//todo: how do i do function pointers
}
func (arg Argument) Val() int{
	if(arg.argType!=Number){
		fmt.Errorf("WANRING: reading value from non-number")
	}
	return arg.value
}
func (arg Argument) Contents() []Argument{
	if(arg.argType!=List){
		fmt.Errorf("WANRING: reading list contents from non-list")
	}
	return arg.data
}
func (arg Argument) String() string{
	var result string
	switch arg.argType{
	case Number:
		result=fmt.Sprint(arg.Val());
	case List:
		result=fmt.Sprint(arg.Contents());
	case Stmt:
		result=arg.stmt.String()
	case Reference:
		result=fmt.Sprintf("'%d",arg.Val())
	}
	return result
}
func (arg Argument) ToGo() string{
	var result string
	switch arg.argType{
	case Number:
		result=fmt.Sprintf("MkArgi(%d)",arg.Val());
	case List:
		result="MkArgl([]int{"
		for i, val:=range arg.data{
			result+=val.ToGo()
			if(i<len(arg.data)-1){
				result+=","
			}
		}
		result+="})"
	case Stmt:
		result=arg.stmt.ToGo()
	case Reference:
		result=fmt.Sprintf("args[%d]",arg.value)
	}
	return result
}

type Statement struct{
	funcName string
	args []Argument
}
func (stmt Statement) String() string{
	var result string
	result="("+stmt.funcName
	for _, v := range stmt.args{
		result+=" "+v.String()
	}
	return result+")"
}
func (stmt Statement) GoFnName() string{
	result,found:=fnMap[stmt.funcName]
	if(found){
		return result
	}
	return stmt.funcName
}
func (stmt Statement) ToGo() string{
	//(myfunction 123 (l 1 2 3))
	//makes result myfunction([]Argument{MkArgi(123),MkArgl([]int{1,2,3})})
	var result string = stmt.GoFnName() + "([]Argument{"
	for i, arg := range stmt.args{
		result+=arg.ToGo()
		if(i<len(stmt.args)-1){
			result+=","
		}
	}
	return result+"})"
}

type Function struct{
	name string
	body []Statement
}
func (fn Function) String() string{
	result:=fmt.Sprintf("(fn %s\n",fn.name)
	for i, stmt:=range fn.body {
		result+=stmt.String()
		if(i<len(fn.body)-1){
			result+="\n"
		}else{
			result+=")"
		}
	}
	return result
}
func (fn Function) ToGo() string{
	result:=fmt.Sprintf("func %s(args []Argument) Argument{\n", fn.name)
	for i,val:=range fn.body {
		if (i<len(fn.body)-1){
			result+=val.ToGo()+"\n"
		}else{
			result+=fmt.Sprintf("return %s\n", val.ToGo())
		}
	}
	return result+"}"
}

func joinFnSlice(slice1 []Function, slice2 []Function) []Function {
	//make a new slice with combined size
	var result=make([]Function, len(slice1)+len(slice2))
	//add slice1
	for i, fn:=range slice1 {
		result[i]=fn
	}
	//add slice2
	offset:=len(slice1)
	for i, fn:=range slice2 {
		result[i+offset]=fn
	}
	return result
}


type Import struct{
	source string
}

func MkFunc(funcName string, bodyStmts []Statement) Function{
	return Function{name:funcName,body:bodyStmts}
}
func MkStmt(fn string, argList[]Argument) Statement{
	return Statement{funcName:fn,args:argList}
}
func MkArgi(val int) Argument{
	return Argument{value:val,argType:Number}
}
func MkArgl(list []Argument) Argument{
	return Argument{data:list,argType:List}
}
func MkArgL(list []int) Argument{
	var result []Argument=make([]Argument,len(list))
	for i, val:=range list {
		result[i]=MkArgi(val)
	}
	return Argument{data:result,argType:List}
}
func MkArgs(statement Statement) Argument{
	return Argument{stmt:statement,argType:Stmt}
}
func MkArgS(fn string, args[]Argument) Argument{
	return MkArgs(MkStmt(fn,args))
}
func MkArgr(to int) Argument{
	return Argument{value:to,argType:Reference}
}

func matchParen(str string) int{
	var parens int=0
	for i,c:=range str {
		switch c {
		case '(':
			parens++
		case ')':
			if (parens==1){
				return i
			}
			parens--
		}
	}
	return -1
}
func extractFromParen(str string) (string, string){
	var first int=strings.Index(str,"(")
	var last int=matchParen(str[first:])
	if(first==-1){
		fmt.Println("WARNING: expected ( in ",str)
		return "", str
	}
	if(last==-1){
		fmt.Println("WARNING: expected ) in ",str)
		//no matching parenthesis
		return "", str
	}
	return str[first+1:last+first], str[last+first+1:]
}
func splitFirst2(str string) (string, string, string){
	space1:=strings.Index(str," ")
	space2:=strings.Index(str[space1+1:]," ")+space1+1
	return str[:space1], str[space1+1:space2], str[space2+1:]	
}

func parseFn(name string, contents string) Function {
	body:=make([]Statement,0)
	//var current string
	var rest string=contents
	for strings.Contains(rest,"(") {
		//current, rest=extractFromParen(rest)
		
		//body=append(body,parseStmt())
	}
	return MkFunc(name,body)
}
func parseStmt(contents []string) Statement{
	var args []Argument=make([]Argument, len(contents)-1)
	fnName:=contents[0]
	for i, val:=range contents[1:] {
		if val[0]=='(' {
			//its a statement
			stmtContents, _:=extractFromParen(val)
			var contentsArray []string=strings.Split(stmtContents," ")
			for i,val:=range contentsArray{
				contentsArray[i]=strings.Trim(val, " \n\t")
			}
			args[i]=MkArgs(parseStmt(contentsArray))
		}else if val[0]=='$' {
			//its a ref
			to,err:=strconv.Atoi(val[1:])
			if(err!=nil){
				fmt.Printf("WARNING: failed to parse Reference %s\n", val)
			}
			args[i]=MkArgr(to)
		}else{
			num,err:=strconv.Atoi(val[1:])
			if(err!=nil){
				fmt.Printf("WARNING: failed to parse Number %s\n", val)
			}
			args[i]=MkArgr(num)
		}
	}
	return MkStmt(fnName,args)
}
func parse(str string) []Function{
	var result []Function=make([]Function,1)
	var current, rest=extractFromParen(str)
	assnType,name,contents:=splitFirst2(current)
	//todo:let expresion support
	if(assnType=="fn"){
		result[0]=parseFn(name, contents)
	}
	if(strings.Contains(rest,"(")){
		return joinFnSlice(result,parse(rest))
	}else {
		return result
	}
}

type TokenType int
const (
	TkParen TokenType = 1 //
	TkKeywordFunc TokenType = 2 //fn
	TkFunc TokenType = 3 //(<func> <rest> <rest>)
	TkReference TokenType = 4 //$num
	TkEnd TokenType = 5 //)
	TkNumber TokenType = 6 //123
)
func (tkt TokenType) String() string{
	var result string = ""
	switch tkt {
	case TkParen:
		result+="("
	case TkEnd:
		result+=")"
	case TkKeywordFunc:
		result+="keyword_fn"
	case TkFunc:
		result+="fn"
	case TkReference:
		result+="$"
	case TkNumber:
		result+="#"
	}
	return result
}

type Token struct{
	content string
	tokenType TokenType
}
func (tk Token) String() string{
	return fmt.Sprintf("[%s,%s]",tk.tokenType.String(),tk.content)
}

type TokenTree struct{
	tk Token
	children []*TokenTree
}
func (tkt TokenTree) String() string{
	return tkt.TreeList(0)
}
func (tkt TokenTree) TreeList(indent int) string{
	var indentStr string
	for i:=0;i<indent;i++ {
		indentStr+="-"
	}
	result:=fmt.Sprintf("%s %s\n",indentStr,tkt.tk.String())
	for _,val:=range tkt.children {
		result+=val.TreeList(indent+1)
	}
	return result
}


func matchParenTk(tks []Token) int{
	var parens int=0
	for i,tk:=range tks {
		switch tk.tokenType {
		case TkParen:
			parens++
		case TkEnd:
			if (parens==1){
				return i
			}
			parens--
		}
	}
	return -1
}
func extractFromParenTk(tks []Token) ([]Token, []Token){
	first:=0
	for ;first<len(tks)&&tks[first].tokenType!=TkParen;first++ {}
	if first==len(tks) {
		fmt.Println("WARNING: expected TkParen at begining but statement ended")
		return make([]Token,0), tks
	}
	//fmt.Println(first)
	var last int=matchParenTk(tks[first:])+first
	if(tks[first].tokenType!=TkParen){
		fmt.Println("WARNING: expected TkParen at begining but got %s",tks[0])
		return make([]Token,0), tks
	}
	if(last==-1){
		fmt.Println("WARNING: expected TkEnd posible mismatched parens &s",tks)
		//no matching parenthesis
		return make([]Token,0), tks
	}
	return tks[first+1:last], tks[last+1:]
}

func parsetkts(tks []Token) []TokenTree {
	var result []TokenTree=make([]TokenTree,1)
	var rest []Token
	result[0], rest=parsetkt(tks)
	if(len(rest)==0){
		return result
	}else{
		return append(result, parsetkts(rest)...)
	}
}

func parsetkt(tks []Token) (TokenTree, []Token) {
	current, rest:=extractFromParenTk(tks);
	//fmt.Println(tks)
	//if(len(current)==0){
	//	return nil, rest
	//}	
	var body []*TokenTree=make([]*TokenTree,0)
	//todo: go through children and make statements trees
	for i:=1;i<len(current);i++ {
		if(current[i].tokenType==TkParen){
			inside,_:=extractFromParenTk(current[i:])
			i+=len(inside)+1//because extractfrom paren drops parens
			parsedInside,_:=parsetkt(inside)
			//if(parsedInside!=nil){
			body=append(body,&parsedInside)
			//}
		}else{
			body=append(body,&TokenTree{tk:current[i]})
		}
	}
	fmt.Printf("tk:%s  |  %s",current[0],current[0:])
	return TokenTree{tk:current[0],children:body}, rest
}

func lex(str string) []Token {
	var result []Token=make([]Token,0)
	var current Token=Token{}
	i:=0
	var number = regexp.MustCompile(`[0-9]`)
	//var letter = regexp.MustCompile(`[a-zA-Z]`)
	for ;i<len(str);i++{
		for strings.ContainsAny(string(str[i])," \n\t") {
			i++
		}
		s:=string(str[i])
		switch {
		case s=="(":
			current.tokenType=TkParen
			break
		case s==")":
			current.tokenType=TkEnd
			break
		case s=="f":
			if(str[i:i+3]=="fn "){
				current.tokenType=TkKeywordFunc
				i+=2
			}
			break
		case s=="$":
			end:=strings.Index(str[i:], " ")
			current.content=str[i+1:end+i]
			current.tokenType=TkReference
			i+=end
			break
		case number.MatchString(string(str[i])):
			end:=strings.IndexAny(str[i:], " )")
			current.content=str[i:end+i]
			current.tokenType=TkNumber
			i+=end-1
			break
		default:
			end:=strings.IndexAny(str[i:], " )")
			current.content=str[i:end+i]
			current.tokenType=TkFunc
			i+=end-1
			
		}
		fmt.Printf("c:%s, i:%d, tk:%s, val:%s\n",string(str[i]),i,current.tokenType.String(),current.content)
		result=append(result,current)
		if(i<len(str)-1){
			current=Token{}
		}
	}
	return result
}

//todo char literal
//todo builds parse tree but chops + function
//unit testing would be nice
func main(){
	fnRewrite("+","REZadd")
	//stmt:=MkStmt("ss", []Argument{MkArgS("print",[]Argument{MkArgi(123)}),MkArgS("lmao", []Argument{MkArgi(1),MkArgi(2)})})
	//stmt2:=MkStmt("+", []Argument{MkArgS("print",[]Argument{MkArgi(123)}),MkArgS("lmao", []Argument{MkArgi(1),MkArgi(2)})})
	//fn:=MkFunc("main",[]Statement{stmt,stmt2})
	
	//fmt.Println(stmt.ToGo())
	s,_:=extractFromParenTk(lex("(+ 1 (+ 3 4))"))
	//fmt.Println(parsetkts(lex("(fn main (+ 1 (+ 3 4)))")))
	fmt.Println(s)
}
