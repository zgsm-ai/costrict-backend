const events = require('events');

class EventEmitter {
    constructor() {
        this.events = {};
    }
    
    on(event, listener) {
        if (!this.events[event]) {
            <｜fim▁hole｜>this.events[event] = [];
        }
        this.events[event].push(listener);
    }
    
    emit(event, ...args) {
        if (this.events[event]) {
            this.events[event].forEach(listener => listener(...args));
        }
    }
}

const emitter = new EventEmitter();
emitter.on('test', (data) => {
    console.log('Received:', data);
});

emitter.emit('test', 'Hello World');