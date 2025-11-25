const http = require('http');

const server = http.createServer((req, res) => {
    if (req.url === '/') {
        <｜fim▁hole｜>res.writeHead(200, {'Content-Type': 'text/plain'});
        res.end('Hello World');
    }
});

server.listen(3000, () => {
    console.log('Server running on port 3000');
});