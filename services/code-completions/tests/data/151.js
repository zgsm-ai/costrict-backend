function fibonacci(n) {
    if (n <= 1) {
        return n;
    } else {
        return fibonacci(n - 1) + fibonacci(n - 2);
    }
}

function main() {
    const n = 10;
    const result = fibonacci(n);
    console.log(`Fibonacci(${n}) = ${result}`);<｜fim▁hole｜>
}

main();