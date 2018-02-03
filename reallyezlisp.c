#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "reallyeztypes.c"
#include "statementbuilder.c"
#include "reallyezprintout.c"
#include "reallyezevalstatement.c"
#include "systemfunctions.c"
#include "systemfunctionsdefinitions.c"
#include "reallyezparser.c"



int main(int argc, char **argv){
  //add_system_functions();
  char *str="(+ 1 3)";
  printf("s");
  Statement s=parse(str);
  print_stmt(s);
  printf("\n");
}

//todo actually compile to c
//todo compile to sql
//todo free the goddamn memory
//todo conditional eval
//tree parser
//make function list into hash table
//defined functions
//todo maybe have it compile directly to asm instead
//errors

//last error segfaults during scanning
