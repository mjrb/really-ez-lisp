#include <stdlib.h>
#include <stdio.h>
#include "reallyeztypes.c"
#include "statementbuilder.c"
#include "reallyezprintout.c"
#include "reallyezevalstatement.c"
int main(int argc, char **argv){
  Argument *args[]={Argi(123), Argi(222)};
  Statement s=Stmt("lmao", 1, args, 3);
  print_stmt(s);
  printf("\n");
}
