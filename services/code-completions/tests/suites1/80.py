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
        return <｜fim▁hole｜>"Error: Division by zero"

def main():
    result1 = add(5, 3)
    result2 = subtract(5, 3)
    result3 = multiply(5, 3)
    result4 = divide(6, 3)
    result5 = divide(5, 0)
    print(f"5 + 3 = {result1}")
    print(f"5 - 3 = {result2}")
    print(f"5 * 3 = {result3}")
    print(f"6 / 3 = {result4}")
    print(f"5 / 0 = {result5}")

if __name__ == "__main__":
    main()