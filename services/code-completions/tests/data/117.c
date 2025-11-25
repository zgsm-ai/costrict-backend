<｜fim▁hole｜>#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define MAX_LENGTH 100

int main() {
    char str[MAX_LENGTH];
    
    printf("Enter a string: ");
    fgets(str, MAX_LENGTH, stdin);
    
    // Remove newline character if present
    str[strcspn(str, "\n")] = '\0';
    
    printf("You entered: %s\n", str);
    printf("Length: %zu\n", strlen(str));
    
    return 0;
}