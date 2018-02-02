#ifndef SYSTEMFUNCTIONS
#define SYSTEMFUNCTIONS

#include "reallyeztypes.c"
#include <stdio.h>
#include "statementbuilder.c"

typedef struct systemfunc Systemfunc;
struct systemfunc{
  char* name;
  Argument (*fn)(Argument*, int);
  Systemfunc *next;
};

Systemfunc *system_funcs=NULL;

//maybe we need to strcpy to keep name if it gets deleted
void add_system_function(char *name, Argument (*fn)(Argument*, int)){
  Systemfunc func;
  func.name=name;
  func.fn=fn;
  if(system_funcs==NULL)
    *system_funcs=func;
  else{
    Systemfunc *system_func=system_funcs;
    for(;system_func->next!=NULL;system_func=system_func->next);
    system_func->next=&func;
  }
}


Argument system_apply(char *name, Argument args[], int argc){
  //todo: make this actually do somthing
  return *Argi(123);
}


Argument add(Argument args[], int argc){
  char argtype='i';
  int sum=0;
  for(int i=0;i<argc;i++){
    if(args[i].argtype=='i')
      sum+=args[i].arg->i;
    else if(args[i].argtype=='p'){
      sum+=args[i].arg->p;
      argtype='p';
    }
  }
  if(argtype=='p')
    return *Argp(sum);
  else if(argtype=='i')
    return *Argi(sum);
}

Argument print(Argument args[], int argc){
  for(int i=0;i<argc;i++){
    if(args[i].argtype=='i')
      printf("%i",args[i].arg->i);
    else if(args[i].argtype=='p')
      printf("POINTER:%i",args[i].arg->p);
  }
  return *Argp(0);
}

Argument printc(Argument args[], int argc){
  for(int i=0;i<argc;i++){
    if(args[i].argtype=='i')
      printf("%c",(char)args[i].arg->i);
  }
  return *Argp(0);
}

void add_system_functions(){
  add_system_function("+", &add);
  add_system_function("print", &print);
  add_system_function("printc", &printc);
}


#endif
