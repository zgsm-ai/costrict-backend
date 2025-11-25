#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include <stdbool.h>

#define MAX_EMPLOYEES 100
#define MAX_NAME_LENGTH 50
#define MAX_DEPARTMENT_LENGTH 30

typedef enum {
    MONDAY,
    TUESDAY,
    WEDNESDAY,
    THURSDAY,
    FRIDAY,
    SATURDAY,
    SUNDAY
} DayOfWeek;

typedef struct {
    int day;
    int month;
    int year;
} Date;

typedef struct {
    DayOfWeek day;
    float hoursWorked;
    bool isHoliday;
} WorkDay;

typedef struct {
    int employeeId;
    char name[MAX_NAME_LENGTH];
    char department[MAX_DEPARTMENT_LENGTH];
    float hourlyWage;
    Date hireDate;
    WorkDay weeklySchedule[7];
    float totalHoursWorked;
    float monthlySalary;
} Employee;

typedef struct {
    Employee employees[MAX_EMPLOYEES];
    int count;
} Company;

const char* dayToString(DayOfWeek day) {
    switch (day) {
        case MONDAY: return "Monday";
        case TUESDAY: return "Tuesday";
        case WEDNESDAY: return "Wednesday";
        case THURSDAY: return "Thursday";
        case FRIDAY: return "Friday";
        case SATURDAY: return "Saturday";
        case SUNDAY: return "Sunday";
        default: return "Unknown";
    }
}

int daysBetweenDates(const Date* date1, const Date* date2) {
    int days1 = date1->year * 365 + date1->month * 30 + date1->day;
    int days2 = date2->year * 365 + date2->month * 30 + date2->day;
    return abs(days2 - days1);
}

void initCompany(Company* company) {
    company->count = 0;
}

int addEmployee(Company* company, const char* name, const char* department,
               float hourlyWage, const Date* hireDate) {
    if (company->count >= MAX_EMPLOYEES) return 0;
    
    Employee* emp = &company->employees[company->count];
    emp->employeeId = 1000 + company->count;
    strncpy(emp->name, name, MAX_NAME_LENGTH - 1);
    emp->name[MAX_NAME_LENGTH - 1] = '\0';
    
    strncpy(emp->department, department, MAX_DEPARTMENT_LENGTH - 1);
    emp->department[MAX_DEPARTMENT_LENGTH - 1] = '\0';
    
    emp->hourlyWage = hourlyWage;
    emp->hireDate = *hireDate;
    emp->totalHoursWorked = 0.0f;
    emp->monthlySalary = 0.0f;
    
    for (int i = 0; i < 7; i++) {
        emp->weeklySchedule[i].day = i;
        emp->weeklySchedule[i].isHoliday = false;
        emp->weeklySchedule[i].hoursWorked = (i < 5) ? 8.0f : 0.0f;
    }
    
    company->count++;
    return 1;
}

void updateWorkSchedule(Employee* emp, DayOfWeek day, float hours, bool isHoliday) {
    emp->weeklySchedule[day].hoursWorked = hours;
    emp->weeklySchedule[day].isHoliday = isHoliday;
}

float calculateMonthlySalary(Employee* emp) {
    float weeklyHours = 0.0f;
    for (int i = 0; i < 7; i++) {
        weeklyHours += emp->weeklySchedule[i].hoursWorked;
    }
    return weeklyHours * 4.33f * emp->hourlyWage;
}

Employee* findEmployeeById(Company* company, int employeeId) {
    for (int i = 0; i < company->count; i++) {
        if (company->employees[i].employeeId == employeeId) {
            return &company->employees[i];
        }
    }
    return NULL;
}

int findEmployeesByDepartment(Company* company, const char* department,
                             Employee* results, int maxResults) {
    int count = 0;
    for (int i = 0; i < company->count && count < maxResults; i++) {
        if (strcmp(company->employees[i].department, department) == 0) {
            results[count] = company->employees[i];
            count++;
        }
    }
    return count;
}

float calculateYearsOfService(const Employee* emp, const Date* currentDate) {
    int daysWorked = daysBetweenDates(&emp->hireDate, currentDate);
    return daysWorked / 365.25f;
}

