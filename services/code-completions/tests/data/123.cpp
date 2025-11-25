#include <iostream>
#include <string>
#include <map>

class Student {
private:
    std::string name;
    int id;
    std::map<std::string, double> grades;
    
public:
    Student(std::string name, int id) : name(name), id(id) {}
    
    void addGrade(const std::string& course, double grade) {
        <｜fim▁hole｜>grades[course] = grade;
    }
    
    double getGPA() const {
        if (grades.empty()) return 0.0;
        
        double sum = 0.0;
        for (const auto& grade : grades) {
            sum += grade.second;
        }
        
        return sum / grades.size();
    }
    
    void printInfo() const {
        std::cout << "Student: " << name << " (ID: " << id << ")" << std::endl;
        std::cout << "GPA: " << getGPA() << std::endl;
    }
};

int main() {
    Student student("John Doe", 12345);
    student.addGrade("Math", 90.5);
    student.addGrade("Physics", 85.0);
    student.printInfo();
    
    return 0;
}