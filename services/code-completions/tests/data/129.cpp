#include <iostream>
#include <vector>
#include <numeric>

int main() {
    std::vector<int> numbers = {1, 2, 3, 4, 5};
    int sum = 0;
    
    // Calculate sum using accumulate
    <｜fim▁hole｜>sum = std::accumulate(numbers.begin(), numbers.end(), 0);
    
    std::cout << "Numbers: ";
    for (int num : numbers) {
        std::cout << num << " ";
    }
    std::cout << std::endl;
    
    std::cout << "Sum: " << sum << std::endl;
    
    return 0;
}