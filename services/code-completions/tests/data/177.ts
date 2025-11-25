function isEven(num: number): boolean {
    return num % 2 === 0;
}

function filterEvenNumbers(numbers: number[]): number[] {
    return numbers.filter(num => isEven(num));
}

const numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
<｜fim▁hole｜>const evenNumbers = filterEvenNumbers(numbers);

console.log('Original numbers:', numbers);
console.log('Even numbers:', evenNumbers);