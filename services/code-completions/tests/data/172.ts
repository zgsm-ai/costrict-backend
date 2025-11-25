interface User {
    id: number;
    name: string;
}

class UserService {
    private users: User[] = [];
    
    addUser(name: string): number {
        const id = this.users.length + 1;
        <｜fim▁hole｜>this.users.push({ id, name });
        return id;
    }
    
    getUser(id: number): User | undefined {
        return this.users.find(user => user.id === id);
    }
}

const service = new UserService();
const userId = service.addUser('John Doe');
console.log('User ID:', userId);
console.log('User:', service.getUser(userId));