#include <stdio.h>
#include <math.h>

double calculateDistance(double x1, double y1, double x2, double y2) {
    double dx = x2 - x1;
    double dy = y2 - y1;
    
    <｜fim▁hole｜>return sqrt(dx * dx + dy * dy);
}

int main() {
    double x1 = 1.0, y1 = 2.0;
    double x2 = 4.0, y2 = 6.0;
    
    double distance = calculateDistance(x1, y1, x2, y2);
    
    printf("Distance between (%.1f, %.1f) and (%.1f, %.1f) is %.2f\n", 
           x1, y1, x2, y2, distance);
    
    return 0;
}