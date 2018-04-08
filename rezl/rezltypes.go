package rezl

import (
	"fmt"
	"strconv"
)

type RezlSlice []interface{}

func (rslice RezlSlice) String() string{
	result:="("
	for i,expr:=range rslice {
		result+=Stringify(expr)
		if i<len(rslice)-1 {
			result+=" "
		}
	}
	return result+")"
}

func dupeLocals(locals map[Symbol]bool) map[Symbol]bool{
	result:=make(map[Symbol]bool)
	for key,val:=range locals{
		result[key]=val
	}
	return result
}
func (rslice RezlSlice) Dependencies(deps Dependencies,locals map[Symbol]bool,ctx RezlContext) {
	if(len(rslice))<1{
		if(len(locals)>0){
			// strange way to get first key of map without reflection
			// hopefuly this is the the name of the function
			var key Symbol
			for k:=range locals{
				key=k
				break
			}
			PanLine("attempt to evaluate empty list near symbol "+key.word+" around line ",key.line)
		}
		panic("attempt to evaluate empty list near ")
	}
	switch rslice[0].(type){
	case Symbol:
		break
	default:
		panic(fmt.Sprintf("Attempt to execute non symbol %T. lists must start with a symbol %s",rslice[0],rslice.String()))
	}
	head:=rslice[0].(Symbol)
	if head.word=="let" {
		moreLocals:=dupeLocals(locals)
		if len(rslice)<3 {
			PanLine("let takes 2 arguments. line ",head.line)
		}
		deps.add(head)
		switch rslice[1].(type){
		case RezlSlice:
			break
		default:
			PanLine("first argument to let must be list. line ",head.line)
		}
		lets:=rslice[1].(RezlSlice)
		switch lets[0].(type){
		case Symbol:
			LetDependencies(lets,deps,moreLocals,head.line,ctx)
		case RezlSlice:
			for _,elem:=range lets{
				switch elem.(type){
				case RezlSlice:
					LetDependencies(elem.(RezlSlice),deps,moreLocals,head.line,ctx)
				default:
					PanLine("Multi let arguments must be list of lists. ex (let ((x 2) (y 3)) (+ x y))). line ",head.line)
				}
			}
		}
		rslice[2:].BlindDependencies(deps, moreLocals,ctx)
		return
	}
	if(head.word=="defn"||head.word=="def"){
		fmt.Println("Warning: you may not be able to do def or defn here. line "+strconv.Itoa(head.line))
	}
	//TODO add context awarenes to make sure you give the right number
	//of args to functions
	if !ctx.Defned(head) {
		PanLine("Symbol "+head.word+" undefined on line ",head.line)
	}
	if head.BuiltIn(){
		deps.add(head)
	}else{
		ctx.defns[head.word].Dependencies(deps,ctx)
	}
	rslice[1:].BlindDependencies(deps,locals,ctx)
	
}
func LetDependencies(pair RezlSlice,deps Dependencies, moreLocals map[Symbol]bool,line int,ctx RezlContext){
	if len(pair)!=2 {
		PanLine("definition inside let must have 2 members ex (let (x 3) (print 3)). line ",line)
	}
	switch pair[0].(type){
	case Symbol:
		break
	default:
		PanLine("attempt to bind value to non symbol. line ",line)
	}
	if !FinalType(pair[1]){
		switch pair[1].(type){
		case Symbol:
			PanLine("cant define a symbol to another symbol. line ",line)
		case RezlSlice:
			pair[1].(RezlSlice).Dependencies(deps,moreLocals,ctx)
		}
	}
	moreLocals[pair[0].(Symbol)]=true
}
func (rslice RezlSlice) BlindDependencies(deps Dependencies, locals map[Symbol]bool,ctx RezlContext){
	if len(rslice)==0{
		return
	}
	for _,elem:=range rslice{
		switch elem.(type){
		case Symbol:
			sym:=elem.(Symbol)
			_,isLocal:=locals[sym]
			_,isGlobal:=ctx.defs[sym.word]
			if(!(isGlobal||isLocal)){
				PanLine("undefined symbol "+sym.word+" on line ",sym.line)
			}
			if(isGlobal&&!isLocal){
				deps.add(sym)
			}
		case RezlSlice:
			elem.(RezlSlice).Dependencies(deps,locals,ctx)
		}
	}
}


