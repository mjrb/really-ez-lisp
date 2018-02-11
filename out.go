package main

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

func read(source string) string {
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
	cont:=read(src)
	return Use{source:src,content:cont}
}

//actual system functions
//rewritten in rezl/rewrite
func RezAdd(args []Argument) Argument{
	sum:=0
	for _,arg:=range args{
		if(arg.argType==Number){
			sum+=arg.value
		}else{
			println("WARNING: attempt to add non-Number "+arg.String()+" has been thwarted")
		}
	}
	return MkArgi(sum)
}
func RezSub(args []Argument) Argument{
	if(len(args)==1){
		return MkArgi(-args[0].Val())
	}else if(len(args)==2){
		return MkArgi(args[0].Val()-args[1].Val())
	}
	println("WARING: - expects 1 or 2 arguments")
	return MkArgi(0)
}
func RezMul(args []Argument) Argument{
	sum:=0
	for _,arg:=range args{
		if(arg.argType==Number){
			sum*=arg.value
		}else{
			println("WARNING: attempt to add non-Number "+arg.String()+" has been thwarted")
		}
	}
	return MkArgi(sum)
}
func RezDiv(args []Argument) Argument{
	if(len(args)==1){
		return MkArgi(1/args[0].Val())
	}else if(len(args)==2){
		return MkArgi(args[0].Val()/args[1].Val())
	}
	println("WARING: - expects 1 or 2 arguments")
	return MkArgi(0)
}

func RezEqual(args []Argument) Argument{
	if(len(args)==2){
		if(args[0].value==args[1].value){
			return MkArgi(1)
		}else{
			return MkArgi(0)
		}
	}
	fmt.Printf("WARNING bad comparison %s = %s",args[0],args[1])
	return MkArgi(0)
}
func RezLesser(args []Argument) Argument{
	if(len(args)==2){
		if(args[0].value<args[1].value){
			return MkArgi(1)
		}else{
			return MkArgi(0)
		}
	}
	fmt.Printf("WARNING bad comparison %s < %s",args[0],args[1])
	return MkArgi(0)
}
func RezGreater(args []Argument) Argument{
	if(len(args)==2){
		if(args[0].value>args[1].value){
			return MkArgi(1)
		}else{
			return MkArgi(0)
		}
	}
	fmt.Printf("WARNING bad comparison %s > %s",args[0],args[1])
	return MkArgi(0)
}

func RezIf(args []Argument) Argument{
	if(len(args)>=2){
		if(args[0].value==1){
			return args[1]
		}
		if(len(args)>2){
			return args[2]
		}
	}
	return MkArgi(0)
	
}

func RezPrint(args []Argument) Argument{
	for _,arg:=range args {
		print(arg.String())
	}
	return MkArgi(0)
}
func RezPrintc(args []Argument) Argument{
	for _,arg:=range args {
		if(arg.argType==Number){
			print(string(rune(arg.value)))
		}else if(arg.argType==List){
			RezPrintc(arg.data)
		}
	}
	return MkArgi(0)
}
func list(args []Argument) Argument{
	return MkArgl(args)
}
func get(args []Argument) Argument{
	if(len(args)==2){
		if(args[1].argType==Number){
			return MkArgi(args[1].value)
		}
		return MkArgi(args[1].data[args[0].value].value)
	}else if(len(args)==3){
		if(args[2].argType==Number){
			return args[2]
		}
		return MkArgl(args[2].data[args[0].value:args[1].value])
	}
	contents:=make([]Argument,len(args)-2)
	for i,arg:=range args[2:] {
		contents[i]=get([]Argument{args[0],args[1],arg})
	}
	return MkArgl(contents)
}
func RezLen(args []Argument) Argument{
	if(len(args)!=1){
		println("WARNING: attempt get lengths from multiple things")
	}
	return MkArgi(len(args[0].Contents()))
}
func join(args []Argument) Argument{
	result:=make([]Argument,0)
	for _,arg:=range args {
		if(arg.argType==List){
			result=append(result,arg.data...)
		}else{
			result=append(result,arg)
		}
	}
	return MkArgl(result)
}
func and(args []Argument) Argument{
	for _,arg:=range args {
		if(arg.Val()==0){
			return MkArgi(0)
		}
	}
	return MkArgi(1);
}
func or(args []Argument) Argument{
	for _,arg:=range args {
		if(arg.Val()!=0){
			return MkArgi(1)
		}
	}
	return MkArgi(0);
}
func tail(args []Argument) Argument{
	if(len(args)==1){
		if(args[0].argType==Number){
			return MkArgi(args[0].value)
		}
		return MkArgl(args[0].data[1:])
	}
	contents:=make([]Argument,len(args))
	for i,arg:=range args {
		contents[i]=tail([]Argument{arg})
	}
	return MkArgl(contents)
}
func head(args []Argument) Argument{
	if(len(args)==1){
		if(args[0].argType==Number){
			return MkArgi(args[0].value)
		}
		return MkArgi(args[0].data[0].value)
	}
	contents:=make([]Argument,len(args))
	for i,arg:=range args {
		contents[i]=head([]Argument{arg})
	}
	return MkArgl(contents)
}

