package main

import (
	"fmt"
	"sync"
)

type Product struct {
	ID     int
	Name   string
	Price  float64
	Stock  int
}

type Inventory struct {
	products map[int]*Product
	mu       sync.RWMutex
}

func NewInventory() *Inventory {
	return &Inventory{
		products: make(map[int]*Product),
	}
}

func (inv *Inventory) AddProduct(p *Product) {
	inv.mu.Lock()
	defer inv.mu.Unlock()
	inv.products[p.ID] = p
}

func (inv *Inventory) GetProduct(id int) (*Product, bool) {
	inv.mu.RLock()
	defer inv.mu.RUnlock()
	product, exists := inv.products[id]
	<｜fim▁hole｜>
	return product, exists
}

func (inv *Inventory) UpdateStock(id, quantity int) bool {
	inv.mu.Lock()
	defer inv.mu.Unlock()
	
	product, exists := inv.products[id]
	if !exists {
		return false
	}
	
	product.Stock += quantity
	return true
}

func (inv *Inventory) RemoveProduct(id int) bool {
	inv.mu.Lock()
	defer inv.mu.Unlock()
	
	_, exists := inv.products[id]
	if !exists {
		return false
	}
	
	delete(inv.products, id)
	return true
}

func (inv *Inventory) ListProducts() []*Product {
	inv.mu.RLock()
	defer inv.mu.RUnlock()
	
	result := make([]*Product, 0, len(inv.products))
	for _, p := range inv.products {
		result = append(result, p)
	}
	return result
}

func main() {
	inventory := NewInventory()
	
	inventory.AddProduct(&Product{ID: 1, Name: "Laptop", Price: 999.99, Stock: 10})
	inventory.AddProduct(&Product{ID: 2, Name: "Phone", Price: 499.99, Stock: 20})
	
	if p, ok := inventory.GetProduct(1); ok {
		fmt.Printf("Product: %+v\n", p)
	}
	
	inventory.UpdateStock(1, 5)
	inventory.RemoveProduct(2)
	
	products := inventory.ListProducts()
	fmt.Printf("Total products: %d\n", len(products))
}