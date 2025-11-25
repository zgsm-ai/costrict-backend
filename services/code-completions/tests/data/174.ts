function calculateArea(width: number, height: number): number {
    return width * height;
}

function calculatePerimeter(width: number, height: number): number {
    return 2 * (width + height);
}

const width = 5;
const height = 10;

<｜fim▁hole｜>const area = calculateArea(width, height);
const perimeter = calculatePerimeter(width, height);

console.log(`Width: ${width}, Height: ${height}`);
console.log(`Area: ${area}`);
console.log(`Perimeter: ${perimeter}`);