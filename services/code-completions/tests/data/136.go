package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func main() {
    file, err := os.Open("input.txt")
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer file.Close()
    
    scanner := bufio.NewScanner(file)
    wordCount := make(map[string]int)
    
    for scanner.Scan() {
        line := scanner.Text()
        words := strings.Fields(line)
        
        for _, word := range words {
            <｜fim▁hole｜>
            word = strings.ToLower(word)
            wordCount[word]++
        }
    }
    
    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading file:", err)
        return
    }
    
    for word, count := range wordCount {
        fmt.Printf("%s: %d\n", word, count)
    }
}