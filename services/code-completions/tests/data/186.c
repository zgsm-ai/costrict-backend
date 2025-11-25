#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include <stdbool.h>

#define MAX_FLIGHTS 100
#define MAX_PASSENGERS 1000
#define MAX_AIRLINE_LENGTH 30
#define MAX_AIRPORT_LENGTH 10
#define MAX_PASSENGER_NAME 50
#define MAX_FLIGHT_NUMBER 10

typedef enum {
    ECONOMY,
    BUSINESS,
    FIRST_CLASS
} SeatClass;

typedef struct {
    int day;
    int month;
    int year;
    int hour;
    int minute;
} DateTime;

typedef struct {
    char flightNumber[MAX_FLIGHT_NUMBER];
    char airline[MAX_AIRLINE_LENGTH];
    char origin[MAX_AIRPORT_LENGTH];
    char destination[MAX_AIRPORT_LENGTH];
    DateTime departureTime;
    DateTime arrivalTime;
    int totalSeats;
    int availableSeats[3]; // Economy, Business, First Class
    float basePrice[3]; // Economy, Business, First Class
    bool isInternational;
} Flight;

typedef struct {
    int passengerId;
    char name[MAX_PASSENGER_NAME];
    char flightNumber[MAX_FLIGHT_NUMBER];
    SeatClass seatClass;
    float pricePaid;
    DateTime bookingTime;
    bool isCheckedIn;
} Booking;

typedef struct {
    Flight flights[MAX_FLIGHTS];
    Booking bookings[MAX_PASSENGERS];
    int flightCount;
    int bookingCount;
    int nextPassengerId;
} AirlineSystem;

const char* seatClassToString(SeatClass seatClass) {
    switch (seatClass) {
        case ECONOMY: return "Economy";
        case BUSINESS: return "Business";
        case FIRST_CLASS: return "First Class";
        default: return "Unknown";
    }
}

void initAirlineSystem(AirlineSystem* system) {
    system->flightCount = 0;
    system->bookingCount = 0;
    system->nextPassengerId = 1000;
}

int addFlight(AirlineSystem* system, const char* flightNumber, const char* airline,
              const char* origin, const char* destination,
              const DateTime* departureTime, const DateTime* arrivalTime,
              int totalSeats, int businessSeats, int firstClassSeats,
              float economyPrice, float businessPrice, float firstClassPrice,
              bool isInternational) {
    if (system->flightCount >= MAX_FLIGHTS) return 0;
    
    Flight* flight = &system->flights[system->flightCount];
    strncpy(flight->flightNumber, flightNumber, MAX_FLIGHT_NUMBER - 1);
    flight->flightNumber[MAX_FLIGHT_NUMBER - 1] = '\0';
    
    strncpy(flight->airline, airline, MAX_AIRLINE_LENGTH - 1);
    flight->airline[MAX_AIRLINE_LENGTH - 1] = '\0';
    
    strncpy(flight->origin, origin, MAX_AIRPORT_LENGTH - 1);
    flight->origin[MAX_AIRPORT_LENGTH - 1] = '\0';
    
    strncpy(flight->destination, destination, MAX_AIRPORT_LENGTH - 1);
    flight->destination[MAX_AIRPORT_LENGTH - 1] = '\0';
    
    flight->departureTime = *departureTime;
    flight->arrivalTime = *arrivalTime;
    flight->totalSeats = totalSeats;
    flight->availableSeats[ECONOMY] = totalSeats;
    flight->availableSeats[BUSINESS] = businessSeats;
    flight->availableSeats[FIRST_CLASS] = firstClassSeats;
    flight->basePrice[ECONOMY] = economyPrice;
    flight->basePrice[BUSINESS] = businessPrice;
    flight->basePrice[FIRST_CLASS] = firstClassPrice;
    flight->isInternational = isInternational;
    
    system->flightCount++;
    return 1;
}

Flight* findFlightByNumber(AirlineSystem* system, const char* flightNumber) {
    for (int i = 0; i < system->flightCount; i++) {
        if (strcmp(system->flights[i].flightNumber, flightNumber) == 0) {
            return &system->flights[i];
        }
    }
    return NULL;
}

