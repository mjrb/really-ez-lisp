package rezl

import "fmt"
import "io/ioutil"

type ReallyEzType int
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
		result=fmt.Sprintf("$%d",arg.Val())
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
	var result string = stmt.GoFnName() + "("
	if(stmt.funcName=="if"){
		//go does not have ternary oeprater so i have to pull of this
		//nonsense of calling an anonyomous func in spot
		//func() Argument {if(arg0!=0){return arg1}; return arg3}()
		if(len(stmt.args)<2){
			fmt.Println("WARNING: if expects 2 or 3 arguments")
		}
		result=fmt.Sprintf("func() Argument{if(%s.value!=0){",stmt.args[0].ToGo())
		result+=fmt.Sprintf("return %s};",stmt.args[1].ToGo())
		if(len(stmt.args)==3){
			result+="return "+stmt.args[2].ToGo()
		}else{
			result+="return MkArgi(0)"
		}
		result+="}()"
		return result
	}
	for i, arg := range stmt.args{
		result+=arg.ToGo()
		if(i<len(stmt.args)-1){
			result+=","
		}
	}
	return result+")"
}

type Keyword interface{
	ToGo() string
}
//turns top level statement slice into a go string
func KeywordsToGo(keywords []Keyword) string{
	var result string=""
	for _,keyword:=range keywords {
		result+=keyword.ToGo()+"\n"
	}
	return result
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
	result:=fmt.Sprintf("func %s(args... Argument) Argument{\n", fn.name)
	if(fn.name=="main"){
		result=fmt.Sprintf("func main(){\n")
	}
	for i,val:=range fn.body {
		if (i<len(fn.body)-1){
			result+=val.ToGo()+"\n"
		}else{
			if(fn.name!="main"){
				result+=fmt.Sprintf("return %s\n", val.ToGo())
			}else{
				result+=fmt.Sprintf("%s\n", val.ToGo())
			}
		}
	}
	return result+"}"
}

func FnsToGo(fns []Function) string{
	var result string=""
	for _,fn:=range fns {
		result+=fn.ToGo()+"\n"
	}
	return result
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

func Read(source string) string {
	bytes,err:=ioutil.ReadFile(source)
	if(err!=nil){
		panic(fmt.Sprintf("FATAL: could not read file at %s. error:%s",source,err))
	}
	return string(bytes)
}

type Use struct {
	source string
	content string
}
func (u Use) ToGo() string {
	return u.content
}

type Import struct {
	source string
	content []Keyword
}
func (i Import) ToGo() string{
	return KeywordsToGo(i.content)
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
func MkUse(src string) Use{
	cont:=Read(src)
	return Use{source:src,content:cont}
}
func MkImport(src string) Import{
	cont:=ParseFile(src)
	return Import{source:src,content:cont}
}
