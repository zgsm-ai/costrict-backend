#include <iostream>
#include <fstream>
#include <string>
#include <vector>

class FileProcessor {
private:
    std::string filename;
    
public:
    FileProcessor(const std::string& file) : filename(file) {}
    
    std::vector<std::string> readLines() {
        std::ifstream file(filename);
        std::vector<std::string> lines;
        std::string line;
        
        if (!file.is_open()) {
            std::cerr << "Error opening file: " << filename << std::endl;
            return lines;
        }
        
        while (std::getline(file, line)) {
            <｜fim▁hole｜>lines.push_back(line);
        }
        
        file.close();
        return lines;
    }
    
    void writeLines(const std::vector<std::string>& lines) {
        std::ofstream file(filename);
        
        if (!file.is_open()) {
            std::cerr << "Error opening file for writing: " << filename << std::endl;
            return;
        }
        
        for (const auto& line : lines) {
            file << line << std::endl;
        }
        
        file.close();
    }
};

int main() {
    FileProcessor processor("example.txt");
    
    std::vector<std::string> linesToWrite = {
        "Hello, World!",
        "This is a test file."
    };
    
    processor.writeLines(linesToWrite);
    
    std::vector<std::string> linesRead = processor.readLines();
    std::cout << "File contents:" << std::endl;
    for (const auto& line : linesRead) {
        std::cout << line << std::endl;
    }
    
    return 0;
}