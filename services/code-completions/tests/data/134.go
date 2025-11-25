package main

import (
    "fmt"
    "time"
)

func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    
    <｜fim▁hole｜>return fibonacci(n-1) + fibonacci(n-2)
}

func main() {
    start := time.Now()
    result := fibonacci(10)
    elapsed := time.Since(start)
    
    fmt.Printf("Fibonacci(10) = %d\n", result)
    fmt.Printf("Time taken: %s\n", elapsed)
}