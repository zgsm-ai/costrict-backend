#include <iostream>
#include <string>

class Student {
private:
    std::string name;
    int age;
    
public:
    Student(std::string n, int a) : name(n), age(a) {}
    
    void display() {
        <｜fim▁hole｜>std::cout << "Name: " << name << ", Age: " << age << std::endl;
    }
};

int main() {
    Student student("Alice", 20);
    student.display();
    return 0;
}