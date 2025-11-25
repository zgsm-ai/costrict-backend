class Person {
    constructor(name, age) {
        this.name = name;
        this.age = age;
    }
    
    introduce() {
        <｜fim▁hole｜>return `Hi, I'm ${this.name} and I'm ${this.age} years old.`;
    }
}

function main() {
    const person = new Person('Alice', 30);
    console.log(person.introduce());
}

main();