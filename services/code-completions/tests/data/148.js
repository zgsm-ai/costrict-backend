<｜fim▁hole｜>const fs = require('fs');

function readFileContent(filePath) {
    try {
        const data = fs.readFileSync(filePath, 'utf8');
        return data;
    } catch (error) {
        console.error('Error reading file:', error);
        return null;
    }
}

function main() {
    const content = readFileContent('example.txt');
    if (content) {
        console.log('File content:');
        console.log(content);
    }
}

main();