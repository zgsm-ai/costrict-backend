function processArray(arr) {
    const result = [];
    arr.forEach(item => {
        <｜fim▁hole｜>if (item % 2 === 0) {
            result.push(item * 2);
        }
    });
    return result;
}

function main() {
    const numbers = [1, 2, 3, 4, 5, 6];
    const processed = processArray(numbers);
    console.log('Original:', numbers);
    console.log('Processed:', processed);
}

main();