int findFlightsByRoute(AirlineSystem* system, const char* origin, const char* destination,
                      Flight* results, int maxResults) {
    int count = 0;
    for (int i = 0; i < system->flightCount && count < maxResults; i++) {
        if (strcmp(system->flights[i].origin, origin) == 0 &&
            strcmp(system->flights[i].destination, destination) == 0) {
            results[count] = system->flights[i];
            count++;
        }
    }
    return count;
}

int bookFlight(AirlineSystem* system, const char* flightNumber, const char* passengerName,
               SeatClass seatClass, const DateTime* bookingTime) {
    Flight* flight = findFlightByNumber(system, flightNumber);
    if (flight == NULL || flight->availableSeats[seatClass] <= 0) return -1;
    
    if (system->bookingCount >= MAX_PASSENGERS) return -1;
    
    Booking* booking = &system->bookings[system->bookingCount];
    booking->passengerId = system->nextPassengerId++;
    strncpy(booking->name, passengerName, MAX_PASSENGER_NAME - 1);
    booking->name[MAX_PASSENGER_NAME - 1] = '\0';
    
    strncpy(booking->flightNumber, flightNumber, MAX_FLIGHT_NUMBER - 1);
    booking->flightNumber[MAX_FLIGHT_NUMBER - 1] = '\0';
    
    booking->seatClass = seatClass;
    booking->pricePaid = flight->basePrice[seatClass];
    booking->bookingTime = *bookingTime;
    booking->isCheckedIn = false;
    
    flight->availableSeats[seatClass]--;
    system->bookingCount++;
    
    return booking->passengerId;
}

int checkInPassenger(AirlineSystem* system, int passengerId) {
    for (int i = 0; i < system->bookingCount; i++) {
        if (system->bookings[i].passengerId == passengerId) {
            system->bookings[i].isCheckedIn = true;
            return 1;
        }
    }
    return 0;
}

int calculateFlightDuration(const Flight* flight) {
    int departureMinutes = flight->departureTime.day * 1440 +
                          flight->departureTime.hour * 60 +
                          flight->departureTime.minute;
    int arrivalMinutes = flight->arrivalTime.day * 1440 +
                        flight->arrivalTime.hour * 60 +
                        flight->arrivalTime.minute;
    
    return arrivalMinutes - departureMinutes;
}

float calculateTotalRevenue(const AirlineSystem* system) {
    float total = 0.0f;
    for (int i = 0; i < system->bookingCount; i++) {
        total += system->bookings[i].pricePaid;
    }
    return total;
}

const char* findMostPopularDestination(AirlineSystem* system) {
    if (system->flightCount == 0) return NULL;
    
    int counts[20] = {0};
    char destinations[20][MAX_AIRPORT_LENGTH];
    int destCount = 0;
    
    for (int i = 0; i < system->flightCount; i++) {
        const char* dest = system->flights[i].destination;
        bool found = false;
        
        for (int j = 0; j < destCount; j++) {
            if (strcmp(destinations[j], dest) == 0) {
                counts[j]++;
                found = true;
                break;
            }
        }
        
        if (!found && destCount < 20) {
            strncpy(destinations[destCount], dest, MAX_AIRPORT_LENGTH - 1);
            destinations[destCount][MAX_AIRPORT_LENGTH - 1] = '\0';
            counts[destCount] = 1;
            destCount++;
        }
    }
    
    int maxIndex = 0;
    for (int i = 1; i < destCount; i++) {
        if (counts[i] > counts[maxIndex]) {
            maxIndex = i;
        }
    }
    
    return destinations[maxIndex];
}

