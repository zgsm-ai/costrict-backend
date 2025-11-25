interface Product {
    name: string;
    price: number;
}

function calculateTotal(products: Product[]): number {
    return products.reduce((total, product) => {
        <｜fim▁hole｜>return total + product.price;
    }, 0);
}

const products: Product[] = [
    { name: 'Apple', price: 1.5 },
    { name: 'Banana', price: 0.8 },
    { name: 'Orange', price: 1.2 }
];

const total = calculateTotal(products);
console.log(`Total: $${total.toFixed(2)}`);