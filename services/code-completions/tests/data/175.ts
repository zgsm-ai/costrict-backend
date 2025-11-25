class Calculator {
    private result: number = 0;
    
    add(value: number): void {
        <｜fim▁hole｜>this.result += value;
    }
    
    subtract(value: number): void {
        this.result -= value;
    }
    
    getResult(): number {
        return this.result;
    }
    
    reset(): void {
        this.result = 0;
    }
}

const calc = new Calculator();
calc.add(10);
calc.subtract(5);
console.log('Result:', calc.getResult());