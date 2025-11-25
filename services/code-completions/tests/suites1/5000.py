#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def add(a, b):
    return a + b

def subtract(a, b):
    return a - b

def multiply(a, b):
    return a * b

def divide(a, b):
    if b != 0:
        return a / b
    else:
        return "Error: Division by zero"

def power(base, exp):
    return base ** exp

def factorial(n):
    if n == 0 or n == 1:
        return 1
    else:
        return n * factorial(n - 1)

def fibonacci(n):
    if n <= 1:
        return n
    else:
        return fibonacci(n - 1) + fibonacci(n - 2)

def is_prime(n):
    if n <= 1:
        return False
    for i in range(2, int(n ** 0.5) + 1):
        if n % i == 0:
            return False
    return True

def gcd(a, b):
    if b == 0:
        return a
    return gcd(b, a % b)

def square_root(n):
    if n < 0:
        return "Error: Negative number"
    return n ** 0.5

def absolute_value(n):
    return abs(n)

def average(numbers):
    if not numbers:
        return 0
    return sum(numbers) / len(numbers)

def sum_list(numbers):
    return sum(numbers)

def reverse_list(numbers):
    return numbers[::-1]

def sort_list(numbers, reverse=False):
    return sorted(numbers, reverse=reverse)

def binary_search(arr, target):
    left, right = 0, len(arr) - 1
    while left <= right:
        mid = (left + right) // 2
        if arr[mid] == target:
            return mid
        elif arr[mid] < target:
            left = mid + 1
        else:
            right = mid - 1
    return -1

