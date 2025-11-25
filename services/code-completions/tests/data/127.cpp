<｜fim▁hole｜>#include <iostream>
#include <vector>
#include <algorithm>

class Calculator {
public:
    int add(int a, int b) {
        return a + b;
    }
    
    int multiply(int a, int b) {
        return a * b;
    }
};

int main() {
    Calculator calc;
    int x = 5, y = 3;
    
    std::cout << "Sum: " << calc.add(x, y) << std::endl;
    std::cout << "Product: " << calc.multiply(x, y) << std::endl;
    
    return 0;
}