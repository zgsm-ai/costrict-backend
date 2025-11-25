def read_file_content(filename):
    try:
        with open(filename, 'r') as file:
            content = file.read()
        return content
    except FileNotFoundError:
        return None

def main():
    filename = "example.txt"
    <｜fim▁hole｜>
    if content:
        print(f"Content of {filename}:")
        print(content)
    else:
        print(f"File {filename} not found.")

if __name__ == "__main__":
    main()