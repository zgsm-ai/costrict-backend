#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define MAX_STUDENTS 50
#define MAX_COURSES 10
#define MAX_NAME_LENGTH 30
#define FILENAME "students.dat"

typedef struct {
    int id;
    char name[MAX_NAME_LENGTH];
    float scores[MAX_COURSES];
    int course_count;
    float average_score;
    char grade;
} Student;

typedef struct {
    Student students[MAX_STUDENTS];
    int count;
} StudentDatabase;

void initialize_database(StudentDatabase *db);
int add_student(StudentDatabase *db, const char *name);
int add_score(StudentDatabase *db, int student_id, int course_index, float score);
void calculate_average(StudentDatabase *db, int student_id);
void assign_grade(StudentDatabase *db, int student_id);
void display_student(const Student *student);
void display_all_students(const StudentDatabase *db);
int find_student_by_id(const StudentDatabase *db, int student_id);
int save_database(const StudentDatabase *db, const char *filename);
int load_database(StudentDatabase *db, const char *filename);
void sort_students_by_average(StudentDatabase *db);
void search_students_by_name(const StudentDatabase *db, const char *name);
void display_statistics(const StudentDatabase *db);
void delete_student(StudentDatabase *db, int student_id);

const char* course_names[MAX_COURSES] = {
    "Mathematics", "Physics", "Chemistry", "Biology", "Computer Science",
    "English", "History", "Geography", "Economics", "Psychology"
};

int main() {
    StudentDatabase db;
    initialize_database(&db);
    
    printf("Student Management System\n");
    printf("1. Add Student\n");
    printf("2. Add Score\n");
    printf("3. Display All Students\n");
    printf("4. Search by ID\n");
    printf("5. Search by Name\n");
    printf("6. Sort by Average\n");
    printf("7. Display Statistics\n");
    printf("8. Delete Student\n");
    printf("0. Exit\n");
    
    int choice;
    do {
        printf("\nEnter your choice: ");
        scanf("%d", &choice);
        
        switch (choice) {
            case 1: {
                char name[MAX_NAME_LENGTH];
                printf("Enter student name: ");
                scanf(" %[^\n]", name);
                int id = add_student(&db, name);
                if (id != -1) printf("Student added with ID: %d\n", id);
                else printf("Failed to add student. Database is full.\n");
                break;
            }
            case 2: {
                int student_id, course_index;
                float score;
                printf("Enter student ID: ");
                scanf("%d", &student_id);
                printf("Enter course index (0-9): ");
                scanf("%d", &course_index);
                printf("Enter score: ");
                scanf("%f", &score);
                if (add_score(&db, student_id, course_index, score)) {
                    printf("Score added successfully.\n");
                    calculate_average(&db, student_id);
                    assign_grade(&db, student_id);
                } else printf("Failed to add score. Student not found.\n");
                break;
            }
            case 3: display_all_students(&db); break;
            case 4: {
                int student_id;
                printf("Enter student ID: ");
                scanf("%d", &student_id);
                int index = find_student_by_id(&db, student_id);
                if (index != -1) display_student(&db.students[index]);
                else printf("Student not found.\n");
                break;
            }
            case 5: {
                char name[MAX_NAME_LENGTH];
                printf("Enter student name to search: ");
                scanf(" %[^\n]", name);
                search_students_by_name(&db, name);
                break;
            }
            case 6: {
                sort_students_by_average(&db);
                printf("Students sorted by average score.\n");
                display_all_students(&db);
                break;
            }
            case 7: display_statistics(&db); break;
            case 8: {
                int student_id;
                printf("Enter student ID to delete: ");
                scanf("%d", &student_id);
                delete_student(&db, student_id);
                break;
            }
            case 0: printf("Exiting program.\n"); break;
            default: printf("Invalid choice. Please try again.\n"); break;
        }
    } while (choice != 0);
    
    return 0;
}

void initialize_database(StudentDatabase *db) {
    db->count = 0;
}

int add_student(StudentDatabase *db, const char *name) {
    if (db->count >= MAX_STUDENTS) return -1;
    
    Student *student = &db->students[db->count];
    student->id = db->count + 1;
    strncpy(student->name, name, MAX_NAME_LENGTH - 1);
    student->name[MAX_NAME_LENGTH - 1] = '\0';
    student->course_count = 0;
    student->average_score = 0.0f;
    student->grade = 'F';
    
    for (int i = 0; i < MAX_COURSES; i++) {
        student->scores[i] = 0.0f;
    }
    
    db->count++;
    return student->id;
}

int add_score(StudentDatabase *db, int student_id, int course_index, float score) {
    int index = find_student_by_id(db, student_id);
    if (index == -1 || course_index < 0 || course_index >= MAX_COURSES) return 0;
    
    db->students[index].scores[course_index] = score;
    if (course_index + 1 > db->students[index].course_count) {
        db->students[index].course_count = course_index + 1;
    }
    return 1;
}

