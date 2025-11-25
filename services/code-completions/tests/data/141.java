import java.util.*;

public class StudentManagementSystem {
    private Map<String, Student> students;
    
    public StudentManagementSystem() {
        students = new HashMap<>();
        <｜fim▁hole｜>
    }
    
    public void addStudent(String id, String name, int age) {
        if (students.containsKey(id)) {
            System.out.println("Student with ID " + id + " already exists");
            return;
        }
        
        Student student = new Student(id, name, age);
        students.put(id, student);
        System.out.println("Added student: " + name);
    }
    
    public void displayAllStudents() {
        System.out.println("All Students:");
        for (Student student : students.values()) {
            System.out.println(student);
        }
    }
    
    public static void main(String[] args) {
        StudentManagementSystem sms = new StudentManagementSystem();
        sms.addStudent("1001", "Alice Johnson", 20);
        sms.addStudent("1002", "Bob Smith", 21);
        sms.displayAllStudents();
    }
}

class Student {
    private String id;
    private String name;
    private int age;
    
    public Student(String id, String name, int age) {
        this.id = id;
        this.name = name;
        this.age = age;
    }
    
    public String getName() {
        return name;
    }
    
    @Override
    public String toString() {
        return "ID: " + id + ", Name: " + name + ", Age: " + age;
    }
}