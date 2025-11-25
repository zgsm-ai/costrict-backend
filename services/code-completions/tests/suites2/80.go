package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

func (p *Person) introduce() {
    <｜fim▁hole｜>
}

func main() {
    person := Person{
        Name: "Alice",
        Age:  30,
    }
    person.introduce()
}