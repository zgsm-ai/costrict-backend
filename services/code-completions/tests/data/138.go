package main

import (
    "fmt"
    "math"
)

func calculateDistance(x1, y1, x2, y2 float64) float64 {
    dx := x2 - x1
    dy := y2 - y1
    <｜fim▁hole｜>distance := math.Sqrt(dx*dx + dy*dy)
    return distance
}

func main() {
    fmt.Printf("Distance: %.2f\n", calculateDistance(0, 0, 3, 4))
}