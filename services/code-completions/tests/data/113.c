#include <stdio.h>
#include <string.h>

#define MAX_LENGTH 100

void reverseString(char *str) {
    int length = strlen(str);
    int i, j;
    char temp;
    
    for (i = 0, j = length - 1; i < j; i++, j--) {
        <｜fim▁hole｜>temp = str[i];
        str[i] = str[j];
        str[j] = temp;
    }
}

int main() {
    char str[MAX_LENGTH];
    
    printf("Enter a string: ");
    fgets(str, MAX_LENGTH, stdin);
    
    // 移除换行符
    str[strcspn(str, "\n")] = '\0';
    
    reverseString(str);
    
    printf("Reversed string: %s\n", str);
    
    return 0;
}