def quick_sort(arr):
    if len(arr) <= 1:
        return arr
    pivot = arr[len(arr) // 2]
    left = [x for x in arr if x < pivot]
    middle = [x for x in arr if x == pivot]
    right = [x for x in arr if x > pivot]
    return quick_sort(left) + middle + quick_sort(right)

def is_even(n):
    return n % 2 == 0

def is_palindrome(s):
    return s == s[::-1]

def count_vowels(s):
    vowels = "aeiouAEIOU"
    return sum(1 for char in s if char in vowels)

def decimal_to_binary(n):
    if n == 0:
        return "0"
    binary = ""
    num = abs(n)
    while num > 0:
        binary = str(num % 2) + binary
        num = num // 2
    return binary if n >= 0 else "-" + binary

def decimal_to_hex(n):
    if n == 0:
        return "0"
    hex_chars = "0123456789ABCDEF"
    hex_str = ""
    num = abs(n)
    while num > 0:
        hex_str = hex_chars[num % 16] + hex_str
        num = num // 16
    return hex_str if n >= 0 else "-" + hex_str

def decimal_to_octal(n):
    if n == 0:
        return "0"
    octal_str = ""
    num = abs(n)
    while num > 0:
        octal_str = str(num % 8) + octal_str
        num = num // 8
    return octal_str if n >= 0 else "-" + octal_str

def is_leap_year(year):
    if year % 4 != 0:
        return False
    elif year % 100 != 0:
        return True
    else:
        return year % 400 == 0

def days_in_month(year, month):
    if month < 1 or month > 12:
        return 0
    if month in [4, 6, 9, 11]:
        return 30
    elif month == 2:
        return 29 if is_leap_year(year) else 28
    else:
        return 31

def distance_between_points(x1, y1, x2, y2):
    return ((x2 - x1) ** 2 + (y2 - y1) ** 2) ** 0.5

def area_of_circle(radius):
    if radius < 0:
        return "Error: Negative radius"
    return 3.14159 * radius ** 2

def circumference_of_circle(radius):
    if radius < 0:
        return "Error: Negative radius"
    return 2 * 3.14159 * radius

def area_of_rectangle(length, width):
    if length < 0 or width < 0:
        return "Error: Negative dimension"
    return length * width

def perimeter_of_rectangle(length, width):
    if length < 0 or width < 0:
        return "Error: Negative dimension"
    return 2 * (length + width)

def area_of_triangle(base, height):
    if base < 0 or height < 0:
        return "Error: Negative dimension"
    return 0.5 * base * height

def perimeter_of_triangle(side1, side2, side3):
    if side1 < 0 or side2 < 0 or side3 < 0:
        return "Error: Negative dimension"
    return side1 + side2 + side3

def volume_of_sphere(radius):
    if radius < 0:
        return "Error: Negative radius"
    return (4/3) * 3.14159 * radius ** 3

def surface_area_of_sphere(radius):
    if radius < 0:
        return "Error: Negative radius"
    return 4 * 3.14159 * radius ** 2

def volume_of_cylinder(radius, height):
    if radius < 0 or height < 0:
        return "Error: Negative dimension"
    return 3.14159 * radius ** 2 * height

def surface_area_of_cylinder(radius, height):
    if radius < 0 or height < 0:
        return "Error: Negative dimension"
    return 2 * 3.14159 * radius * (radius + height)

def volume_of_cone(radius, height):
    if radius < 0 or height < 0:
        return "Error: Negative dimension"
    return (1/3) * 3.14159 * radius ** 2 * height

def surface_area_of_cone(radius, height):
    if radius < 0 or height < 0:
        return "Error: Negative dimension"
    slant_height = (radius ** 2 + height ** 2) ** 0.5
    return 3.14159 * radius * (radius + slant_height)

def matrix_addition(matrix1, matrix2):
    if len(matrix1) != len(matrix2) or len(matrix1[0]) != len(matrix2[0]):
        return "Error: Matrices must have same dimensions"
    result = []
    for i in range(len(matrix1)):
        row = []
        for j in range(len(matrix1[0])):
            row.append(matrix1[i][j] + matrix2[i][j])
        result.append(row)
    return result

def matrix_subtraction(matrix1, matrix2):
    if len(matrix1) != len(matrix2) or len(matrix1[0]) != len(matrix2[0]):
        return "Error: Matrices must have same dimensions"
    result = []
    for i in range(len(matrix1)):
        row = []
        for j in range(len(matrix1[0])):
            row.append(matrix1[i][j] - matrix2[i][j])
        result.append(row)
    return result

def matrix_multiplication(matrix1, matrix2):
    if len(matrix1[0]) != len(matrix2):
        return "Error: Number of columns in first matrix must equal number of rows in second matrix"
    result = []
    for i in range(len(matrix1)):
        row = []
        for j in range(len(matrix2[0])):
            sum_val = 0
            for k in range(len(matrix2)):
                sum_val += matrix1[i][k] * matrix2[k][j]
            row.append(sum_val)
        result.append(row)
    return result

def matrix_transpose(matrix):
    result = []
    for j in range(len(matrix[0])):
        row = []
        for i in range(len(matrix)):
            row.append(matrix[i][j])
        result.append(row)
    return result

def matrix_determinant(matrix):
    if len(matrix) != len(matrix[0]):
        return "Error: Matrix must be square"
    n = len(matrix)
    if n == 1:
        return matrix[0][0]
    if n == 2:
        return matrix[0][0] * matrix[1][1] - matrix[0][1] * matrix[1][0]
    
    det = 0
    for j in range(n):
        minor = [row[:j] + row[j+1:] for row in matrix[1:]]
        det += ((-1) ** j) * matrix[0][j] * matrix_determinant(minor)
    return det

def main():
    # Basic operations
    print(f"5 + 3 = {add(5, 3)}")
    print(f"5 - 3 = {subtract(5, 3)}")
    print(f"5 * 3 = {multiply(5, 3)}")
    print(f"6 / 3 = {divide(6, 3)}")
    print(f"5 / 0 = {divide(5, 0)}")
    print(f"2 ^ 3 = {power(2, 3)}")
    
    # Math functions
    print(f"5! = {factorial(5)}")
    print(f"Fibonacci(7) = {fibonacci(7)}")
    print(f"Is 17 prime? {is_prime(17)}")
    print(f"GCD of 48 and 18 = {gcd(48, 18)}")
    print(f"Square root of 16 = {square_root(16)}")
    print(f"Absolute value of -10 = {absolute_value(-10)}")
    
    # List operations
    print(f"Average of [1, 2, 3, 4, 5] = {average([1, 2, 3, 4, 5])}")
    print(f"Sum of [1, 2, 3, 4, 5] = {sum_list([1, 2, 3, 4, 5])}")
    print(f"Reverse of [1, 2, 3, 4, 5] = {reverse_list([1, 2, 3, 4, 5])}")
    print(f"Sorted [3, 1, 4, 2, 5] = {sort_list([3, 1, 4, 2, 5])}")
    
    # Search and sort
    print(f"Binary search for 3 in [1, 2, 3, 4, 5] = {binary_search([1, 2, 3, 4, 5], 3)}")
    print(f"Quick sort [3, 1, 4, 2, 5] = {quick_sort([3, 1, 4, 2, 5])}")
    
    # String operations
    print(f"Is 10 even? {is_even(10)}")
    print(f"Is 'madam' a palindrome? {is_palindrome('madam')}")
    print(f"Vowels in 'hello world' = {count_vowels('hello world')}")
    
    # Conversions
    print(f"Decimal 42 to binary = {decimal_to_binary(42)}")
    print(f"Decimal 255 to hex = {decimal_to_hex(255)}")
    print(f"Decimal 65 to octal = {decimal_to_octal(65)}")
    
    # Date and geometry
    print(f"Is 2020 a leap year? {is_leap_year(2020)}")
    print(f"Days in February 2020 = {days_in_month(2020, 2)}")
    print(f"Distance between (1,2) and (4,6) = {distance_between_points(1, 2, 4, 6)}")
    print(f"Area of circle with radius 5 = {area_of_circle(5)}")
    print(f"Circumference of circle with radius 5 = {circumference_of_circle(5)}")
    print(f"Area of rectangle (5,10) = {area_of_rectangle(5, 10)}")
    print(f"Perimeter of rectangle (5,10) = {perimeter_of_rectangle(5, 10)}")
    print(f"Area of triangle with base 10 and height 5 = {area_of_triangle(10, 5)}")
    print(f"Perimeter of triangle with sides 3,4,5 = {perimeter_of_triangle(3, 4, 5)}")
    print(f"Volume of sphere with radius 5 = {volume_of_sphere(5)}")
    print(f"Surface area of sphere with radius 5 = {surface_area_of_sphere(5)}")
    print(f"Volume of cylinder (3,10) = {volume_of_cylinder(3, 10)}")
    print(f"Surface area of cylinder (3,10) = {surface_area_of_cylinder(3, 10)}")
    print(f"Volume of cone (3,10) = {volume_of_cone(3, 10)}")
    print(f"Surface area of cone (3,10) = {surface_area_of_cone(3, 10)}")
    
    # Matrix operations
    m1 = [[1, 2], [3, 4]]
    m2 = [[5, 6], [7, 8]]
    print(f"Matrix addition: {matrix_addition(m1, m2)}")
    print(f"Matrix subtraction: {matrix_subtraction(m1, m2)}")
    print(f"Matrix multiplication: {matrix_multiplication(m1, m2)}")
    print(f"Matrix transpose: {matrix_transpose(m1)}")
    print(f"Matrix determinant: {matrix_determinant(m1)}")

if __name__ == "__main__":
    main()