void printEmployee(const Employee* emp, const Date* currentDate) {
    printf("Employee ID: %d\n", emp->employeeId);
    printf("Name: %s\n", emp->name);
    printf("Department: %s\n", emp->department);
    printf("Hourly Wage: $%.2f\n", emp->hourlyWage);
    printf("Hire Date: %d/%d/%d\n", emp->hireDate.day, emp->hireDate.month, emp->hireDate.year);
    printf("Years of Service: %.1f\n", calculateYearsOfService(emp, currentDate));
    printf("Weekly Schedule:\n");
    
    for (int i = 0; i < 7; i++) {
        printf("  %s: %.1f hours", dayToString(emp->weeklySchedule[i].day),
               emp->weeklySchedule[i].hoursWorked);
        if (emp->weeklySchedule[i].isHoliday) printf(" (Holiday)");
        printf("\n");
    }
    
    printf("Monthly Salary: $%.2f\n", calculateMonthlySalary(emp));
}

void printAllEmployees(const Company* company, const Date* currentDate) {
    for (int i = 0; i < company->count; i++) {
        printf("--- Employee %d ---\n", i + 1);
        printEmployee(&company->employees[i], currentDate);
        printf("\n");
    }
}

float calculateTotalPayroll(Company* company) {
    float total = 0.0f;
    for (int i = 0; i < company->count; i++) {
        total += calculateMonthlySalary(&company->employees[i]);
    }
    return total;
}

float calculateAverageHourlyWage(const Company* company) {
    if (company->count == 0) return 0.0f;
    
    float total = 0.0f;
    for (int i = 0; i < company->count; i++) {
        total += company->employees[i].hourlyWage;
    }
    
    return total / company->count;
}

int main() {
    Company techCompany;
    initCompany(&techCompany);
    
    Date currentDate = {12, 11, 2025};
    
    // Add employees
    Date hireDate1 = {15, 1, 2020};
    addEmployee(&techCompany, "John Smith", "Engineering", 45.50f, &hireDate1);
    
    Date hireDate2 = {10, 3, 2019};
    addEmployee(&techCompany, "Jane Doe", "Marketing", 38.75f, &hireDate2);
    
    Date hireDate3 = {5, 7, 2021};
    addEmployee(&techCompany, "Mike Johnson", "Engineering", 52.00f, &hireDate3);
    
    Date hireDate4 = {20, 11, 2022};
    addEmployee(&techCompany, "Sarah Williams", "HR", 41.25f, &hireDate4);
    
    // Update work schedules
    Employee* john = findEmployeeById(&techCompany, 1000);
    <｜fim▁hole｜>updateWorkSchedule(john, MONDAY, 8.0f, false);
    updateWorkSchedule(john, TUESDAY, 8.0f, false);
    updateWorkSchedule(john, WEDNESDAY, 7.5f, false);
    updateWorkSchedule(john, THURSDAY, 8.5f, false);
    updateWorkSchedule(john, FRIDAY, 8.0f, false);
    updateWorkSchedule(john, SATURDAY, 0.0f, false);
    updateWorkSchedule(john, SUNDAY, 0.0f, false);
    
    Employee* jane = findEmployeeById(&techCompany, 1001);
    updateWorkSchedule(jane, MONDAY, 9.0f, false);
    updateWorkSchedule(jane, TUESDAY, 8.0f, false);
    updateWorkSchedule(jane, WEDNESDAY, 8.0f, false);
    updateWorkSchedule(jane, THURSDAY, 8.0f, false);
    updateWorkSchedule(jane, FRIDAY, 7.0f, false);
    updateWorkSchedule(jane, SATURDAY, 4.0f, false);
    updateWorkSchedule(jane, SUNDAY, 0.0f, false);
    
    // Print all employees
    printAllEmployees(&techCompany, &currentDate);
    
    // Print company statistics
    printf("Company Statistics:\n");
    printf("Total Employees: %d\n", techCompany.count);
    printf("Total Monthly Payroll: $%.2f\n", calculateTotalPayroll(&techCompany));
    printf("Average Hourly Wage: $%.2f\n", calculateAverageHourlyWage(&techCompany));
    
    return 0;
}