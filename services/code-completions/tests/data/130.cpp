#include <iostream>
#include <fstream>
#include <string>

int main() {
    std::string filename = "example.txt";
    std::ofstream outFile(filename);
    
    if (outFile.is_open()) {
        outFile << "Hello, C++!" << std::endl;
        outFile << "This is a test file." << std::endl;
        outFile.close();
        std::cout << "File created successfully." << std::endl;
    } else {
        std::cout << "Unable to create file." << std::endl;
    }
    
    <｜fim▁hole｜>std::ifstream inFile(filename);
    if (inFile.is_open()) {
        std::string line;
        std::cout << "File content:" << std::endl;
        while (getline(inFile, line)) {
            std::cout << line << std::endl;
        }
        inFile.close();
    }
    
    return 0;
}