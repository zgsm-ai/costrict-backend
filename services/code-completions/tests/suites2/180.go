package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

type WordCounter struct {
    words map[string]int
}

func NewWordCounter() *WordCounter {
    return &WordCounter{
        words: make(map[string]int),
    }
}

func (wc *WordCounter) ProcessText(text string) {
    scanner := bufio.NewScanner(strings.NewReader(text))
    scanner.Split(bufio.ScanWords)
    
    for scanner.Scan() {
        word := strings.ToLower(scanner.Text())
        <｜fim▁hole｜>
    }
}

func (wc *WordCounter) GetWordCount(word string) int {
    return wc.words[strings.ToLower(word)]
}

func (wc *WordCounter) GetTopWords(n int) []string {
    type wordCount struct {
        word  string
        count int
    }
    
    var wordCounts []wordCount
    for word, count := range wc.words {
        wordCounts = append(wordCounts, wordCount{word, count})
    }
    
    for i := 0; i < len(wordCounts)-1; i++ {
        for j := i + 1; j < len(wordCounts); j++ {
            if wordCounts[i].count < wordCounts[j].count {
                wordCounts[i], wordCounts[j] = wordCounts[j], wordCounts[i]
            }
        }
    }
    
    var result []string
    for i := 0; i < n && i < len(wordCounts); i++ {
        result = append(result, fmt.Sprintf("%s: %d", wordCounts[i].word, wordCounts[i].count))
    }
    
    return result
}

func main() {
    counter := NewWordCounter()
    
    text := "Go is an open source programming language that makes it easy to build simple, reliable, and efficient software."
    counter.ProcessText(text)
    
    fmt.Println("Word 'go' count:", counter.GetWordCount("go"))
    fmt.Println("Top 3 words:")
    for _, word := range counter.GetTopWords(3) {
        fmt.Println(word)
    }
}