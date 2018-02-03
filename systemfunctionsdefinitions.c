#ifndef SYSTEM_FUNCTIONS_DEFINITIONS
#define SYSTEM_FUNCTIONS_DEFINITIONS

#include "reallyeztypes.c"
#include "systemfunctions.c"
#include <stdio.h>

int numbertize(Argument arg){
  if(arg.argtype=='p')
    return arg.arg->p;
  else
    return arg.arg->i;
}

//printing
Argument print(Argument* args, int argc){
  //print out leadin elements with spaces
  for(int i=0;i<argc-1;i++){
    if(args[i].argtype=='i')
      printf("%d ",args[i].arg->i);
    else if(args[i].argtype=='p')
      printf("POINTER:%i ",args[i].arg->p);
  }
  //print last element with newline
  if(args[argc-1].argtype=='i')
    printf("%d\n",args[argc-1].arg->i);
  else if(args[argc-1].argtype=='p')
    printf("POINTER:%i\n",args[argc-1].arg->p);
  //return null
  return *Argp(0);
}

Argument printc(Argument args[], int argc){
  for(int i=0;i<argc;i++){
    if(args[i].argtype=='i')
      printf("%c",(char)args[i].arg->i);
  }
  printf("\n");
  return *Argp(0);
}

//arithmetic
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

Argument sub(Argument args[], int argc){
  char argtype='i';
  int sum=0;
  for(int i=0;i<argc;i++){
    if(args[i].argtype=='i')
      sum-=args[i].arg->i;
    else if(args[i].argtype=='p'){
      sum-=args[i].arg->p;
      argtype='p';
    }
  }
  if(argtype=='p')
    return *Argp(sum);
  else if(argtype=='i')
    return *Argi(sum);
}

Argument mult(Argument args[], int argc){
  char argtype='i';
  int sum=0;
  for(int i=0;i<argc;i++){
    if(args[i].argtype=='i')
      sum*=args[i].arg->i;
    else if(args[i].argtype=='p'){
      sum*=args[i].arg->p;
      argtype='p';
    }
  }
  if(argtype=='p')
    return *Argp(sum);
  else if(argtype=='i')
    return *Argi(sum);
}

Argument divide(Argument args[], int argc){
  char argtype='i';
  int sum=0;
  for(int i=0;i<argc;i++){
    if(args[i].argtype=='i')
      sum-=args[i].arg->i;
    else if(args[i].argtype=='p'){
      sum-=args[i].arg->p;
      argtype='p';
    }
  }
  if(argtype=='p')
    return *Argp(sum);
  else if(argtype=='i')
    return *Argi(sum);
}

//comparison
Argument equals(Argument args[], int argc){
  if(argc<2)
    printf("WARNING: = expects 2 inputs\n");
  int e1=numbertize(args[0]);
  int e2=numbertize(args[1]);
  if(e1==e2)
    return *Argi(1);
  else 
    return *Argi(0);
}

Argument lessthan(Argument args[], int argc){
  if(argc<2)
    printf("WARNING: = expects 2 inputs\n");
  int e1=numbertize(args[0]);
  int e2=numbertize(args[1]);
  if(e1<e2)
    return *Argi(1);
  else 
    return *Argi(0);
}

Argument lequals(Argument args[], int argc){
  if(argc<2)
    printf("WARNING: = expects 2 inputs\n");
  int e1=numbertize(args[0]);
  int e2=numbertize(args[1]);
  if(e1<=e2)
    return *Argi(1);
  else 
    return *Argi(0);
}

Argument if_func(Argument args[], int argc){
  if(argc<3)
    printf("WARNING: if expects 3 inputs\n");
  int condition=numbertize(args[0]);
  if(condition==1)
    return args[1];
  else
    return args[2];
}





//addfunctions at runtime
void add_system_functions(){
  add_system_function("+", &add);
  add_system_function("-", &sub);
  add_system_function("*", &mult);
  add_system_function("/", &divide);
  add_system_function("=", &equals);
  add_system_function("<", &lessthan);
  add_system_function("=", &lequals);
  add_system_function("if", &if_func);
  add_system_function("print", &print);
  add_system_function("printc", &printc);
}

#endif
