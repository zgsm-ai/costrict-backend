package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
)

type Document struct {
	ID      int      `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type SearchEngine struct {
	documents    []*Document
	invertedIndex map[string][]int
	mutex        sync.RWMutex
}

type SearchResult struct {
	DocumentID int     `json:"document_id"`
	Score      float64 `json:"score"`
	Snippet    string  `json:"snippet"`
}

func NewSearchEngine() *SearchEngine {
	return &SearchEngine{
		documents:    make([]*Document, 0),
		invertedIndex: make(map[string][]int),
		mutex:        sync.RWMutex{},
	}
}

func (se *SearchEngine) AddDocument(doc *Document) {
	se.mutex.Lock()
	defer se.mutex.Unlock()
	
	doc.ID = len(se.documents)
	se.documents = append(se.documents, doc)
	
	// Index document
	terms := se.tokenize(doc.Title + " " + doc.Content)
	for _, term := range terms {
		se.invertedIndex[term] = append(se.invertedIndex[term], doc.ID)
	}
}

func (se *SearchEngine) tokenize(text string) []string {
	text = strings.ToLower(text)
	
	// Remove punctuation
	for _, punct := range []string{".", ",", "!", "?", ";", ":", "(", ")", "[", "]", "{", "}", "\"", "'"} {
		text = strings.ReplaceAll(text, punct, " ")
	}
	
	words := strings.Fields(text)
	return words
}

func (se *SearchEngine) calculateTF(docID int) map[string]float64 {
	doc := se.documents[docID]
	terms := se.tokenize(doc.Title + " " + doc.Content)
	
	termCount := make(map[string]int)
	for _, term := range terms {
		termCount[term]++
	}
	
	tf := make(map[string]float64)
	totalTerms := len(terms)
	
	for term, count := range termCount {
		tf[term] = float64(count) / float64(totalTerms)
	}
	
	return tf
}

func (se *SearchEngine) calculateIDF() map[string]float64 {
	idf := make(map[string]float64)
	totalDocs := len(se.documents)
	
	for term, docIDs := range se.invertedIndex {
		docFreq := len(docIDs)
		idf[term] = math.Log(float64(totalDocs) / float64(docFreq))
	}
	
	return idf
}

func (se *SearchEngine) calculateTFIDF(query string) []*SearchResult {
	se.mutex.RLock()
	defer se.mutex.RUnlock()
	
	queryTerms := se.tokenize(query)
	idf := se.calculateIDF()
	
	// Calculate TF-IDF scores for each document
	scores := make(map[int]float64)
	
	for term, idfValue := range idf {
		for _, docID := range se.invertedIndex[term] {
			tf := se.calculateTF(docID)
			if tfValue, exists := tf[term]; exists {
				scores[docID] += tfValue * idfValue
			}
		}
	}
	
	// Boost scores for documents that contain more query terms
	for docID, score := range scores {
		doc := se.documents[docID]
		docTerms := se.tokenize(doc.Title + " " + doc.Content)
		
		termMatches := 0
		for _, queryTerm := range queryTerms {
			for _, docTerm := range docTerms {
				if queryTerm == docTerm {
					termMatches++
					break
				}
			}
		}
		
		// Boost score based on term matches
		boost := float64(termMatches) / float64(len(queryTerms))
		scores[docID] = score * (1 + boost)
	}
	
	// Convert to slice and sort
	results := make([]*SearchResult, 0)
	for docID, score := range scores {
		snippet := se.generateSnippet(se.documents[docID], queryTerms)
		results = append(results, &SearchResult{
			DocumentID: docID,
			Score:      score,
			Snippet:    snippet,
		})
	}
	
	// Sort by score (descending)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	
	return results
}

func (se *SearchEngine) generateSnippet(doc *Document, queryTerms []string) string {
	content := doc.Content
	if len(content) > 200 {
		content = content[:200] + "..."
	}
	
	// Highlight query terms in snippet
	for _, term := range queryTerms {
		content = strings.ReplaceAll(content, term, "<b>"+term+"</b>")
	}
	
	return content
}

func (se *SearchEngine) Search(query string, limit int) []*SearchResult {
	if query == "" {
		return []*SearchResult{}
	}
	
	results := se.calculateTFIDF(query)
	
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}
	
	return results
}

func (se *SearchEngine) GetDocument(id int) (*Document, bool) {
	se.mutex.RLock()
	defer se.mutex.RUnlock()
	
	if id < 0 || id >= len(se.documents) {
		return nil, false
	}
	
	return se.documents[id], true
}

func (se *SearchEngine) GetStatistics() map[string]interface{} {
	se.mutex.RLock()
	defer se.mutex.RUnlock()
	
	stats := make(map[string]interface{})
	stats["total_documents"] = len(se.documents)
	stats["indexed_terms"] = len(se.invertedIndex)
	
	// Calculate average document length
	totalLength := 0
	for _, doc := range se.documents {
		if doc != nil {
			totalLength += len(doc.Title) + len(doc.Content)
		}
	}
	
	if len(se.documents) > 0 {
		stats["avg_document_length"] = totalLength / len(se.documents)
	}
	
	// Calculate most frequent terms
	termFreq := make(map[string]int)
	for term, docIDs := range se.invertedIndex {
		termFreq[term] = len(docIDs)
	}
	
	type termCount struct {
		term  string
		count int
	}
	
	var termCounts []termCount
	for term, count := range termFreq {
		termCounts = append(termCounts, termCount{term, count})
	}
	
	sort.Slice(termCounts, func(i, j int) bool {
		return termCounts[i].count > termCounts[j].count
	})
	
	var topTerms []string
	for i := 0; i < 5 && i < len(termCounts); i++ {
		topTerms = append(topTerms, termCounts[i].term)
	}
	
	stats["top_terms"] = topTerms
	
	return stats
}

func generateSampleDocuments() []*Document {
	titles := []string{
		"Introduction to Go Programming",
		"Web Development with Go",
		"Concurrency in Go",
		"Go Microservices",
		"Testing Go Applications",
	}
	
	contents := []string{
		"Go is a statically typed, compiled programming language designed at Google. It is syntactically similar to C, but with memory safety, garbage collection, structural typing, and CSP-style concurrency.",
		"Building web applications in Go is efficient and scalable. The standard library provides excellent support for HTTP servers and clients.",
		"Go makes concurrency easy with goroutines and channels. These primitives allow developers to write concurrent programs that can utilize multiple CPU cores effectively.",
		"Microservices architecture is a popular pattern for building scalable systems. Go's lightweight goroutines make it an excellent choice for microservices.",
		"Testing is a critical part of software development. Go includes a built-in testing framework that makes it easy to write unit tests, benchmarks, and examples.",
	}
	
	tags := [][]string{
		{"go", "programming", "basics"},
		{"go", "web", "development"},
		{"go", "concurrency", "performance"},
		{"go", "microservices", "architecture"},
		{"go", "testing", "quality"},
	}
	
	documents := make([]*Document, len(titles))
	for i := 0; i < len(titles); i++ {
		documents[i] = &Document{
			Title:   titles[i],
			Content: contents[i],
			Tags:    tags[i],
		}
	}
	
	return documents
}

func main() {
	fmt.Println("Search Engine Demo")
	fmt.Println("=================")
	
	// Create a new search engine
	engine := NewSearchEngine()
	
	// Add sample documents
	documents := generateSampleDocuments()
	for _, doc := range documents {
		engine.AddDocument(doc)
	}
	
	fmt.Printf("Indexed %d documents\n", len(documents))
	
	// Display statistics
	stats := engine.GetStatistics()
	fmt.Printf("\nSearch Engine Statistics:\n")
	fmt.Printf("Total documents: %d\n", stats["total_documents"])
	fmt.Printf("Indexed terms: %d\n", stats["indexed_terms"])
	fmt.Printf("Top terms: %v\n", stats["top_terms"])
	
	// Perform searches
	queries := []string{
		"Go programming",
		"concurrency",
		"web development",
		"performance",
		"testing",
	}
	
	for _, query := range queries {
		fmt.Printf("\nSearch query: \"%s\"\n", query)
		results := engine.Search(query, 3)
		
		if len(results) == 0 {
			fmt.Println("No results found")
			continue
		}
		
		for i, result := range results {
			doc, exists := engine.GetDocument(result.DocumentID)
			if !exists {
				continue
			}
			
			fmt.Printf("\nResult %d (Score: %.2f):\n", i+1, result.Score)
			fmt.Printf("Title: %s\n", doc.Title)
			fmt.Printf("Snippet: %s\n", result.Snippet)
			fmt.Printf("Tags: %v\n", doc.Tags)
		}
	}
	
	// Save index
	<｜fim▁hole｜>
	
	fmt.Println("\nSearch engine demo completed")
}