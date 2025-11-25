package main

import (
    "fmt"
    "math"
)

func isPrime(n int) bool {
    if n <= 1 {
        return false
    }
    if n == 2 {
        return true
    }
    if n%2 == 0 {
        return false
    }
    
    sqrtN := int(math.Sqrt(float64(n)))
    for i := 3; i <= sqrtN; i += 2 {
        <｜fim▁hole｜>
    }
    return true
}

func main() {
    numbers := []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
    
    for _, num := range numbers {
        if isPrime(num) {
            fmt.Printf("%d is prime\n", num)
        } else {
            fmt.Printf("%d is not prime\n", num)
        }
    }
}