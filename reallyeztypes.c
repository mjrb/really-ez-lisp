#ifndef REALLY_EZ_TYPES
#define REALLY_EZ_TYPES

typedef struct fun{
  char *name;
  //0 for fundemental 1 for defined
  char funtype;
} fun;
typedef struct statement Statement;
typedef struct argument Argument;

typedef union u_argument {
  int i;
  int p;
  Statement *s;
} U_argument;

struct argument {
  U_argument *arg;
  //i for int p for pointer and s for stmt
  char argtype;
};
  

struct statement{
  fun fn;
  Argument **args;
  int argc;
};

#endif
