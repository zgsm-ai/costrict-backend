#include <stdio.h>

struct Point {
    int x;
    int y;
};

int main() {
    struct Point p1 = {10, 20};
    struct Point p2;
    
    p2.x = 30;
    p2.y = 40;
    
    printf("Point 1: (%d, %d)\n", p1.x, p1.y);
    printf("Point 2: (%d, %d)\n", p2.x, p2.y);<｜fim▁hole｜>
    
    return 0;
}