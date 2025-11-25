function calculateArea(shape: 'circle' | 'rectangle', dimensions: any): number {
    if (shape === 'circle') {
        const radius = dimensions.radius as number;
        <｜fim▁hole｜>return Math.PI * radius * radius;
    } else {
        const width = dimensions.width as number;
        const height = dimensions.height as number;
        return width * height;
    }
}

const circleArea = calculateArea('circle', { radius: 5 });
const rectangleArea = calculateArea('rectangle', { width: 4, height: 6 });

console.log(`Circle area: ${circleArea.toFixed(2)}`);
console.log(`Rectangle area: ${rectangleArea}`);