interface User {
    id: number;
    name: string;
    email: string;
}

class UserService {
    private users: User[] = [];
    
    addUser(user: User): void {
        <｜fim▁hole｜>this.users.push(user);
    }
    
    findUserById(id: number): User | undefined {
        return this.users.find(user => user.id === id);
    }
}

const service = new UserService();
service.addUser({ id: 1, name: 'John', email: 'john@example.com' });
const user = service.findUserById(1);
console.log(user);