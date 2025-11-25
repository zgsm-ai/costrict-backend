def process_list(items):
    result = []
    for item in items:
        <｜fim▁hole｜>if item % 2 == 0:
            result.append(item * 2)
    return result

def main():
    numbers = [1, 2, 3, 4, 5]
    processed = process_list(numbers)
    print(f"Original: {numbers}")
    print(f"Processed: {processed}")

if __name__ == "__main__":
    main()