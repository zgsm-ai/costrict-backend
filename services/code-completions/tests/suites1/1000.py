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

def lcm(a, b):
    return a * b // gcd(a, b)

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

def mode(numbers):
    if not numbers:
        return 0
    counts = {}
    for num in numbers:
        counts[num] = counts.get(num, 0) + 1
    max_count = max(counts.values())
    for num, count in counts.items():
        if count == max_count:
            return num
    return 0

def sum_list(numbers):
    total = 0
    for num in numbers:
        total += num
    return total

def product_list(numbers):
    product = 1
    for num in numbers:
        product *= num
    return product

def reverse_list(numbers):
    return numbers[::-1]

def sort_list(numbers, reverse=False):
    return sorted(numbers, reverse=reverse)

def count_occurrences(numbers, target):
    count = 0
    for num in numbers:
        if num == target:
            count += 1
    return count

def find_max(numbers):
    if not numbers:
        return None
    max_num = numbers[0]
    for num in numbers:
        if num > max_num:
            max_num = num
    return max_num

def find_min(numbers):
    if not numbers:
        return None
    min_num = numbers[0]
    for num in numbers:
        if num < min_num:
            min_num = num
    return min_num

def is_even(n):
    return n % 2 == 0

def is_odd(n):
    return n % 2 != 0

def is_palindrome(s):
    return s == s[::-1]

def capitalize_words(s):
    return s.title()

def count_vowels(s):
    vowels = "aeiouAEIOU"
    <｜fim▁hole｜>
    count = 0
    for char in s:
        if char in vowels:
            count += 1
    return count

def count_consonants(s):
    consonants = "bcdfghjklmnpqrstvwxyzBCDFGHJKLMNPQRSTVWXYZ"
    count = 0
    for char in s:
        if char in consonants:
            count += 1
    return count

def main():
    result1 = add(5, 3)
    result2 = subtract(5, 3)
    result3 = multiply(5, 3)
    result4 = divide(6, 3)
    result5 = divide(5, 0)
    result6 = power(2, 3)
    result7 = factorial(5)
    result8 = fibonacci(7)
    result9 = is_prime(17)
    result10 = gcd(48, 18)
    result11 = lcm(12, 18)
    result12 = square_root(16)
    result13 = absolute_value(-10)
    result14 = maximum(10, 20)
    result15 = minimum(10, 20)
    result16 = average([1, 2, 3, 4, 5])
    result17 = median([1, 2, 3, 4, 5])
    result18 = mode([1, 2, 2, 3, 4])
    result19 = sum_list([1, 2, 3, 4, 5])
    result20 = product_list([1, 2, 3, 4])
    result21 = reverse_list([1, 2, 3, 4, 5])
    result22 = sort_list([3, 1, 4, 2, 5])
    result23 = sort_list([3, 1, 4, 2, 5], True)
    result24 = count_occurrences([1, 2, 2, 3, 2, 4], 2)
    result25 = find_max([1, 2, 3, 4, 5])
    result26 = find_min([1, 2, 3, 4, 5])
    result27 = is_even(10)
    result28 = is_odd(11)
    result29 = is_palindrome("madam")
    result30 = capitalize_words("hello world")
    result31 = count_vowels("hello world")
    result32 = count_consonants("hello world")
    
    print(f"5 + 3 = {result1}")
    print(f"5 - 3 = {result2}")
    print(f"5 * 3 = {result3}")
    print(f"6 / 3 = {result4}")
    print(f"5 / 0 = {result5}")
    print(f"2 ^ 3 = {result6}")
    print(f"5! = {result7}")
    print(f"Fibonacci(7) = {result8}")
    print(f"Is 17 prime? {result9}")
    print(f"GCD of 48 and 18 = {result10}")
    print(f"LCM of 12 and 18 = {result11}")
    print(f"Square root of 16 = {result12}")
    print(f"Absolute value of -10 = {result13}")
    print(f"Maximum of 10 and 20 = {result14}")
    print(f"Minimum of 10 and 20 = {result15}")
    print(f"Average of [1, 2, 3, 4, 5] = {result16}")
    print(f"Median of [1, 2, 3, 4, 5] = {result17}")
    print(f"Mode of [1, 2, 2, 3, 4] = {result18}")
    print(f"Sum of [1, 2, 3, 4, 5] = {result19}")
    print(f"Product of [1, 2, 3, 4] = {result20}")
    print(f"Reverse of [1, 2, 3, 4, 5] = {result21}")
    print(f"Sorted [3, 1, 4, 2, 5] = {result22}")
    print(f"Reverse sorted [3, 1, 4, 2, 5] = {result23}")
    print(f"Count of 2 in [1, 2, 2, 3, 2, 4] = {result24}")
    print(f"Max in [1, 2, 3, 4, 5] = {result25}")
    print(f"Min in [1, 2, 3, 4, 5] = {result26}")
    print(f"Is 10 even? {result27}")
    print(f"Is 11 odd? {result28}")
    print(f"Is 'madam' a palindrome? {result29}")
    print(f"Capitalize 'hello world' = {result30}")
    print(f"Vowels in 'hello world' = {result31}")
    print(f"Consonants in 'hello world' = {result32}")

if __name__ == "__main__":
    main()