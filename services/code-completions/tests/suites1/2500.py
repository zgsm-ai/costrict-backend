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

def maximum(a, b):
    return a if a > b else b

def minimum(a, b):
    return a if a < b else b

def average(numbers):
    if not numbers:
        return 0
    return sum(numbers) / len(numbers)

def median(numbers):
    if not numbers:
        return 0
    sorted_numbers = sorted(numbers)
    n = len(sorted_numbers)
    if n % 2 == 0:
        return (sorted_numbers[n//2 - 1] + sorted_numbers[n//2]) / 2
    else:
        return sorted_numbers[n//2]

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

def bubble_sort(arr):
    n = len(arr)
    for i in range(n):
        for j in range(0, n - i - 1):
            if arr[j] > arr[j + 1]:
                arr[j], arr[j + 1] = arr[j + 1], arr[j]
    return arr

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
    print(f"Maximum of 10 and 20 = {maximum(10, 20)}")
    print(f"Minimum of 10 and 20 = {minimum(10, 20)}")
    
    # List operations
    print(f"Average of [1, 2, 3, 4, 5] = {average([1, 2, 3, 4, 5])}")
    print(f"Median of [1, 2, 3, 4, 5] = {median([1, 2, 3, 4, 5])}")
    print(f"Sum of [1, 2, 3, 4, 5] = {sum_list([1, 2, 3, 4, 5])}")
    print(f"Reverse of [1, 2, 3, 4, 5] = {reverse_list([1, 2, 3, 4, 5])}")
    print(f"Sorted [3, 1, 4, 2, 5] = {sort_list([3, 1, 4, 2, 5])}")
    
    # Search and sort
    print(f"Binary search for 3 in [1, 2, 3, 4, 5] = {binary_search([1, 2, 3, 4, 5], 3)}")
    print(f"Bubble sort [3, 1, 4, 2, 5] = {bubble_sort([3, 1, 4, 2, 5])}")
    print(f"Quick sort [3, 1, 4, 2, 5] = {quick_sort([3, 1, 4, 2, 5])}")
    
    # String operations
    print(f"Is 10 even? {is_even(10)}")
    print(f"Is 'madam' a palindrome? {is_palindrome('madam')}")
    print(f"Vowels in 'hello world' = {count_vowels('hello world')}")
    
    # Conversion
    print(f"Decimal 42 to binary = {decimal_to_binary(42)}")

if __name__ == "__main__":
    main()