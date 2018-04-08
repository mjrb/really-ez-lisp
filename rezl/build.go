package rezl
import "fmt"
type FuncGenorator interface{
	ToGo() string
}
type Func struct{
	name Symbol
	params []Symbol
	body RezlSlice
	// eventualy add infered typing instead of full dynamic
	// returnType RezlType
	// add context for local vars?
}
func (fn Func) String() string{
	result:="(def "
	result+=fn.name.word
	result+=" "
	result+="("
	for i,param:=range fn.params{
		result+=param.word
		if i<len(fn.params)-1{
			result+=" "
		}
	}
	result+=")\n"
	for i,stmt:=range fn.body{
		result+="   "+Stringify(stmt)
		if i<len(fn.body)-1{
			result+="\n"
		}
	}
	return result+")"
}

type Dependencies map[string]bool
func (deps Dependencies) add(sym Symbol){
	if !deps.exists(sym) {
		deps[sym.word]=true
	}
}
func (deps Dependencies) exists(sym Symbol) bool{
	_,exists:=deps[sym.word]
	return exists
}
func NewDependencies() Dependencies{
	return Dependencies(make(map[string]bool))
}

func (fn Func) Dependencies(deps Dependencies, ctx RezlContext){
	println(fn.name.word)
	if deps.exists(fn.name) {
		return
	}
	locals:=make(map[Symbol]bool)
	locals[fn.name]=true
	for _,param:=range fn.params{
		locals[param]=true
	}
	for _,stmt:=range fn.body{
		switch stmt.(type){
		case RezlSlice:
			stmt.(RezlSlice).Dependencies(deps,locals,ctx)
		case Symbol:
			PanLine("Member of function body cannot be Symbol. line",stmt.(Symbol).line)
		}
	}
	deps.add(fn.name)
}

type Var struct {
	name Symbol
	value interface{}
	// to be added if type inference is added
	// varType RezlType
	
}
func (v Var) String() string{
	return fmt.Sprintf("(def %s %s)",v.name.word,Stringify(v.value))
}

type RezlContext struct {
	defns map[string]Func
	defs map[string]Var
}
func (ctx RezlContext) Defns() map[string]Func{
	return ctx.defns
}
func (ctx RezlContext) Defs() map[string]Var{
	return ctx.defs
}
func (ctx RezlContext) Defned(sym Symbol) bool{
	_,funcDefined:=ctx.defns[sym.word]
	if(funcDefined){
		return true
	}else if(sym.BuiltIn()){
		return true
	}
	return false
		
}
func (ctx RezlContext) ToGo(deps Dependencies) string{
	result:="package main\n"
	//var write bool 
	for word,def:=range ctx.defs {
		_,write:=deps[word]
		if write {
			result+=ToGo(def)+"\n"
		}
	}
	for word,defn:=range ctx.defns {
		_,write:=deps[word]
		if write {
			result+=ToGo(defn)+"\n"
		}
	}
	return result
	
}
// unused until infered typing. may be added
// also the nmumber may grow
const nativeFuncs=24
func NativeFuncs() RezlContext{
	ctx:=RezlContext{defns:make(map[string]Func)}
	return ctx
}

func (ctx RezlContext) String() string{
	result:="(defcontext\n"
	result+=":defs [\n"
	for _,def:=range ctx.defs{
		result+=def.String()
		result+="\n"
	}
	result+="]\n:defns [\n"
	for _,defn:=range ctx.defns{
		result+=defn.String()
		result+="\n"
	}
	return result+"])"
}

func ScanSymbols(tree []RezlSlice) RezlContext{
	// main starts out as -1 not found
	ctx:=RezlContext{defns:make(map[string]Func),defs:make(map[string]Var)}
	def:=func(name Symbol,value interface{}){
		ctx.defs[name.word]=Var{name,value}
	}
	defn:=func(name Symbol,params []Symbol,body RezlSlice){
		ctx.defns[name.word]=Func{name,params,body}
	}
	for _,topLevelStmt:=range tree {
		var body RezlSlice
		var params []Symbol
		var name Symbol
		var value interface{}
		var sym Symbol
		switch topLevelStmt[0].(type){
		case Symbol:
			sym=topLevelStmt[0].(Symbol)	
		default:
			panic("first element of top level statment is not symbol")
		}
		switch sym.word {
		case "defn":
			if len(topLevelStmt)<4 {
				PanLine("defn takes three arguments. line ",sym.line)
			}
			switch topLevelStmt[1].(type){
			case Symbol:
				name=topLevelStmt[1].(Symbol)
			default:
				PanLine("function name must be symbol. line ",sym.line)
			}
			switch topLevelStmt[2].(type){
			case RezlSlice:
				params=make([]Symbol,len(topLevelStmt[2].(RezlSlice)))
				// convert RezlSlice to []Symbol
				for i,param:=range topLevelStmt[2].(RezlSlice){
					switch param.(type){
					case Symbol:
						params[i]=param.(Symbol)
					default:
						fmt.Printf("%T",param)
						PanLine("Parameter to function must be symbol. line ",sym.line)
					}
				}
			default:
				PanLine("Parameters to a function must be in the form of RezlSlice. line ",sym.line)
			}
			body=topLevelStmt[3:]
			defn(name,params,body)
		case "def":
			if len(topLevelStmt)!=3 {
				PanLine("def takes two arguments. line ",sym.line)
			}
			switch topLevelStmt[1].(type){
			case Symbol:
				name=topLevelStmt[1].(Symbol)
			default:
				PanLine("function name must be symbol. line ",sym.line)
			}
			value=topLevelStmt[2]
			def(name,value)
		default:
			PanLine("toplevel statment is not def or defn on line ",sym.line)
		}
	}
	return ctx	
}
// staging for type inference
func (ctx RezlContext) ScanTypes(){
	/*for i,def:=range ctx.defs{
		ctx.defs[i].varType=GetType(def.value)
	}*/
}

