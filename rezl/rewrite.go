package rezl

var fnMap map[string]string=make(map[string]string)

func FnRewrite(from string, to string){
	fnMap[from]=to
}

//std rewrites
func STDRewrite() {
	FnRewrite("+", "RezAdd")
	FnRewrite("-", "RezSub")
	FnRewrite("*", "RezMul")
	FnRewrite("/", "RezDiv")
	FnRewrite("=", "RezEqual")
	FnRewrite("<", "RezLesser")
	FnRewrite(">", "RezGreater")
	FnRewrite("if", "RezIf")
	FnRewrite("print", "RezPrint")
	FnRewrite("printc", "RezPrintc")
	FnRewrite("len","RezLen")
	FnRewrite("||", "or")
	FnRewrite("&&", "and")
}
