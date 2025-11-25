#include <stdio.h>
#include <string.h>

void reverseString(char *str) {
    int length = strlen(str);
    int i, j;
    char temp;
    
    for (i = 0, j = length - 1; i < j; i++, j--) {
        temp = str[i];
        str[i] = str[j];
        str[j] = temp;
    }
}

int main() {
    char text[] = "Hello, World!";
    <｜fim▁hole｜>reverseString(text);
    printf("Original: Hello, World!\n");
    printf("Reversed: %s\n", text);
    return 0;
}