void calculate_average(StudentDatabase *db, int student_id) {
    int index = find_student_by_id(db, student_id);
    if (index == -1) return;
    
    Student *student = &db->students[index];
    float sum = 0.0f;
    int valid_scores = 0;
    
    for (int i = 0; i < MAX_COURSES; i++) {
        if (student->scores[i] > 0) {
            sum += student->scores[i];
            valid_scores++;
        }
    }
    
    student->average_score = valid_scores > 0 ? sum / valid_scores : 0.0f;
}

void assign_grade(StudentDatabase *db, int student_id) {
    int index = find_student_by_id(db, student_id);
    if (index == -1) return;
    
    float average = db->students[index].average_score;
    
    if (average >= 90.0f) db->students[index].grade = 'A';
    else if (average >= 80.0f) db->students[index].grade = 'B';
    else if (average >= 70.0f) db->students[index].grade = 'C';
    else if (average >= 60.0f) db->students[index].grade = 'D';
    else db->students[index].grade = 'F';
}

void display_student(const Student *student) {
    printf("\nStudent ID: %d\n", student->id);
    printf("Name: %s\n", student->name);
    printf("Average Score: %.2f\n", student->average_score);
    printf("Grade: %c\n", student->grade);
    printf("Course Scores:\n");
    
    for (int i = 0; i < MAX_COURSES; i++) {
        if (student->scores[i] > 0) {
            printf("  %s: %.2f\n", course_names[i], student->scores[i]);
        }
    }
}

void display_all_students(const StudentDatabase *db) {
    printf("\n=== All Students ===\n");
    printf("Total Students: %d\n\n", db->count);
    
    for (int i = 0; i < db->count; i++) {
        display_student(&db->students[i]);
        printf("------------------------\n");
    }
}

int find_student_by_id(const StudentDatabase *db, int student_id) {
    for (int i = 0; i < db->count; i++) {
        if (db->students[i].id == student_id) return i;
    }
    return -1;
}

int save_database(const StudentDatabase *db, const char *filename) {
    FILE *file = fopen(filename, "wb");
    if (file == NULL) return 0;
    
    fwrite(&db->count, sizeof(int), 1, file);
    fwrite(db->students, sizeof(Student), db->count, file);
    fclose(file);
    return 1;
}

int load_database(StudentDatabase *db, const char *filename) {
    FILE *file = fopen(filename, "rb");
    if (file == NULL) return 0;
    
    fread(&db->count, sizeof(int), 1, file);
    fread(db->students, sizeof(Student), db->count, file);
    fclose(file);
    return 1;
}

void sort_students_by_average(StudentDatabase *db) {
    for (int i = 0; i < db->count - 1; i++) {
        for (int j = 0; j < db->count - i - 1; j++) {
            if (db->students[j].average_score < db->students[j + 1].average_score) {
                Student temp = db->students[j];
                db->students[j] = db->students[j + 1];
                db->students[j + 1] = temp;
            }
        }
    }
}

void search_students_by_name(const StudentDatabase *db, const char *name) {
    printf("\nSearch Results for '%s':\n", name);
    int found = 0;
    
    for (int i = 0; i < db->count; i++) {
        if (strstr(db->students[i].name, name) != NULL) {
            display_student(&db->students[i]);
            printf("------------------------\n");
            found++;
        }
    }
    
    if (found == 0) {
        printf("No students found with name containing '%s'.\n", name);
    }
}

void display_statistics(const StudentDatabase *db) {
    if (db->count == 0) {
        printf("No students in database.\n");
        return;
    }
    
    float total_average = 0.0f;
    int grade_counts[5] = {0};
    
    for (int i = 0; i < db->count; i++) {
        total_average += db->students[i].average_score;
        switch (db->students[i].grade) {
            case 'A': grade_counts[0]++; break;
            case 'B': grade_counts[1]++; break;
            case 'C': grade_counts[2]++; break;
            case 'D': grade_counts[3]++; break;
            case 'F': grade_counts[4]++; break;
        }
    }
    
    printf("\n=== Statistics ===\n");
    printf("Total Students: %d\n", db->count);
    printf("Class Average: %.2f\n", total_average / db->count);
    printf("Grade Distribution:\n");
    printf("  A: %d (%.1f%%)\n", grade_counts[0], (float)grade_counts[0] / db->count * 100);
    printf("  B: %d (%.1f%%)\n", grade_counts[1], (float)grade_counts[1] / db->count * 100);
    printf("  C: %d (%.1f%%)\n", grade_counts[2], (float)grade_counts[2] / db->count * 100);
    printf("  D: %d (%.1f%%)\n", grade_counts[3], (float)grade_counts[3] / db->count * 100);
    printf("  F: %d (%.1f%%)\n", grade_counts[4], (float)grade_counts[4] / db->count * 100);
}

void delete_student(StudentDatabase *db, int student_id) {
    int index = find_student_by_id(db, student_id);
    if (index == -1) {
        printf("Student not found.\n");
        return;
    }
    
    for (int i = index; i < db->count - 1; i++) {
        db->students[i] = db->students[i + 1];
    }
    
    db->count--;
    printf("Student with ID %d deleted successfully.\n", student_id);
}