void printFlight(const Flight* flight) {
    printf("Flight: %s (%s)\n", flight->flightNumber, flight->airline);
    printf("Route: %s → %s\n", flight->origin, flight->destination);
    printf("Departure: %02d/%02d/%d %02d:%02d\n",
           flight->departureTime.day, flight->departureTime.month, flight->departureTime.year,
           flight->departureTime.hour, flight->departureTime.minute);
    printf("Arrival: %02d/%02d/%d %02d:%02d\n",
           flight->arrivalTime.day, flight->arrivalTime.month, flight->arrivalTime.year,
           flight->arrivalTime.hour, flight->arrivalTime.minute);
    printf("Duration: %d minutes\n", calculateFlightDuration(flight));
    printf("Type: %s\n", flight->isInternational ? "International" : "Domestic");
    
    printf("Available Seats:\n");
    printf("  Economy: %d ($%.2f)\n", flight->availableSeats[ECONOMY], flight->basePrice[ECONOMY]);
    printf("  Business: %d ($%.2f)\n", flight->availableSeats[BUSINESS], flight->basePrice[BUSINESS]);
    printf("  First Class: %d ($%.2f)\n", flight->availableSeats[FIRST_CLASS], flight->basePrice[FIRST_CLASS]);
}

void printBooking(const Booking* booking) {
    printf("Passenger ID: %d\n", booking->passengerId);
    printf("Name: %s\n", booking->name);
    printf("Flight: %s\n", booking->flightNumber);
    printf("Class: %s\n", seatClassToString(booking->seatClass));
    printf("Price Paid: $%.2f\n", booking->pricePaid);
    printf("Booking Time: %02d/%02d/%d %02d:%02d\n",
           booking->bookingTime.day, booking->bookingTime.month, booking->bookingTime.year,
           booking->bookingTime.hour, booking->bookingTime.minute);
    printf("Status: %s\n", booking->isCheckedIn ? "Checked In" : "Not Checked In");
}

int main() {
    AirlineSystem airline;
    initAirlineSystem(&airline);
    
    // Create some date times
    DateTime dep1 = {15, 12, 2025, 8, 30};
    DateTime arr1 = {15, 12, 2025, 11, 45};
    DateTime dep2 = {15, 12, 2025, 14, 15};
    DateTime arr2 = {15, 12, 2025, 18, 30};
    DateTime dep3 = {16, 12, 2025, 7, 0};
    DateTime arr3 = {16, 12, 2025, 10, 30};
    
    DateTime bookTime1 = {12, 11, 2025, 10, 30};
    DateTime bookTime2 = {12, 11, 2025, 11, 45};
    DateTime bookTime3 = {12, 11, 2025, 14, 20};
    
    // Add flights
    addFlight(&airline, "AA123", "American Airlines", "JFK", "LAX", 
              &dep1, &arr1, 150, 30, 10, 250.0f, 750.0f, 1500.0f, true);
    addFlight(&airline, "UA456", "United Airlines", "JFK", "LAX", 
              &dep2, &arr2, 120, 25, 8, 230.0f, 680.0f, 1350.0f, true);
    addFlight(&airline, "DL789", "Delta Airlines", "JFK", "BOS", 
              &dep3, &arr3, 100, 20, 6, 120.0f, 360.0f, 720.0f, false);
    
    // Print all flights
    for (int i = 0; i < airline.flightCount; i++) {
        printf("--- Flight %d ---\n", i + 1);
        printFlight(&airline.flights[i]);
        printf("\n");
    }
    
    // Book some flights
    int passenger1 = bookFlight(&airline, "AA123", "John Smith", ECONOMY, &bookTime1);
    int passenger2 = bookFlight(&airline, "AA123", "Jane Doe", BUSINESS, &bookTime2);
    int passenger3 = bookFlight(&airline, "DL789", "Mike Johnson", FIRST_CLASS, &bookTime3);
    <｜fim▁hole｜>
    // Check in some passengers
    checkInPassenger(&airline, passenger1);
    checkInPassenger(&airline, passenger3);
    
    // Print all bookings
    for (int i = 0; i < airline.bookingCount; i++) {
        printf("--- Booking %d ---\n", i + 1);
        printBooking(&airline.bookings[i]);
        printf("\n");
    }
    
    // Print system statistics
    printf("System Statistics:\n");
    printf("Total Flights: %d\n", airline.flightCount);
    printf("Total Bookings: %d\n", airline.bookingCount);
    printf("Total Revenue: $%.2f\n", calculateTotalRevenue(&airline));
    
    const char* popularDest = findMostPopularDestination(&airline);
    if (popularDest) {
        printf("Most Popular Destination: %s\n", popularDest);
    }
    
    return 0;
}