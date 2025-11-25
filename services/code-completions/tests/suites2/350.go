package main

import (
	"container/heap"
	"fmt"
	"math"
	"strings"
)

type HuffmanNode struct {
	char  rune
	freq  int
	left  *HuffmanNode
	right *HuffmanNode
	index int
}

type PriorityQueue []*HuffmanNode

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].freq < pq[j].freq
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*HuffmanNode)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func buildFrequencyTable(text string) map[rune]int {
	freq := make(map[rune]int)
	for _, char := range text {
		freq[char]++
	}
	return freq
}

func buildHuffmanTree(freq map[rune]int) *HuffmanNode {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	for char, count := range freq {
		node := &HuffmanNode{
			char:  char,
			freq:  count,
			left:  nil,
			right: nil,
		}
		heap.Push(&pq, node)
	}

	for pq.Len() > 1 {
		left := heap.Pop(&pq).(*HuffmanNode)
		right := heap.Pop(&pq).(*HuffmanNode)

		internal := &HuffmanNode{
			char:  0,
			freq:  left.freq + right.freq,
			left:  left,
			right: right,
		}

		heap.Push(&pq, internal)
	}

	return heap.Pop(&pq).(*HuffmanNode)
}

func generateCodes(root *HuffmanNode, code string, codes map[rune]string) {
	if root == nil {
		return
	}

	if root.char != 0 {
		codes[root.char] = code
		return
	}

	generateCodes(root.left, code+"0", codes)
	generateCodes(root.right, code+"1", codes)
}

func encodeText(text string, codes map[rune]string) string {
	var encoded strings.Builder
	for _, char := range text {
		encoded.WriteString(codes[char])
	}
	return encoded.String()
}

func decodeText(encoded string, root *HuffmanNode) string {
	var decoded strings.Builder
	current := root

	for _, bit := range encoded {
		if bit == '0' {
			current = current.left
		} else {
			current = current.right
		}

		if current.char != 0 {
			decoded.WriteRune(current.char)
			current = root
		}
	}

	return decoded.String()
}

func calculateCompressionRatio(original, encoded string) float64 {
	originalBits := len(original) * 8
	encodedBits := len(encoded)
	<｜fim▁hole｜>
	return float64(originalBits) / float64(encodedBits)
}

func main() {
	text := "this is an example of a huffman tree"
	fmt.Printf("Original text: \"%s\"\n", text)
	fmt.Printf("Original size: %d bytes (%d bits)\n", len(text), len(text)*8)

	// Build frequency table
	freq := buildFrequencyTable(text)
	fmt.Println("\nCharacter frequencies:")
	for char, count := range freq {
		fmt.Printf("%c: %d\n", char, count)
	}

	// Build Huffman tree
	huffmanTree := buildHuffmanTree(freq)

	// Generate Huffman codes
	codes := make(map[rune]string)
	generateCodes(huffmanTree, "", codes)

	// Print Huffman codes
	fmt.Println("\nHuffman Codes:")
	for char, code := range codes {
		fmt.Printf("%c: %s\n", char, code)
	}

	// Encode text
	encoded := encodeText(text, codes)
	fmt.Printf("\nEncoded text: %s\n", encoded)
	fmt.Printf("Encoded size: %d bits\n", len(encoded))

	// Calculate and display compression ratio
	ratio := calculateCompressionRatio(text, encoded)
	fmt.Printf("Compression ratio: %.2f:1\n", ratio)

	// Decode text
	decoded := decodeText(encoded, huffmanTree)
	fmt.Printf("\nDecoded text: \"%s\"\n", decoded)

	// Verify that original and decoded texts match
	if text == decoded {
		fmt.Println("Compression and decompression successful!")
	} else {
		fmt.Println("Error: Original and decoded texts do not match!")
	}

	// Calculate entropy
	var entropy float64
	totalChars := len(text)
	for _, count := range freq {
		probability := float64(count) / float64(totalChars)
		entropy -= probability * math.Log2(probability)
	}
	fmt.Printf("\nEntropy: %.2f bits per character\n", entropy)
}
