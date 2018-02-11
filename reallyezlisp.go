package main

import "fmt"
import "./rezl"
import "os"
import "flag"

func WriteToFile (content string, destination string) {
	f, err:=os.Create(destination)
	defer f.Close()
	if(err!=nil){
		panic(fmt.Sprintf("FATAL: could not create file %s error:%s",destination,err))
	}
	_,err =f.Write([]byte(content))
	if(err!=nil){
		panic(fmt.Sprintf("FATAL: could not write to file %s error:%s",destination,err))
	}
}

//todo char literal
//todo builds parse tree but chops + function
//unit testing would be nice
func main(){
	srcFile:=flag.String("i", "pls name input file", "the input file name")
	outFile:=flag.String("o", "out.go", "the output file name")
	flag.Parse()
	rezl.STDRewrite()
	keywords:=rezl.ParseFile(*srcFile)
	//keywords:=rezl.ParseString("(use std)(fn main (printc 72 101 108 108 111 44 32 87 111 114 100 108 33 10))")
	fmt.Println(rezl.KeywordsToGo(keywords))
	WriteToFile(rezl.KeywordsToGo(keywords), *outFile)
}
