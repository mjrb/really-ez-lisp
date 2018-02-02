CC=gcc

default: reallyezlisp.c
	$(CC) -o reallyezlisp reallyezlisp.c -Wno-format-security

clean:
	rm *~
