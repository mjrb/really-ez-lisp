CC=gcc

default: reallyezlisp.c
	$(CC) -o reallyezlisp reallyezlisp.c -Wno-format-security -std=gnu99

clean:
	rm *~
