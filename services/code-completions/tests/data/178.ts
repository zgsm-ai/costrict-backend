<｜fim▁hole｜>interface Product {
    id: number;
    name: string;
    price: number;
}

class ShoppingCart {
    private items: Product[] = [];
    
    addItem(product: Product): void {
        this.items.push(product);
    }
    
    getTotal(): number {
        return this.items.reduce((total, item) => total + item.price, 0);
    }
}

const cart = new ShoppingCart();
cart.addItem({ id: 1, name: 'Book', price: 20 });
cart.addItem({ id: 2, name: 'Pen', price: 5 });
console.log(`Total: $${cart.getTotal()}`);