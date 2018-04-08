package rezl

import (
	"strings"
	"regexp"
	"io/ioutil"
	"errors"
	"strconv"
)

func PanOnErr(err error,message string,line int){
	if err!=nil {
		panic(message+strconv.Itoa(line))
	}
}

func PanLine(message string, line int){
	panic(message+strconv.Itoa(line))
}

func ReadFile(filename string) []RezlSlice{
	contents,err:=ioutil.ReadFile(filename)
	if err!=nil{
		panic("could not find file: "+filename)
	}
	return ReadString(string(contents))
}

func ReadString(str string) []RezlSlice{
	result:=make([]RezlSlice,0)
	line:=1
	ctx:=RezlContext{}
	for i:=0;i<len(str);i++{
		if str[i]=='\n' {
			line++
		}else if str[i]=='(' {
			match,linesConsumed,err:=MatchParen(str,i)
			if err!=nil{
				panic("Mismatched Paren on line "+strconv.Itoa(line))
			}
			expr,_:=ReadExpression(str[i:match+1],line,ctx)
			result=append(result,expr)
			line+=linesConsumed
			i=match
			
		}
	}
	return result
}

// begin should be the first opening paren to match
func MatchParen(str string,begin int) (int,int,error){
	depth:=0
	linesConsumed:=0
	for i:=begin;i<len(str);i++{
		if str[i]=='\n' {
			linesConsumed++
		}else if str[i]=='(' {
			depth++
		}else if str[i]==')' {
			depth--
			if depth==0 {
				return i,linesConsumed,nil
			}
		}
	}
	return -1,linesConsumed,errors.New("Mismatched Parenthesis")
}

// string should start and end with paren
func ReadExpression(str string, line int, ctx RezlContext) (RezlSlice, int){
	result:=make(RezlSlice,0)
	add:=func(element interface{}){
		result=append(result,element)
	}
	// start point to the begining of what were looking at
	start:=1
	number:=regexp.MustCompile("[0-9]")
	//whitespace:=regexp.MustCompile("^[ \n\t,\\)]+$")
	//trimable:=regexp.MustCompile("^[ \n\t,\\)]+")
	for i:=1;i<len(str);i++{
		//TODO: add vector literal support
		if str[i]=='('{
			// recurse on sub expression
			match,_,err:=MatchParen(str,i)
			PanOnErr(err,"mismatched paren on line ",line)
			exprStr:=str[i:match+1]
			expr,lineAfterRead:=ReadExpression(exprStr,line,ctx)
			add(expr)
			line=lineAfterRead
			i+=len(exprStr)
			start=i+1
		}else if str[i]=='"'{
			// string
			j:=strings.Index(str[i+1:],"\"")
			if(j==-1){
				PanOnErr(errors.New(""),"dunterminated string literal on line ",line)
			}
			j+=i+1
			add(str[i+1:j])
			start=j+1
			i=j+1
		}else if strings.ContainsAny(string(str[i]),") \n\t,") {
			if i-start!=0{
				word:=strings.TrimSpace(str[start:i])
				if len(word)==0{
					// empty
				}else if word[0]==':'{
					// keyword
					add(Keyword(word[1:]))
				}else if word[0]=='\\'{
					add(word[1])
				}else if number.MatchString(string(word[0])){
					if strings.Contains(word,".") {
						//float
						value,err:=strconv.ParseFloat(word,64)
						PanOnErr(err,"malformed float on line ",line)
						add(value)
					}else{
						//int
						value,err:=strconv.Atoi(word)
						PanOnErr(err,"malformed int on line ",line)
						add(value)
					}
				}else{
					// symbol
					add(Symbol{word,line})
				}
				if(str[i]=='\n'){
					line++
				}
				start=i+1
			}
		}
	}
	return result, line
}
