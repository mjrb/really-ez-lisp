#ifndef SYSTEMFUNCTIONS
#define SYSTEMFUNCTIONS

#include "reallyeztypes.c"
#include <stdio.h>
#include <string.h>
#include "statementbuilder.c"

typedef struct systemfunc Systemfunc;
struct systemfunc{
  char* name;
  Argument (*fn)(Argument*, int);
  Systemfunc *next;
};

Systemfunc *system_funcs;

//maybe we need to strcpy to keep name if it gets deleted
void add_system_function(char *name, Argument (*fn)(Argument*, int)){
  Systemfunc *func=malloc(sizeof(Systemfunc));
  func->name=name;
  func->fn=fn;
  if(system_funcs==NULL)
    system_funcs=func;
  else{
    Systemfunc *system_func=system_funcs;
    while(system_func->next!=NULL)
      system_func=system_func->next;
    system_func->next=func;
  }
}


Argument system_apply(char *name, Argument args[], int argc){
  Systemfunc *system_func=system_funcs;
  
  while(system_func!=NULL&&strcmp(system_func->name,name)!=0){
    system_func=system_func->next;  
  }
  return system_func->fn(args,argc);
}


#endif
