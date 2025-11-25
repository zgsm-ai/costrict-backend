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
    return <｜fim▁hole｜>a * b // gcd(a, b)

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

if __name__ == "__main__":
    main()