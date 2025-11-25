#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>

#define MAX_PRODUCTS 200
#define MAX_NAME_LENGTH 60
#define MAX_CATEGORY_LENGTH 30
#define MAX_SUPPLIER_LENGTH 40

typedef struct {
    int day;
    int month;
    int year;
} Date;

typedef struct {
    char name[MAX_NAME_LENGTH];
    char category[MAX_CATEGORY_LENGTH];
    char supplier[MAX_SUPPLIER_LENGTH];
    int productId;
    float price;
    int stockQuantity;
    int minStockLevel;
    Date lastRestockDate;
    int timesSold;
} Product;

typedef struct {
    Product products[MAX_PRODUCTS];
    int count;
} Inventory;

void initInventory(Inventory* inv) {
    inv->count = 0;
}

int addProduct(Inventory* inv, const char* name, const char* category, const char* supplier,
              int productId, float price, int stockQuantity, int minStockLevel,
              const Date* lastRestockDate) {
    if (inv->count >= MAX_PRODUCTS) return 0;
    
    Product* product = &inv->products[inv->count];
    strncpy(product->name, name, MAX_NAME_LENGTH - 1);
    product->name[MAX_NAME_LENGTH - 1] = '\0';
    
    strncpy(product->category, category, MAX_CATEGORY_LENGTH - 1);
    product->category[MAX_CATEGORY_LENGTH - 1] = '\0';
    
    strncpy(product->supplier, supplier, MAX_SUPPLIER_LENGTH - 1);
    product->supplier[MAX_SUPPLIER_LENGTH - 1] = '\0';
    
    product->productId = productId;
    product->price = price;
    product->stockQuantity = stockQuantity;
    product->minStockLevel = minStockLevel;
    product->lastRestockDate = *lastRestockDate;
    product->timesSold = 0;
    
    inv->count++;
    return 1;
}

Product* findProductById(Inventory* inv, int productId) {
    for (int i = 0; i < inv->count; i++) {
        if (inv->products[i].productId == productId) {
            return &inv->products[i];
        }
    }
    return NULL;
}

int findProductsByCategory(Inventory* inv, const char* category,
                          Product* results, int maxResults) {
    int count = 0;
    for (int i = 0; i < inv->count && count < maxResults; i++) {
        if (strcmp(inv->products[i].category, category) == 0) {
            results[count] = inv->products[i];
            count++;
        }
    }
    return count;
}

int findProductsBySupplier(Inventory* inv, const char* supplier,
                          Product* results, int maxResults) {
    int count = 0;
    for (int i = 0; i < inv->count && count < maxResults; i++) {
        if (strcmp(inv->products[i].supplier, supplier) == 0) {
            results[count] = inv->products[i];
            count++;
        }
    }
    return count;
}

int sellProduct(Inventory* inv, int productId, int quantity, const Date* currentDate) {
    Product* product = findProductById(inv, productId);
    if (product == NULL || product->stockQuantity < quantity) return 0;
    
    product->stockQuantity -= quantity;
    product->timesSold += quantity;
    
    if (product->stockQuantity < product->minStockLevel) {
        printf("WARNING: Product %s (%d) needs restocking! Current stock: %d\n",
               product->name, productId, product->stockQuantity);
    }
    
    return 1;
}

int restockProduct(Inventory* inv, int productId, int quantity, const Date* restockDate) {
    Product* product = findProductById(inv, productId);
    if (product == NULL) return 0;
    
    product->stockQuantity += quantity;
    product->lastRestockDate = *restockDate;
    
    return 1;
}

int updateProductPrice(Inventory* inv, int productId, float newPrice) {
    Product* product = findProductById(inv, productId);
    if (product == NULL) return 0;
    
    product->price = newPrice;
    return 1;
}

float calculateTotalInventoryValue(const Inventory* inv) {
    float total = 0.0f;
    for (int i = 0; i < inv->count; i++) {
        total += inv->products[i].price * inv->products[i].stockQuantity;
    }
    return total;
}

int findProductsNeedingRestock(Inventory* inv, Product* results, int maxResults) {
    int count = 0;
    for (int i = 0; i < inv->count && count < maxResults; i++) {
        if (inv->products[i].stockQuantity < inv->products[i].minStockLevel) {
            results[count] = inv->products[i];
            count++;
        }
    }
    return count;
}

