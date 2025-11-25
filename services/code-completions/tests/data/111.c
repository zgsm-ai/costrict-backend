#include <stdio.h>
#include <stdlib.h>

int main() {
    int *numbers;
    int size = 5;
    int i;
    
    // 分配内存
    numbers = (int *)malloc(size * sizeof(int));
    if (numbers == NULL) {
        printf("Memory allocation failed\n");
        return 1;
    }
    
    // 初始化数组
    for (i = 0; i < size; i++) {
        <｜fim▁hole｜>numbers[i] = i * 2;
    }
    
    // 打印数组
    for (i = 0; i < size; i++) {
        printf("%d ", numbers[i]);
    }
    printf("\n");
    
    // 释放内存
    free(numbers);
    
    return 0;
}