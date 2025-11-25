<｜fim▁hole｜>package main

import (
    "fmt"
    "time"
)

func main() {
    current := time.Now()
    fmt.Printf("Current time: %s\n", current.Format("2006-01-02 15:04:05"))
}