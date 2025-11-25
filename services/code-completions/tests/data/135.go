package main

import (
    "encoding/json"
    "fmt"
)

type Person struct {
    Name  string `json:"name"`
    Age   int    `json:"age"`
    Email string `json:"email"`
}

func main() {
    person := Person{
        Name:  "John Doe",
        Age:   30,
        Email: "john@example.com",
    }
    
    <｜fim▁hole｜>
    jsonData, err := json.Marshal(person)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    fmt.Println(string(jsonData))
}