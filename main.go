package main
import (
	"./rezl"
	"fmt"
)

func main(){
	//stuff:=rezl.ReadString("\n\n(def fewre 3)(defn lmao (a b) (print a (print 3 b)))\n(defn main () (print \"lmao\" (let ((d 4) (g 6)) (print d g))))")
	stuff:=rezl.ReadString("(defn sayHello () (print \"Hello, World!\"))(defn main () (sayHello))")
	for _,expr:=range stuff {
		fmt.Println(expr)
	}
	ctx:=rezl.ScanSymbols(stuff)
	deps:=rezl.NewDependencies()
	ctx.Defns()["main"].Dependencies(deps,ctx)
	fmt.Println(ctx)
	fmt.Println(deps)
	fmt.Println(ctx.ToGo(deps))
}
