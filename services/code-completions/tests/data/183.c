#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

#define MAX_BOOKS 1000
#define MAX_TITLE_LENGTH 100
#define MAX_AUTHOR_LENGTH 50
#define MAX_ISBN_LENGTH 20

typedef enum {
    FICTION,
    NON_FICTION,
    SCIENCE,
    TECHNOLOGY,
    ART,
    HISTORY,
    BIOGRAPHY,
    OTHER
} BookCategory;

typedef struct {
    char title[MAX_TITLE_LENGTH];
    char author[MAX_AUTHOR_LENGTH];
    char isbn[MAX_ISBN_LENGTH];
    BookCategory category;
    int publicationYear;
    float price;
    int quantity;
    int timesBorrowed;
} Book;

typedef struct {
    Book books[MAX_BOOKS];
    int count;
} Library;

/**
 * Convert category enum to string
 * @param category Book category enum
 * @return String representation of category
 */
const char* categoryToString(BookCategory category) {
    switch (category) {
        case FICTION: return "Fiction";
        case NON_FICTION: return "Non-Fiction";
        case SCIENCE: return "Science";
        case TECHNOLOGY: return "Technology";
        case ART: return "Art";
        case HISTORY: return "History";
        case BIOGRAPHY: return "Biography";
        default: return "Other";
    }
}

/**
 * Initialize a new library
 * @param lib Pointer to the library to initialize
 */
void initLibrary(Library* lib) {
    lib->count = 0;
    srand(time(NULL));
}

/**
 * Add a new book to the library
 * @param lib Pointer to the library
 * @param title Book title
 * @param author Book author
 * @param isbn Book ISBN
 * @param category Book category
 * @param year Publication year
 * @param price Book price
 * @param quantity Initial quantity
 * @return 1 if successful, 0 if library is full
 */
int addBook(Library* lib, const char* title, const char* author, const char* isbn, 
           BookCategory category, int year, float price, int quantity) {
    if (lib->count >= MAX_BOOKS) {
        return 0;
    }
    
    Book* book = &lib->books[lib->count];
    strncpy(book->title, title, MAX_TITLE_LENGTH - 1);
    book->title[MAX_TITLE_LENGTH - 1] = '\0';
    
    strncpy(book->author, author, MAX_AUTHOR_LENGTH - 1);
    book->author[MAX_AUTHOR_LENGTH - 1] = '\0';
    
    strncpy(book->isbn, isbn, MAX_ISBN_LENGTH - 1);
    book->isbn[MAX_ISBN_LENGTH - 1] = '\0';
    
    book->category = category;
    book->publicationYear = year;
    book->price = price;
    book->quantity = quantity;
    book->timesBorrowed = 0;
    
    lib->count++;
    return 1;
}

/**
 * Find a book by ISBN
 * @param lib Pointer to the library
 * @param isbn ISBN to find
 * @return Pointer to book if found, NULL otherwise
 */
Book* findBookByISBN(Library* lib, const char* isbn) {
    for (int i = 0; i < lib->count; i++) {
        if (strcmp(lib->books[i].isbn, isbn) == 0) {
            return &lib->books[i];
        }
    }
    return NULL;
}

/**
 * Find books by author
 * @param lib Pointer to the library
 * @param author Author name to search
 * @param results Array to store results
 * @param maxResults Maximum number of results to store
 * @return Number of books found
 */
int findBooksByAuthor(Library* lib, const char* author, Book* results, int maxResults) {
    int count = 0;
    for (int i = 0; i < lib->count && count < maxResults; i++) {
        if (strcmp(lib->books[i].author, author) == 0) {
            results[count] = lib->books[i];
            count++;
        }
    }
    return count;
}

/**
 * Find books by category
 * @param lib Pointer to the library
 * @param category Category to search
 * @param results Array to store results
 * @param maxResults Maximum number of results to store
 * @return Number of books found
 */
int findBooksByCategory(Library* lib, BookCategory category, Book* results, int maxResults) {
    int count = 0;
    for (int i = 0; i < lib->count && count < maxResults; i++) {
        if (lib->books[i].category == category) {
            results[count] = lib->books[i];
            count++;
        }
    }
    return count;
}