Product* findMostSoldProduct(Inventory* inv) {
    if (inv->count == 0) return NULL;
    
    Product* mostSold = &inv->products[0];
    for (int i = 1; i < inv->count; i++) {
        if (inv->products[i].timesSold > mostSold->timesSold) {
            mostSold = &inv->products[i];
        }
    }
    
    return mostSold;
}

void printProduct(const Product* product) {
    printf("Product ID: %d\n", product->productId);
    printf("Name: %s\n", product->name);
    printf("Category: %s\n", product->category);
    printf("Supplier: %s\n", product->supplier);
    printf("Price: $%.2f\n", product->price);
    printf("Stock Quantity: %d\n", product->stockQuantity);
    printf("Minimum Stock Level: %d\n", product->minStockLevel);
    printf("Last Restock Date: %d/%d/%d\n",
           product->lastRestockDate.day,
           product->lastRestockDate.month,
           product->lastRestockDate.year);
    printf("Times Sold: %d\n", product->timesSold);
    
    if (product->stockQuantity < product->minStockLevel) {
        printf("STATUS: NEEDS RESTOCKING\n");
    } else {
        printf("STATUS: IN STOCK\n");
    }
}

void printAllProducts(const Inventory* inv) {
    for (int i = 0; i < inv->count; i++) {
        printf("--- Product %d ---\n", i + 1);
        printProduct(&inv->products[i]);
        printf("\n");
    }
}

int main() {
    Inventory storeInventory;
    initInventory(&storeInventory);
    
    Date today = {12, 11, 2025};
    Date lastWeek = {5, 11, 2025};
    Date lastMonth = {15, 10, 2025};
    
    // Add products
    addProduct(&storeInventory, "Laptop Computer", "Electronics", "TechSupply", 1001, 899.99f, 25, 10, &lastWeek);
    addProduct(&storeInventory, "Wireless Mouse", "Electronics", "TechSupply", 1002, 29.99f, 50, 15, &lastWeek);
    addProduct(&storeInventory, "Office Chair", "Furniture", "ComfortFurnishings", 2001, 149.99f, 12, 5, &lastMonth);
    addProduct(&storeInventory, "Desk Lamp", "Furniture", "ComfortFurnishings", 2002, 39.99f, 30, 10, &lastMonth);
    addProduct(&storeInventory, "Notebook Set", "Stationery", "PaperGoods", 3001, 12.99f, 100, 25, &lastWeek);
    addProduct(&storeInventory, "Pen Set", "Stationery", "PaperGoods", 3002, 8.99f, 150, 30, &lastWeek);
    
    // Simulate sales
    sellProduct(&storeInventory, 1001, 5, &today);
    sellProduct(&storeInventory, 1002, 20, &today);
    sellProduct(&storeInventory, 2001, 4, &today);
    sellProduct(&storeInventory, 3001, 30, &today);
    sellProduct(&storeInventory, 3002, 60, &today);
    
    // Restock some products
    <｜fim▁hole｜>restockProduct(&storeInventory, 1002, 25, &today);
    restockProduct(&storeInventory, 3002, 50, &today);
    
    // Update prices
    updateProductPrice(&storeInventory, 1001, 849.99f);
    updateProductPrice(&storeInventory, 2001, 139.99f);
    
    // Print all products
    printAllProducts(&storeInventory);
    
    // Print inventory statistics
    printf("Inventory Statistics:\n");
    printf("Total Products: %d\n", storeInventory.count);
    printf("Total Inventory Value: $%.2f\n", calculateTotalInventoryValue(&storeInventory));
    
    Product* mostSold = findMostSoldProduct(&storeInventory);
    if (mostSold) {
        printf("Most Sold Product: %s (%d units)\n", mostSold->name, mostSold->timesSold);
    }
    
    // Find products needing restock
    Product restockNeeded[10];
    int restockCount = findProductsNeedingRestock(&storeInventory, restockNeeded, 10);
    if (restockCount > 0) {
        printf("\nProducts Needing Restock (%d):\n", restockCount);
        for (int i = 0; i < restockCount; i++) {
            printf("- %s (ID: %d, Stock: %d)\n", 
                   restockNeeded[i].name, 
                   restockNeeded[i].productId, 
                   restockNeeded[i].stockQuantity);
        }
    }
    
    return 0;
}