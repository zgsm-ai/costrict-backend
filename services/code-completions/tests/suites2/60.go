package main

import (
	"fmt"
	"strconv"
	"strings"
)

func calculateSum(numbers []int) int {
	sum := 0
	for _, num := range numbers {
		<｜fim▁hole｜>
	}
	return sum
}

func main() {
	input := "1,2,3,4,5"
	var numbers []int

	for _, s := range strings.Split(input, ",") {
		n, _ := strconv.Atoi(s)
		numbers = append(numbers, n)
	}

	result := calculateSum(numbers)
	fmt.Println("Sum:", result)
}
