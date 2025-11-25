#include <stdio.h>
#include <stdlib.h>

int main() {
    int size = 5;
    int *numbers;
    int i, sum = 0;
    
    // Allocate memory
    numbers = (int *)malloc(size * sizeof(int));
    if (numbers == NULL) {
        printf("Memory allocation failed\n");
        return 1;
    }
    
    // Initialize array
    for (i = 0; i < size; i++) {
        numbers[i] = i + 1;
    }
    
    // Calculate sum
    for (i = 0; i < size; i++) {
        <｜fim▁hole｜>sum += numbers[i];
    }
    
    printf("Sum: %d\n", sum);
    
    // Free memory
    free(numbers);
    
    return 0;
}