func Stringify(expr interface{}) string{
	switch expr.(type){
	case string:
		return "\""+expr.(string)+"\""
	case Symbol:
		return expr.(Symbol).word
	case int:
		return fmt.Sprintf("%d",expr.(int))
	case float64:
		return fmt.Sprintf("%f",expr.(float64))
	case Keyword:
		return ":"+string(expr.(Keyword))
	case uint8:
		return "\\"+string(expr.(uint8))
	case RezlSlice:
		return expr.(RezlSlice).String()
	default:
		return fmt.Sprintf("%T",expr)
	}
}

func GetType(expr interface{}) RezlType{
	switch expr.(type){
	case string:
		return RString
	case Symbol:
		return RSymbol
	case int:
		return RInt
	case float64:
		return RFloat
	case Keyword:
		return RKeyword
	case uint8:
		return RChar
	case RezlSlice:
		return RRezlSlice
	default:
		return REmpty
	}
}

func FinalType(expr interface{}) bool{
	switch expr.(type){
	case string:
		return true
	case Symbol:
		return false
	case int:
		return true
	case float64:
		return true
	case Keyword:
		return true
	case uint8:
		return true
	case RezlSlice:
		return false
	default:
		return false
	}
}

// symbols refer to somthing else so for debuging its
// useful to know where a symbol is being used
// if it is not defined
type Symbol struct {
	word string
	line int
}
func (sym Symbol) BuiltIn() bool{
	switch sym.word {
	case "print":
	case "println":
	case "+":
	case "-":
	case "/":
	case "*":
		return true
	case "default":
		return true
	default:
		return false
	}
}

type Keyword string

// protocol!
func ToGo(expr interface{}) string{
	switch expr.(type){
	case Var:
		return VarToGo(expr.(Var))
	case Func:
		return FuncToGo(expr.(Func))
	case string:
		return "\""+expr.(string)+"\""
	case Symbol:
		return expr.(Symbol).word
	case int:
		return fmt.Sprintf("%d",expr.(int))
	case float64:
		return fmt.Sprintf("%f",expr.(float64))
	case Keyword:
		panic("keywords not yet implemented")
		//return ":"+string(expr.(Keyword))
	case uint8:
		return "'"+string(expr.(uint8))+"'"
	case RezlSlice:
		return CallToGo(expr.(RezlSlice))
	default:
		panic(fmt.Sprintf("futile attempt to turn %T into go", expr))
	}
}
func VarToGo(def Var) string{
	//screen for reserved word
	return def.name.word+":="+ToGo(def.value)
}
func FuncToGo(defn Func) string{
	result:="func "
	result+=defn.name.word
	main:=defn.name.word=="main"
	result+="("
	for i,param:=range defn.params {
		result+=param.word+" interface{}"
		if(i<len(defn.params)-1){
			result+=", "
		}
	}
	if main {
		result+=") {\n"
	}else{
		result+=") interface{}{\n"
	}
	for i,stmt:=range defn.body {
		if i==len(defn.body)-1 {
			if !main{
				result+=CallToGo(stmt.(RezlSlice), true)
			}else{
				result+=CallToGo(stmt.(RezlSlice))
			}
		}else{
			result+=CallToGo(stmt.(RezlSlice))
		}
		result+="\n"
	}
	return result+"}"
	
}
func CallToGo(expr RezlSlice, ret ...bool) string{
	result:=""
	// name is given as the first item in expr
	// it has already been checked to be a symbol in the
	// anylize phase we just need to do a type assertion
	fnsym:=expr[0].(Symbol)
	// funky stuff for optional return argument
	if len(ret)>0{
		if ret[0] {
			result+="return"
		}
	}
	if fnsym.BuiltIn() {
		if len(ret)>0 {
			result+=BuiltinToGo(expr,ret[0])
		}else{
			result+=BuiltinToGo(expr,false)
		}
			
	}else{
		result+=fnsym.word
		result+="("
		for i,arg:=range expr[1:]{
			result+=ToGo(arg)
			if i<len(expr[1:])-1 {
				result+=", "
			}
		}
		result+=")"
	}
	return result
}



type RezlType int
const (
	RInt RezlType=iota
	RFloat
	RChar
	RString
	RRezlSlice
	RKeyword
	RSymbol
	REmpty
)
