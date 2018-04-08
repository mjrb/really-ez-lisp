func BuiltinToGo(expr RezlSlice, ret bool) string{
	// name is given as the first item in expr
	// it has already been checked to be a symbol in the
	// anylize phase we just need to do a type assertion
	name:=expr[0].(Symbol).word
	switch name {
	case "print":
		return PrintToGo(expr, ret)
	case "println":
		return PrintlnToGo(expr, ret)
	case "+":
		return PlusToGo(expr)
	case "let":
		return LetToGo(expr)
		
	}
	PanLine("unknown builtin "+name+"on line ",expr[0].(Symbol).line)
	panic("unreachable")
}
func PrintToGo(expr RezlSlice, ret bool) string{
	result:=""
	if ret{
		result+=" func() int{"
	}
	result+="print("
	for i,arg:=range expr[1:]{
		if FinalType(arg) {
			result+=ToGo(arg)
		}else{
			switch arg.(type) {
			case Symbol:
				result+=ToGo(arg)
			case RezlSlice:
				fmt.Println("Warning: posible attempt to print non printable "+Stringify(arg)+" on line "+strconv.Itoa(expr[0].(Symbol).line))
				result+=ToGo(arg)
			default:
				panic(fmt.Sprintf("attempt to print unknown type %s (type %T) on line %d",Stringify(arg),arg,expr[0].(Symbol).line))
			}
		}
		if i<len(expr[1:])-1 {
			result+=" "
		}
	}
	result+=")"
	if ret{
		result+=";return 0}()"
	}
	return result
}
func PrintlnToGo(expr RezlSlice, ret bool) string{
	result:=""
	if ret{
		result+=" func() int{"
	}
	result+="println("
	for i,arg:=range expr[1:]{
		if FinalType(arg) {
			result+=ToGo(arg)
		}else{
			switch arg.(type) {
			case Symbol:
				result+=ToGo(arg)
			case RezlSlice:
				fmt.Println("Warning: posible attempt to print non printable "+Stringify(arg)+" on line "+strconv.Itoa(expr[0].(Symbol).line))
				result+=ToGo(arg)
			default:
				panic(fmt.Sprintf("attempt to print unknown type %s (type %T) on line %d",Stringify(arg),arg,expr[0].(Symbol).line))
			}
		}
		if i<len(expr[1:])-1 {
			result+=" "
		}
	}
	result+=")"
	if ret{
		result+=";return 0}()"
	}
	return result
}
func PlusToGo(expr RezlSlice) string{
	result:=""
	for i,elem:=range expr[1:] {
		switch elem.(type){
		case int:
		case float64:
		case uint8:
		case Symbol:
		case Keyword:
			
		default:
			
		}
	}
	return result
}
func LetToGo(expr RezlSlice) string{
	result:="{"
	return result+"}"
}
