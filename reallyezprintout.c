#ifndef REALLY_EZ_PRINTOUT
#define REALLY_EZ_PRINTOUT
#include <stdio.h>

void print_stmt(Statement s);
void print_arg(Argument *arg){
  switch(arg->argtype){
  case 'i':
    printf("%i",arg->arg->i);
    break;
  case 'p':
    printf("p%i",arg->arg->p);
    break;
  case 's':
    print_stmt(*(arg->arg->s));
  }
}

void print_stmt(Statement s){
  printf("(");
  printf(s.fn.name);
  for(int i=0;i<s.argc;i++){
    printf(" ");
    print_arg(s.args[i]);
  }
  printf(")");
}

#endif
