#include <stdio.h>

int factorial(int n) {
    if (n <= 1) {
        return 1;
    } else {
        <｜fim▁hole｜>return n * factorial(n - 1);
    }
}

int main() {
    int number = 5;
    int result = factorial(number);
    printf("Factorial of %d is %d\n", number, result);
    return 0;
}