const crypto = require('crypto');

function hashPassword(password, salt) {
    const hash = crypto.createHmac('sha256', salt);
    hash.update(password);
    <｜fim▁hole｜>return hash.digest('hex');
}

const password = 'mypassword';
const salt = 'mysalt';
const hashedPassword = hashPassword(password, salt);

console.log('Password:', password);
console.log('Salt:', salt);
console.log('Hashed Password:', hashedPassword);