func merge(args []Argument) Argument{
return func() Argument{if(and([]Argument{RezEqual([]Argument{RezLen([]Argument{args[0]}),MkArgi(0)}),RezGreater([]Argument{RezLen([]Argument{args[1]}),MkArgi(0)})}).value!=0){return func() Argument{if(RezEqual([]Argument{RezLen([]Argument{args[1]}),MkArgi(1)}).value!=0){return get([]Argument{MkArgi(0),args[1]})};return join([]Argument{get([]Argument{MkArgi(0),args[1]}),merge([]Argument{args[0],tail([]Argument{args[1]})})})}()};return func() Argument{if(RezEqual([]Argument{RezLen([]Argument{args[1]}),MkArgi(0)}).value!=0){return func() Argument{if(RezEqual([]Argument{RezLen([]Argument{args[0]}),MkArgi(1)}).value!=0){return get([]Argument{MkArgi(0),args[0]})};return func() Argument{if(RezGreater([]Argument{RezLen([]Argument{args[0]}),MkArgi(1)}).value!=0){return join([]Argument{head([]Argument{args[0]}),merge([]Argument{tail([]Argument{args[0]}),args[1]})})};return RezPrint([]Argument{args[0],args[1],MkArgi(10)})}()}()};return func() Argument{if(RezLesser([]Argument{get([]Argument{MkArgi(0),args[0]}),get([]Argument{MkArgi(0),args[1]})}).value!=0){return join([]Argument{head([]Argument{args[0]}),merge([]Argument{tail([]Argument{args[0]}),args[1]})})};return join([]Argument{head([]Argument{args[1]}),merge([]Argument{args[0],tail([]Argument{args[1]})})})}()}()}()
}
func mergeSort(args []Argument) Argument{
return func() Argument{if(RezEqual([]Argument{MkArgi(2),RezLen([]Argument{args[0]})}).value!=0){return func() Argument{if(RezGreater([]Argument{get([]Argument{MkArgi(1),args[0]}),head([]Argument{args[0]})}).value!=0){return list([]Argument{head([]Argument{args[0]}),get([]Argument{MkArgi(1),args[0]})})};return list([]Argument{get([]Argument{MkArgi(1),args[0]}),head([]Argument{args[0]})})}()};return func() Argument{if(RezEqual([]Argument{MkArgi(1),RezLen([]Argument{args[0]})}).value!=0){return list([]Argument{head([]Argument{args[0]})})};return merge([]Argument{mergeSort([]Argument{get([]Argument{MkArgi(0),RezDiv([]Argument{RezLen([]Argument{args[0]}),MkArgi(2)}),args[0]})}),mergeSort([]Argument{get([]Argument{RezDiv([]Argument{RezLen([]Argument{args[0]}),MkArgi(2)}),RezLen([]Argument{args[0]}),args[0]})})})}()}()
}
func main(){
RezPrint([]Argument{mergeSort([]Argument{list([]Argument{MkArgi(7),MkArgi(4),MkArgi(3),MkArgi(8),MkArgi(9),MkArgi(6),MkArgi(2),MkArgi(1),MkArgi(5)})})})
RezPrintc([]Argument{MkArgi(10)})
}
