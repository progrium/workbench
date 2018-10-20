package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

typedef struct Book {
   char  title[50]; // 50
   int   book_id;	// 4
   float price;		// 4
} Book;

void printBook( void* book ) {
	struct Book *b;
	b = (struct Book*) book;
	printf( "Book title : %s\n", b->title);
	printf( "Book price : %f\n", b->price);
	printf( "Book book_id : %d\n", b->book_id);
	printf( "C sizeof : %lu\n", sizeof(struct Book));
}

void *newBook() {
	struct Book *book;
	book = (struct Book*) malloc(sizeof(struct Book));
	strcpy( book->title, "C Programming");
	book->book_id = 6495407;
	book->price = 10.50;
	return book;
}

void *printBookCopy( void* book, size_t n) {
	char *b = malloc(n);
	memcpy(b, book, n);
	return b;
}

void run() {
	printBook(newBook());
}

*/
import "C"
import (
	"fmt"
	"unsafe"
)

type Book struct {
	Title [50]uint8
	ID    int32
	Price float32
}

func main() {
	C.run()
	var t [50]uint8
	b := Book{
		Title: t,
		ID:    6495407,
		Price: 10.50,
	}
	fmt.Println("Go sizeof :", unsafe.Sizeof(b))
	C.printBook(unsafe.Pointer(&b))
}
