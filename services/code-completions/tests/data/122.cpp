#include <iostream>
#include <memory>

class Resource {
public:
    Resource() {
        std::cout << "Resource acquired" << std::endl;
    }
    
    ~Resource() {
        std::cout << "Resource released" << std::endl;
    }
    
    void doSomething() {
        std::cout << "Doing something with resource" << std::endl;
    }
};

int main() {
    {
        <｜fim▁hole｜>std::unique_ptr<Resource> resource = std::make_unique<Resource>();
        resource->doSomething();
    } // Resource is automatically released here
    
    std::cout << "Program ended" << std::endl;
    return 0;
}