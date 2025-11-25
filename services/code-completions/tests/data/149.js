function calculateSum(numbers) {
    let sum = 0;
    for (let i = 0; i < numbers.length; i++) {
        <｜fim▁hole｜>sum += numbers[i];
    }
    return sum;
}

function main() {
    const numbers = [1, 2, 3, 4, 5];
    const result = calculateSum(numbers);
    console.log(`Sum of [${numbers}] is ${result}`);
}

main();