/**
 * Borrow a book
 * @param lib Pointer to the library
 * @param isbn ISBN of the book to borrow
 * @return 1 if successful, 0 if book not found or out of stock
 */
int borrowBook(Library* lib, const char* isbn) {
    Book* book = findBookByISBN(lib, isbn);
    if (book == NULL || book->quantity <= 0) {
        return 0;
    }
    
    book->quantity--;
    book->timesBorrowed++;
    return 1;
}

/**
 * Return a book
 * @param lib Pointer to the library
 * @param isbn ISBN of the book to return
 * @return 1 if successful, 0 if book not found
 */
int returnBook(Library* lib, const char* isbn) {
    Book* book = findBookByISBN(lib, isbn);
    if (book == NULL) {
        return 0;
    }
    
    book->quantity++;
    return 1;
}

/**
 * Print book information
 * @param book Pointer to the book
 */
void printBook(const Book* book) {
    printf("Title: %s\n", book->title);
    printf("Author: %s\n", book->author);
    printf("ISBN: %s\n", book->isbn);
    printf("Category: %s\n", categoryToString(book->category));
    printf("Year: %d\n", book->publicationYear);
    printf("Price: $%.2f\n", book->price);
    printf("Quantity: %d\n", book->quantity);
    printf("Times Borrowed: %d\n", book->timesBorrowed);
}

/**
 * Print all books in the library
 * @param lib Pointer to the library
 */
void printAllBooks(const Library* lib) {
    for (int i = 0; i < lib->count; i++) {
        printf("--- Book %d ---\n", i + 1);
        printBook(&lib->books[i]);
        printf("\n");
    }
}

/**
 * Calculate total value of all books
 * @param lib Pointer to the library
 * @return Total value
 */
float calculateTotalValue(const Library* lib) {
    float total = 0.0f;
    for (int i = 0; i < lib->count; i++) {
        total += lib->books[i].price * lib->books[i].quantity;
    }
    return total;
}

/**
 * Find most popular book (most times borrowed)
 * @param lib Pointer to the library
 * @return Pointer to most popular book, NULL if library is empty
 */
Book* findMostPopularBook(Library* lib) {
    if (lib->count == 0) {
        return NULL;
    }
    
    Book* mostPopular = &lib->books[0];
    for (int i = 1; i < lib->count; i++) {
        if (lib->books[i].timesBorrowed > mostPopular->timesBorrowed) {
            mostPopular = &lib->books[i];
        }
    }
    
    return mostPopular;
}

int main() {
    Library library;
    initLibrary(&library);
    
    // Add books
    addBook(&library, "The Great Gatsby", "F. Scott Fitzgerald", "978-0-7432-7356-5", FICTION, 1925, 12.99, 5);
    addBook(&library, "To Kill a Mockingbird", "Harper Lee", "978-0-06-112008-4", FICTION, 1960, 14.99, 3);
    addBook(&library, "1984", "George Orwell", "978-0-452-28423-4", FICTION, 1949, 13.99, 7);
    addBook(&library, "A Brief History of Time", "Stephen Hawking", "978-0-553-38016-3", SCIENCE, 1988, 18.99, 4);
    addBook(&library, "The Art of War", "Sun Tzu", "978-1-59030-225-7", HISTORY, -500, 9.99, 6);
    
    // Borrow some books
    borrowBook(&library, "978-0-7432-7356-5");
    borrowBook(&library, "978-0-7432-7356-5");
    borrowBook(&library, "978-0-452-28423-4");
    
    // Print all books
    printAllBooks(&library);
    
    // Print library statistics
    printf("Total Books: %d\n", library.count);
    printf("Total Value: $%.2f\n", calculateTotalValue(&library));
    
    Book* popular = findMostPopularBook(&library);
    if (popular) {
        printf("Most Popular Book: %s (borrowed %d times)\n", 
               popular->title, popular->timesBorrowed);
    }
    
    return 0;
}