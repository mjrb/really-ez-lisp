#ifndef REALLY_EZ_PARSER
#define REALLY_EZ_PARSER

#include "reallyeztypes.c"
#include "statementbuilder.c"
#include <string.h>
#include <stdlib.h>

Statement parse(char* input_str){
  //strbrk can find ( or ' '
  printf("s");
  int length=strlen(input_str);
  if(length<2) printf("WARNING: probably not a well formed lisp expression");
  char *fn_name;
  //ignore first (
  int last=1;
  printf("len %d ",length);
  //search for " " and set fn name 
  int i=(int)(strchr(input_str,' ')-input_str);
  fn_name=malloc(sizeof(char)*(i-last));
  strncpy(fn_name,input_str+last,i-last);
  //continue searching for first non space
  //while *iterator!=input_string+length
  //if not ( look for next space and make this an arg 
  //if it is ( capture all the way to the matching ) and parse this
  //return statement
  Argument *args[]={Argp(0)};
  return Stmt(fn_name,1,args,1);
}

#endif
