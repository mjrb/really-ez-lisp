#ifndef REALLY_EZ_EVAL_STATEMENT
#define REALLY_EZ_EVAL_STATEMENT

#include "statementbuilder.c"
#include "systemfunctions.c"
#include <stdlib.h>

Argument apply(char *name, Argument args[], int argc){
  return *Argi(123);//make argument
}
Argument eval(Statement s){
  Argument *args=malloc(sizeof(Argument)*s.argc);
  for(int i=0;i<s.argc;i++){
    if(s.args[i]->argtype=='s')
      args[i]=eval(*s.args[i]->arg->s);
    else
      args[i]=*s.args[i];
  }
  if(s.fn.funtype==0)
    return system_apply(s.fn.name, args, s.argc);
  else
    return apply(s.fn.name, args, s.argc);
}

#endif
