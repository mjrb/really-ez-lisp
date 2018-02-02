#ifndef STATEMENT_BUILDER
#define STATEMENT_BUILDER

#include <stdlib.h>
#include "reallyeztypes.c"

Statement Stmt(char *name, char funtype, Argument *args[], int argc){
  fun fn;
  fn.name=name;
  fn.funtype=funtype;
  Statement s;
  s.fn=fn;
  s.args=args;
  s.argc=argc;
  return s;
}
Argument *Argi(int i){
  U_argument *u_arg=malloc(sizeof(U_argument));
  u_arg->i=i;
  Argument *arg_p=malloc(sizeof(Argument));
  arg_p->arg=u_arg;
  arg_p->argtype='i';
  return arg_p;
}
Argument *Argp(int p){
  U_argument *u_arg=malloc(sizeof(U_argument));
  u_arg->p=p;
  Argument *arg_p=malloc(sizeof(Argument));
  arg_p->arg=u_arg;
  arg_p->argtype='p';
  return arg_p;
}
Argument *Args(Statement *s){
  U_argument *u_arg=malloc(sizeof(U_argument));
  u_arg->s=s;
  Argument *arg_p=malloc(sizeof(Argument));
  arg_p->arg=u_arg;
  arg_p->argtype='s';
  return arg_p;
}

#endif
