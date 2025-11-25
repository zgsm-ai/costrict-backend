function filterArray<T>(items: T[], predicate: (item: T) => boolean): T[] {
    const result: T[] = [];
    for (const item of items) {
        <｜fim▁hole｜>if (predicate(item)) {
            result.push(item);
        }
    }
    return result;
}

const numbers = [1, 2, 3, 4, 5, 6];
const evenNumbers = filterArray(numbers, (n) => n % 2 === 0);
console.log('Even numbers:', evenNumbers);

const words = ['apple', 'banana', 'cherry', 'date'];
const longWords = filterArray(words, (w) => w.length > 5);
console.log('Long words:', longWords);