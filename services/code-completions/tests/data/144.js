const fs = require('fs');

class FileReader {
    constructor(filePath) {
        this.filePath = filePath;
    }
    
    read() {
        <｜fim▁hole｜>return fs.readFileSync(this.filePath, 'utf8');
    }
}

const reader = new FileReader('./example.txt');
console.log(reader.read());