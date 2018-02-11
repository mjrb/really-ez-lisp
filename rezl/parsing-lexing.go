package rezl

import "strings"
import "fmt"
import "strconv"
import "regexp"

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
	TkUse = 7 //use
	TkImport = 8 //import 
)
func (tkt TokenType) String() string{
	var result string = ""
	switch tkt {
	case TkParen:
		result="("
	case TkEnd:
		result=")"
	case TkKeywordFunc:
		result="keyword_fn"
	case TkFunc:
		result="fn"
	case TkReference:
		result="$"
	case TkNumber:
		result="#"
	case TkUse:
		result="keyword_use"
	case TkImport:
		result="keyword_import"
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
func (tkt TokenTree) ToArgument() Argument{
	var result Argument=Argument{}
	var err error
	switch tkt.tk.tokenType {
		case TkNumber:
		result.argType=Number;
		result.value,err=strconv.Atoi(tkt.tk.content)
		if(err!=nil){
			fmt.Printf("failed to parse token as number %s",tkt.tk.String())
		}
		case TkReference:
		result.argType=Reference
		result.value,err=strconv.Atoi(tkt.tk.content)
		if(err!=nil){
			fmt.Printf("failed to parse token as reference %s",tkt.tk.String())
		}
		case TkFunc:
		result.argType=Stmt
		result.stmt=tkt.ToStatement()
	}
	return result
}
func (tkt TokenTree) ToStatement() Statement{
	if(tkt.tk.tokenType!=TkFunc){
		fmt.Printf("warning expected function name (TkFunc) but got %s\n",tkt.tk.String())
	}

	var args []Argument=make([]Argument,0)
	result:=Statement{funcName:tkt.tk.content}
	for _,val:=range tkt.children{
		if((*val).tk.tokenType!=TkEnd){
			args=append(args,(*val).ToArgument())
		}
	}
	result.args=args
	return result
}
func (tkt TokenTree) ToFunction() Function{
	if(tkt.tk.tokenType!=TkKeywordFunc){
		fmt.Printf("warning expected function keyword (TkKeywordFunc) but got %s\n",tkt.tk.tokenType.String())
	}
	//first child should be the name of function to define
	fnNameTk:=tkt.children[0].tk

	result:=Function{name:fnNameTk.content}
	var stmts []Statement=make([]Statement,0)

	for _,val:=range tkt.children[1:] {
		stmts=append(stmts,(*val).ToStatement())
	}

	result.body=stmts
	return result
}
func (tkt TokenTree) ToUse() Use{
	//the token that has the source filename is the first child of import
	//due to abuse of the lexer this should be of tokenType fn
	srcTk:=(*tkt.children[0]).tk
	if(srcTk.tokenType!=TkFunc){
		fmt.Printf("WARNING: trying to use from invalid source: %s",srcTk.String())
	}
	source:=srcTk.content
	return MkUse(source)
}
func (tkt TokenTree) ToImport() Import{
	//the token that has the source filename is the first child of import
	//due to abuse of the lexer this should be of tokenType fn
	srcTk:=(*tkt.children[0]).tk
	if(srcTk.tokenType!=TkFunc){
		fmt.Printf("WARNING: trying to import from invalid source: %s",srcTk.String())
	}
	source:=srcTk.content
	return MkImport(source)
}
func TktsToFunction(tkts []TokenTree) []Function{
	funcs:=make([]Function,len(tkts))
	for i,tkt:=range tkts {
		funcs[i]=tkt.ToFunction()
	}
	return funcs
}
func TktsToKeywords(tkts []TokenTree) []Keyword{
	keywords:=make([]Keyword,len(tkts))
	for i,tkt:=range tkts {
		switch tkt.tk.tokenType{
		case TkKeywordFunc:
			keywords[i]=tkt.ToFunction()
		case TkUse:
			keywords[i]=tkt.ToUse()
		case TkImport:
			keywords[i]=tkt.ToImport()
		default:
			fmt.Printf("WARNING: non keyword token at top level :%s",tkt.tk.String())
		}
	}
	return keywords
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

func ParseTkts(tks []Token) []TokenTree {
	var result []TokenTree=make([]TokenTree,1)
	var rest []Token
	result[0], rest=parseTkt(tks)
	if(len(rest)==0){
		return result
	}else{
		return append(result, ParseTkts(rest)...)
	}
}

func parseTkt(tks []Token) (TokenTree, []Token) {
	current, rest:=extractFromParenTk(tks);
	//fmt.Println(tks)
	//if(len(current)==0){
	//	return nil, rest
	//}	
	var body []*TokenTree=make([]*TokenTree,0)
	//todo: go through children and make statements trees
	for i:=1;i<len(current);i++ {
		if(current[i].tokenType==TkParen){
			endParen:=matchParenTk(current[i:])
			parsedInside,_:=parseTkt(current[i:i+endParen+1])
			i+=endParen//because extractfrom paren drops parens
			//if(parsedInside!=nil){
			body=append(body,&parsedInside)
			//}
		}else{
			body=append(body,&TokenTree{tk:current[i]})
		}
	}
	fmt.Printf("tk:%s  |  %s\n",current[0],current[0:])
	return TokenTree{tk:current[0],children:body}, rest
}

func makeIntStr(str string) string{
	if(str[0]=='\\'){
		switch str[1]{
		case 'n':
			return strconv.Itoa(int('\n'))
		case '\\':
			return strconv.Itoa(int('\\'))
		case '\t':
			return strconv.Itoa(int('\t'))
		default:
			return strconv.Itoa(int(str[1]))
		}
	}
	return strconv.Itoa(int(str[0]))
				
}

func Lex(str string) []Token {
	println(str)
	println(len(str))
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
		case s=="-":
			if (number.MatchString(string(str[i+1]))){
				end:=strings.IndexAny(str[i:], "\n\t )")
				current.content=str[i:end+i]
				current.tokenType=TkNumber
				i+=end-1
				break
			}
		case s=="'":
			end:=2
			if(string(str[i+1])=="\\"){
				end=3
			}
			if(string(str[i+end])!="'"){
				fmt.Printf("WARNING: unterminated char near character %d\n",i)
			}
			current.content=makeIntStr(str[i+1:i+3])
			current.tokenType=TkNumber
			i+=end
			break
		case s=="\"":
			end:=strings.IndexByte(str[i+1:], '"')+1
			if(end==0){
				fmt.Printf("WARNING: unterminated string near char %d\n",i)
			}
			//instert (list
			current.tokenType=TkParen
			result=append(result,current)
			current=Token{}
			
			current.tokenType=TkFunc
			current.content="list"
			result=append(result,current)
			current=Token{}
			
			for j:=i+1;j<i+end;j++ {
				current.tokenType=TkNumber
				current.content=makeIntStr(str[j:j+2])
				fmt.Printf("c:%s, i:%d, tk:%s, val:%s\n",string(str[i]),i,current.tokenType.String(),current.content)
				if(str[j]=='\\'){
					j++
				}
				result=append(result,current)
				current=Token{}
			}
			//instert )
			current.tokenType=TkEnd
			
			i+=end
		case s=="f":
			if(str[i:i+3]=="fn "){
				current.tokenType=TkKeywordFunc
				i+=2
				break
			}
			fallthrough
		case s=="u":
			if(str[i:i+4]=="use "){
				current.tokenType=TkUse
				i+=3
				break
			}
			fallthrough
		case s=="i":
			if(str[i:i+7]=="import "){
				current.tokenType=TkImport
				i+=6
				break
			}
			fallthrough
		default:
			end:=strings.IndexAny(str[i:], "\n\t )")
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
	println("done lexing")
	return result
}

func ParseString(str string) []Keyword {
	tks:=Lex(str)
	tkts:=ParseTkts(tks)
	keywords:=TktsToKeywords(tkts)
	return keywords
}

func ParseFile(src string) []Keyword {
	content:=read(src)
	return ParseString(content)
}

