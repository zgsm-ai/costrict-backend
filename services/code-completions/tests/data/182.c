#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>

#define MAX_STUDENTS 100
#define MAX_NAME_LENGTH 50
#define MAX_COURSES 10

typedef struct {
    int id;
    char name[MAX_NAME_LENGTH];
    float grades[MAX_COURSES];
    int courseCount;
    float average;
} Student;

typedef struct {
    Student students[MAX_STUDENTS];
    int count;
} StudentDatabase;

/**
 * Initialize a new student database
 * @param db Pointer to the database to initialize
 */
void initDatabase(StudentDatabase* db) {
    db->count = 0;
}

/**
 * Add a new student to the database
 * @param db Pointer to the database
 * @param name Student name
 * @param id Student ID
 * @return 1 if successful, 0 if database is full
 */
int addStudent(StudentDatabase* db, const char* name, int id) {
    if (db->count >= MAX_STUDENTS) {
        return 0;
    }
    
    Student* student = &db->students[db->count];
    student->id = id;
    strncpy(student->name, name, MAX_NAME_LENGTH - 1);
    student->name[MAX_NAME_LENGTH - 1] = '\0';
    student->courseCount = 0;
    student->average = 0.0f;
    
    db->count++;
    return 1;
}

/**
 * Add a grade for a student
 * @param student Pointer to the student
 * @param grade Grade value (0-100)
 * @return 1 if successful, 0 if course limit reached
 */
int addGrade(Student* student, float grade) {
    if (student->courseCount >= MAX_COURSES) {
        return 0;
    }
    
    if (grade < 0 || grade > 100) {
        return 0;
    }
    
    student->grades[student->courseCount] = grade;
    student->courseCount++;
    
    // Calculate average
    float sum = 0.0f;
    for (int i = 0; i < student->courseCount; i++) {
        sum += student->grades[i];
    }
    student->average = sum / student->courseCount;
    
    return 1;
}

/**
 * Find a student by ID
 * @param db Pointer to the database
 * @param id Student ID to find
 * @return Pointer to student if found, NULL otherwise
 */
Student* findStudentById(StudentDatabase* db, int id) {
    for (int i = 0; i < db->count; i++) {
        if (db->students[i].id == id) {
            return &db->students[i];
        }
    }
    return NULL;
}

/**
 * Print student information
 * @param student Pointer to the student
 */
void printStudent(const Student* student) {
    printf("ID: %d\n", student->id);
    printf("Name: %s\n", student->name);
    printf("Grades: ");
    for (int i = 0; i < student->courseCount; i++) {
        printf("%.1f ", student->grades[i]);
    }
    printf("\nAverage: %.2f\n", student->average);
}

/**
 * Print all students in the database
 * @param db Pointer to the database
 */
void printAllStudents(const StudentDatabase* db) {
    for (int i = 0; i < db->count; i++) {
        printf("--- Student %d ---\n", i + 1);
        printStudent(&db->students[i]);
        printf("\n");
    }
}

/**
 * Calculate class average
 * @param db Pointer to the database
 * @return Class average grade
 */
float calculateClassAverage(const StudentDatabase* db) {
    if (db->count == 0) {
        return 0.0f;
    }
    
    float totalSum = 0.0f;
    for (int i = 0; i < db->count; i++) {
        totalSum += db->students[i].average;
    }
    
    return totalSum / db->count;
}

/**
 * Find the student with highest average grade
 * @param db Pointer to the database
 * @return Pointer to student with highest average
 */
Student* findTopStudent(StudentDatabase* db) {
    if (db->count == 0) {
        return NULL;
    }
    
    Student* topStudent = &db->students[0];
    for (int i = 1; i < db->count; i++) {
        if (db->students[i].average > topStudent->average) {
            topStudent = &db->students[i];
        }
    }
    
    return topStudent;
}

int main() {
    StudentDatabase database;
    initDatabase(&database);
    
    // Add students
    addStudent(&database, "Alice Johnson", 1001);
    addStudent(&database, "Bob Smith", 1002);
    addStudent(&database, "Charlie Brown", 1003);
    
    // Add grades for first student
    Student* alice = findStudentById(&database, 1001);
    <｜fim▁hole｜>addGrade(alice, 85.5f);
    addGrade(alice, 92.0f);
    addGrade(alice, 78.5f);
    
    // Add grades for second student
    Student* bob = findStudentById(&database, 1002);
    addGrade(bob, 76.0f);
    addGrade(bob, 88.5f);
    addGrade(bob, 91.0f);
    
    // Add grades for third student
    Student* charlie = findStudentById(&database, 1003);
    addGrade(charlie, 95.0f);
    addGrade(charlie, 87.5f);
    addGrade(charlie, 82.0f);
    
    // Print all students
    printAllStudents(&database);
    
    // Print class statistics
    printf("Class Average: %.2f\n", calculateClassAverage(&database));
    
    Student* top = findTopStudent(&database);
    if (top) {
        printf("Top Student: %s with average %.2f\n", top->name, top->average);
    }
    
    return 0;
}